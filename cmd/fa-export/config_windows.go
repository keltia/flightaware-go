// +build !unix,windows

package main

import (
	"os"
	"path/filepath"
)

var (
	baseDir = filepath.Join(os.Getenv("%APPLOCALDATA%"),
		"flightaware",
		MyName,
	)
)
