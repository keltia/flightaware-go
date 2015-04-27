// config.go
//
// My homemade configuration class

import (
	"github.com/go-yaml/yaml"
)

type Config struct{
	user string
	password string
	site string
	port int
	dests map[]
}

func LoadConfig(file string) map[string]interface{} {

}
