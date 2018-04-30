package main

import (
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
