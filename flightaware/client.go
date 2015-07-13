// client.go
//

// Flightaware client package
package flightaware

import (
	"../config"
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"
	"errors"
	"strconv"
	"strings"
)

const (
	AUTHSTR = "%s version 4.0 username %s password %s events \"position\"\n"
)

type FAClient struct {
	Host     config.Config
	Bytes    int64
	Pkts     int32
	Conn     *tls.Conn
	Feed_one func([]byte)
	Verbose  bool
	FeedType string
	// For range event type
	RangeT   []int64
}

// Private functions
// Default callback
func defaultFeed(buf []byte) { fmt.Println(string(buf)) }

// Transform N:M into a array
func stringtoRange(s string) ([]int64, error) {
	begEnd := strings.Split(s, ":")

	if len(begEnd) != 2 {
		return []int64{}, errors.New("only one value")
	}
	var (
		beginT int64
		endT   int64
		err    error
	)

	if beginT, err = strconv.ParseInt(begEnd[0], 10, 64); err != nil {
		return []int64{}, err
	}

	if endT, err = strconv.ParseInt(begEnd[1], 10, 64); err != nil {
		return []int64{}, err
	}

	if beginT >= endT {
		return []int64{}, errors.New("begin > end")
	}
	return []int64{beginT, endT}, nil
}

// Send banner to FA
func (cl *FAClient) authClient(conn *tls.Conn) error {
	var authStr string = ""

	rc := cl.Host
	switch cl.FeedType {
		case "live":
			authStr = cl.FeedType
			if cl.Verbose {
				log.Println("Live traffic feed")
			}
		case "pitr":
			authStr = fmt.Sprintf("%s %d", cl.FeedType, cl.RangeT[0])
			if cl.Verbose {
				log.Println("Live traffic replay starting at %v",
					time.Unix(cl.RangeT[0], 0))
			}
		case "range":
			authStr = fmt.Sprintf("%s %d %d", cl.FeedType, cl.RangeT[0], cl.RangeT[1])
			if cl.Verbose {
				log.Printf("Replay traffic from %v to %v\n",
					time.Unix(cl.RangeT[0], 0),
					time.Unix(cl.RangeT[1], 0))
			}
	}

	conf := fmt.Sprintf(AUTHSTR, authStr, cl.FeedType, rc.User, rc.Password)
	_, err := conn.Write([]byte(conf))
	if err != nil {
		log.Println("Error configuring feed", err.Error())
		return err
	}
	return nil
}

// Public functions

// Create new instance of the client
func NewClient(rc config.Config) *FAClient {
	cl := new(FAClient)
	cl.Host = rc
	cl.Feed_one = defaultFeed

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
func (client *FAClient) CheckEvents(feedType, feedTimings string) error {
	// -t live and -T are mutually exclusive
	if feedType == "live" && feedTimings != "" {
		return errors.New("Error: can't use -t live and -T")
	}

	// Check when -t pitr that -T is single valued
	if feedType == "pitr" {
		if feedTimings == "" {
			return errors.New("Error: you must specify a value with -T")
		}

		// Allow only one value to -T for -t pitr
		if strings.Index(feedTimings, ":") != -1 {
			return errors.New("Error: only one value for -t pitr and -T")
		}

		// Check value
		restart, err := strconv.ParseInt(feedTimings, 10, 64)
		if err != nil {
			return err
		}
		if restart >= time.Now().Unix() {
			return errors.New(fmt.Sprintf("Error: -T %d is in the future", restart))
		}
		// Store out final value
		client.RangeT[0] = restart
	}

	if feedType == "range" {
		rangeT, err := stringtoRange(feedTimings)
		if err != nil {
			return errors.New(fmt.Sprintf("Bad range specified in %s\n", feedTimings))
		}
		// Store out final values
		client.RangeT = rangeT
	}
	return nil
}

// consumer part of the FA client
func (cl *FAClient) StartWriter() (chan []byte, error) {
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
				log.Printf("Read %d bytes\n", len(buf))
			}
			(cl.Feed_one)(buf)

			cl.Bytes += int64(len(buf))
			cl.Pkts++
		}
	}()
	return ch, nil
}

// This is the main function here:
// - starts the consumer in the background
// - reads data from FA and send it to the consumer
func (cl *FAClient) Start() error {
	var rc config.Config = cl.Host

	str := rc.Site + ":" + rc.Port
	if cl.Verbose {
		log.Printf("Connecting to %v with TLS\n", str)
	}

	conn, err := tls.Dial("tcp", str, &tls.Config{
		RootCAs:            nil,
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Println("failed to connect: " + err.Error())
		return err
	}

	if cl.Verbose {
		log.Println("TLS negociation done.")
	}

	if err := cl.authClient(conn); err != nil {
		log.Printf("Error: auth error for %s\n", rc.User)
		return err
	}

	if cl.Verbose {
		log.Println("Flightaware init done.")
	}
	cl.Conn = conn

	// Starting here everything is flowing from that connection
	ch, err := cl.StartWriter()
	if err != nil {
		log.Printf("Error: starting writer - %s\n", err.Error())
		return err
	}

	// Loop over chunks of data
	sc := bufio.NewScanner(cl.Conn)
	for {
		if cl.Verbose {
			log.Println("Now waiting for data")
		}
		for sc.Scan() {
			buf := sc.Text()

			if nb := len(buf); nb != 0 {
				if cl.Verbose {
					log.Printf("Sending %d bytes\n", nb)
				}
				ch <- []byte(buf)
			}
		}
	}
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
