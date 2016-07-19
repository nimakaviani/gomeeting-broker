package models

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	URL string `json:"url"`
}

func LoadConfig(filePath string) (Config, error) {
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		return Config{}, e
	}
	var config Config
	json.Unmarshal(file, &config)
	return config, nil
}
