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
	"../config"
	"../flightaware"
	"../utils"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	"runtime/pprof"
)

const PPROF_PATH = "/tmp/fa-export.prof"

var (
	RcFile    = "flightaware"
	client    *flightaware.FAClient
	fOutputFH *os.File

	RangeT []time.Time
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
		if fPProf {
			log.Printf("Stopping profiling…")
			log.Printf(`
Profiling mode was enabled.
Please use go tool pprof %s %s to read profiling data`,
			flag.Arg(0),
			PPROF_PATH)
			pprof.StopCPUProfile()
		}
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

	// Check for output filter

	// Transform the value if present
	if fHexid != "" {
		hexid := fmt.Sprintf("hexid:\"%s\"", fHexid)
		client.AddOutputFilter(hexid)
	}

	// Replace defaults by anything on the CLI
	if fUserName != "" {
		client.Host.DefUser = fUserName
	}
	if fDest != "" {
		client.Host.DefDest = fDest
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

		// Do no invert if using "pitr" as B>E in this case
		if tFeedEnd.Before(tFeedBegin) && fFeedType != "pitr" {
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

	if fPProf {
		pp, err := os.Create(PPROF_PATH)
		if err != nil {
			log.Fatalf("Can't create profiling file: %v\n", err)
		}
		pprof.StartCPUProfile(pp)
	}

	c, err := config.LoadConfig(RcFile)
	if err != nil {
		log.Fatalf("Error loading %s: %s\n", RcFile, err.Error())
	}

	client = flightaware.NewClient(*c)
	client.FeedType = fFeedType

	checkCommandLine()

	// Open output file
	if fOutput != "" {
		if fVerbose {
			log.Printf("Output file is %s\n", fOutput)
		}

		// Check if the file already exist
		if fi, err := os.Stat(fOutput); err != nil {
			if fVerbose {
				log.Printf("Warning: %s (%v) already exists!", fOutput, fi.ModTime())
				if fOverwrite {
					log.Println("… overwriting it.")
				}
			}
			// Default for fOverwrite is false so we save the file
			if !fOverwrite {
				newFile := fmt.Sprintf("%s.old", fOutput)
				os.Rename(fOutput, newFile)
				if fVerbose {
					log.Printf("Info: %s renamed into %s\n", fOutput, newFile)
				}
			}
		}

		//Default for Create() is to overwrite if already exist
		if fOutputFH, err = os.Create(fOutput); err != nil {
			log.Printf("Error: can not create/overwrite %s\n", fOutput)
			panic(err)
		}

		client.AddHandler(fileOutput)
		// XXX FIXME Handle fAutoRotate
	} else {
		if fAutoRotate {
			log.Println("Warning: -A needs -o to work, ignoring")
			fAutoRotate = false
		}
	}

	// Check if we want a live stream or a more specialized one
	if fFeedType != "" {
		if err := client.SetFeed(fFeedType, RangeT); err != nil {
			log.Printf("%s", err.Error())
		}
	}

	// Check the various possible input filters
	client.AddInputFilter(flightaware.FILTER_EVENT, fEventType)
	client.AddInputFilter(flightaware.FILTER_AIRLINE, fAirlineFilter)
	client.AddInputFilter(flightaware.FILTER_IDENT, fIdentFilter)
	client.AddInputFilter(flightaware.FILTER_AIRPORT, fAirportFilter)
	client.AddInputFilter(flightaware.FILTER_LATLONG, fLatLongFilter)

	// Check if we did specify a timeout with -i
	if fsTimeout != "" {
		fTimeout = utils.CheckTimeout(fsTimeout)

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
