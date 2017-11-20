// datalog.go

package flightaware

/*
Implements the DataLog(), a payload-specific version of log.Printf().
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// Payload is our main object
type Payload struct {
	Clock string
	Rest  interface{}
}

// DataLog is a clone of log.Printf() with data-specific time
func DataLog(buf []byte, str string) {
	var data Payload

	// Parse json payload
	if err := json.Unmarshal(buf, &data); err != nil {
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
