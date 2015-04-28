// client.go
//
// Flightaware client package

package flightaware

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"../config"
)

type Client struct {
	Host	config.Config
	Bytes	int64
	Pkts	int32
}

func HelloWorld() {
	fmt.Println("Hello package")
}

func NewClient(rc config.Config) (Client, error) {
	cl := Client{rc, 0, 0}

	str := rc.Site + ":" + rc.Port
	fmt.Printf("Connecting to %v\n", str)

	roots := x509.NewCertPool()
	conn, err := tls.Dial("tcp", str, &tls.Config{
		RootCAs: roots,
		InsecureSkipVerify: true,	// XXX FIXME
	})
	if err != nil {
		panic("failed to connect: " + err.Error())
	}

	conn.Close()
	return cl, err
}
