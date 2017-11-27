package news

import (
	"github.com/l1z2g9/go-quote2/util"
	"bytes"
	_ "crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
)

const (
	host = "https://www.oschina.net/action"
)

func OschinaApi(action string, token string, id string) string {
	var result string
	var urlstring string

	if action == "oauth" {
		return getApiAuthorization()
	}

	path := host + "/openapi/" + action
	client := http.Client{}

	if action == "favorite_list" || action == "news_list" {
		urlstring = fmt.Sprintf(path+"?access_token=%s&pageSize=500&page=1", token)
	}

	if action == "news_detail" {
		urlstring = fmt.Sprintf(path+"?access_token=%s&id=%s", token, id)
	}

	util.Info.Println("action", action, "urlstring", urlstring)

	resp, _ := client.Get(urlstring)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	result = string(body)

	return result
}

func getApiAuthorization() string {
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: cookieJar}

	login(client)
	touchAuthUrl(client)
	token := passInfoByForm(client)
	return token
}

func GetAccessToken(code string) string {
	path := host + "/openapi/token"

	client := http.Client{}
	resp, _ := client.Get(fmt.Sprintf(path+"?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s&grant_type=%s", "pHAnlp34IsGn9yZMYCZR", "http://quote-querychinesesto.rhcloud.com/oschina", "X7h1K4C20OYUt5ncdwpsdDo1IO99Mpys", code, "authorization_code"))

	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)

	defer resp.Body.Close()

	var dat map[string]interface{}
	err := json.Unmarshal([]byte(content), &dat)
	if err != nil {
		util.Error.Println(err)
	}

	token := dat["access_token"].(string)
	expire := strconv.FormatFloat(dat["expires_in"].(float64), 'f', -1, 64)

	util.Info.Println("token | expire : ", token+"|"+expire)

	return token + "|" + expire
}

func passInfoByCustom(client *http.Client) {
	urlstring := host + "/oauth2/authorize"
	urlData := url.Values{}
	urlData.Set("response_type", "code")
	urlData.Set("client_id", "pHAnlp34IsGn9yZMYCZR")
	urlData.Set("redirect_uri", "http://quote-querychinesesto.rhcloud.com/oschina")
	urlData.Set("scope", "user_api")
	urlData.Set("user_oauth_approval", "true")
	urlData.Set("state", "")
	util.Info.Println("urlData.Encode()", urlData.Encode())
	r, _ := http.NewRequest("POST", urlstring, bytes.NewBufferString(urlData.Encode()))
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:31.0) Gecko/20100101 Firefox/31.0")
	r.Header.Add("Referer", `https://www.oschina.net/action/oauth2/authorize? \
                    response_type=code&client_id=pHAnlp34IsGn9yZMYCZR& \
                    redirect_uri=http://quote-querychinesesto.rhcloud.com/oschina`)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(urlData.Encode())))

	//resp, _ := client.Do(r)
	resp, err := http.DefaultTransport.RoundTrip(r) // just get the result, do not redirect a new page
	if err != nil {
		util.Error.Println("Error RoundTrip", urlstring)
	}

	util.Info.Println(resp.Status)
	util.Info.Println(resp.StatusCode)
}

func touchAuthUrl(client *http.Client) {
	path := host + "/oauth2/authorize?response_type=code&client_id=pHAnlp34IsGn9yZMYCZR&&redirect_uri=http://quote-querychinesesto.rhcloud.com/oschina"

	resp, _ := client.Get(path)
	util.Info.Println("touchAuthUrl http status", resp.Status)
	defer resp.Body.Close()
}

func passInfoByForm(client *http.Client) string {
	path := host + "/oauth2/authorize"
	urlData := url.Values{}
	urlData.Set("response_type", "code")
	urlData.Set("client_id", "pHAnlp34IsGn9yZMYCZR")
	urlData.Set("redirect_uri", "http://quote-querychinesesto.rhcloud.com/oschina")
	urlData.Set("scope", "news_api,project_api,user_api,favorite_api")
	urlData.Set("user_oauth_approval", "true")
	urlData.Set("state", "")

	resp, err := client.PostForm(path, urlData)
	if err != nil {
		util.Error.Fatal(err)
	}

	util.Info.Println("passInfoByForm http status", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)

	defer resp.Body.Close()
	return content
}

func login(client *http.Client) {
	//b := sha1.Sum([]byte(""))
	//pwd := fmt.Sprintf("%x", b)
	pwd := "48d19eeaf7fc00265dc884f2863624e0cdb3ee8a"

	path := host + "/user/hash_login"
	urlData := url.Values{}
	urlData.Set("email", "querychinesesto@gmail.com")
	urlData.Set("pwd", pwd)
	urlData.Set("verifyCode", "")

	resp, _ := client.PostForm(path, urlData)

	defer resp.Body.Close()
}
