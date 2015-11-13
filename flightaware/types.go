// types.go

/*
  This file implements the types for the flightaware package
*/
package flightaware

import (
	"flightaware-go/config"
	"crypto/tls"
	"regexp"
)

type FAClient struct {
	Started      bool
	Host         config.Config
	Bytes        int64
	Pkts         int32
	Conn         *tls.Conn
	Feed_one     func([]byte)
	Filter       func(*FAClient, []byte) bool
	InputFilters []string
	OutputFilters []*regexp.Regexp
	Verbose      bool
	FeedType     string
	// For range event type
	RangeT []int64
}

// FA records types
//
// https://fr.flightaware.com/commercial/firehose/firehose_documentation.rvt

type FAgeneric struct {
	Type    string
	Payload []byte
}

type FApoint struct {
	// mandatory
	Lat float32
	Lon float32
	// common w/ 4D
	Clock string
	Name  string
	Alt   string
	// optional
	Gs    string
	AirspeedKts  string
	AirspeedMach string
}

type FAflightplan struct {
	// mandatory
	Type         string
	Ident        string
	Status       string
	Orig         string
	Dest         string
	Edt          string
	Eta          string
	Ete          string
	Id           string
	// common
	AircraftType string
	Suffix       string
	Reg          string
	Speed        string
	Alt          string
	FacilityHash string
	FacilityName string
	// optional
    Prefix       string
	Waypoints    []interface{}
	FDWaypoints  []interface{}
	Route        string
	Atcident     string
}

type FAdeparture struct {
	// mandatory
	Type         string
	Ident        string
	Orig         string
	Dest         string
	Adt          string
	Eta          string
	Id           string
	// common
	AircraftType string
	FacilityHash string
	FacilityName string
	// optional
	Synthetic    string
	Atcident     string
}

type FAarrival struct {
	// mandatory
	Type         string
	Ident        string
	Orig         string
	Dest         string
	Aat          string
	timeType     string
	Id           string
	// common
	FacilityHash string
	FacilityName string
	// optional
	Synthetic    string
	Atcident     string
}

type FAcancellation struct {
	// mandatory
	Type         string
	Ident        string
	Orig         string
	Dest         string
	Id           string
	// common
	FacilityHash string
	FacilityName string
	// optional
	Atcident     string
}

type FAposition struct {
	// mandatory
	Type         string
	Ident        string
	Lat          string
	Lon          string
	Clock        string
	Id           string
	UpdateType   string
	AirGround    string
	FacilityHash string
	FacilityName string
	// common
	Alt          string
	Gs           string
	Heading      string
	Rp1Lat       string
	Rp1Lon       string
	Rp1Alt       string
	Rp1Clock     string
	Squawk     string
	Hexid        string
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
	AltChange  string
	Reg        string
}
