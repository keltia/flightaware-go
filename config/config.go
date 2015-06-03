// config.go
//
// My homemade configuration class

package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"log"
)

type Dest struct {
	Broker string
	Name   string
}

type Config struct {
	User     string
	Password string
	Site     string
	Port     string
	Dests    map[string]Dest
	Default  string
	Feed_one func([]byte)
}

func (dest *Dest) String() string {
	return fmt.Sprintf("%v: %v", dest.Broker, dest.Name)
}

func LoadConfig(file string) (Config, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return Config{}, err
	}

	c := new(Config)
	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		log.Println("Error parsing yaml")
		return Config{}, err
	}

	c.Default = "mine"
	c.Feed_one = func(buf []byte) { fmt.Println(buf)}
	return *c, err
}
