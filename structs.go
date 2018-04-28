package main

import (
	"github.com/bwmarrin/discordgo"
)

type song struct {
	Name      string
	URL       string
	VDURL     string
	Thumbnail string
	Time      int
	User      *discordgo.User
}
