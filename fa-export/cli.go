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
	"../utils"
	"time"
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

	RangeT       []time.Time
)

// my usage string
const (
	cliUsage = `
Usage: %s [-o FILE] [-A] [-i N(s|mn|h|d)] [-f live|pitr|range [-B date [-E date]] [-v]
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
	flag.BoolVar(&fVerbose, "v", false, "Set verbose flag.")

	// Default is "live", incompatible with -B/-E
	if fFeedType == "live" && (fFeedBegin != "" || fFeedEnd != "") {
		fmt.Printf("Error: -B & -E are incompatible with -f live (the default)\n")
		os.Exit(1)
	}

	// When using -f pitr, we need -B starttime
	if fFeedType == "pitr" && fFeedBegin == "" {
		fmt.Printf("Error: you MUST use -B to specify starting time with -f pitr\n")
		os.Exit(1)
	}

	// When using -f range, we need -B starttime & -E endtime
	if fFeedType == "range" && (fFeedBegin == "" || fFeedEnd == "") {
		fmt.Printf("Error: you MUST use both -B & -E to specify times with -f range\n")
		os.Exit(1)
	}

	// Now parse them
	var (
		tFeedBegin, tFeedEnd time.Time
		err                  error
	)

	if fFeedBegin != "" {
		tFeedBegin, err = utils.ParseDate(fFeedBegin)
		if err != nil {
			fmt.Printf("Error: bad date format %v\n", fFeedBegin)
			os.Exit(1)
		}
	}

	if fFeedEnd != "" {
		tFeedEnd, err = utils.ParseDate(fFeedEnd)
		if err != nil {
			fmt.Printf("Error: bad date format %v\n", fFeedEnd)
			os.Exit(1)
		}
	} else {
		tFeedEnd = time.Time{}
	}

	if tFeedEnd.Before(tFeedBegin) {
		fmt.Printf("Warning: reversed date range, inverting.")
		tFeedBegin, tFeedEnd = tFeedEnd, tFeedBegin
	}
	RangeT[0] = tFeedBegin
	RangeT[1] = tFeedEnd
}
