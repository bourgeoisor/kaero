package application

import (
	"encoding/json"
	"io"
	"kaero/utils"
	"os"
)

func readConfig(path string) (*utils.Config, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	configBytes, err := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	config := &utils.Config{}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, err
	}

	setDefaultConfig(config)

	return config, nil
}

func setDefaultConfig(config *utils.Config) {
	if config.MaxNickLength == 0 {
		config.MaxNickLength = 16
	}
}
