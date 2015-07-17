// utils.go
//
// Misc. utility functions
//
// Copyright 2015 © by Ollivier Robert for the EEC
//

// Misc. utility functions
package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

// Transform N:M into a array
func StringtoRange(s string) ([]int64, error) {
	begEnd := strings.Split(s, ":")

	if len(begEnd) != 2 {
		return []int64{}, errors.New("only one value")
	}
	var (
		beginT int64
		endT   int64
		err    error
	)

	if beginT, err = strconv.ParseInt(begEnd[0], 10, 64); err != nil {
		return []int64{}, errors.New("Can't parse beginT")
	}

	if endT, err = strconv.ParseInt(begEnd[1], 10, 64); err != nil {
		return []int64{}, errors.New("Can't parse endT")
	}

	if beginT >= endT {
		return []int64{}, errors.New("begin > end")
	}
	return []int64{beginT, endT}, nil
}

