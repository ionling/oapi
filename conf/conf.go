package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Server *Server
	Abbr   *Abbr
}

type Server struct {
	Addr string
}

type Abbr struct {
	UID     string
	TokenID string
}

func LoadBytes(bs []byte) (conf *Config, err error) {
	conf = &Config{}
	err = json.Unmarshal(bs, conf)
	return
}

func LoadFile(name string) (conf *Config, err error) {
	bs, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return LoadBytes(bs)
}
