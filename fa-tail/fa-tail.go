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
	BSIZE = 1024
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

func main() {
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

	_, err = fh.Seek(fileStat.Size() - BSIZE, 2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to seek into the file %s\n", fn)
		os.Exit(1)
	}

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
		nbRecords++
	}

	var lastFA FArecord

	if err := json.Unmarshal([]byte(lastRecord), &lastFA); err != nil {
		fmt.Printf("Unable to decode %v\n", lastRecord)
		os.Exit(1)
	}
	iClock, err := strconv.ParseInt(lastFA.Clock, 10, 64)
	fmt.Printf("%s: size %d bytes\n", fn, fileStat.Size())
	fmt.Printf("# records: %d - Last record: %v\n", nbRecords, time.Unix(iClock, 0))
}
