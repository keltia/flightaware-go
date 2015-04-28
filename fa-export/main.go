// main.go
//
//

package main

import (
	"../flightaware"
	"../config"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
)

var (
	RcFile = filepath.Join(os.Getenv("HOME"), ".flightaware", "config.yml")
)

func main() {
	app := cli.NewApp()
	app.Name = "fa-export"
	app.Author = "Ollivier Robert"
	app.Version = "0.0.1"
	app.Usage = "fa-export"

	fmt.Println("Hello world\n", app.Name)
	flightaware.HelloWorld()

	c, err := config.LoadConfig(RcFile)
	if err != nil {
		fmt.Println("Error loading")
	}
	fmt.Println(c.Dests)
	fmt.Println(c.Default, c.Dests[c.Default])
}
