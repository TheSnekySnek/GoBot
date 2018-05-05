package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

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

func usersInVC() int {
	i := 0
	for _, guild := range session.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.ChannelID == config.VC {
				i++
			}
		}
	}
	return i
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func checkForListeners() {
	if usersInVC() < 2 && shouldPlay {
		fmt.Println("NOPPL")
		shouldPlay = false
		skip <- true
	} else if usersInVC() > 1 && !shouldPlay && !isPlaying {
		shouldPlay = true
		fmt.Println("PPL")
		if len(queue) > 0 {
			fmt.Println("Playing Queue")
			curSong, queue = queue[0], queue[1:]
			time.Sleep(2 * time.Second)
			go play(curSong.VDURL, modifier)
			//session.ChannelMessageSend(config.TC, "Playing "+curSong.Name)
		} else {
			fmt.Println("Playing Pl")
			pl, err := getPlaylist()
			if err != nil {
				fmt.Println("Error Getting PL")
				return
			}
			if len(pl.Songs) > 0 {
				playYT(pl.Songs[rand.Intn(len(pl.Songs))].URL, true, nil, func(sn song) {})
			} else {
				session.ChannelMessageSend(config.TC, "Playlist is empty. Use !add [url] to add a song")
			}

		}
	}
}

func presenceHandler() {
	for {
		time.Sleep(10 * time.Second)
		go checkForListeners()
	}
}
