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

func regQueue(q []song) {
	sv, err2 := json.Marshal(q)
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println("Q SAVE")
	err3 := ioutil.WriteFile("./config/queue.json", sv, 0644)
	if err3 != nil {
		fmt.Println(err3)
	}
	return
}

func initQueue() ([]song, error) {
	var q []song
	file, err := ioutil.ReadFile("./config/queue.json")
	if err != nil {
		return q, err
	}
	json.Unmarshal(file, &q)
	return q, nil
}

func clearQueue() {
	err3 := ioutil.WriteFile("./config/queue.json", make([]byte, 0), 0644)
	if err3 != nil {
		fmt.Println(err3)
	}
	return
}
