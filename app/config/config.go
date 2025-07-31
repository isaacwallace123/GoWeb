package config

import (
	"github.com/isaacwallace123/GoUtils/jsonutil"
	"github.com/isaacwallace123/GoWeb/app/types"
	"os"
	"strconv"
)

type ServerConfig struct {
	Port int `json:"port"`
}

var (
	Server ServerConfig
	Static []types.StaticConfig
)

type appConfig struct {
	Server ServerConfig         `json:"server"`
	Static []types.StaticConfig `json:"static"`
}

func LoadConfig(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var loaded appConfig
	if err := jsonutil.FromBytes(file, &loaded); err != nil {
		return err
	}
	Server = loaded.Server
	Static = loaded.Static
	return nil
}

// PortString returns ":8080" (or whatever garbage port you configured)
func PortString() string {
	return ":" + strconv.Itoa(Server.Port)
}
