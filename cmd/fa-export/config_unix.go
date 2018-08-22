// +build !windows

package main

import (
	"os"
	"path/filepath"
)

var (
	baseDir = filepath.Join(os.Getenv("HOME"),
		".config",
		RcFile,
	)
)
