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
	Host	config.Config
	Bytes	int64
	Pkts	int32
	Conn	*tls.Conn
	ch		chan []byte
}

func NewClient(rc config.Config) *FAClient {
	cl := new(FAClient)
	cl.Host = rc

	return cl
}

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
			fmt.Println(string(buf))

			cl.Bytes += int64(len(buf))
			cl.Pkts++
		}
	}()
	return ch, nil
}

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

	cl.ch = make(chan []byte, 100)
	conf := fmt.Sprintf("live version 4.0 username %s password %s events \"position\"", rc.User, rc.Password)
	_, err = conn.Write([]byte(conf))
	if err != nil {
		log.Println("Error configuring feed", err.Error())
		return err
	}

	log.Println("Flightaware init done.")
	cl.Conn = conn

	// Starting here everything is flowing from that connection
	ch, err := cl.StartWriter()

	//var	buf []byte

	log.Println("Loop")
	b := bufio.NewReader(cl.Conn)
	buf := make([]byte, 100)
	for {
		nb, err := b.Read(buf)
		if err != nil {
			log.Fatalf("error reading socket %v", nb)
		}
		ch <- buf
	}
	/*
	// Loop over chunks of data
	sc := bufio.NewScanner(cl.Conn)
	for {
		log.Println("Now waiting for data")
		for sc.Scan() {
			log.Println("in Scan()")
			buf := sc.Text()
			nb := len(buf)
			if err == nil && nb != 0 {
				log.Printf("Sending %d bytes\n", nb)
				ch <- []byte(buf)
			}
		}
	}
	*/
	return nil
}

func (cl *FAClient) Close() error {
	var err error

	if err := cl.Conn.Close(); err != nil {
		log.Println("Flightaware client shutdown.")
	}
	return err
}
