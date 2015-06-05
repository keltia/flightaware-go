// fa-export.go
//
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
	Client *flightaware.FAClient
)

// Starts here.
func main() {
	// Handle SIGINT
	go func() {
	    sigint := make(chan os.Signal, 3)
	    signal.Notify(sigint, os.Interrupt)
	    <-sigint

		log.Printf("FA client stopped:")
		log.Printf("  %d pkts %ld bytes", Client.Pkts, Client.Bytes)
		if err := Client.Close(); err != nil {
			log.Println("Error closing connection:", err)
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}()

	flag.Parse()

	log.Println("Hello world\n")

	c, err := config.LoadConfig(RcFile)
	if err != nil {
		log.Fatal("Error loading")
	}
	log.Println(c.Dests)
	log.Println(c.Default, c.Dests[c.Default])

	Client = flightaware.NewClient(c)

	// Get the flow running
	if err := Client.Start(); err != nil {
		log.Fatalln("Error: unable to connect:", err)
	}
}
