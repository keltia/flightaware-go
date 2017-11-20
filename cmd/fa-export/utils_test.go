package main

import (
	"testing"
	"errors"
)


func TestAnalyzeFormat(t *testing.T) {
	_, err := AnalyzeFormat("foo")
	if err == nil {
		t.Errorf("Error: %s", err.Error())
	}

	_, err = AnalyzeFormat("%L")
	if err == nil {
		t.Errorf("Error: %s", err.Error())
	}

	day  := Rotation{RT_DAY}
	hour := Rotation{RT_HOUR}

	ret, err := AnalyzeFormat("%d")
	if ret != day {
		t.Errorf("Error: %s", err.Error())
	}

	ret, err = AnalyzeFormat("%h")
	if ret != hour {
		t.Errorf("Error: %s", err.Error())
	} else {
		t.Log("AnalyzeFormat passed")
	}
}

func TestParseDate(t *testing.T) {
	tm, err := ParseDate("2015-06-09 00:00:00")
	if err != nil {
		t.Errorf("Error: bad parsing for 2015-06-09 00:00:00: %v", err.Error())
	}
	if tm != 1433800800 {
		t.Errorf("Error: bad parsing for 2015-06-09 00:00:00: %v", errors.New("bad value"))
	}
}