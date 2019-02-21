package flightaware

import (
	"fmt"

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
			if cl.Verbose {
				DataLog(buf, fmt.Sprintf("Read %d bytes\n", len(buf)))
			}

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
