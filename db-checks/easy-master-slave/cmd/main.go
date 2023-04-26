package main

import (
	"flag"

	"easy-master-slave/internal/apiserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "config/config.yaml", "Path to config file")
}

func main() {
	flag.Parse()

	if err := apiserver.Start(configPath); err != nil {
		panic(err)
	}
}
