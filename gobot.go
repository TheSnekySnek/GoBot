package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var curVC *discordgo.VoiceConnection
var isPlaying = false
var skip = make(chan bool)
var stop = make(chan bool)
var session *discordgo.Session
var curSong song
var config configuration
var queue []song
var pl playlist
var modifier = 0
var firstBoot = true
var skipMan []string
var clearMan []string
var listeners = false
var shouldPlay = false
var volume = 100.0

func main() {

	//We load the config file
	fmt.Println("Loading Config...")
	conf, err := loadConfig()
	if err != nil {
		fmt.Println("Error reading Config")
		return
	}
	config = conf

	//Create a new Session
	fmt.Println("Starting Session...")
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Error creating Session")
		return
	}
	session = dg

	//Create a new Handler for messages
	dg.AddHandler(messageCreate)

	//Open the Session
	fmt.Println("Openning Session...")
	err = dg.Open()
	if err != nil {
		fmt.Println("Error starting Session")
		return
	}

	//Join the Voice Channel

	fmt.Println("Joining VC...")
	vc, err := dg.ChannelVoiceJoin(config.Guild, config.VC, false, false)
	if err != nil {
		fmt.Println("Error joining Voice")
		return
	}

	//Assign our current Voice Connection
	curVC = vc

	//Load the Queue
	qu, err := initQueue()
	if err == nil {
		queue = qu
	}

	//load the playlist
	fmt.Println("Loading Playlist...")
	//Set the Seed
	rand.Seed(time.Now().UTC().UnixNano())
	pl, err = getPlaylist()
	/*if len(queue) > 0 {
		fmt.Println("Playing Queue")
		curSong, queue = queue[0], queue[1:]
		regQueue(queue)
		go play(curSong.VDURL, modifier)
	} else if len(pl.Songs) > 0 {
		go playYT(pl.Songs[rand.Intn(len(pl.Songs))].URL, true, pl.Songs[0].User, func(sn song) {})
	} else {
		session.ChannelMessageSend(config.TC, "Playlist is empty. Use !add [url] to add a song")
	}*/

	go presenceHandler()

	fmt.Println("Bot is now running")
	<-stop
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.ChannelID != config.TC || !strings.HasPrefix(m.Content, "!") || !isUserInVC(s, m.Author.ID) {
		return
	}
	commandHandler(m)
}
