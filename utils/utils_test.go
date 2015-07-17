package utils

import (
	"testing"
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
