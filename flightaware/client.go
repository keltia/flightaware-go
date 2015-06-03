// client.go
//
// Flightaware client package

package flightaware

import (
	"crypto/tls"
	"log"
	"fmt"
	"../config"
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

func (cl *FAClient) WriteData() (int, error) {
	buf, ok := <-cl.ch
	if !ok {

	}
	// Do something

	return len(buf), nil
}

func (cl *FAClient) Start() error {
	var rc	config.Config = cl.Host

	str := rc.Site + ":" + rc.Port
	log.Printf("Connecting to %v with TLS\n", str)

	conn, err := tls.Dial("tcp", str, &tls.Config{
		RootCAs: nil,
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
	go cl.WriteData()
	return nil
}

func (cl *FAClient) Close() error {
	var err error

	if err := cl.Conn.Close(); err != nil {
		log.Println("Flightaware client shutdown.")
	}
	return err
}
