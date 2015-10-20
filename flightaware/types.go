// types.go

/*
  This file implements the types for the flightaware package
*/
package flightaware

import (
	"../config"
	"crypto/tls"
)

// Our main struct to move state around
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

// Generic type that will allow to partially decode only what I need
type FArecord struct {
	Type		string
	Data		interface{}
}
