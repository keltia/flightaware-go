// types.go

/*
  This file implements the types for the flightaware package
*/
package flightaware

import (
	"../config"
	"crypto/tls"
)

type FAClient struct {
	Started      bool
	Host         config.Config
	Bytes        int64
	Pkts         int32
	Conn         *tls.Conn
	Feed_one     func([]byte)
	Filter       func([]byte) bool
	InputFilters []string
	Verbose      bool
	FeedType     string
	// For range event type
	RangeT []int64
}
