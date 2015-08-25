// fa-tail.go
//
// tail(1)-like utility that gets the last record and display the clock element
//
// Copyright 2015 © by Ollivier Robert for the EEC

// Implement the fa-tail application, a FA-aware tail(1) clone.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	BSIZE = 2048
)

type FArecord struct {
	Type       string
	AirGround  string
	AltChange  string
	Clock      string
	Gs         string
	Heading    string
	Hexid      string
	Id         string
	Ident      string
	Lat        string
	Lon        string
	Reg        string
	Squawk     string
	UpdateType string
}

var (
	fVerbose bool
	fileStat os.FileInfo
)

func main() {
	flag.BoolVar(&fVerbose, "v", false, "Be verbose")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Error: You must specify a file!\n")
	}
	fn := flag.Arg(0)

	// Get the size to seek near the end
	fileStat, err := os.Stat(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error stat(2) on %s: %v\n", fn, err.Error())
	}

	//Seek and read
	fh, err := os.Open(fn)
	if err != nil {
		fmt.Printf("Unable to open %s\n", fn)
		os.Exit(1)
	}

	// Go forward fast
	_, err = fh.Seek(fileStat.Size() - BSIZE, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to seek into the file %s at %d\n", fn, fileStat.Size() - BSIZE)
		os.Exit(1)
	}

	// Then we go with the usual scanning thingy
	scanner := bufio.NewScanner(fh)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read %s\n", fn)
		os.Exit(1)
	}

	var (
		lastRecord string
		nbRecords  int64
	)

	for scanner.Scan() {
		lastRecord = scanner.Text()

		// In verbose mode we display the last few records read
		if fVerbose {
			fmt.Printf("%d: %s\n", nbRecords, lastRecord)
		}
		nbRecords++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading: %v", err)
	} else {
		// EOF
		if fVerbose {
			fmt.Printf("Last record: %s\n", lastRecord)
		}
	}

	var lastFA FArecord

	if err := json.Unmarshal([]byte(lastRecord), &lastFA); err != nil {
		fmt.Printf("Unable to decode %v: %v\n", lastRecord, err)
		os.Exit(1)
	}
	iClock, err := strconv.ParseInt(lastFA.Clock, 10, 64)
	fmt.Printf("%s: size %d bytes\n", fn, fileStat.Size())
	fmt.Printf("Last record: %v\n", time.Unix(iClock, 0))
}
