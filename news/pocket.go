package news

import (
	"../util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

// Get idea from http://www.jamesfmackenzie.com/getting-started-with-the-pocket-developer-api/

func MyPocketItems(tag string) []byte {
	path := "https://getpocket.com/v3/get"
	urlData := url.Values{}
	urlData.Set("consumer_key", "54638-6333f424d3695345246bea9d")
	urlData.Set("access_token", "c68f2fe2-b655-f89f-cb4d-acca07")
	urlData.Set("detailType", "complete")
	//urlData.Set("sort", "newest")
	if tag != "" {
		t, _ := url.QueryUnescape(tag)
		util.Info.Println("show me the tag", t)
		urlData.Set("tag", tag)
	}

	client := http.Client{}
	resp, _ := client.PostForm(path, urlData)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	//ioutil.WriteFile("getpocket2.html", body, 0644)
	/*body, _ := ioutil.ReadFile("getpocket.html")*/

	var dat map[string]interface{}
	err := json.Unmarshal(body, &dat)
	if err != nil {
		util.Error.Fatal(err)
	}

	var items []Item
	list := dat["list"].(map[string]interface{})
	for _, item := range list {
		row := item.(map[string]interface{})
		timeAdded, _ := strconv.ParseInt(row["time_added"].(string), 10, 64)
		title := row["given_title"].(string)
		givenUrl := row["given_url"].(string)
		var tags []string
		if row["tags"] != nil {
			pocketTags := row["tags"].(map[string]interface{})

			for i, _ := range pocketTags {
				tags = append(tags, i)
			}
		}

		if len(title) == 0 {
			title = row["resolved_title"].(string)
			if len(title) == 0 {
				title = givenUrl
			}
		}

		i := Item{row["item_id"].(string), title, givenUrl, timeAdded, row["excerpt"].(string), tags}
		items = append(items, i)
	}

	sort.Sort(ByTimeUpdate(items)) // need to sort items as it is extracted from a map

	_items, _ := json.Marshal(items)

	return _items
}

// For reference
func ConvertTime(unixTimeStamp string) string {
	// unixTimeStamp := "1434508579"
	unixIntValue, err := strconv.ParseInt(unixTimeStamp, 10, 64)

	if err != nil {
		util.Error.Println(err)
	}

	t := time.Unix(unixIntValue, 0)
	time := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	//fmt.Println("time ", time)
	return time
}

type ByTimeUpdate []Item

func (b ByTimeUpdate) Len() int {
	return len(b)
}

func (b ByTimeUpdate) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByTimeUpdate) Less(i, j int) bool {
	return b[i].TimeAdded > b[j].TimeAdded
}

type Item struct {
	ItemId    string
	Title     string
	Url       string
	TimeAdded int64
	Excerpt   string
	Tags      []string
}
