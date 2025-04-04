package config

import (
	"encoding/json"
	"fmt"
	"gorgon/pkg/schemas"
	"os"
	"reflect"
)

func InitializeOrUpdateConfigFile() error {
	var updated bool
	config := schemas.ConfigFile{
		ProwlarrApiKey:             "",
		ProwlarrHost:               "",
		ProwlarrPort:               "",
		QBittorrentHost:            "",
		QBittorrentPort:            "",
		QBittorrentUsername:        "",
		QBittorrentPassword:        "",
		QBittorrentDownloadFolder:  "",
		DefaultShowInfoFolder:      "assets/shows",
		DefaultInstalledShowFolder: "$HOME/Videos/show",
	}

	// Verifica se existe o arquivo
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		// Não existe: cria com os valores default
		fmt.Println("Creating new config file...")
		return saveConfig(config)
	}

	// Já existe: carrega o existente
	existingConfig, err := LoadConfig()
	if err != nil {
		return err
	}

	// Preenche os campos faltantes
	existingVal := reflect.ValueOf(existingConfig).Elem()
	defaultVal := reflect.ValueOf(config)
	for i := 0; i < existingVal.NumField(); i++ {
		field := existingVal.Field(i)
		if field.Kind() == reflect.String && field.String() == "" {
			field.SetString(defaultVal.Field(i).String())
			updated = true
		}
	}

	if updated {
		fmt.Println("Updating config file with missing fields...")
		return saveConfig(*existingConfig)
	}

	fmt.Println("Config file already up to date.")
	return nil
}

func saveConfig(config schemas.ConfigFile) error {
	file, err := os.Create(ConfigPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
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
