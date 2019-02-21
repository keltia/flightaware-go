// utils.go
//
// Copyright 2018 © by Ollivier Robert <roberto@keltia.net>

package flightaware

// debug displays only if fDebug is set
func (cl *FAClient) debug(str string, a ...interface{}) {
	if cl.level >= 2 {
		cl.Log.Printf(str, a...)
	}
}

// debug displays only if fVerbose is set
func (cl *FAClient) verbose(str string, a ...interface{}) {
	if cl.level >= 1 {
		cl.Log.Printf(str, a...)
	}
}
