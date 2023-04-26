package config

import (
	"bufio"
	"io"
	"os"

	"gopkg.in/yaml.v3"

	"easy-master-slave/pkg/model"
)

func ParseConfig(configPath string) (*model.Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var config model.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
