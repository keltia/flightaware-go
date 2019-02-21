// auth.go

/*
Package flightaware This file contains the authentication & connection functions
*/
package flightaware

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

const (
	faAuthStr = "%s username %s password %s %s\n"
)

// Send banner to FA
func (cl *FAClient) authClient(conn *tls.Conn) error {
	var authStr string

	rc := cl.Host
	switch cl.FeedType {
	case "live":
		authStr = fmt.Sprintf("%s", cl.FeedType)
		cl.verbose("Live traffic feed")
	case "pitr":
		authStr = fmt.Sprintf("%s %d", cl.FeedType, cl.RangeT[0])
		cl.verbose("Live traffic replay starting at %v",
			time.Unix(cl.RangeT[0], 0))
	case "range":
		authStr = fmt.Sprintf("%s %d %d", cl.FeedType, cl.RangeT[0], cl.RangeT[1])
		cl.verbose("Replay traffic from %v to %v\n",
			time.Unix(cl.RangeT[0], 0),
			time.Unix(cl.RangeT[1], 0))
	}

	cl.verbose("Using username %s", rc.User)
	cl.verbose("Using %s as prefix.", authStr)
	cl.verbose("Adding input filters: %s\n", setInputFilters(cl.InputFilters))

	// Set connection string including filters if any
	conf := fmt.Sprintf(faAuthStr, authStr,
		rc.User,
		rc.Password,
		setInputFilters(cl.InputFilters))

	_, err := conn.Write([]byte(conf))
	if err != nil {
		return errors.Wrap(err, "auth/feed")
	}
	return nil
}

// Connection handling, manage both initial and reconnections
func (cl *FAClient) connectFA(str string, initial bool) (*tls.Conn, error) {
	var rc = cl.Host

	if initial {
		cl.verbose("Connecting to %s with TLS\n", str)
	} else {
		cl.verbose("Reconnecting to %s...\n", str)
	}

	// XXX
	conn, err := tls.Dial("tcp", str, &tls.Config{
		RootCAs:            nil,
		InsecureSkipVerify: true,
	})
	if err != nil {
		return &tls.Conn{}, errors.Wrap(err, "connectFA")
	}

	cl.verbose("TLS negotiation done.")

	if err := cl.authClient(conn); err != nil {
		return &tls.Conn{}, errors.Wrapf(err, "connectFA/%s", rc.User)
	}

	cl.verbose("Flightaware init done.")
	return conn, nil
}
