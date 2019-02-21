// client.go
//

/*
Package flightaware implements the cl part to the FA API

 You start by creating a cl instance with your credentials passed as Config
 struct.

 	cl := flightaware.NewClient(Config)

 Then you can configure the feed type with

 	cl.SetFeed(string, []time.Time)

 You can also set a timeout time with a value in seconds

 	cl.SetTimeout(int64)

 You can add one or more different input filters:

    cl.AddInputFilter(<type>, <value>)

 where type can be one of

     FILTER_EVENT
     FILTER_AIRLINE
     FILTER_IDENT
     FILTER_AIRPORT
     FILTER_LATLONG

 The filters you specify will be checked remotely by FlightAware according to the
 documentation available at
 https://fr.flightaware.com/commercial/firehose/firehose_documentation.rvt

 You can specify output filters with using cl.AddOutputFilter(string)

 The default handler is to display all packets.  You can change the default handler
 with

 	cl.AddHandler(func([]byte)

 Last action is to start the consuming/producer loop with

 	cl.Start()

 Reading will be closed either though getting an EOF from FA or being will killed either
 manually or through the timeout value.

 You can then use

 	cl.Close()

 to properly close the reading channel.
*/
package flightaware

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
)

// Public functions

// NewClient creates new instance of the cl
func NewClient(rc Config) *FAClient {
	cl := &FAClient{
		Host:          rc,
		FeedOne:       defaultFeed,
		Filter:        defaultFilter,
		RangeT:        make([]int64, 2),
		Started:       false,
		InputFilters:  []string{},
		OutputFilters: []*regexp.Regexp{},
	}
	return cl
}

// AddHandler changes default callback
func (cl *FAClient) AddHandler(fn func([]byte)) {
	cl.FeedOne = fn
}

// SetTimer allows run of specified duration
func (cl *FAClient) SetTimer(timer int64) {
	// Sleep for fTimeout seconds then sends Interrupt
	go func() {
		time.Sleep(time.Duration(timer) * time.Second)
		if cl.Verbose {
			log.Println("Timer off, time to kill")
		}
		myself, _ := os.FindProcess(os.Getpid())
		myself.Signal(os.Interrupt)
	}()
}

// SetFeed adds a given feed
func (cl *FAClient) SetFeed(feedType string, RangeT []time.Time) error {
	// Check when -t pitr that -T is single valued
	if feedType == "pitr" {
		// Check value
		restart := RangeT[0]
		if restart.After(time.Now()) {
			return fmt.Errorf("Error: -B %v is in the future", restart)
		}
		// Store out final value
		cl.RangeT[0] = restart.Unix()
	}

	if feedType == "range" {
		// Store out final values in UNIX epoch format
		cl.RangeT[0] = RangeT[0].Unix()
		cl.RangeT[1] = RangeT[1].Unix()
	}
	return nil
}

// Start run tha who thing.
func (cl *FAClient) Start() (err error) {
	var rc = cl.Host

	// This is the main function here:
	// - starts the consumer in the background
	// - reads data from FA and send it to the consumer

	// Build the connection string
	str := rc.Site + ":" + fmt.Sprintf("%d", rc.Port)

	// Do the actual connection
	conn, err := cl.connectFA(str, true)
	if err != nil {
		log.Fatalf("Error: can not connect with %s: %v", str, err)
	}
	cl.Conn = conn

	// Starting here everything is flowing from that connection
	ch, err := cl.startWriter()
	if err != nil {
		log.Printf("Error: starting writer - %s\n", err.Error())
		return err
	}

	cl.Started = true
	// Loop over chunks of data
	sc := bufio.NewScanner(cl.Conn)
	for {
		for sc.Scan() {
			buf := sc.Text()

			if nb := len(buf); nb != 0 {
				if cl.Verbose {
					DataLog([]byte(buf), fmt.Sprintf("Sending %d bytes\n", nb))
				}
				ch <- []byte(buf)
			}
		}
		if err = sc.Err(); err != nil {
			log.Println("Error reading:", err)

			// Reconnect
			conn, err = cl.connectFA(str, false)
			cl.Conn = conn
			sc = bufio.NewScanner(cl.Conn)
		} else {
			// Got EOF
			break
		}
	}
	return nil
}

// Close properly close the TLS connection
func (cl *FAClient) Close() error {
	var err error

	if err = cl.Conn.Close(); err != nil {
		log.Println("Error closing connection " + err.Error())
	}
	if cl.Verbose {
		log.Println("Flightaware cl shutdown.")
	}
	return err
}
