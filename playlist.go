package main

import (
	"encoding/json"
	"io/ioutil"
)

type playlist struct {
	Songs []song `json:"Songs"`
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

func addPlaylist(s song) {
	var pl playlist
	file, err := ioutil.ReadFile("./config/playlist.json")
	if err != nil {
		return
	}
	json.Unmarshal(file, &pl)
	pl.Songs = append(pl.Songs, s)
	sv, err2 := json.Marshal(pl)
	if err2 != nil {
		return
	}
	err3 := ioutil.WriteFile("./config/playlist.json", sv, 0644)
	if err3 != nil {
		return
	}
	return
}
