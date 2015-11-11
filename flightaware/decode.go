// decode.go

/*
  This package implements functions to decode the various different types
  of data sent by Flightaware.
 */
package flightaware

import (
	"encoding/json"
	"fmt"
	"errors"
)

func getType (record []byte) (string, error) {
	var generic FAgeneric

	// Check input record
	if err := json.Unmarshal(record, &generic); err != nil {
		return "", errors.New(fmt.Sprintf("Unable to decode %v: %v\n", record, err))
	}
	return generic.Type, nil
}

func decodeRecord(record []byte) (interface{}) {
	var generic interface{}

	// Check input record
	if err := json.Unmarshal(record, &generic); err != nil {
		return errors.New(fmt.Sprintf("Unable to decode %v: %v\n", record, err))
	}

	payload := generic.(map[string]interface{})
	switch payload["type"] {
	case "position":


	}
	return nil
}

func decodePosition(record map[string]interface{}) (FAposition, error) {

	generic := FAposition{}
	for k, v := range record {
		switch k {
		// mandatory fields
		case "type": generic.Type = v.(string)
		case "ident": generic.Ident = v.(string)
		case "lat": generic.Lat = v.(string)
		case "lon": generic.Lon = v.(string)
		case "clock": generic.Clock = v.(string)
		case "id": generic.Id = v.(string)
		case "updateType": generic.UpdateType = v.(string)
		case "air_ground": generic.AirGround = v.(string)
		case "facility_hash": generic.FacilityHash = v.(string)
		case "facility_name": generic.FacilityName = v.(string)
		case "alt": generic.Alt = v.(string)
		case "gs": generic.Gs = v.(string)
		case "heading": generic.Heading = v.(string)
		case "rp1lat": generic.Rp1Lat = v.(string)
		case "rp1lon": generic.Rp1Lon = v.(string)
		case "rp1alt": generic.Rp1Alt = v.(string)
		case "rp1clock": generic.Rp1Clock = v.(string)
		case "squawk": generic.Squawk = v.(string)
		case "hexid": generic.Hexid = v.(string)
		case "fob": generic.Fob = v.(string)
		case "oat": generic.Oat = v.(string)
		case "airspeed_kts": generic.AirspeedKts = v.(string)
		case "airspeed_mach": generic.AirspeedMach = v.(string)
		case "winds": generic.Winds = v.(string)
		case "eta": generic.Eta = v.(string)
		case "baro_alt": generic.BaroAlt = v.(string)
		case "gps_alt": generic.GpsAlt = v.(string)
		case "atcident": generic.Atcident = v.(string)
		}
	}
	return generic, nil
}

func decodeDeparture(record []byte) (FAdeparture, error) {

	generic := FAdeparture{}
	// Check input record
	if err := json.Unmarshal(record, &generic); err != nil {
		return generic, errors.New(fmt.Sprintf("Unable to decode %v: %v\n", record, err))
	}
	return generic, nil
}

func decodeArrival(record []byte) (FAarrival, error) {

	generic := FAarrival{}
	// Check input record
	if err := json.Unmarshal(record, &generic); err != nil {
		return generic, errors.New(fmt.Sprintf("Unable to decode %v: %v\n", record, err))
	}
	return generic, nil
}

func decodeCancellation(record []byte) (FAcancellation, error) {

	generic := FAcancellation{}
	// Check input record
	if err := json.Unmarshal(record, &generic); err != nil {
		return generic, errors.New(fmt.Sprintf("Unable to decode %v: %v\n", record, err))
	}
	return generic, nil
}

func decodeFlightplan(record []byte) (FAflightplan, error) {

	generic := FAflightplan{}
	// Check input record
	if err := json.Unmarshal(record, &generic); err != nil {
		return generic, errors.New(fmt.Sprintf("Unable to decode %v: %v\n", record, err))
	}
	return generic, nil
}
