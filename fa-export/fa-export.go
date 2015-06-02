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
)

// Starts here.
func main() {
	// Handle SIGINT
	go func() {
	    sigint := make(chan os.Signal, 3)
	    signal.Notify(sigint, os.Interrupt)
	    <-sigint
	    log.Println("Program killed !")

		//doShutdown()

	    os.Exit(0)
	}()


	flag.Parse()

	log.Println("Hello world\n")

	c, err := config.LoadConfig(RcFile)
	if err != nil {
		log.Fatal("Error loading")
	}
	log.Println(c.Dests)
	log.Println(c.Default, c.Dests[c.Default])


	if client, err := flightaware.NewClient(c); err != nil {
		log.Printf("Bytes=%v\n", client.Bytes)
	}
}
