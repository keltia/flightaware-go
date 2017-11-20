// config.go
//
// Copyright 2015 © by Ollivier Robert <roberto@keltia.net>
//

/*
Package config implement my homemade configuration class

Looks into a TOML file for configuration options and returns a config.Config
struct.

	import "config"

	rc := config.LoadConfig("foo.toml")

rc will be serialized from TOML.

TOML: https://github.com/naoina/toml
*/
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/naoina/toml"
)

// Dest is a output
type Dest struct {
	Broker string
	Name   string
	Type   string
}

// User holds credentials
type User struct {
	User     string
	Password string
}

// Config is the main object
type Config struct {
	Site    string
	Port    int
	DefUser string
	DefDest string
	Users   map[string]User
	Dests   map[string]Dest
}

// String is a basic Stringer for Dest
func (dest *Dest) String() string {
	return fmt.Sprintf("%v: %v", dest.Broker, dest.Name)
}

// LoadConfig loads a file as a YAML document and return the structure
func LoadConfig(file string) (*Config, error) {
	var sFile string

	// Check for tag
	if !strings.HasSuffix(file, ".toml") {
		// file must be a tag so add a "."
		sFile = filepath.Join(os.Getenv("HOME"),
			fmt.Sprintf(".%s", file),
			"config.toml")
	} else {
		sFile = file
	}

	c := new(Config)
	buf, err := ioutil.ReadFile(sFile)
	if err != nil {
		return c, errors.New(fmt.Sprintf("Can not read %s file.", sFile))
	}

	err = toml.Unmarshal(buf, &c)
	if err != nil {
		return c, errors.New(fmt.Sprintf("Can not parse %s: %v", sFile, err))
	}

	return c, err
}
