// fa-export.go
//
// Small FlightAware Go client
//
// Copyright 2015 © by Ollivier Robert for the EEC
//

package main

import (
	"../flightaware"
	"../config"
	"flag"
	"os"
	"path/filepath"
	"os/signal"
	"log"
	"fmt"
	"regexp"
	"strconv"
	"time"
	"syscall"
)

const (
	SigALRM = 14
)

var (
	RcFile = filepath.Join(os.Getenv("HOME"), ".flightaware", "config.yml")
	client        *flightaware.FAClient
	fOutputFH	*os.File

	timeMods	= map[string]int64{
		"mn": 60,
		"h": 3600,
		"d": 3600*24,
	}
)

// fOutput callback
func fileOutput(buf []byte) {
	nb, err := fmt.Fprintln(fOutputFH, string(buf));
	if err != nil {
		log.Fatalf("Error writing %d bytes: %v", nb, err)
	}
}

// Proper shutdown
func stopEverything() {
	log.Printf("FA client stopped:")
	log.Printf("  %d pkts %d bytes", client.Pkts, client.Bytes)
	if err := client.Close(); err != nil {
		log.Println("Error closing connection:", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

// Check for specific modifiers, returns seconds
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

	// Handle SIGALRM
	go func() {
		sigalrm := make(chan os.Signal, SigALRM)
  		signal.Notify(sigalrm, os.Interrupt)
  		<-sigalrm

		log.Println("Timeout reached.")
		stopEverything()
	}()

	flag.Parse()

	c, err := config.LoadConfig(RcFile)
	if err != nil {
		log.Fatalf("Error loading configuration %f: %v\n", RcFile, err)
	}

	client = flightaware.NewClient(c)

	// Propagate this to the Client struct
	if fVerbose {
		client.Verbose = true
	}

	// Open output file
	if (fOutput != "") {
		if (fVerbose) {
			log.Printf("Output file is %s\n", fOutput)
		}

		if fOutputFH, err = os.Create(fOutput); err != nil {
			log.Printf("Error creating %s\n", fOutput)
			panic(err)
		}

		client.AddHandler(fileOutput)
	}

	// Check if we did specify a timeout with -i
	if fsTimeout != "" {
		fTimeout = checkTimeout(fsTimeout)

		log.Printf("Run for %ds\n", fTimeout)
		// Sleep for fTimeout seconds then sends SIGALRM
		go func() {
			time.Sleep(time.Duration(fTimeout) * time.Second)
			syscall.Kill(syscall.Getpid(), SigALRM)
		}()
	}

	// Get the flow running
	if err := client.Start(); err != nil {
		log.Fatalln("Error: unable to connect:", err)
	}
}
