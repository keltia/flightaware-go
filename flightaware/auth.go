// auth.go

/*
 This file contains the authentication & connection functions
*/
package flightaware

import (
	"crypto/tls"
	"fmt"
	"github.com/keltia/flightaware-go/config"
	"log"
	"time"
)

const (
	FA_AUTHSTR = "%s username %s password %s %s\n"
)

// Send banner to FA
func (cl *FAClient) authClient(conn *tls.Conn) error {
	var authStr string = ""

	rc := cl.Host
	switch cl.FeedType {
	case "live":
		authStr = fmt.Sprintf("%s", cl.FeedType)
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
		log.Printf("Using username %s", rc.DefUser)
		log.Printf("Using %s as prefix.", authStr)
		log.Printf("Adding input filters: %s\n", setInputFilters(cl.InputFilters))
	}

	// Set connection string including filters if any
	conf := fmt.Sprintf(FA_AUTHSTR, authStr,
		rc.Users[rc.DefUser].User,
		rc.Users[rc.DefUser].Password,
		setInputFilters(cl.InputFilters))

	_, err := conn.Write([]byte(conf))
	if err != nil {
		log.Println("Error configuring feed", err.Error())
		return err
	}
	return nil
}

// Connection handling, manage both initial and reconnections
func (cl *FAClient) connectFA(str string, initial bool) (*tls.Conn, error) {
	var rc config.Config = cl.Host

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
		log.Printf("Error: auth error for %s\n", rc.Users[rc.DefUser].User)
		return &tls.Conn{}, err
	}

	if cl.Verbose {
		log.Println("Flightaware init done.")
	}
	return conn, nil
}
