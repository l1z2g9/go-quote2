package rthk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
)

func ExportFeedForCitySnap() string {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "RTHK City Snapshot",
		Link:        &feeds.Link{Href: "http://www.rthk.hk/radio/radio1/programme/City_Snapshot"},
		Description: "我們邀請旅居世界各地的名人為節目撰稿及以廣東話聲音演繹，以感性角度去分析他們身處的國家時事。 ",
		Author:      &feeds.Author{Name: "RTHK"},
		Created:     now,
	}

	var items []*feeds.Item

	urlTmpl := "http://www.rthk.hk/radio/catchUp?c=radio1&p=City_Snapshot&page=%d&m="
	for i := 1; i < 2; i++ {
		url := fmt.Sprintf(urlTmpl, i)
		fmt.Println(url)
		resp, _ := http.Get(url)
		defer resp.Body.Close()

		dat, _ := ioutil.ReadAll(resp.Body)

		var audioRow AudioRow
		if err := json.Unmarshal(dat, &audioRow); err != nil {
			panic(err)
		}

		for _, i := range audioRow.Content {
			text := fmt.Sprintf("%s %s %s %s %s\n", i.ID, i.Title, i.Date, i.Part, i.Photos)
			fmt.Println(text)

			t, err := time.Parse("02/01/2006", i.Date)
			if err != nil {
				log.Fatal(err)
			}

			item := &feeds.Item{
				Title:   i.Title,
				Link:    &feeds.Link{Href: fmt.Sprintf("http://stmw3.rthk.hk/aod/_definst_/radio/archive/radio1/City_Snapshot/mp3/mp3:%s.mp3/playlist.m3u8", t.Format("20060102"))},
				Created: t,
			}
			items = append(items, item)
		}
	}

	feed.Items = items

	atom, err := feed.ToRss()
	if err != nil {
		log.Fatal(err)
	}

	return atom
}

type AudioRow struct {
	Status  string `json:"status"`
	Content []struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Date   string `json:"date"`
		Photos struct {
			Photo []interface{} `json:"photo"`
		} `json:"photos"`
		Part []string `json:"part"`
	} `json:"content"`
	NextPage int `json:"nextPage"`
}
