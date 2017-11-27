package news

import (
	"github.com/l1z2g9/go-quote2/util"
	"encoding/json"
	_ "fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

/*func initDB() *sql.DB {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	db_path := dir + "/../stock.db"

	db, _ := sql.Open("sqlite3", db_path)
	return db
}*/

const translate_url = "http://translate.googleapis.com/translate_a/single?client=gtx&sl=en&tl=zh-CN&dt=t&q="

func GetNewsList(offset string, searchText string) string {
	db := util.GetDB()
	defer db.Close()

	var titleList []string

	_offset := ""

	criteria := ""
	if searchText != "" {
		criteria = "where title like '%" + searchText + "%'"
	} else if offset != "" {
		n, _ := strconv.Atoi(offset)
		_offset = "limit 20 offset " + strconv.Itoa((n-1)*20)
	}

	stmt := "select Id, Title, datetime(Last_Update) from NHK_WORLD_Daily_News " + criteria + " Order by Last_Update desc " + _offset
	rows, _ := db.Query(stmt)
	util.Info.Println(stmt)

	defer rows.Close()

	for rows.Next() {
		var id string
		var title string
		var date string
		rows.Scan(&id, &title, &date)
		title = util.Trim(title)

		news := &News{id, title, "", "", date, ""}
		item, _ := json.Marshal(news)
		titleList = append(titleList, string(item))
	}

	return "[" + strings.Join(titleList, ",") + "]"
}

func GetNewsContent(id string) string {
	db := util.GetDB()
	defer db.Close()

	var content string
	var like string
	db.QueryRow("select Content, Like from NHK_WORLD_Daily_News where id = ?", id).Scan(&content, &like)
	content = util.Trim(content)
	news := &News{Content: content, Like: like}
	item, _ := json.Marshal(news)

	util.Info.Println("newsId ", id, " is read")
	return string(item)
}

func GetNewsAudioList() string {
	db := util.GetDB()
	defer db.Close()

	var titleList []string
	stmt := "select Url, Listend from NHK_WORLD_News_mp3 Order by Last_Update desc limit 20"
	rows, _ := db.Query(stmt)
	util.Info.Println(stmt)

	defer rows.Close()

	for rows.Next() {
		var url string
		var listened string

		rows.Scan(&url, &listened)

		news := &News{Id: url, Listened: listened}
		item, _ := json.Marshal(news)
		titleList = append(titleList, string(item))
	}

	return "[" + strings.Join(titleList, ",") + "]"
}

func NewsFavor(id string, favor string) {
	db := util.GetDB()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		util.Error.Fatal(err)
	}

	sqlStmt := "update NHK_WORLD_Daily_News set like = ? where id = ?"
	stmt, _ := tx.Prepare(sqlStmt)

	util.Info.Println("id", id, "like", favor)

	_, err = stmt.Exec(favor, id)
	if err != nil {
		util.Error.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	err = tx.Commit()
	if err != nil {
		util.Error.Fatal(err)
	}

	stmt.Close()

	util.Info.Println("Article of Id ["+id+"] has been liked: ", favor)
}

func TranslateWord(id string, word string) string {
	db := util.GetDB()
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		util.Error.Fatal(err)
	}

	var words string
	db.QueryRow("select translated_words from NHK_WORLD_Daily_News where id = ?", id).Scan(&words)

	if words != "" {
		if !strings.Contains(words, word) {
			words = words + ", " + word
		}
	} else {
		words = word
	}

	sqlStmt := "update NHK_WORLD_Daily_News set translated_words = ? where id = ?"
	stmt, _ := tx.Prepare(sqlStmt)

	util.Info.Println("id", id, "translated_words ", words)

	_, err = stmt.Exec(words, id)
	if err != nil {
		util.Error.Printf("%q: %s\n", err, sqlStmt)
	}

	err = tx.Commit()
	if err != nil {
		util.Error.Fatal(err)
	}

	stmt.Close()

	util.Info.Println("Translated words in article of Id ["+id+"] : ", words)

	resp, _ := http.Get(translate_url + word)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func GetBookmarkNews() string {
	db := util.GetDB()
	defer db.Close()

	var titleList []string
	stmt := "select Id, Title, datetime(Last_Update) from NHK_WORLD_Daily_News where like = 'T' Order by Last_Update desc"
	rows, _ := db.Query(stmt)
	util.Info.Println(stmt)

	defer rows.Close()

	for rows.Next() {
		var id string
		var title string
		var date string
		rows.Scan(&id, &title, &date)
		title = util.Trim(title)

		news := &News{id, title, "", "", date, ""}
		item, _ := json.Marshal(news)
		titleList = append(titleList, string(item))
	}

	return "[" + strings.Join(titleList, ",") + "]"
}

type News struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Like        string `json:"like"`
	Last_Update string `json:"lastUpdate"`
	Listened    string `json:"listened"`
}
