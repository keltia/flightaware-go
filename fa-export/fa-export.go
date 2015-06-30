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
)

var (
	RcFile = filepath.Join(os.Getenv("HOME"), ".flightaware", "config.yml")
	Client 		*flightaware.FAClient
	fOutputFH	*os.File
)

// fOutput callback
func fileOutput(buf []byte) {
	if nb, err := fOutputFH.Write(buf); err != nil {
		log.Fatalf("Error writing %d bytes: %v", nb, err)
	}
}

// Starts here.
func main() {
	// Handle SIGINT
	go func() {
	    sigint := make(chan os.Signal, 3)
	    signal.Notify(sigint, os.Interrupt)
	    <-sigint

		log.Printf("FA client stopped:")
		log.Printf("  %d pkts %d bytes", Client.Pkts, Client.Bytes)
		if err := Client.Close(); err != nil {
			log.Println("Error closing connection:", err)
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}()

	flag.Parse()

	c, err := config.LoadConfig(RcFile)
	if err != nil {
		log.Fatalf("Error loading configuration %f: %v\n", RcFile, err)
	}
	if fVerbose {
		log.Println(c.Dests)
		log.Println(c.Default, c.Dests[c.Default])
	}

	Client = flightaware.NewClient(c)

	// Open output file
	if (fOutput != "") {
		if (fVerbose) {
			log.Printf("Output file is %s\n", fOutput)
		}

		if fOutputFH, err = os.Create(fOutput); err != nil {
			log.Printf("Error creating %s\n", fOutput)
			panic(err)
		}

		Client.AddHandler(fileOutput)
	}


	// Get the flow running
	if err := Client.Start(); err != nil {
		log.Fatalln("Error: unable to connect:", err)
	}
}
