package news

import (
	"../util"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// example
// curl -X POST -H "content-type:application/json" https://cloud.feedly.com/v3/feeds/.mget --data '["feed/http://dave.cheney.net/feed","feed/http://www.yatzer.com/feed/index.php"]' -v
// curl -H "Authorization: AzWD1rUVVW2jr3vjo8qHhz4nSaaG1HIStUWTxRqptzPtdi4_WBGK0AaU3qIy-uDeY0gXTrQZzLxdJJzyAZLMfuSoWRtBjAD1lk_lM9tTA-FbPvQDQQMkdy6aUmMiHQSw31K3nyaeiIenNY_aAGcpOpYxLOGiKRLpb8iZ6oYqxBYhgllC8Bn_a6pqwwWg3Gd5OmLaWCQ0LfrDFp9iCq7hWc57O9_8iA:feedlydev" https://cloud.feedly.com/v3/subscriptions

const (
	urlPath = "https://cloud.feedly.com/v3/"
)

type Feedly struct {
}

func (f Feedly) GetProfile() []byte {
	return callApi("profile", true)
}

func (f Feedly) GetCategories() []byte {
	return callApi("categories", true)
}

func (f Feedly) GetSubscriptions() []byte {
	var dat Subscriptions
	byt := callApi("subscriptions", true)
	if err := json.Unmarshal(byt, &dat); err != nil {
		util.Error.Fatal("Error on Subscriptions: ", err)
	}

	var minSubscriptions []MinSubscription
	for _, s := range dat {
		cat := ""
		if len(s.Categories) > 0 {
			cat = s.Categories[0].Label
		}
		minSubscription := MinSubscription{s.ID, s.Title, s.Website, cat}
		minSubscriptions = append(minSubscriptions, minSubscription)
	}

	minSubscriptions_, _ := json.Marshal(minSubscriptions)
	return minSubscriptions_
}

func (f Feedly) GetListByFeed(streamId string) []byte {
	var dat Stream
	byt := callApi("streams/contents?count=50&streamId="+streamId, false)
	if err := json.Unmarshal(byt, &dat); err != nil {
		util.Error.Fatal("Error on GetListByFeed: ", err)
	}

	var items []FeedItem
	for _, i := range dat.Items {
		item := FeedItem{i.ID, i.Title, i.Summary.Content, i.Published, i.Content.Content, i.OriginID}
		items = append(items, item)
	}

	items_, _ := json.Marshal(items)
	return items_
}

func (f Feedly) GetEntryContent(entryId, alternateHref string) []byte {
	var dat Entries
	byt := callApi("entries/"+entryId, false)
	if err := json.Unmarshal(byt, &dat); err != nil {
		util.Error.Fatal("Error on GetEntryContent: ", err)
	}

	for _, s := range dat {
		c := s.Content.Content
		if len(c) > 0 {
			return []byte(c)
		}
	}

	body, _ := ioutil.ReadAll(util.GetLink(alternateHref))

	return body
}

func callApi(action string, authRequired bool) []byte {
	client := http.Client{}
	path := urlPath + action
	util.Info.Println("call url ", path)

	dir := util.GetExecutePath()
	secret_code, _ := ioutil.ReadFile(dir + "/../news/feedly_code")

	req, err := http.NewRequest("GET", path, nil)
	if authRequired {
		req.Header.Add("Authorization", "OAuth "+string(secret_code))
	}

	resp, err := client.Do(req)
	if err != nil {
		util.Error.Println("Error on profile: ", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body
}

type Subscriptions []struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Website    string `json:"website"`
	Categories []struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	} `json:"categories"`
	Updated     int64   `json:"updated"`
	Subscribers int     `json:"subscribers"`
	Velocity    float64 `json:"velocity"`
	ContentType string  `json:"contentType,omitempty"`
	CoverURL    string  `json:"coverUrl,omitempty"`
	IconURL     string  `json:"iconUrl,omitempty"`
	Partial     bool    `json:"partial,omitempty"`
	VisualURL   string  `json:"visualUrl,omitempty"`
	CoverColor  string  `json:"coverColor,omitempty"`
	State       string  `json:"state,omitempty"`
}

type MinSubscription struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Website  string `json:"website"`
	Category string `json:"category"`
}

type Stream struct {
	ID        string `json:"id"`
	Direction string `json:"direction"`
	Updated   int64  `json:"updated"`
	Title     string `json:"title"`
	Alternate []struct {
		Href string `json:"href"`
		Type string `json:"type"`
	} `json:"alternate"`
	Continuation string `json:"continuation"`
	Items        []struct {
		ID          string   `json:"id"`
		Keywords    []string `json:"keywords"`
		OriginID    string   `json:"originId"`
		Fingerprint string   `json:"fingerprint"`
		Recrawled   int64    `json:"recrawled"`
		Origin      struct {
			StreamID string `json:"streamId"`
			Title    string `json:"title"`
			HTMLURL  string `json:"htmlUrl"`
		} `json:"origin"`
		Content struct {
			Content   string `json:"content"`
			Direction string `json:"direction"`
		} `json:"content"`
		Title     string `json:"title"`
		Published int64  `json:"published"`
		Crawled   int64  `json:"crawled"`
		Alternate []struct {
			Type string `json:"type"`
			Href string `json:"href"`
		} `json:"alternate"`
		Author  string `json:"author"`
		Summary struct {
			Content   string `json:"content"`
			Direction string `json:"direction"`
		} `json:"summary"`
		Unread     bool `json:"unread"`
		Categories []struct {
			ID    string `json:"id"`
			Label string `json:"label"`
		} `json:"categories"`
		Engagement int `json:"engagement"`
	} `json:"items"`
}

type FeedItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Published int64  `json:"published"`
	Content   string `json:"content"`
	OriginID  string `json:"originId"`
}

type Entries []struct {
	ID          string   `json:"id"`
	Keywords    []string `json:"keywords"`
	OriginID    string   `json:"originId"`
	Fingerprint string   `json:"fingerprint"`
	Recrawled   int64    `json:"recrawled"`
	Content     struct {
		Content   string `json:"content"`
		Direction string `json:"direction"`
	} `json:"content"`
	Title     string `json:"title"`
	Published int64  `json:"published"`
	Crawled   int64  `json:"crawled"`
	Summary   struct {
		Content   string `json:"content"`
		Direction string `json:"direction"`
	} `json:"summary"`
	Alternate []struct {
		Href string `json:"href"`
		Type string `json:"type"`
	} `json:"alternate"`
	Author string `json:"author"`
	Origin struct {
		StreamID string `json:"streamId"`
		Title    string `json:"title"`
		HTMLURL  string `json:"htmlUrl"`
	} `json:"origin"`
	Unread     bool `json:"unread"`
	Categories []struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	} `json:"categories"`
	Engagement int `json:"engagement"`
}
