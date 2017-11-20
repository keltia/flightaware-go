// types.go

package flightaware

/*
  This file implements the types for the flightaware package
*/

import (
	"crypto/tls"
	"regexp"
)

type Config struct {
	Site     string
	Port     int
	User     string
	Password string
	FeedType string
}

// FAClient holds most data & configuration for a given client
type FAClient struct {
	Started       bool
	Host          Config
	Bytes         int64
	Pkts          int32
	Conn          *tls.Conn
	FeedOne       func([]byte)
	Filter        func(*FAClient, []byte) bool
	InputFilters  []string
	OutputFilters []*regexp.Regexp
	Verbose       bool
	FeedType      string
	// For range event type
	RangeT []int64
}

// FA records types
//
// https://fr.flightaware.com/commercial/firehose/firehose_documentation.rvt

// FAgeneric is just for finding the type
type FAgeneric struct {
	Type    string
	Payload []byte
}

// FApoint is used for holding points
type FApoint struct {
	// mandatory
	Lat float32
	Lon float32
	// common w/ 4D
	Clock string
	Name  string
	Alt   string
	// optional
	Gs           string
	AirspeedKts  string
	AirspeedMach string
}

// FAflightplan is a flight plan
type FAflightplan struct {
	// mandatory
	Type   string
	Ident  string
	Status string
	Orig   string
	Dest   string
	Edt    string
	Eta    string
	Ete    string
	ID     string
	// common
	AircraftType string
	Suffix       string
	Reg          string
	Speed        string
	Alt          string
	FacilityHash string
	FacilityName string
	// optional
	Prefix      string
	Waypoints   []interface{}
	FDWaypoints []interface{}
	Route       string
	Atcident    string
}

// FAdeparture is for departure events
type FAdeparture struct {
	// mandatory
	Type  string
	Ident string
	Orig  string
	Dest  string
	Adt   string
	Eta   string
	ID    string
	// common
	AircraftType string
	FacilityHash string
	FacilityName string
	// optional
	Synthetic string
	Atcident  string
}

// FAarrival is the corresponding type for arrivals
type FAarrival struct {
	// mandatory
	Type     string
	Ident    string
	Orig     string
	Dest     string
	Aat      string
	timeType string
	ID       string
	// common
	FacilityHash string
	FacilityName string
	// optional
	Synthetic string
	Atcident  string
}

// FAcancellation is for cancelled flights
type FAcancellation struct {
	// mandatory
	Type  string
	Ident string
	Orig  string
	Dest  string
	ID    string
	// common
	FacilityHash string
	FacilityName string
	// optional
	Atcident string
}

// FAposition is an ADS-B position
type FAposition struct {
	// mandatory
	Type         string
	Ident        string
	Lat          string
	Lon          string
	Clock        string
	ID           string
	UpdateType   string
	AirGround    string
	FacilityHash string
	FacilityName string
	// common
	Alt      string
	Gs       string
	Heading  string
	Rp1Lat   string
	Rp1Lon   string
	Rp1Alt   string
	Rp1Clock string
	Squawk   string
	Hexid    string
	// optional
	Fob          string
	Oat          string
	AirspeedKts  string
	AirspeedMach string
	Winds        string
	Eta          string
	BaroAlt      string
	GpsAlt       string
	Atcident     string
	// unknown
	AltChange string
	Reg       string
}
