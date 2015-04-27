// main.go
//
//

package main

import (
	"../flightaware"
	"fmt"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "fa-export"
	app.Author = "Ollivier Robert"
	app.Version = "0.0.1"
	app.Usage = "fa-export"
	fmt.Println("Hello world\n", app.Name)
	flightaware.HelloWorld()
}
