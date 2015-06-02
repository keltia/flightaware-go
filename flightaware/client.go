// client.go
//
// Flightaware client package

package flightaware

import (
	"crypto/tls"
	"fmt"

	"../config"
)

type Client struct {
	Host	config.Config
	Bytes	int64
	Pkts	int32
}

func NewClient(rc config.Config) (*Client, error) {
	cl := new(Client)
	cl.Host = rc

	str := rc.Site + ":" + rc.Port
	fmt.Printf("Connecting to %v with TLS\n", str)

	conn, err := tls.Dial("tcp", str, &tls.Config{
		RootCAs: nil,
	})
	if err != nil {
		panic("failed to connect: " + err.Error())
	}

	fmt.Println("TLS negociation done.")

	conf := fmt.Sprintf("live version 4.0 username %s password %s events \"position\"", rc.User, rc.Password)
	conn.Write([]byte(conf))

	fmt.Println("Flightaware init done.")

	// Insert here the io.Reader code

	conn.Close()
	return cl, err
}
