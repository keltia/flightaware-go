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
	fVerbose   bool
	fOutput    string
	fFeedType string
	fTimeout   int64
	fsTimeout  string
	fAutoRotate bool
	fRestart	string
)

// my usage string
const (
	cliUsage = `
Usage: %s [-o FILE] [-A] [-i N(s|mn|h|d)] [-f live|pitr|range -F p1[:p2]] [-v]
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
	flag.StringVar(&fRestart, "F", "", "Parameters for -f pitr|range")
	flag.StringVar(&fsTimeout, "i", "60s", "Stop after N s/mn/h/days")
	flag.BoolVar(&fAutoRotate, "A", false, "Autorotate output file")
	flag.BoolVar(&fVerbose, "v", false, "Set verbose flag.")
}
