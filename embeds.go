package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func nowPlaying(song song) *discordgo.MessageEmbed {
	var username string
	var avatar string
	var min int
	var sec int
	var tm string

	min = song.Time / 60
	sec = song.Time - min*60

	if min < 10 && sec < 10 {
		tm = "0" + strconv.Itoa(min) + ":0" + strconv.Itoa(sec)
	} else if min < 10 {
		tm = "0" + strconv.Itoa(min) + ":" + strconv.Itoa(sec)
	} else {
		tm = strconv.Itoa(min) + ":" + strconv.Itoa(sec)
	}

	if song.User != nil {
		username = song.User.Username
		avatar = song.User.AvatarURL("1024x1024")
	} else {
		username = "Playlist"
		avatar = "https://cdn.discordapp.com/avatars/344604921485721610/eadc5467d4981aa9051356cdd3ee3673.png?size=2048"
	}
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Song added by " + username,
			IconURL: avatar,
		},
		Color: 0x2eaae5,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Timestamp",
				Value:  tm,
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: song.Thumbnail,
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Title:     "**" + song.Name + "**",
		URL:       song.URL,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "GoBot OpenSource MB by TheSnekySnek",
		},
	}
	fmt.Println("Send")
	return embed
}

func getQueue(songs []song) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0x2eaae5,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: songs[0].Thumbnail,
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "GoBot OpenSource MB by TheSnekySnek",
		},
	}

	for i := 0; i < len(songs); i++ {
		if i < 23 {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   strconv.Itoa(i) + ". " + songs[i].Name,
				Value:  "Added by " + songs[i].User.Username + "\n",
				Inline: false,
			})
		} else {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "And " + strconv.Itoa(len(songs)-i) + " more",
				Value:  "...",
				Inline: true,
			})
		}
	}

	fmt.Println("Send")
	return embed
}
