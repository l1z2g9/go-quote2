package rthk

import (
	"github.com/l1z2g9/go-quote2/util"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const (
	urlPrefix_ch = "http://news.rthk.hk/rthk/ch/component/k2/"
	urlPrefix_en = "http://news.rthk.hk/rthk/en/component/k2/"
	prefix       = "http://programme.rthk.hk/channel/radio/"
)

func GetRthkNews(cat string) []byte {
	// NULL:all, 1:tou tiao, 2:greater china, 3:local, 4:world news, 5:finance, 6:sport
	var url string
	if cat == "1" {
		url = "http://m.rthk.hk/index.htm"
	} else {
		url = fmt.Sprintf("http://news.rthk.hk/rthk/webpageCache/services/loadModNewsShowSp2List.php?lang=zh-TW&cat=%s&newsCount=60&dayShiftMode=1", cat)
		util.Info.Println("cat", cat)
	}

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	//body, _ := ioutil.ReadFile("local.html")
	page := string(body)

	var news []RthkNews

	if cat == "1" {
		r := regexp.MustCompile(`(?s)<li><a href='(.*?)'>(.*?)</a>`)

		for _, x := range r.FindAllStringSubmatch(page, -1) {
			var lang string
			if strings.Contains(x[1], "/ch/") {
				url = strings.Replace(x[1], urlPrefix_ch, "", -1)
				lang = "ch"
			} else {
				url = strings.Replace(x[1], urlPrefix_en, "", -1)
				lang = "en"
			}

			news = append(news, RthkNews{x[2], url, "", lang})
		}
	} else {
		r := regexp.MustCompile(`(?s)<a href='(.+?)'>(.*?)</a>`)
		r2 := regexp.MustCompile(`(?s)<div class='ns2-created'>(.*?)</div>`)
		y := r2.FindAllStringSubmatch(page, -1)

		for i, x := range r.FindAllStringSubmatch(page, -1) {
			// url = http://news.rthk.hk/rthk/ch/component/k2/1270672-20160705.htm => 1270672-20160705.htm
			url := strings.Replace(x[1], urlPrefix_ch, "", -1)
			news = append(news, RthkNews{x[2], url, y[i][1], "ch"})
		}
	}

	d, _ := json.Marshal(&news)
	return d
}

/*func GetRthkNewsDetail(suffix string, lang string) string {
	var url string
	if lang == "ch" {
		url = urlPrefix_ch + suffix
	} else {
		url = urlPrefix_en + suffix
	}

	util.Info.Println("url", url)
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	page := string(body)
	titleReg := regexp.MustCompile(`<meta property="og:title" content="(.+?)" />`)
	textReg := regexp.MustCompile(`(?s)<div class="itemFullText">(.+?)</div>`)

	title := titleReg.FindStringSubmatch(page)
	text := textReg.FindStringSubmatch(page)

	d := ""
	if len(title) > 0 {
		d = strings.Replace(text[1], "<br />", "", -1)
		d = strings.Replace(d, "\t", "", -1)
		d = fmt.Sprintf(`<font size="4">%s</font><p><font size="4">%s</font>`, strings.Replace(title[1], " ", "", -1), d)
	}

	return d
}*/

type RthkNews struct {
	Title string
	Url   string
	Date  string
	Lang  string
}

// UpdateRthkRadio...
func UpdateRthkRadio(cat string) {
	util.Info.Println("Updating category", cat)
	weeks := "6"
	var url, category, category_ch string
	if cat == "parents_are_no_aliens" {
		url = "http://programme.rthk.hk/channel/radio/programme.php?name=radio1/parents_are_no_aliens&p=5286&m=archive&page=1&item=" + weeks
		category = "parents_are_no_aliens"
		category_ch = "我們不是怪獸"
	} else if cat == "musicinaction" {
		url = "http://programme.rthk.hk/channel/radio/programme.php?name=radio2/musicinaction&p=4548&m=archive&page=1&item=" + weeks
		category = "musicinaction"
		category_ch = "騷動音樂"
	} else if cat == "City_Snapshot" {
		url = "http://programme.rthk.hk/channel/radio/programme.php?name=radio1/City_Snapshot&p=3172&m=archive&page=1&item=30"
		category = "City_Snapshot"
		category_ch = "大城小事"
	} else if cat == "intlpopchart" {
		url = "http://programme.rthk.hk/channel/radio/programme.php?name=radio2/intlpopchart&p=415&m=archive&page=1&item=" + weeks
		category = "intlpopchart"
		category_ch = "環球榜"
	} else if cat == "People" {
		url = "http://programme.rthk.hk/channel/radio/programme.php?name=radio1/People&p=4819&m=archive&page=1&item=" + weeks
		category = "People"
		category_ch = "古今風雲人物"
		updateRTHKPodCast(category)
	} else if cat == "World_in_a_Nutshell" {
		url = "http://programme.rthk.hk/channel/radio/programme.php?name=radio1/World_in_a_Nutshell&p=4469&m=archive&page=1&item=" + weeks
		category = "World_in_a_Nutshell"
		category_ch = "十萬八千里"
	}
	radio := readAllRTHKRadioArchive(url)

	for _, r := range radio {
		title := r.title
		url := prefix + r.url
		date := r.date

		title, desc := readRTHKRadio(category, category_ch, url)
		saveRTHKRadio(title, desc, date, category)
	}
}

func saveRTHKRadio(title string, desc string, date string, category string) {
	db := util.GetDB()
	defer db.Close()

	var c int
	db.QueryRow("select count(1) from RTHK_Radio where type = ? and strftime('%Y-%m-%d', date) = ?", category, date).Scan(&c)

	if c > 0 {
		util.Info.Printf("Date: %s, Title: %s was saved, skip! \n", date, title)
		return
	}

	sql := "insert into RTHK_Radio (Title, Description, Date, Last_Update, Type) values (?, ?, ?, date('now', '+8 hour'), ?)"

	tx, err := db.Begin()
	if err != nil {
		util.Error.Fatal(err)
	}

	stmt, err := tx.Prepare(sql)
	if err != nil {
		util.Error.Fatal(err)
	}

	_, err = stmt.Exec(date+" "+title, desc, date, category)
	if err != nil {
		util.Error.Fatal(err)
	}
	stmt.Close()

	util.Info.Printf("Date: %s, Title: %s has been saved. \n", date, title)
	tx.Commit()
}

func readRTHKRadio(category string, text string, url string) (string, string) {
	util.Info.Println("Get", category, "of date", url)

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	page := string(body)

	r := regexp.MustCompile(`(?s)<meta name="title" content="rthk.hk 香港電台網站: ` + text + `:(.+?)" />`)
	r2 := regexp.MustCompile(`(?s)<meta name="description" content="(.+?)" />`)

	title := r.FindStringSubmatch(page)
	desc := r2.FindStringSubmatch(page)

	var a, b string
	if len(title) > 0 {
		a = title[1]
	} else {
		util.Info.Println("Fail to get title", title)
	}

	if len(desc) > 0 {
		b = desc[1]
	} else {
		util.Info.Println("Fail to get description", desc)
	}

	a = unescapeHTML(a)
	b = unescapeHTML(b)

	return a, b
}

func unescapeHTML(html string) string {
	html = strings.Replace(html, "&amp;", "&", -1)
	html = strings.Replace(html, "&quot;", "\"", -1)
	html = strings.Replace(html, "&lt;", "<", -1)
	html = strings.Replace(html, "&gt;", ">", -1)
	html = strings.Replace(html, "&nbsp;", " ", -1)
	return html
}

func readAllRTHKRadioArchive(url string) []RTHKRadio {
	var output []RTHKRadio
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	page := string(body)

	titleReg := regexp.MustCompile(`<div class="title"><a href="(.+m=episode)">((\d{4}-\d{2}-\d{2}).+?) </a></div>`) // get title
	result := titleReg.FindAllStringSubmatch(page, -1)
	for _, res := range result {
		output = append(output, RTHKRadio{title: res[2], url: res[1], date: res[3]})
	}
	return output
}

type RTHKRadio struct {
	title string
	desc  string
	url   string
	date  string
}

// RTHK PodCast
type Rss struct {
	Channel []Channel `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title    string `xml:"title"`
	Link     string `xml:"link"`
	Subtitle string `xml:"subtitle"`
	Summary  string `xml:"summary"`
	Guid     string `xml:"guid"`
	PubDate  string `xml:"pubDate"`
	Duration string `xml:"duration"`
}

func updateRTHKPodCast(cat string) {
	util.Info.Println("updateRTHKPodCast", cat)
	var url string
	if cat == "People" {
		url = "http://podcast.rthk.hk/podcast/people.xml"
		fetchData(cat, url)
	}
}

func fetchData(cat, url string) {
	rss := Rss{}

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	xmlContent, _ := ioutil.ReadAll(resp.Body)

	err := xml.Unmarshal(xmlContent, &rss)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%#v\n", rss)
	fmt.Println(len(rss.Channel[0].Items))

	db := util.GetDB()
	defer db.Close()

	selectSql := "select count(1) from RTHK_PODCAST where Title = ?"
	insertSql := "insert into RTHK_PODCAST (Title, Link, Subtitle, Summary, Guid, PubDate, Duration, Last_Update, Type) values (?, ?, ?, ?, ?, ?, ?, date('now', '+8 hour'), ?)"

	tx, err := db.Begin()
	if err != nil {
		util.Error.Println(err)
	}

	if err != nil {
		util.Error.Println(err)
	}

	stmt, err := tx.Prepare(insertSql)
	for _, item := range rss.Channel[0].Items {
		title := item.Title

		var count int
		db.QueryRow(selectSql, title).Scan(&count)
		if count == 0 {
			fmt.Println("Save title", title, item.Link)

			_, err = stmt.Exec(item.Title, item.Link, item.Subtitle, item.Summary, item.Guid, item.PubDate, item.Duration, cat)
			if err != nil {
				util.Error.Fatal(err)
			}

			util.Info.Println("Fetch Detail", item.Link)
			updateDesc(tx, item.Link)
		}
	}

	stmt.Close()
	tx.Commit()
}

func updateDesc(tx *sql.Tx, link string) {
	updateSql := "update RTHK_PODCAST set Description = ? where link = ?"

	stmt, err := tx.Prepare(updateSql)

	desc := readDetail(link)

	_, err = stmt.Exec(desc, link)
	if err != nil {
		util.Error.Println(err)
	}

	stmt.Close()
}

func readDetail(url string) string {
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	re := regexp.MustCompile(`(?s)<meta property="og:description" content="(.+?)" />`)
	desc := re.FindStringSubmatch(string(body))
	//fmt.Println(desc)

	return desc[1]
}
