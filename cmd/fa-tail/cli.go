package main

import "flag"

func init() {
	flag.BoolVar(&fVersion, "version,V", false, "Display version & quit.")
	flag.BoolVar(&fCount, "c", false, "Count records.")
	flag.BoolVar(&fVerbose, "v", false, "Be verbose")
}
