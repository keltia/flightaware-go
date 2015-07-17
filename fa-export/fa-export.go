// fa-export.go
//
// Small FlightAware Go client
//
// Copyright 2015 © by Ollivier Robert for the EEC

// Implement the fa-export client.
package main

import (
	"../config"
	"../flightaware"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
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
	if fVerbose {
		log.Printf("FA client stopped:")
		log.Printf("  %d pkts %d bytes", client.Pkts, client.Bytes)
	}
	if err := client.Close(); err != nil {
		log.Println("Error closing connection:", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
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

	// Propagate this to the Client struct
	if fVerbose {
		client.Verbose = true
	}

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
}
