// filters.go

package flightaware

/*
This file implements the various filter-related functions.
*/

import (
	"fmt"
	"log"
	"regexp"
)

const (
	FILTER_EVENT = iota
	FILTER_AIRLINE
	FILTER_IDENT
	FILTER_LATLONG
	FILTER_AIRPORT
)

var (
	filterTypes = map[int]string{
		FILTER_EVENT:   "events \"%s\"",
		FILTER_AIRLINE: "filter \"%s\"",
		FILTER_IDENT:   "idents \"%s\"",
		FILTER_LATLONG: "latlong \"%s\"",
		FILTER_AIRPORT: "airport_filter \"%s\"",
	}
)

// Private functions

// Generate the proper argument for a given filter
func generateFilter(fType int, str string) string {
	return fmt.Sprintf(filterTypes[fType], str)
}

// Generate the filter list for FA
func setInputFilters(inputFilters []string) string {
	result := ""

	for _, str := range inputFilters {
		result = result + " " + str
	}
	return result
}

// Generate a regex from a simple pattern
func generateRegex(str string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^.*%s.*$", str))
}

// Public functions

// AddInputFilter adds an input filter to the list
func (cl *FAClient) AddInputFilter(fType int, str string) {
	if str != "" {
		if cl.Verbose {
			log.Printf("Adding input filter type %d on %s\n", fType, str)
		}
		cl.InputFilters = append(cl.InputFilters, generateFilter(fType, str))
	}
}

// AddOutputFilter adds an output filter
func (cl *FAClient) AddOutputFilter(str string) {
	if str != "" {
		of := generateRegex(str)
		if cl.Verbose {
			log.Printf("Adding output filter on %s: %v\n", str, of)
		}
		cl.OutputFilters = append(cl.OutputFilters, of)
	}
}
