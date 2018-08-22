// config.go
//
// Copyright 2015-2018 © by Ollivier Robert <roberto@keltia.net>
//

/*
Looks into a TOML file for configuration options and returns a Config
struct.

	rc := config.LoadConfig("foo.toml")

rc will be serialized from TOML.

TOML: https://github.com/naoina/toml
*/
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/naoina/toml"
	"github.com/pkg/errors"
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
	c := Config{}

	if file == "" {
		file = filepath.Join(baseDir, configName)
	}

	// Check if there is any config file
	if _, err := os.Stat(file); err != nil {
		return nil, errors.Wrap(err, "LoadConfig")
	}

	verbose("file=%s, found it", file)

	// Read it
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrapf(err, "reading %s", file)
	}

	err = toml.Unmarshal(buf, &c)
	if err != nil {
		return nil, errors.Wrapf(err, "parse error %s", file)
	}

	return &c, err
}
