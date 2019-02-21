package flightaware

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Private functions

// Default callback
func defaultFeed(buf []byte) { fmt.Println(string(buf)) }

// Default filter
func defaultFilter(cl *FAClient, buf []byte) bool {
	if len(cl.OutputFilters) != 0 {
		for index, flt := range cl.OutputFilters {
			// First match so behaviour is OR
			if flt.Match(buf) {
				cl.verbose("%d", index)
				return true
			}
			cl.verbose(".")
		}
		// If no match from any of the filters, no cigar
		return false
	}
	return true
}

// consumer part of the FA cl
func (cl *FAClient) startWriter() (chan []byte, error) {
	var err error

	cl.verbose("Waiting for data…")
	ch := make(chan []byte, 1000)
	go func() {
		for {
			buf, ok := <-ch
			if !ok {
				err = errors.Errorf("reading %s: %v", string(buf), ok)
				return
			}
			// Do something
			cl.dataLog(buf, fmt.Sprintf("Read %d bytes\n", len(buf)))

			// Insert filter call
			if ok = (cl.Filter)(cl, buf); ok {
				(cl.FeedOne)(buf)
			}

			cl.Bytes += int64(len(buf))
			cl.Pkts++
		}
	}()
	return ch, err
}

// Payload is our main object
type Payload struct {
	Clock string
	Rest  interface{}
}

// dataLog is a clone of log.Printf() with data-specific time
func (cl *FAClient) dataLog(buf []byte, str string) {
	// Only if verbose or more
	if cl.level >= 1 {
		var data Payload

		// Parse json payload
		if err := json.Unmarshal(buf, &data); err != nil {
			cl.Log.Printf("Error: decoding %v: %v\n", data, err)
		}

		// string -> []byte
		datePkt, err := strconv.ParseInt(data.Clock, 10, 64)
		if err != nil {
			cl.Log.Printf("Error: parsing %v: %v\n", data.Clock, err)
		}

		// Now log
		pktTime := time.Unix(datePkt, 0)
		strTime := pktTime.Format("2006/01/02 15:04:05")
		cl.Log.Printf("%s %s", strTime, str)
	}
}
