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
}

func NewClient(rc config.Config) (*FAClient, error) {
	cl := new(FAClient)
	cl.Host = rc

	str := rc.Site + ":" + rc.Port
	log.Printf("Connecting to %v with TLS\n", str)

	conn, err := tls.Dial("tcp", str, &tls.Config{
		RootCAs: nil,
	})
	if err != nil {
		panic("failed to connect: " + err.Error())
	}

	log.Println("TLS negociation done.")

	conf := fmt.Sprintf("live version 4.0 username %s password %s events \"position\"", rc.User, rc.Password)
	conn.Write([]byte(conf))

	log.Println("Flightaware init done.")

	// Insert here the io.Reader code

	conn.Close()
	return cl, err
}
