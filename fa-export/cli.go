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
	fVerbose     bool
	fOutput      string
	fFeedType    string
	fTimeout     int64
	fsTimeout    string
	fAutoRotate  bool
	fFeedBegin   string
	fFeedEnd     string
	fUserName    string
	fDest        string
)

// my usage string
const (
	cliUsage = `
Usage: %s [-o FILE] [-A] [-i N(s|mn|h|d)] [-f live|pitr|range [-B date [-E date]] [-v] [-u user]
`
)

// Redefine Usage
var Usage = func() {
	fmt.Fprintf(os.Stderr, cliUsage, os.Args[0])
	flag.PrintDefaults()
}

// called by flag.Parse()
func init() {
	// cli
	flag.StringVar(&fOutput, "o", "", "Specify output FILE.")
	flag.StringVar(&fFeedType, "f", "live", "Specify which feed we want (default live)")
	flag.StringVar(&fFeedBegin, "B", "", "Begin time for -f pitr|range")
	flag.StringVar(&fFeedEnd, "E", "", "End time for -f range")
	flag.StringVar(&fsTimeout, "i", "", "Stop after N s/mn/h/days")
	flag.BoolVar(&fAutoRotate, "A", false, "Autorotate output file")
	flag.StringVar(&fDest, "d", "", "Default destination (NOT IMPL)")
	flag.StringVar(&fUserName, "u", "", "Username to connect with")
	flag.BoolVar(&fVerbose, "v", false, "Set verbose flag.")
}
