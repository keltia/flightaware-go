// fa-tail.go
//
// tail(1)-like tuility that gets the last record and display the clock element
//
// Copyright 2015 © by Ollivier Robert for the EEC

// Implement the fa-export client.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"os"
	"time"
	"strconv"
	"fmt"
)

type FArecord struct {
	Type		string
	AirGround	string
	AltChange	string
	Clock		string
	Gs			string
	Heading		string
	Hexid		string
	Id			string
	Ident		string
	Lat			string
	Lon			string
	Reg			string
	Squawk		string
	UpdateType	string
}

func main() {
	flag.Parse()
	fh, err := os.Open(flag.Arg(0))
	scanner := bufio.NewScanner(fh)
	if err != nil {
		fmt.Printf("Unable to read %s\n", flag.Arg(0))
		os.Exit(1)
	}

	var (
		lastRecord string
		nbRecords	int64
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
	fmt.Printf("# records: %d - Last record: %v\n", nbRecords, time.Unix(iClock, 0))
}
