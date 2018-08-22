// utils.go
//
// Misc. utility functions
//
// Copyright 2015 © by Ollivier Robert for the EEC
//

/*
 Package utils implements misc. utility functions
*/
package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

const (
	RT_HOUR = 1
	RT_DAY  = 2

	TIMEFMT = "2006-01-02 15:04:05"
)

var (
	timeMods = map[string]int64{
		"mn": 60,
		"h":  3600,
		"d":  3600 * 24,
	}
)

type Rotation struct {
	Type int
}

// Check the format for the logs to be rotated
func AnalyzeFormat(sFmt string) (Rotation, error) {
	format := []rune(sFmt)
	if format[0] != '%' {
		return Rotation{}, errors.New("Badly formatted string w/o %")
	}

	switch format[1] {
	case 'h':
		return Rotation{RT_HOUR}, nil
	case 'd':
		return Rotation{RT_DAY}, nil
	default:
		return Rotation{}, errors.New(fmt.Sprintf("Unknown modifier %s\n", string(format[1])))
	}
}

// Parse date into UNIX epoch-style int64
func ParseDate(date string) (time.Time, error) {
	tDate, err := time.Parse(TIMEFMT, date)
	if err != nil {
		return time.Time{}, err
	}
	return tDate, nil
}

// Check for specific modifiers, returns seconds
//
//XXX could use time.ParseDuration except it does not support days.
func CheckTimeout(value string) int64 {
	mod := int64(1)
	re := regexp.MustCompile(`(?P<time>\d+)(?P<mod>(s|mn|h|d)*)`)
	match := re.FindStringSubmatch(value)
	if match == nil {
		return 0
	} else {
		// Get the base tm
		tm, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return 0
		}

		// Look for meaningful modifier
		if match[2] != "" {
			mod = timeMods[match[2]]
			if mod == 0 {
				mod = 1
			}
		}

		// At the worst, mod == 1.
		return tm * mod
	}
}

// debug displays only if fDebug is set
func debug(str string, a ...interface{}) {
	if fDebug {
		log.Printf(str, a...)
	}
}

// verbose displays only if fVerbose is set
func verbose(str string, a ...interface{}) {
	if fVerbose {
		log.Printf(str, a...)
	}
}
