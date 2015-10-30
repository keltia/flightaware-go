// cli.go
//
// Everything related to command-line flag handling
//
// Copyright 2015 © by Ollivier Robert for the EEC
//

// Command-line Interface for fa-export
package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// cli
	fVerbose       bool
	fOutput        string
	fEventType     string
	fFeedType      string
	fTimeout       int64
	fsTimeout      string
	fAutoRotate    bool
	fFeedBegin     string
	fFeedEnd       string
	fOverwrite     bool
	fUserName      string
	fDest          string
	fAirlineFilter string
	fIdentFilter   string
	fLatLongFilter string
	fAirportFilter string
	fHexid         string
	fPProf         bool
)

// my usage string
const (
	cliUsage = `
%s version %s
Usage: %s [-o FILE] [-d N(s|mn|h|d)][-f live|pitr|range [-B date [-E date]] [-v] [-u user]

       Filters (OR is implied if multiple):
          [-e type] [-F airline] [-I plane-ident] [-L lat/lon] [-P airport-glob]
       Output filter (not on theFA command line)
          [-X hexid]
`
)

// Redefine Usage
var Usage = func() {
	fmt.Fprintf(os.Stderr, cliUsage, os.Args[0], FA_VERSION, os.Args[0])
	flag.PrintDefaults()
}

// called by flag.Parse()
func init() {
	// cli
	flag.StringVar(&fsTimeout, "d", "", "Stop after N s/mn/h/days")
	flag.StringVar(&fEventType, "e", "", "Events to stream")
	flag.StringVar(&fFeedType, "f", "live", "Specify which feed we want")
	flag.StringVar(&fOutput, "o", "", "Specify output FILE.")
	flag.BoolVar(&fOverwrite, "O", false, "Overwrite existing file?")
	flag.BoolVar(&fPProf, "p", false, "Enable profiling")
	flag.StringVar(&fUserName, "u", "", "Username to connect with")
	flag.BoolVar(&fVerbose, "v", false, "Set verbose flag.")
	flag.BoolVar(&fAutoRotate, "A", false, "Autorotate output file")
	flag.StringVar(&fFeedBegin, "B", "", "Begin time for -f pitr|range")
	flag.StringVar(&fDest, "D", "", "Default destination (NOT IMPL)")
	flag.StringVar(&fFeedEnd, "E", "", "End time for -f range")
	flag.StringVar(&fAirlineFilter, "F", "", "Airline filter")
	flag.StringVar(&fIdentFilter, "I", "", "Aircraft Ident filter")
	flag.StringVar(&fLatLongFilter, "L", "", "Lat/Long filter")
	flag.StringVar(&fAirportFilter, "P", "", "Airport filter")
	flag.StringVar(&fHexid, "X", "", "Hexid output filter")
}
