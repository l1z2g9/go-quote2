package util

import (
	_ "errors"
	"html"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func GetLink(urlstring string) io.ReadCloser {
	urlstring = html.EscapeString(urlstring)
	Info.Println("Get link", urlstring)
	req, err := http.NewRequest("GET", urlstring, nil)

	if err != nil {
		Error.Println("Error access", urlstring)
	}

	var body io.ReadCloser
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		Error.Println("Error RoundTrip", urlstring)
	}

	defer resp.Body.Close()
	if resp.StatusCode == 302 {
		downloadUrl := resp.Header["Location"][0]
		Info.Println("Redirect to link", downloadUrl)

		req, _ = http.NewRequest("GET", downloadUrl, nil)

		req.URL = &url.URL{
			Scheme: "https",
			Host:   req.Host,
			// fix the problem of asterisk character in url while getting google drive url properly
			Opaque: strings.Replace(downloadUrl, "https:", "", -1),
		}

		cookieJar, _ := cookiejar.New(nil)
		client := http.Client{Jar: cookieJar, CheckRedirect: func(req *http.Request, via []*http.Request) error {
			//fmt.Println("req", req, "via", via)
			//return errors.New("Don't redirect!")
			return nil
		}}

		resp, err = client.Do(req)

		if err != nil {
			Error.Println("Error access link", downloadUrl)
		}

		//body, _ = ioutil.ReadAll(resp.Body)
		//ioutil.WriteFile("abc", body, 0644)

		body = resp.Body
	} else {
		client := http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			Error.Println("Error handle url", urlstring)
		}

		body = resp.Body
	}

	return body
}
