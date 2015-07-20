// client.go
//

// Flightaware client package
package flightaware

import (
	"flightaware-go/config"
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"
	"errors"
	"io"
)

const (
	FA_AUTHSTR = "%s username %s password %s events \"position\"\n"
	FA_VERSION = "version 4.0"
)

type FAClient struct {
	Host     config.Config
	Bytes    int64
	Pkts     int32
	Conn     *tls.Conn
	Feed_one func([]byte)
	Filter   func([]byte) bool
	Verbose  bool
	FeedType string
	// For range event type
	RangeT   []int64
}

// Private functions
// Default callback
func defaultFeed(buf []byte) { fmt.Println(string(buf)) }

// Default filter
func defaultFilter(buf []byte) bool { return true }

// Send banner to FA
func (cl *FAClient) authClient(conn *tls.Conn) error {
	var authStr string = ""

	rc := cl.Host
	switch cl.FeedType {
		case "live":
			authStr = fmt.Sprintf("%s %s", cl.FeedType, FA_VERSION)
			if cl.Verbose {
				log.Println("Live traffic feed")
			}
		case "pitr":
			authStr = fmt.Sprintf("%s %d", cl.FeedType, cl.RangeT[0])
			if cl.Verbose {
				log.Printf("Live traffic replay starting at %v",
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

	if cl.Verbose {
		log.Printf("Using %s as prefix.", authStr)
	}
	conf := fmt.Sprintf(FA_AUTHSTR, authStr, rc.User, rc.Password)
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
	cl.Filter = defaultFilter
	cl.RangeT = make([]int64, 2)

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
			return errors.New(fmt.Sprintf("Error: -B %d is in the future", restart))
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

			// Insert filter call
			if ok = (cl.Filter)(buf); ok {
				(cl.Feed_one)(buf)
			}

			cl.Bytes += int64(len(buf))
			cl.Pkts++
		}
	}()
	return ch, nil
}

// Connection handling, manage both initial and reconnections
func (cl *FAClient) ConnectFA(initial bool) (*tls.Conn, error) {
	var rc config.Config = cl.Host

	str := rc.Site + ":" + rc.Port
	if initial {
		if cl.Verbose {
			log.Printf("Connecting to %s with TLS\n", str)
		}
	} else {
		if cl.Verbose {
			log.Printf("Reconnecting to %s...\n", str)
		}
	}

	conn, err := tls.Dial("tcp", str, &tls.Config{
		RootCAs:            nil,
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Println("failed to connect: " + err.Error())
		return &tls.Conn{}, err
	}

	if cl.Verbose {
		log.Println("TLS negociation done.")
	}

	if err := cl.authClient(conn); err != nil {
		log.Printf("Error: auth error for %s\n", rc.User)
		return &tls.Conn{}, err
	}

	if cl.Verbose {
		log.Println("Flightaware init done.")
	}
	return conn, nil
}

// This is the main function here:
// - starts the consumer in the background
// - reads data from FA and send it to the consumer
func (cl *FAClient) Start() error {
	var rc config.Config = cl.Host

	conn, err := cl.ConnectFA(true)
	cl.Conn = conn

	str := rc.Site + ":" + rc.Port
	if cl.Verbose {
		log.Printf("Connecting to %v with TLS\n", str)
	}

	// Starting here everything is flowing from that connection
	ch, err := cl.StartWriter()
	if err != nil {
		log.Printf("Error: starting writer - %s\n", err.Error())
		return err
	}

	var stopped bool = false

	// Loop over chunks of data
	sc := bufio.NewScanner(cl.Conn)
	for {
		for sc.Scan() {
			buf := sc.Text()

			if nb := len(buf); nb != 0 {
				if cl.Verbose {
					log.Printf("Sending %d bytes\n", nb)
				}
				ch <- []byte(buf)
			}
		}
		if err := sc.Err(); err != nil {
			log.Println("Error reading:", err)

			// Reconnect
			conn, err = cl.ConnectFA(false)
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
