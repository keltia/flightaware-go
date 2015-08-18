// fa-export.go
//
// Small FlightAware Go client
//
// Copyright 2015 © by Ollivier Robert for the EEC

/*
 Package main implements the fa-export client.
 */
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"time"
	"flightaware-go/config"
	"flightaware-go/flightaware"
	"flightaware-go/utils"
)

var (
	RcFile    = filepath.Join(os.Getenv("HOME"), ".flightaware", "config.yml")
	client    *flightaware.FAClient
	fOutputFH *os.File

	timeMods = map[string]int64{
		"mn": 60,
		"h":  3600,
		"d":  3600 * 24,
	}

	RangeT       []time.Time
)

// fOutput callback
func fileOutput(buf []byte) {
	nb, err := fmt.Fprintln(fOutputFH, string(buf))
	if err != nil {
		log.Fatalf("Error writing %d bytes: %v", nb, err)
	}
}

// Proper shutdown
func stopEverything() {
	if client.Started {
		if fVerbose {
			log.Printf("FA client stopped:")
			log.Printf("  %d pkts %d bytes", client.Pkts, client.Bytes)
		}
		if err := client.Close(); err != nil {
			log.Println("Error closing connection:", err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}

// Check for specific modifiers, returns seconds
//
//XXX could use time.ParseDuration except it does not support days.
func checkTimeout(value string) int64 {
	mod := int64(1)
	re := regexp.MustCompile(`(?P<time>\d+)(?P<mod>(s|mn|h|d)*)`)
	match := re.FindStringSubmatch(value)
	if match == nil {
		return 0
	} else {
		// Get the base time
		time, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return 0
		}

		// Look for meaningful modifier
		if match[2] != "" {
			mod = timeMods[match[2]]
			if mod == 0 {
				mod = 1
			}
		}

		// At the worst, mod == 1.
		return time * mod
	}
}

// Check various parameters
func checkCommandLine() {
	// Propagate this to the Client struct
	if fVerbose {
		client.Verbose = true
	}

	// Default is "live", incompatible with -B/-E
	if fFeedType == "live" && (fFeedBegin != "" || fFeedEnd != "") {
		log.Printf("Error: -B & -E are incompatible with -f live (the default)\n")
		os.Exit(1)
	}

	// When using -f pitr, we need -B starttime
	if fFeedType == "pitr" && fFeedBegin == "" {
		log.Printf("Error: you MUST use -B to specify starting time with -f pitr\n")
		os.Exit(1)
	}

	// When using -f range, we need -B starttime & -E endtime
	if fFeedType == "range" && (fFeedBegin == "" || fFeedEnd == "") {
		log.Printf("Error: you MUST use both -B & -E to specify times with -f range\n")
		os.Exit(1)
	}

	// Now parse them
	var (
		tFeedBegin, tFeedEnd time.Time
		err                  error
	)

	if fFeedType != "live" {
		RangeT = make([]time.Time, 2)
		if fFeedBegin != "" {
			tFeedBegin, err = utils.ParseDate(fFeedBegin)
			if err != nil {
				log.Printf("Error: bad date format %v\n", fFeedBegin)
				os.Exit(1)
			}
		}

		if fFeedEnd != "" {
			tFeedEnd, err = utils.ParseDate(fFeedEnd)
			if err != nil {
				log.Printf("Error: bad date format %v\n", fFeedEnd)
				os.Exit(1)
			}
		} else {
			tFeedEnd = time.Time{}
		}

		if tFeedEnd.Before(tFeedBegin) {
			log.Printf("Warning: reversed date range, inverting.")
			tFeedBegin, tFeedEnd = tFeedEnd, tFeedBegin
		}
		RangeT[0] = tFeedBegin
		RangeT[1] = tFeedEnd
		if fVerbose {
			log.Printf("tFeedBegin: %v - tFeedEnd: %v\n", tFeedBegin, tFeedEnd)
		}
	}
}

// Starts here.
func main() {
	// Handle SIGINT
	go func() {
		sigint := make(chan os.Signal, 3)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		stopEverything()
	}()

	flag.Parse()

	c, err := config.LoadConfig(RcFile)
	if err != nil {
		log.Fatalf("Error loading configuration %s: %s\n", RcFile, err.Error())
	}

	client = flightaware.NewClient(c)
	client.FeedType = fFeedType

	checkCommandLine()

	// Open output file
	if fOutput != "" {
		if fVerbose {
			log.Printf("Output file is %s\n", fOutput)
		}

		if fOutputFH, err = os.Create(fOutput); err != nil {
			log.Printf("Error creating %s\n", fOutput)
			panic(err)
		}

		client.AddHandler(fileOutput)
		// XXX FIXME Handle fAutoRotate
	} else {
		if fAutoRotate {
			log.Println("Warning: -A needs -O to work, ignoring")
			fAutoRotate = false
		}
	}

	// Check if we want a live stream or a more specialized one
	if fFeedType != "" {
		if err := client.SetFeed(fFeedType, RangeT); err != nil {
			log.Printf("%s", err.Error())
		}
	}

	// Check if we did specify a timeout with -i
	if fsTimeout != "" {
		fTimeout = checkTimeout(fsTimeout)

		if fVerbose {
			log.Printf("Run for %ds\n", fTimeout)
		}
		client.SetTimer(fTimeout)
	}

	// Get the flow running
	if err := client.Start(); err != nil {
		log.Fatalln("Error: unable to connect:", err)
	}

	// Got an EOF from Scan(), finish gracefully
	log.Printf("Stream finished, good bye\n")
	stopEverything()
	// NOTREACHED
}
