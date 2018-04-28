package main

import (
	"encoding/json"
	"io/ioutil"
)

type playlist struct {
	Songs []string `json:"Songs"`
}

func getPlaylist() (playlist, error) {
	var pl playlist
	file, err := ioutil.ReadFile("./config/playlist.json")
	if err != nil {
		return pl, err
	}
	json.Unmarshal(file, &pl)
	return pl, nil
}
