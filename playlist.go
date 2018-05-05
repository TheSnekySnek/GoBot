package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/bwmarrin/discordgo"
)

type song struct {
	Name      string
	URL       string
	VDURL     string
	Thumbnail string
	Time      time.Time
	Duration  time.Duration
	User      *discordgo.User
}

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
	fmt.Println(pl.Songs[0].User.ID)
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

func removePlaylist(num int) {
	var pl playlist
	file, err := ioutil.ReadFile("./config/playlist.json")
	if err != nil {
		return
	}
	json.Unmarshal(file, &pl)
	pl.Songs = pl.Songs[:num+copy(pl.Songs[num:], pl.Songs[num+1:])] // Removes song from slice
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
