// filters.go

/*
  This file implements the various filter-related functions.
*/
package flightaware

import (
	"fmt"
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
	return regexp.MustCompile(fmt.Sprintf(".*%s.*", str))
}

// Public functions

// Add an input filter to the list
func (cl *FAClient) AddInputFilter(fType int, str string) {
	if str != "" {
		cl.InputFilters = append(cl.InputFilters, generateFilter(fType, str))
	}
}

// Add an output filter
func (cl *FAClient) AddOutputFilter(str string) {
	if str != "" {
		of := generateRegex(str)
		cl.OutputFilters = append(cl.OutputFilters, of)
	}
}
