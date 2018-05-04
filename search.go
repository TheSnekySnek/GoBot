package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly"
)

func searchYT(query string, callback func(string)) {
	query = strings.Replace(query, " ", "%20", -1)
	fmt.Println(query)
	c := colly.NewCollector()
	c.OnHTML("#results > ol > li > ol > li:nth-of-type(2) > div", func(e *colly.HTMLElement) {
		callback(e.Attr("data-context-item-id"))
		c = nil
	})

	c.Visit("https://www.youtube.com/results?search_query=" + query)
}

func searchSpotify(query string, m *discordgo.MessageCreate) {
	rep := false
	t1 := strings.LastIndex(query, "?")
	if t1 > 0 {
		query = query[0:t1]
	}
	fmt.Println(query)
	c := colly.NewCollector()
	c.SetRequestTimeout(90 * time.Second)
	c.OnHTML("#tubetify-generate tr > td > a", func(e *colly.HTMLElement) {
		playYT(strings.Replace(e.Text, "#", "", -1), false, m.Author, func(s song) {
			queue = append(queue, s)
			if !rep {
				session.ChannelMessageSend(config.TC, "Playlist has been added")
			}
			rep = true
		})
	})
	err := c.Post("https://tubetify.com/generate", map[string]string{"spotify-tracks": query})
	if err != nil {
		log.Fatal(err)
	}
}

func searchLyrics(query string, callback func(string, string)) {
	callTime := time.Now().UnixNano() / int64(time.Millisecond)
	step := 0
	found := 0
	lyrics := ""
	subLyr := 999
	t1 := strings.Index(query, "[")
	if t1 > 0 {
		query = query[0 : t1-1]
	}
	t2 := strings.Index(query, "(")
	if t2 > 0 {
		query = query[0 : t2-1]
	}
	t3 := strings.Index(query, "-")
	t4 := strings.LastIndex(query, "-")
	if t3 != t4 {
		query = query[0 : t4-1]
	}

	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		if step == 1 {
			aORb := regexp.MustCompile("mxm-lyrics__content")
			matches := aORb.FindAllStringIndex(string(r.Body), -1)
			subLyr = len(matches)
		}
	})
	c.OnHTML(".title[href^='/lyrics']", func(e *colly.HTMLElement) {
		if step == 0 {
			step++
			c.Visit("https://www.musixmatch.com" + e.Attr("href"))
		}
	})
	c.OnHTML("p.mxm-lyrics__content", func(e *colly.HTMLElement) {
		if step == 1 {
			lyrics += e.Text
			found++
			if subLyr == found {
				diff := (time.Now().UnixNano() / int64(time.Millisecond)) - callTime
				callback(lyrics, strconv.Itoa(int(diff)))
			}
		}
	})

	c.Visit("https://www.musixmatch.com/search/" + query)
	c = nil
}
