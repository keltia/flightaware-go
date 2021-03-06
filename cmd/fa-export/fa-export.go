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
	"github.com/keltia/flightaware-go"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"time"
)

const (
	pprofPath = "/tmp/fa-export.prof"
)

var (
	RcFile    = "flightaware"
	client    *flightaware.FAClient
	fOutputFH *os.File

	// Us
	MyName = filepath.Base(os.Args[0])

	configName = "config.toml"

	// Our version
	FAversion = "1.6.0"

	RangeT []time.Time

	cnf Config
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
				pprofPath)
			pprof.StopCPUProfile()
		}
		verbose("FA client stopped:")
		verbose("  %d pkts %d bytes", client.Pkts, client.Bytes)
		if err := client.Close(); err != nil {
			log.Fatal("Error closing connection:", err)
		}
	}
	os.Exit(0)
}

// Check various parameters
func checkCommandLine() {
	// Propagate this to the Client struct
	if fVerbose {
		client.Verbose = true
		fmt.Printf("%s version %s API version: %s\n",
			MyName, FAversion, flightaware.FAVersion)
	}

	// Default is "live", incompatible with -B/-E
	if fFeedType == "live" && (fFeedBegin != "" || fFeedEnd != "") {
		log.Fatalf("Error: -B & -E are incompatible with -f live (the default)\n")
	}

	// When using -f pitr, we need -B starttime
	if fFeedType == "pitr" && fFeedBegin == "" {
		log.Fatalf("Error: you MUST use -B to specify starting time with -f pitr\n")
	}

	// When using -f range, we need -B starttime & -E endtime
	if fFeedType == "range" && (fFeedBegin == "" || fFeedEnd == "") {
		log.Fatalf("Error: you MUST use both -B & -E to specify times with -f range\n")
		os.Exit(1)
	}

	// Check for output filter

	// Transform the value if present
	if fHexid != "" {
		hexid := fmt.Sprintf("\"hexid\":\"%s\"", fHexid)
		client.AddOutputFilter(hexid)
	}

	// Replace defaults by anything on the CLI
	if fUserName != "" {
		cnf.DefUser = fUserName
	}
	if fDest != "" {
		cnf.DefDest = fDest
	}

	// Now parse them
	var (
		tFeedBegin, tFeedEnd time.Time
		err                  error
	)

	if fFeedType != "live" {
		RangeT = make([]time.Time, 2)
		if fFeedBegin != "" {
			tFeedBegin, err = ParseDate(fFeedBegin)
			if err != nil {
				log.Fatalf("Error: bad date format %v\n", fFeedBegin)
			}
		}

		if fFeedEnd != "" {
			tFeedEnd, err = ParseDate(fFeedEnd)
			if err != nil {
				log.Fatalf("Error: bad date format %v\n", fFeedEnd)
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
		verbose("tFeedBegin: %v - tFeedEnd: %v\n", tFeedBegin, tFeedEnd)
	}
}

// Starts here.
func main() {
	var (
		cnf *Config
		err error
	)

	// Handle SIGINT
	go func() {
		sigint := make(chan os.Signal, 3)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		stopEverything()
	}()

	flag.Usage = Usage
	flag.Parse()

	if fPProf {
		pp, err := os.Create(pprofPath)
		if err != nil {
			log.Fatalf("Can't create profiling file: %v\n", err)
		}
		pprof.StartCPUProfile(pp)
	}

	cnf, err = LoadConfig(fConfig)
	if err != nil {
		log.Fatalf("Error loading %s: %v\n", baseDir, err)
	}

	checkCommandLine()

	client = flightaware.NewClient(flightaware.Config{
		Site:     cnf.Site,
		Port:     cnf.Port,
		User:     cnf.Users[cnf.DefUser].User,
		Password: cnf.Users[cnf.DefUser].Password,
		FeedType: fFeedType,
	})

	// Open output file
	if fOutput != "" {
		verbose("Output file is %s\n", fOutput)

		// Check if the file already exist
		if fi, err := os.Stat(fOutput); err == nil {
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
				verbose("Info: %s renamed into %s\n", fOutput, newFile)
			}
		}

		//Default for Create() is to overwrite if already exist
		if fOutputFH, err = os.Create(fOutput); err != nil {
			log.Fatalf("can not create/overwrite %s: %v\n", fOutput, err)
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
			log.Printf("%v", err)
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
		fTimeout = CheckTimeout(fsTimeout)

		verbose("Run for %ds\n", fTimeout)
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
