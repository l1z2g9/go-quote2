package soundcloud

import (
	"github.com/l1z2g9/go-quote2/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	//"bytes"
)

// curl -vv -X POST --data 'username=&password=&client_id=8xSl5dwv0YD0kkIs2NUbaw5dmD2Tk9fK&redirect_uri=https://myapp-demo965238.rhcloud.com/scd/callback&response_type=token' https://soundcloud.com/connect/login
// Show playlist:   https://api.soundcloud.com/playlists/305138956?client_id=8xSl5dwv0YD0kkIs2NUbaw5dmD2Tk9fK
// Play my track:   https://api.soundcloud.com/me/playlists?oauth_token=1-274177-278918219-b916411cb1ac3
// Play track:      https://api.soundcloud.com/tracks/177671751/stream?client_id=8xSl5dwv0YD0kkIs2NUbaw5dmD2Tk9fK
//                  https://api.soundcloud.com/tracks/177671751/stream?oauth_token=1-274177-278918219-e242a1332bef8

const (
	hostPath = "https://soundcloud.com/"
	apiPath  = "https://api.soundcloud.com/"
	clientId = "8xSl5dwv0YD0kkIs2NUbaw5dmD2Tk9fK"
)

type Soundcloud struct {
	Oauth_token string
}

func (s Soundcloud) Login() string {
	urlData := url.Values{}

	info := util.Decrypt("C5UYDC77LSt2TgNNd05qupmteuxbMhNRL8XOzEIxjAE=")
	i := strings.Split(string(info), "X")

	urlData.Set("username", i[0])
	urlData.Set("password", i[1])
	urlData.Set("client_id", "8xSl5dwv0YD0kkIs2NUbaw5dmD2Tk9fK")
	urlData.Set("redirect_uri", "https://myapp-demo965238.rhcloud.com/scd/callback")
	urlData.Set("response_type", "token")

	request, err := http.NewRequest("POST", hostPath+"connect/login", strings.NewReader(urlData.Encode())) // or bytes.NewBufferString
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:31.0) Gecko/20100101 Firefox/31.0")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(urlData.Encode())))

	//client := http.Client{}
	resp, err := http.DefaultTransport.RoundTrip(request) // just get the result, do not redirect a new page
	if err != nil {
		util.Error.Println("Error RoundTrip", hostPath+"connect/login")
	}

	fragment := ""
	util.Info.Println(resp.StatusCode)
	if resp.StatusCode == 302 {
		loc, _ := resp.Location()
		util.Info.Println("loc", loc)
		fragment = loc.Fragment
	}

	defer resp.Body.Close()
	return fragment
}

func (s Soundcloud) GetPlaylists() []byte {
	var dat playlists
	byt := s.callApi("me/playlists")

	if err := json.Unmarshal(byt, &dat); err != nil {
		util.Error.Fatal("Error on GetPlaylists: ", err)
	}

	var myPlaylists []MyPlaylist
	for _, i := range dat {
		//myPlaylist := assemblePlaylist(i) will throw error 'cannot use i (type struct { playlist }) as type playlist in argument'
		myPlaylist := MyPlaylist{}
		myPlaylist.ID = i.ID
		myPlaylist.Title = i.Title

		var myTracks []MyTrack
		for _, j := range i.Tracks {
			myTrack := MyTrack{j.ID, j.Title, j.Description, j.StreamURL + "?client_id=" + clientId}
			myTracks = append(myTracks, myTrack)
		}

		myPlaylist.MyTracks = myTracks

		myPlaylists = append(myPlaylists, myPlaylist)
	}

	items_, _ := json.Marshal(myPlaylists)
	return items_

}

func (s Soundcloud) GetPlaylist(p string) []byte {
	var dat playlist
	byt := s.callApi(fmt.Sprintf("playlists/%s", p))

	if err := json.Unmarshal(byt, &dat); err != nil {
		util.Error.Fatal("Error on GetPlaylist: ", err)
	}

	myPlaylist := assemblePlaylist(dat)
	items_, _ := json.Marshal(myPlaylist.MyTracks)
	return items_
}

func assemblePlaylist(dat playlist) MyPlaylist {
	myPlaylist := MyPlaylist{}
	myPlaylist.ID = dat.ID
	myPlaylist.Title = dat.Title

	var myTracks []MyTrack
	for _, j := range dat.Tracks {
		myTrack := MyTrack{j.ID, j.Title, j.Description, j.StreamURL + "?client_id=" + clientId}
		myTracks = append(myTracks, myTrack)
	}

	myPlaylist.MyTracks = myTracks
	return myPlaylist
}

