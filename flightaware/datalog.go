// datalog.go

/*
  This package implements the DataLog(), a payload-specific version of log.Printf().
 */
package flightaware

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)


type payload struct {
	Clock string
	Rest  interface{}
}

// clone of log.Printf() with data-specific time
func DataLog(payload []byte, str string) {
	var data payload

	// Parse json payload
	if err := json.Unmarshal(payload, &data); err != nil {
		log.Printf("Error: decoding %v: %v\n", data, err)
	}

	// string -> []byte
	datePkt, err := strconv.ParseInt(data.Clock, 10, 64)
	if err != nil {
		log.Printf("Error: parsing %v: %v\n", data.Clock, err)
	}

	// Now log
	pktTime := time.Unix(datePkt, 0)
	strTime := pktTime.Format("2006/01/02 15:04:05")
	fmt.Fprintf(os.Stderr, "%s %s", strTime, str)
}
