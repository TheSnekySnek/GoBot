package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func commandHandler(m *discordgo.MessageCreate) {
	var args = strings.Split(m.Content, " ")
	cont := strings.Replace(m.Content, args[0]+" ", "", -1)

	if args[0] == "!play" {
		if _, err := strconv.Atoi(cont); err == nil {
			i, _ := strconv.Atoi(cont)
			if i < len(pl.Songs) {
				playYT(pl.Songs[i].URL, true, m.Author, func(sn song) {})
			} else {
				session.ChannelMessageSend(m.ChannelID, "Invalid Song ID")
			}
		} else if strings.Contains(cont, "spotify.com") {
			session.ChannelMessageSend(m.ChannelID, "Converting Spotify playlist... This can take a while")
			searchSpotify(cont, m)
		} else if strings.HasPrefix(cont, "http") {
			playYT(cont, true, m.Author, func(sn song) {})
		} else {
			session.ChannelMessageSend(m.ChannelID, "Searching for "+cont)
			searchYT(cont, func(link string) {
				playYT(link, true, m.Author, func(sn song) {})
			})
		}
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
		session.ChannelMessageSend(m.ChannelID, "Skipping song...")
		skip <- true
	}
	if args[0] == "!np" || args[0] == "!song" {
		em := nowPlaying(curSong)
		_, err := session.ChannelMessageSendEmbed(m.ChannelID, em)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
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
		if !contains(clearMan, m.Author.ID) {
			clearMan = append(clearMan, m.Author.ID)
			p := usersInVC()
			if len(clearMan) > (p-1)/2 {
				queue = queue[:0]
				session.ChannelMessageSend(m.ChannelID, "The queue has been cleared")
				clearMan = clearMan[:0]

			} else {
				session.ChannelMessageSend(m.ChannelID, strconv.Itoa(len(skipMan))+"/"+strconv.Itoa((p-1)/2)+" Votes to clear")
			}
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
				lyrArr := strings.Split(lyrics, "\n\n")
				for i := 0; i < len(lyrArr); i++ {
					session.ChannelMessageSend(m.ChannelID, lyrArr[i])
				}
				session.ChannelMessageSend(m.ChannelID, "\n\nQuery took **"+execTime+"ms** to execute")
			}
		})
	}
	if args[0] == "!volume" {
		if val, err := strconv.Atoi(cont); err == nil {
			if val > -1 && val < 101 {
				volume = float64(val)
				fmt.Println(volume)
			} else {
				session.ChannelMessageSend(m.ChannelID, "Volume must be between 0 and 100")
			}
		} else {
			session.ChannelMessageSend(m.ChannelID, "Volume must be a number between 0 and 100")
		}
	}
	if args[0] == "!add" {
		if isMod(m.Author.ID) {
			fmt.Println(cont)
			if strings.HasPrefix(cont, "http") {
				playYT(cont, false, m.Author, func(sn song) {
					addPlaylist(sn)
					session.ChannelMessageSend(m.ChannelID, sn.Name+" has been added to the playlist")
				})
			} else {
				searchYT(cont, func(link string) {
					playYT(link, false, m.Author, func(sn song) {
						addPlaylist(sn)
						session.ChannelMessageSend(m.ChannelID, sn.Name+" has been added to the playlist")
					})
				})
			}
		}
	}
	if args[0] == "!delete" {
		if isMod(m.Author.ID) {
			num, err := strconv.Atoi(cont)
			if err != nil {
				session.ChannelMessageSend(m.ChannelID, "Please provide a number")
				return
			}
			removePlaylist(num)
		}
	}

	if args[0] == "!playlist" {
		displayPlaylist(session, m)
	}

}
