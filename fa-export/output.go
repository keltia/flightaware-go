// output.go

/*
  Implements the various output file(s)-related functions
 */
package main

import (
	"log"
	"os"
	"fmt"
)

// Open output file(s)
func openOutputFiles(openType int) {
	if openType == OPEN_SINGLE {
		if fVerbose {
			log.Printf("Output file is %s\n", fOutput)
		}

		// Check if the file already exist
		if fi, err := os.Stat(fOutput); err != nil {
			if fVerbose {
				log.Printf("Warning: %s (%v) already exists!", fOutput, fi.ModTime())
				if fOverwrite {
					log.Println("… overwriting it.")
				}
			}
			// Default for fOverwrite is false so we save the file
			if !fOverwrite {
				newFile := fmt.Sprintf("%s.old", fOutput)
				os.Rename(fOutput, newFile)
				if fVerbose {
					log.Printf("Info: %s renamed into %s\n", fOutput, newFile)
				}
			}
		}

		var err error

		//Default for Create() is to overwrite if already exist
		OutputFH, err = os.Create(fOutput);
		if err != nil {
			log.Printf("Error: can not create/overwrite %s\n", fOutput)
			panic(err)
		}

		client.AddHandler(fileOutput)
		// XXX FIXME Handle fAutoRotate
	} else {
		client.AddHandler(multipleFileOutput)
	}
}

