package news

import (
	"../util"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var (
	db *sql.DB
)

func ReadNHK() {
	t := time.Now()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t = t.In(loc)
	time := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	util.Info.Println("Update Time in Asia/Shanghai zone: ", time)

	db = util.GetDB()
	defer db.Close()

	for i := 1; i <= 9; i++ {
		saveNews("index" + strconv.Itoa(i) + ".html")
	}

	saveMp3URL()
}

func saveNews(index string) {
	url := "http://k.nhk.jp/daily/" + index

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	page := string(body)

	m := regexp.MustCompile(`<hr>\d\. (.+)<br>`).FindStringSubmatch(page)
	title := m[1]

	sqlStmt := "select count(1) from NHK_WORLD_Daily_News where title = ? and date(Last_update) = date('now', '+8 hour')" // Asia/China timezone
	var count int

	db.QueryRow(sqlStmt, title).Scan(&count)

	if count > 0 {
		util.Info.Println("Article [" + title + "] was saved, skip !!")
	} else {
		tx, err := db.Begin()

		if err != nil {
			util.Error.Fatal(err)
		}

		sqlStmt = "insert into NHK_WORLD_Daily_News(Title, Content, Url, Last_Update) values(?, ?, ?, datetime('now', '+8 hour'))"
		stmt, _ := tx.Prepare(sqlStmt)

		_, err = stmt.Exec(title, page, index)
		if err != nil {
			util.Error.Printf("%q: %s\n", err, sqlStmt)
			return
		}

		err = tx.Commit()
		if err != nil {
			util.Error.Fatal(err)
		}

		stmt.Close()

		util.Info.Println("Article [" + title + "] has been saved")
	}
}

func getURL() string {
	url := "http://www.nhk.or.jp/rj/podcast/rss/english.xml"
	mp3url := ""

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	page := string(body)

	m := regexp.MustCompile(`enclosure url="(.+?)".+`).FindStringSubmatch(page)
	if len(m) > 0 {
		mp3url = m[1]
	}

	return mp3url
}

func saveMp3URL() {
	mp3url := getURL()
	sqlStmt := "select count(1) from NHK_WORLD_News_mp3 where url = ?"
	var count int

	db.QueryRow(sqlStmt, mp3url).Scan(&count)

	if count == 0 {
		tx, err := db.Begin()

		if err != nil {
			util.Error.Fatal(err)
		}

		sqlStmt = "insert into NHK_WORLD_News_mp3(Url, Last_Update) values(?, datetime('now', '+8 hour'))"
		stmt, _ := tx.Prepare(sqlStmt)

		_, err = stmt.Exec(mp3url)
		if err != nil {
			util.Error.Printf("%q: %s\n", err, sqlStmt)
			return
		}

		err = tx.Commit()
		if err != nil {
			util.Error.Fatal(err)
		}

		stmt.Close()

		util.Info.Println("mp3url [", mp3url, "] has been saved")
	} else {
		util.Info.Println("mp3url [", mp3url, "] was saved, skip !!")
	}
}
