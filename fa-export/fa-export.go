// fa-export.go
//
//

package main

import (
	"../flightaware"
	"../config"
	"flag"
	"fmt"
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

	fmt.Println("Hello world\n")
	flightaware.HelloWorld()

	c, err := config.LoadConfig(RcFile)
	if err != nil {
		fmt.Println("Error loading")
	}
	fmt.Println(c.Dests)
	fmt.Println(c.Default, c.Dests[c.Default])

	client, err := flightaware.NewClient(c)
	fmt.Printf("Bytes=%v\n", client.Bytes)
}
