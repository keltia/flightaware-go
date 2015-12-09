// client.go
//

/*
 Package Flightaware implements the client part to the FA API

 You start by creating a client instance with your credentials passed as config.Config
 struct, previously generated by calling config.LoadConfig() or manually.

 	client := flightaware.NewClient(config.Config)

 Then you can configure the feed type with

 	client.SetFeed(string, []time.Time)

 You can also set a timeout time with a value in seconds

 	client.SetTimeout(int64)

 You can add one or more different input filters:

    client.AddInputFilter(<type>, <value>)

 where type can be one of

     FILTER_EVENT
     FILTER_AIRLINE
     FILTER_IDENT
     FILTER_AIRPORT
     FILTER_LATLONG

 The filters you specify will be checked remotely by FlightAware according to the
 documentation available at
 https://fr.flightaware.com/commercial/firehose/firehose_documentation.rvt

 You can specify output filters with using client.AddOutputFilter(string)

 The default handler is to display all packets.  You can change the default handler
 with

 	client.AddHandler(func([]byte)

 Last action is to start the consuming/producer loop with

 	client.Start()

 Reading will be closed either though getting an EOF from FA or being will killed either
 manually or through the timeout value.

 You can then use

 	client.Close()

 to properly close the reading channel.
*/
package flightaware

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
	"github.com/keltia/flightaware-go/config"
)

// Private functions

// Default callback
func defaultFeed(buf []byte) { fmt.Println(string(buf)) }

// Default filter
func defaultFilter(cl *FAClient, buf []byte) bool {
	if len(cl.OutputFilters) != 0 {
		for index, flt := range cl.OutputFilters {
			// First match so behaviour is OR
			if flt.Match(buf) {
				if cl.Verbose {
					log.Printf("%d", index)
				}
				return true
			}
			if cl.Verbose {
				log.Print(".")
			}
		}
		// If no match from any of the filters, no cigar
		return false
	}
	return true
}

// consumer part of the FA client
func (cl *FAClient) startWriter() (chan []byte, error) {
	if cl.Verbose {
		log.Println("Waiting for data…")
	}
	ch := make(chan []byte, 1000)
	go func() {
		for {
			buf, ok := <-ch
			if !ok {
				log.Fatalf("Error: reading data: %s: %v", string(buf), ok)
			}
			// Do something
			if cl.Verbose {
				DataLog(buf, fmt.Sprintf("Read %d bytes\n", len(buf)))
			}

			// Insert filter call
			if ok = (cl.Filter)(cl, buf); ok {
				(cl.Feed_one)(buf)
			}

			cl.Bytes += int64(len(buf))
			cl.Pkts++
		}
	}()
	return ch, nil
}

// Public functions

// Create new instance of the client
func NewClient(rc config.Config) *FAClient {
	cl := new(FAClient)
	cl.Host = rc
	cl.Feed_one = defaultFeed
	cl.Filter = defaultFilter
	cl.RangeT = make([]int64, 2)
	cl.Started = false
	cl.InputFilters = []string{}
	cl.OutputFilters = []*regexp.Regexp{}

	return cl
}

// Change default callback
func (cl *FAClient) AddHandler(fn func([]byte)) {
	cl.Feed_one = fn
}

// Allow run of specified duration
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

// Check if parameters for the event type are consistent
// Check that -t has also -T and the right parameters
func (client *FAClient) SetFeed(feedType string, RangeT []time.Time) error {
	// Check when -t pitr that -T is single valued
	if feedType == "pitr" {
		// Check value
		restart := RangeT[0]
		if restart.After(time.Now()) {
			return errors.New(fmt.Sprintf("Error: -B %v is in the future", restart))
		}
		// Store out final value
		client.RangeT[0] = restart.Unix()
	}

	if feedType == "range" {
		// Store out final values in UNIX epoch format
		client.RangeT[0] = RangeT[0].Unix()
		client.RangeT[1] = RangeT[1].Unix()
	}
	return nil
}

// This is the main function here:
// - starts the consumer in the background
// - reads data from FA and send it to the consumer
func (cl *FAClient) Start() error {
	var rc config.Config = cl.Host

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
		if err := sc.Err(); err != nil {
			log.Println("Error reading:", err)

			// Reconnect
			conn, err = cl.connectFA(str, false)
			sc = bufio.NewScanner(cl.Conn)
		} else {
			// Got EOF
			break
		}
	}
	return nil
}

// Properly close the TLS connection
func (cl *FAClient) Close() error {
	var err error

	if err := cl.Conn.Close(); err != nil {
		log.Println("Error closing connection " + err.Error())
	}
	if cl.Verbose {
		log.Println("Flightaware client shutdown.")
	}
	return err
}
