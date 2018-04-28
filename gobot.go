package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/rylio/ytdl"
)

var curVC *discordgo.VoiceConnection
var isPlaying = false
var skip = make(chan bool)
var session *discordgo.Session
var curSong song
var config configuration
var queue []song
var pl playlist

func main() {
	fmt.Println("Loading Config...")
	conf, err := loadConfig()
	if err != nil {
		fmt.Println("Error reading Config")
		return
	}
	config = conf

	fmt.Println("Starting Session...")
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Error creating Session")
		return
	}
	session = dg

	dg.AddHandler(messageCreate)

	fmt.Println("Openning Session...")
	err = dg.Open()
	if err != nil {
		fmt.Println("Error starting Session")
		return
	}

	fmt.Println("Joining VC...")
	vc, err := dg.ChannelVoiceJoin(config.Guild, config.VC, false, false)
	if err != nil {
		fmt.Println("Error joining Voice")
		return
	}

	curVC = vc

	fmt.Println("Loading Playlist...")
	pl, err = getPlaylist()

	playYT(pl.Songs[0], nil)

	fmt.Println("Bot is now running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.ChannelID != config.TC {
		return
	}
	var args = strings.Split(m.Content, " ")

	if args[0] == "!play" {

		playYT(args[1], m)
	}
	if args[0] == "!close" {
		if isMod(m.Author.ID) {
			session.Close()
			os.Exit(0)
		} else {
			session.ChannelMessageSend(m.ChannelID, "Denied")
		}
	}
	if args[0] == "!skip" {
		skip <- true
	}
	if args[0] == "!np" || args[0] == "!song" {
		em := nowPlaying(curSong)
		session.ChannelMessageSendEmbed(m.ChannelID, em)
	}
	if args[0] == "!queue" {
		if len(queue) > 0 {
			em := getQueue(queue)
			session.ChannelMessageSendEmbed(m.ChannelID, em)
		} else {
			session.ChannelMessageSend(m.ChannelID, "The queue is empty\n\nType !play [YT Url] to play a song")
		}

	}
	if args[0] == "!clear" {
		if isMod(m.Author.ID) {
			queue = queue[:0]
			session.ChannelMessageSend(m.ChannelID, "The queue has been cleared")
		} else {
			session.ChannelMessageSend(m.ChannelID, "Denied")
		}
	}

}

func playYT(link string, m *discordgo.MessageCreate) {
	video, err := ytdl.GetVideoInfo(link)

	if err != nil {
		session.ChannelMessageSend(m.ChannelID, "Error while accessing video")
		return
	}

	for _, format := range video.Formats {
		if format.AudioEncoding == "opus" || format.AudioEncoding == "aac" || format.AudioEncoding == "vorbis" {
			var nSong song
			nSong.URL = link
			data, err := video.GetDownloadURL(format)
			if err != nil {
				session.ChannelMessageSend(config.TC, err.Error())
			}
			url := data.String()
			nSong.VDURL = url
			data1 := video.GetThumbnailURL("default")
			nSong.Thumbnail = data1.String()

			nSong.Name = video.Title
			if m != nil {
				nSong.User = m.Author
			} else {
				nSong.User = nil
			}

			nSong.Time = 0
			if isPlaying {
				queue = append(queue, nSong)
				session.ChannelMessageSend(config.TC, nSong.Name+" has been added to the queue")
			} else {
				curSong = nSong
				go play(url)
				session.ChannelMessageSend(config.TC, "Playing "+curSong.Name)
			}
			watchTime()

			return
		}
	}
}

func watchTime() {
	if isPlaying {
		curSong.Time++
	}
	time.Sleep(time.Second)
	watchTime()
}

func play(url string) {
	if !isPlaying {
		fmt.Println("Speak")
		isPlaying = true
		dgvoice.PlayAudioFile(curVC, url, skip)
		isPlaying = false
		time.Sleep(time.Second)
		curSong.Time = 0
		if len(queue) > 0 {
			curSong, queue = queue[0], queue[1:]
			go play(curSong.VDURL)
			session.ChannelMessageSend(config.TC, "Playing "+curSong.Name)
		} else {
			playYT(pl.Songs[rand.Intn(len(pl.Songs))], nil)
		}
	}
}

func isMod(userID string) bool {
	for i := 0; i < len(config.Mods); i++ {
		if userID == config.Mods[i] {
			return true
		}
	}
	return false
}
