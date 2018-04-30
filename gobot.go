package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rylio/ytdl"
)

var curVC *discordgo.VoiceConnection
var isPlaying = false
var skip = make(chan bool)
var skop = make(chan bool)
var session *discordgo.Session
var curSong song
var config configuration
var queue []song
var pl playlist
var modifier = 0
var firstBoot = true

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
	<-skop
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.ChannelID != config.TC {
		return
	}
	if !isUserInVC(s, m.Author.ID) {
		return
	}
	var args = strings.Split(m.Content, " ")
	cont := strings.Replace(m.Content, "!play ", "", -1)

	if args[0] == "!play" {
		if strings.HasPrefix(cont, "http") {
			playYT(cont, m)
		} else {
			session.ChannelMessageSend(m.ChannelID, "Searching for "+cont)
			searchYT(cont, func(link string) {
				playYT(link, m)
			})
		}
	}
	if args[0] == "!close" {
		if isMod(m.Author.ID) {
			s.Close()
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
			session.ChannelMessageSend(m.ChannelID, "The queue is empty\n\nType !play [song] to play a song")
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
	if args[0] == "!nightcore" {
		modifier = 1
		session.ChannelMessageSend(m.ChannelID, "NIGHTCORE MODE - ACTIVATED")
		session.ChannelMessageSend(m.ChannelID, "https://media2.giphy.com/media/A6LvcKNhL4fbG/giphy.gif")
	}
	if args[0] == "!daycore" {
		modifier = 2
		session.ChannelMessageSend(m.ChannelID, "DAYCORE MODE - ACTIVATED")
		session.ChannelMessageSend(m.ChannelID, "https://media1.giphy.com/media/3NtY188QaxDdC/giphy.gif")
	}
	if args[0] == "!midcore" {
		modifier = 0
		session.ChannelMessageSend(m.ChannelID, "SYSTEMS BACK TO NORMAL CAPTAIN")
	}
	if args[0] == "!lyrics" {
		searchLyrics(curSong.Name, func(lyrics string, execTime string) {
			if len(lyrics) < 2000 {
				session.ChannelMessageSend(m.ChannelID, lyrics)
				session.ChannelMessageSend(m.ChannelID, "\n\nQuery took **"+execTime+"ms** to execute")
			} else {
				fmt.Println(lyrics)
				lyrArr := strings.Split(lyrics, "\n\n")
				for i := 0; i < len(lyrArr); i++ {
					session.ChannelMessageSend(m.ChannelID, lyrArr[i])
				}
				session.ChannelMessageSend(m.ChannelID, "\n\nQuery took **"+execTime+"ms** to execute")
			}
		})
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

			nSong.Duration = video.Duration

			nSong.Name = video.Title
			if m != nil {
				nSong.User = m.Author
			} else {
				nSong.User = nil
			}
			if isPlaying {
				queue = append(queue, nSong)
				session.ChannelMessageSend(config.TC, nSong.Name+" has been added to the queue")
			} else {
				curSong = nSong
				go play(url, modifier)
			}
			return
		}
	}
}

func play(url string, mod int) {
	if !isPlaying {
		isPlaying = true
		curSong.Time = time.Now()
		PlayAudioFile(curVC, url, mod, skip)
		fmt.Println("Player stopped")
		isPlaying = false
		time.Sleep(2 * time.Second)
		if len(queue) > 0 {
			fmt.Println("Playing Queue")
			curSong, queue = queue[0], queue[1:]
			time.Sleep(2 * time.Second)
			go play(curSong.VDURL, modifier)
			//session.ChannelMessageSend(config.TC, "Playing "+curSong.Name)
		} else {
			fmt.Println("Playing Pl")
			playYT(pl.Songs[rand.Intn(len(pl.Songs))], nil)
		}
	}
}

func isMod(userID string) bool {
	fmt.Println(config.Mods[0])
	for i := 0; i < len(config.Mods); i++ {
		if userID == config.Mods[i] {
			return true
		}
	}
	return false
}

func isUserInVC(session *discordgo.Session, userid string) bool {
	for _, guild := range session.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == userid {
				if vs.ChannelID == config.VC {
					return true
				}
				return false
			}
		}
	}
	return false
}
