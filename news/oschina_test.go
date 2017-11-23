package news

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestOschinaApi(t *testing.T) {
	OschinaApi("", "", "")
}

func TestUrlValues(t *testing.T) {
	v := url.Values{}
	v.Set("name", "Ava")
	v.Add("friend", "Jess")
	v.Add("friend", "Sarah")
	v.Add("friend", "Zoe")
	// v.Encode() == "name=Ava&friend=Jess&friend=Sarah&friend=Zoe"
	fmt.Println(v.Get("name"))
	fmt.Println(v.Get("friend"))
	fmt.Println(v["friend"])
	fmt.Println(v.Get("aa"))
	fmt.Println(v["aa"])
}

func TestParseToken(t *testing.T) {
	token := `{"access_token":"e237184f-cf74-4849-9a49-84b0de4bbc2e","refresh_token":"5b43af96-b088-4c72-957d-ed52d3bd82df","uid":52891,
"token_type":"bearer","expires_in":604799}`
	var dat map[string]interface{}
	err := json.Unmarshal([]byte(token), &dat)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(dat["access_token"])
}

func TestGetUrl(t *testing.T) {
	client := http.Client{}
	resp, _ := client.Get("http://quote-querychinesesto.rhcloud.com/oschina/oauth")
	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)
	fmt.Println(content)

	resp.Body.Close()
}
