package models

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	URL      string `json:"url"`
	Timezone string `json:"timezone"`
	Rooms    []Room
}

type Room struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Floor    string `json:"floor"`
	Capacity int    `json:"capacity"`
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
