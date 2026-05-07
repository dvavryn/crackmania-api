package config

import (
	"encoding/json"
	"io"
	"os"
)

func GetConfig() (Config, error) {
	jsonFile, err := os.Open("./config.json")
	if err != nil {
		return Config{}, err
	}
	defer jsonFile.Close()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return Config{}, err
	}

	var out Config
	if err = json.Unmarshal(bytes, &out); err != nil {
		return Config{}, err
	}

	return out, nil
}
