// client.go
//
// Flightaware client package

package flightaware

import (
	"crypto/tls"
	"log"
	"fmt"
	"../config"
	"bufio"
)

type FAClient struct {
	Host		config.Config
	Bytes		int64
	Pkts		int32
	Conn		*tls.Conn
	Feed_one	func([]byte)
}

// Create new instance of the client
func NewClient(rc config.Config) *FAClient {
	cl := new(FAClient)
	cl.Host = rc
	cl.Feed_one = defaultFeed

	return cl
}

// Default callback
func defaultFeed(buf []byte) { fmt.Println(string(buf)) }

// Change default callback
func (cl *FAClient) AddHandler(fn func([]byte)) {
	cl.Feed_one = fn
}

// consumer part of the FA client
func (cl *FAClient) StartWriter() (chan []byte, error) {
	log.Println("Waiting for data…")
	ch := make(chan []byte, 1000)
	go func() {
		for {
			buf, ok := <-ch
			if !ok {
				log.Fatalf("Error: reading data: %s: %v", string(buf), ok)
			}
			// Do something
			log.Printf("Read %d bytes\n", len(buf))
			(cl.Feed_one)(buf)

			cl.Bytes += int64(len(buf))
			cl.Pkts++
		}
	}()
	return ch, nil
}

// Send banner to FA
func authClient(conn *tls.Conn, rc config.Config) error {
	conf := fmt.Sprintf("live version 4.0 username %s password %s events \"position\"\n", rc.User, rc.Password)
	_, err := conn.Write([]byte(conf))
	if err != nil {
		log.Println("Error configuring feed", err.Error())
		return err
	}
	return nil
}

// This is the main function here:
// - starts the consumer in the background
// - reads data from FA and send it to the consumer
func (cl *FAClient) Start() error {
	var rc	config.Config = cl.Host

	str := rc.Site + ":" + rc.Port
	log.Printf("Connecting to %v with TLS\n", str)

	conn, err := tls.Dial("tcp", str, &tls.Config{
		RootCAs: nil,
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Println("failed to connect: " + err.Error())
		return err
	}

	log.Println("TLS negociation done.")

	if err := authClient(conn, rc); err != nil {
		log.Printf("Error: auth error for %s\n", rc.User)
		return err
	}

	log.Println("Flightaware init done.")
	cl.Conn = conn

	// Starting here everything is flowing from that connection
	ch, err := cl.StartWriter()
	if err != nil {
		log.Println("Error: starting writer - %v", err)
		return err
	}

	// Loop over chunks of data
	sc := bufio.NewScanner(cl.Conn)
	for {
		log.Println("Now waiting for data")
		for sc.Scan() {
			buf := sc.Text()

			if nb := len(buf); nb != 0 {
				log.Printf("Sending %d bytes\n", nb)
				ch <- []byte(buf)
			}
		}
	}

	return nil
}

func (cl *FAClient) Close() error {
	var err error

	if err := cl.Conn.Close(); err != nil {
		log.Println("Flightaware client shutdown.")
	}
	return err
}
