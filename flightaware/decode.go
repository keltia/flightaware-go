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

func decodeRecord(record []byte) (interface{}, error) {
	var generic interface{}

	// Check input record
	if err := json.Unmarshal(record, &generic); err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to decode %v: %v\n", record, err))
	}

	payload := generic.(map[string]interface{})
	switch payload["type"] {
	case "position":
        return decodePosition(payload)
	case "flightplan":
		return decodeFlightplan(payload)
	case "departure":
		return decodeDeparture(payload)
	case "arrival":
		return decodeArrival(payload)
	case "cancellation":
		return decodeCancellation(payload)
	default:
		return nil, errors.New(fmt.Sprintf("Unknown record type %s", payload["type"]))
	}
	return nil, nil
}

func decodeFlightplan(record map[string]interface{}) (FAflightplan, error) {

	generic := FAflightplan{}
	for k, v := range record {
		switch k {
		// mandatory fields
		case "type": generic.Type = v.(string)
		case "ident": generic.Ident = v.(string)
		case "status": generic.Status = v.(string)
		case "orig": generic.Orig = v.(string)
		case "dest": generic.Dest = v.(string)
		case "edt": generic.Edt = v.(string)
		case "eta": generic.Eta = v.(string)
		case "ete": generic.Ete = v.(string)
		case "id": generic.Id = v.(string)
		// common
		case "aircrafttype": generic.AircraftType = v.(string)
		case "suffix": generic.Suffix = v.(string)
		case "reg": generic.Reg = v.(string)
		case "speed": generic.Speed = v.(string)
		case "alt": generic.Alt = v.(string)
		case "facility_hash": generic.FacilityHash = v.(string)
		case "facility_name": generic.FacilityName = v.(string)
		// optional
		case "prefix": generic.Prefix = v.(string)
		case "waypoints": generic.Waypoints = v.([]FApoint)
		case "FDwaypoints": generic.FDWaypoints = v.([]FApoint)
		case "route": generic.Route = v.(string)
		case "atcident": generic.Atcident = v.(string)
		}
	}
	return generic, nil
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

func decodeDeparture(record map[string]interface{}) (FAdeparture, error) {

	generic := FAdeparture{}
	for k, v := range record {
		switch k {
		// mandatory fields
		case "type": generic.Type = v.(string)
		case "ident": generic.Ident = v.(string)
		case "orig": generic.Orig = v.(string)
		case "dest": generic.Dest = v.(string)
		case "adt": generic.Eta = v.(string)
		case "eta": generic.Eta = v.(string)
		case "id": generic.Id = v.(string)
		// common fields
		case "aircrafttype": generic.AircraftType = v.(string)
		case "facility_hash": generic.FacilityHash = v.(string)
		case "facility_name": generic.FacilityName = v.(string)
		// optional
		case "synthetic": generic.Synthetic = v.(string)
		case "atcident": generic.Atcident = v.(string)
		}
	}
	return generic, nil
}

func decodeArrival(record map[string]interface{}) (FAarrival, error) {

	generic := FAarrival{}
	for k, v := range record {
		switch k {
		// mandatory fields
		case "type": generic.Type = v.(string)
		case "ident": generic.Ident = v.(string)
		case "orig": generic.Orig = v.(string)
		case "dest": generic.Dest = v.(string)
		case "aat": generic.Aat = v.(string)
		case "timeType": generic.timeType = v.(string)
		case "id": generic.Id = v.(string)
		// common
		case "facility_hash": generic.FacilityHash = v.(string)
		case "facility_name": generic.FacilityName = v.(string)
		// optional
		case "synthetic": generic.Synthetic = v.(string)
		case "atcident": generic.Atcident = v.(string)
		}
	}
	return generic, nil
}

func decodeCancellation(record map[string]interface{}) (FAcancellation, error) {

	generic := FAcancellation{}
	for k, v := range record {
		switch k {
		// mandatory fields
		case "type": generic.Type = v.(string)
		case "ident": generic.Ident = v.(string)
		case "orig": generic.Orig = v.(string)
		case "dest": generic.Dest = v.(string)
		case "id": generic.Id = v.(string)
		// common
		case "facility_hash": generic.FacilityHash = v.(string)
		case "facility_name": generic.FacilityName = v.(string)
		// optional
		case "atcident": generic.Atcident = v.(string)
		}
	}
	return generic, nil
}
