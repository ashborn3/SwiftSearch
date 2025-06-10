package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	HomePath      string `json:"homePath"`
	CachePath     string `json:"cachePath"`
	EncryptionKey string `json:"encryptionKey"`
	Ip            string `json:"ip"`
	Port          int    `json:"port"`
	SyncTime      int    `json:"syncTime"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
