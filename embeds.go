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
	var out1 time.Time
	diff := time.Now().Sub(song.Time)
	if diff > song.Duration {
		out1 = time.Time{}.Add(song.Duration)
	} else {
		out1 = time.Time{}.Add(diff)
	}
	out2 := time.Time{}.Add(song.Duration)

	if song.User != nil {
		username = song.User.Username
		avatar = song.User.AvatarURL("1024x1024")
	} else {
		username = "Playlist"
		avatar = "https://cdn.discordapp.com/avatars/344604921485721610/ed0b0758b187107d09a9fecfe243ec9d.jpg?size=1024"
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
				Value:  out1.Format("04:05") + " / " + out2.Format("04:05"),
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

func displayPlaylist(s *discordgo.Session, m *discordgo.MessageCreate) {

	pl, err := getPlaylist()
	if err != nil {
		return
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0x2eaae5,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: pl.Songs[0].Thumbnail,
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "GoBot OpenSource MB by TheSnekySnek",
		},
	}

	for i := 0; i < len(pl.Songs); i++ {
		if i%23 == 0 && i != 0 {
			session.ChannelMessageSendEmbed(m.ChannelID, embed)
			embed = &discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{},
				Color:  0x2eaae5,
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: pl.Songs[0].Thumbnail,
				},
				Timestamp: time.Now().Format(time.RFC3339),
				Footer: &discordgo.MessageEmbedFooter{
					Text: "GoBot OpenSource MB by TheSnekySnek",
				},
			}
		} else {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   strconv.Itoa(i) + ". " + pl.Songs[i].Name,
				Value:  "Added by " + pl.Songs[i].User.Username + "\n",
				Inline: false,
			})
		}
	}
	session.ChannelMessageSendEmbed(m.ChannelID, embed)
	return
}
