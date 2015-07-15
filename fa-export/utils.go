// utils.go
//
// Misc. utility functions
//
// Copyright 2015 © by Ollivier Robert for the EEC
//

// Misc. utility functions
package main

import (
	"errors"
	"fmt"
)

const (
	RT_HOUR = 1
	RT_DAY  = 2
)

type Rotation struct {
	Type	int

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
