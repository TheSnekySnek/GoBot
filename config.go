package main

import (
	"encoding/json"
	"io/ioutil"
)

type configuration struct {
	Guild string `json:"Guild"`
	VC    string `json:"VC"`
	TC    string `json:"TC"`
	Token string `json:"Token"`
}

func loadConfig() (configuration, error) {
	var conf configuration
	file, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		return conf, err
	}
	json.Unmarshal(file, &conf)

	return conf, nil
}
