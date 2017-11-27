package rthk

import (
	"github.com/l1z2g9/go-quote2/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// GetCitySnapshots ...
func GetCitySnapshots(category string) []byte {
	tmpl := "https://drive.google.com/folderview?id=%s&usp=sharing&tid=0B1tCMe1zn1gHTjZVVEk0YWZzYm8"
	tmpl2 := "https://drive.google.com/drive/folders/%s?tid=0B1tCMe1zn1gHTjZVVEk0YWZzYm8"

	var urlItems []*UrlItem
	if "India" == category {
		url := fmt.Sprintf(tmpl, "0B1tCMe1zn1gHUUlmcGU4R1dHS1k")
		url2 := fmt.Sprintf(tmpl2, "0B1tCMe1zn1gHUVNTd0pYaDl2SGM")
		urlItems = getDrivePath(url, url2)
	} else if "Japan" == category {
		url := fmt.Sprintf(tmpl, "0B1tCMe1zn1gHRDFiVjFBck9KTWM")
		url2 := fmt.Sprintf(tmpl2, "0B1tCMe1zn1gHUFZ2LVllcFVzdzg")
		urlItems = getDrivePath(url, url2)
	} else if "England" == category {
		url := fmt.Sprintf(tmpl, "0B1tCMe1zn1gHRUxPWUtUanJCYVk")
		url2 := fmt.Sprintf(tmpl2, "0B1tCMe1zn1gHbXRxNlJvZzBMRTQ")
		urlItems = getDrivePath(url, url2)
	} else if "Australia" == category {
		url := fmt.Sprintf(tmpl, "0B1tCMe1zn1gHWDdwYXBQMWdNOVE")
		urlItems = getDrivePath(url)
	} else if "New_Zealand" == category {
		url := fmt.Sprintf(tmpl, "0B1tCMe1zn1gHYkRuajNnWjFxb3c")
		urlItems = getDrivePath(url)
	} else if "USA" == category {
		url := fmt.Sprintf(tmpl, "0B1tCMe1zn1gHeFNiT1gxNFlVeUU")
		urlItems = getDrivePath(url)
	}

	//urlItems := getDrivePath(urls)

	db := util.GetDB()
	defer db.Close()

	re := regexp.MustCompile("City_Snapshot-(.+?)-.+") // Date

	var title, description string
	for _, item := range urlItems {
		match := re.FindStringSubmatch(item.Title)

		db.QueryRow("select title, description from RTHK_Radio where type = 'City_Snapshot' and strftime('%Y%m%d',date) = ?",
			match[1]).Scan(&title, &description)
		item.Description = description
	}

	dat, _ := json.Marshal(urlItems)
	return dat
}

// UrlItem ...
type UrlItem struct {
	Seq         int
	Title       string
	Description string
	URL         string
}

func getDrivePath(urls ...string) []*UrlItem {
	var urlItems []*UrlItem
	for _, url := range urls {
        util.Info.Println("Get path: ", url)

		resp, _ := http.Get(url)
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		list := string(body)

		re := regexp.MustCompile(`window\['_DRIVE_ivd'\] = '\[\[\[(.+)1\]\\n';`)
		re2 := regexp.MustCompile(`\\x22(.+?)\\x22`)

		items := re.FindStringSubmatch(list)[1]
		i := 0
		var prevTitle string
		for _, item := range strings.Split(items, `]\n,[`) {
			x := re2.FindAllStringSubmatch(item, -1)
			if len(x) > 2 {
				title := x[2][1]
				url := fmt.Sprintf("https://drive.google.com/uc?id=%s&authuser=0&export=download", x[0][1])
				if strings.Contains(title, "mp3") {
					i = i + 1
					urlItem := &UrlItem{i, title, "", url}
					urlItems = append(urlItems, urlItem)
				}
				if strings.Contains(title, "(link).txt") {
					resp, _ := http.Get(url)
					defer resp.Body.Close()

					body, _ := ioutil.ReadAll(resp.Body)

					urlItem := &UrlItem{i, title, "", string(body)}
					urlItems = append(urlItems, urlItem)
				}
				if strings.Contains(title, "txt") {
					title = strings.Replace(title, "txt", "", -1)
					if strings.Contains(prevTitle, title) {
						urlItems[len(urlItems)-1].Description = url
					}
				}
				prevTitle = title
			}
		}
	}

	return urlItems
}

func GetYCantonese() []byte {
	url := "https://drive.google.com/folderview?id=0B1tCMe1zn1gHbUxHcDdBUkV0aTQ&usp=sharing"
	url2 := "https://drive.google.com/drive/folders/0B1tCMe1zn1gHRjQyZDBCM3ZUSnM?usp=sharing"
	url3 := "https://drive.google.com/folderview?id=0B1tCMe1zn1gHaERIZWdGbXdNYzA&usp=sharing"
	url4 := "https://drive.google.com/folderview?id=0B1tCMe1zn1gHdEVRdlZUajNCYWc&usp=sharing"

	urlItems := getDrivePath(url, url2, url3, url4)

	dat, _ := json.Marshal(urlItems)
	return dat
}

func GetRthkRadio(category string) []byte {
	db := util.GetDB()
	defer db.Close()

	// no limit of 200 for this radio list due to the compresss with gzip
	rows, _ := db.Query("select title, description, strftime('%Y%m%d',date) from RTHK_Radio where type = ? order by date desc", category)
	defer rows.Close()

	var urlItems []UrlItem
	for rows.Next() {
		var title, description, date string
		rows.Scan(&title, &description, &date)
		urlItem := UrlItem{Title: title, Description: description, URL: date}
		urlItems = append(urlItems, urlItem)
	}
	dat, _ := json.Marshal(urlItems)
	return dat
}

func GetRthkPodcast(category string) []byte {
	db := util.GetDB()
	defer db.Close()

	rows, _ := db.Query("select Seq, Title, Subtitle, Summary, Guid from RTHK_PODCAST where type = ? order by seq desc", category)
	defer rows.Close()

	var urlItems []UrlItem
	for rows.Next() {
		var seq int
		var title, subtitle, summary, guid string
		rows.Scan(&seq, &title, &subtitle, &summary, &guid)
		desc := subtitle
		if len(desc) == 0 {
			desc = summary
		}

		urlItem := UrlItem{Seq: seq, Title: title, Description: desc, URL: guid}
		urlItems = append(urlItems, urlItem)
	}
	dat, _ := json.Marshal(urlItems)
	return dat
}

func GetRTHKPodcasDesc(seq string) string {
	db := util.GetDB()
	defer db.Close()

	var desc string
	db.QueryRow("select Description from RTHK_PODCAST where seq = ?", seq).Scan(&desc)

	return desc
}
