// cli.go
//
// Everything related to command-line flag handling
//
// Copyright 2015 © by Ollivier Robert for the EEC
//

package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// cli
	fVerbose	bool
	fOutput		string
	fEventType	string
)

// my usage string
const (
	cliUsage	= `
Usage: %s [-o FILE] [-t live|pitr|range] [-v]
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
	flag.StringVar(&fEventType, "t", "live", "Specify which events we want (default live)")
	flag.BoolVar(&fVerbose, "v", false, "Set verbose flag.")
}