func (s Soundcloud) callApi(action string) []byte {
	client := http.Client{}
	path := apiPath + action

	if s.Oauth_token != "" {
		path = path + "?oauth_token=" + s.Oauth_token
	} else {
		path = path + "?client_id=" + clientId
	}

	req, err := http.NewRequest("GET", path, nil)

	util.Info.Println("call url ", path)

	resp, err := client.Do(req)
	if err != nil {
		util.Error.Println("Error on callApi: ", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body
}

type MyPlaylist struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	MyTracks []MyTrack
}

type MyTrack struct {
	ID          int `json:"Seq"`
	Title       string
	Description string
	StreamURL   string `json:"URL"`
}

type playlists []struct {
	playlist
}

type playlist struct {
	Duration     int         `json:"duration"`
	ReleaseDay   interface{} `json:"release_day"`
	PermalinkURL string      `json:"permalink_url"`
	Genre        string      `json:"genre"`
	Permalink    string      `json:"permalink"`
	PurchaseURL  interface{} `json:"purchase_url"`
	ReleaseMonth interface{} `json:"release_month"`
	Description  interface{} `json:"description"`
	URI          string      `json:"uri"`
	LabelName    interface{} `json:"label_name"`
	TagList      string      `json:"tag_list"`
	ReleaseYear  interface{} `json:"release_year"`
	TrackCount   int         `json:"track_count"`
	UserID       int         `json:"user_id"`
	LastModified string      `json:"last_modified"`
	License      string      `json:"license"`
	Tracks       []struct {
		Kind                string      `json:"kind"`
		ID                  int         `json:"id"`
		CreatedAt           string      `json:"created_at"`
		UserID              int         `json:"user_id"`
		Duration            int         `json:"duration"`
		Commentable         bool        `json:"commentable"`
		State               string      `json:"state"`
		OriginalContentSize int         `json:"original_content_size"`
		LastModified        string      `json:"last_modified"`
		Sharing             string      `json:"sharing"`
		TagList             string      `json:"tag_list"`
		Permalink           string      `json:"permalink"`
		Streamable          bool        `json:"streamable"`
		EmbeddableBy        string      `json:"embeddable_by"`
		Downloadable        bool        `json:"downloadable"`
		PurchaseURL         string      `json:"purchase_url"`
		LabelID             interface{} `json:"label_id"`
		PurchaseTitle       string      `json:"purchase_title"`
		Genre               string      `json:"genre"`
		Title               string      `json:"title"`
		Description         string      `json:"description"`
		LabelName           string      `json:"label_name"`
		Release             interface{} `json:"release"`
		TrackType           interface{} `json:"track_type"`
		KeySignature        interface{} `json:"key_signature"`
		Isrc                interface{} `json:"isrc"`
		VideoURL            interface{} `json:"video_url"`
		Bpm                 interface{} `json:"bpm"`
		ReleaseYear         interface{} `json:"release_year"`
		ReleaseMonth        interface{} `json:"release_month"`
		ReleaseDay          interface{} `json:"release_day"`
		OriginalFormat      string      `json:"original_format"`
		License             string      `json:"license"`
		URI                 string      `json:"uri"`
		User                struct {
			ID           int    `json:"id"`
			Kind         string `json:"kind"`
			Permalink    string `json:"permalink"`
			Username     string `json:"username"`
			LastModified string `json:"last_modified"`
			URI          string `json:"uri"`
			PermalinkURL string `json:"permalink_url"`
			AvatarURL    string `json:"avatar_url"`
		} `json:"user"`
		AttachmentsURI   string      `json:"attachments_uri"`
		PermalinkURL     string      `json:"permalink_url"`
		ArtworkURL       interface{} `json:"artwork_url"`
		WaveformURL      string      `json:"waveform_url"`
		StreamURL        string      `json:"stream_url"`
		PlaybackCount    int         `json:"playback_count"`
		DownloadCount    int         `json:"download_count"`
		FavoritingsCount int         `json:"favoritings_count"`
		CommentCount     int         `json:"comment_count"`
	} `json:"tracks"`
	PlaylistType  interface{} `json:"playlist_type"`
	ID            int         `json:"id"`
	Downloadable  interface{} `json:"downloadable"`
	Sharing       string      `json:"sharing"`
	CreatedAt     string      `json:"created_at"`
	Release       interface{} `json:"release"`
	Kind          string      `json:"kind"`
	Title         string      `json:"title"`
	Type          interface{} `json:"type"`
	PurchaseTitle interface{} `json:"purchase_title"`
	ArtworkURL    interface{} `json:"artwork_url"`
	Ean           interface{} `json:"ean"`
	Streamable    bool        `json:"streamable"`
	User          struct {
		PermalinkURL string `json:"permalink_url"`
		Permalink    string `json:"permalink"`
		Username     string `json:"username"`
		URI          string `json:"uri"`
		LastModified string `json:"last_modified"`
		ID           int    `json:"id"`
		Kind         string `json:"kind"`
		AvatarURL    string `json:"avatar_url"`
	} `json:"user"`
	EmbeddableBy string      `json:"embeddable_by"`
	LabelID      interface{} `json:"label_id"`
}
