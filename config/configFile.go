package config

import (
	"encoding/json"
	"fmt"
	"gorgon/pkg/schemas"
	"os"
)

func InitializeConfigFile() error {
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		config := schemas.ConfigFile{
			ProwlarrApiKey: "",
		}

		file, err := os.Create(ConfigPath)
		if err != nil {
			return err
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		return encoder.Encode(config)
	}

	fmt.Println("Configuration file already exist.")
	return nil
}

func LoadConfig() (*schemas.ConfigFile, error) {
	file, err := os.Open(ConfigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var configFile schemas.ConfigFile
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configFile)
	if err != nil {
		return nil, err
	}

	return &configFile, nil
}
