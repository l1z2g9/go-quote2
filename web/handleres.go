package web

import (
	"../news"
	"../rthk"
	"../securities"
	"../soundcloud"
	"../util"
	_ "archive/zip"
	_ "bytes"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	_ "net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

func wrapHTML(content string) string {
	content = strings.Replace(content, "\n", "<br>", -1)
	return "<body bgcolor='#c0cdc0'>" + content + "<body>"
}

func getQuoteList(req *http.Request, isEncrypt string) string {
	ret := ""
	category := ""
	if cate, ok := req.URL.Query()["c"]; ok {
		category = cate[0]
	}

	if sf, ok := req.URL.Query()["sf"]; ok {
		if len(sf[0]) > 0 {
			util.Info.Println("sortField:", sf[0])
			securities.SetSortField(sf[0])
		} else {
			securities.SetSortField("PercentChange")
		}
	}

	if isEncrypt == "T" {
		ret = util.Encrypt(securities.GetQuoteWithFormat(category))
	} else {
		ret = securities.GetQuoteWithFormat(category)
	}

	util.Info.Println("category:", category, "isEncrypt:", isEncrypt)
	return ret
}

// Home is the home page
func Home(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "Home")
	fmt.Fprint(res, "Hello My Go Application")
}

// Index is the home page
func Index(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "Index")

	//data := getQuoteList(req, "N")
	//util.CompressData(res, []byte(data))
	fmt.Fprint(res, getQuoteList(req, "T"))
}

// Mobile is the mobile page
func Mobile(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "Mobile")

	data := wrapHTML(getQuoteList(req, "N"))
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	util.CompressData(res, []byte(data))
	//fmt.Fprint(res, wrapHTML(getQuoteList(req, "N")))
}

// ShowBookInfo ...
func ShowBookInfo(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "ShowBookInfo")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	body, _ := ioutil.ReadAll(req.Body)
	util.CompressData(res, util.ShowBookInfo(body))
	//fmt.Fprint(res, util.ShowBookInfo(body))
}

// LoginAndShowBookInfo ...
func LoginAndShowBookInfo(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "LoginAndShowBookInfo")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	data := util.LoginAndShowBookInfo()
	util.CompressData(res, data)
	//fmt.Fprint(res, util.LoginAndShowBookInfo())
}

// RenewBooks ...
func RenewBooks(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "Renew Library Books")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Fprint(res, util.RenewBooks(body))
}

func BookCheckoutHist(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "BookCheckoutHist")
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	util.CompressData(res, []byte(util.BookCheckoutHist()))
	//fmt.Fprint(res, util.BookCheckoutHist())
}

func BookSearchHist(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "BookSearchHist")
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	util.CompressData(res, []byte(util.BookSearchHist()))
	//fmt.Fprint(res, util.BookSearchHist())
}

// SearchBooks ...
func SearchBooks(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "SearchBooks")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := getVars(req)
	fmt.Fprint(res, util.SearchBooks(vars["queryText"], vars["location"], vars["pageNum"]))
}

// GetBookDetail ...
func GetBookDetail(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "GetBookDetail")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := getVars(req)

	util.CompressData(res, util.GetBookDetail(vars["chamoId"]))
	//fmt.Fprint(res, util.GetBookDetail(vars["chamoId"]))
}


// NewsList ...
func NewsList(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "NewsList")
	res.Header().Set("Content-Type", "application/json")
	vars := getVars(req)

	data := news.GetNewsList(vars["offset"], "")
	util.CompressData(res, []byte(data))
	//fmt.Fprint(res, news.GetNewsList(vars["offset"], ""))
}

// BookmarkNews ...
func BookmarkNews(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "BookmarkNews")
	res.Header().Set("Content-Type", "application/json")

	data := news.GetBookmarkNews()
	util.CompressData(res, []byte(data))
	//fmt.Fprint(res, news.GetBookmarkNews())
}

// NewsAudioList ...
func NewsAudioList(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "NewsAudioList")
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprint(res, news.GetNewsAudioList())
}

// NewsQuery ...
func NewsQuery(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "NewsQuery")
	res.Header().Set("Content-Type", "application/json")
	vars := getVars(req)

	data := news.GetNewsList("", vars["query"])
	util.CompressData(res, []byte(data))
	//fmt.Fprint(res, news.GetNewsList("", vars["query"]))
}

// NewsContent ...
func NewsContent(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "NewsContent")
	vars := getVars(req)
	fmt.Fprint(res, news.GetNewsContent(vars["id"]))
}

// NewsFavor ...
func NewsFavor(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "NewsFavor")
	vars := getVars(req)

	favor := vars["favor"]
	if "like" == favor {
		news.NewsFavor(vars["id"], "T")
	} else {
		news.NewsFavor(vars["id"], "F")
	}
	res.WriteHeader(http.StatusOK)
}

// TranslateWord ...
func TranslateWord(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "TranslateWord")
	vars := getVars(req)
	id := vars["id"]
	word := vars["word"]

	fmt.Fprint(res, news.TranslateWord(id, word))
}

// CitySnapShotList ...
func CitySnapShotList(res http.ResponseWriter, req *http.Request) {
	vars := getVars(req)
	saveAccessInfo(res, req, "CitySnapShotList "+vars["cate"])
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	data := rthk.GetCitySnapshots(vars["cate"])
	util.CompressData(res, data)

	//fmt.Fprint(res, rthk.GetCitySnapshots(vars["cate"]))
}

// YCantoneseList ...
func YCantoneseList(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "YCantoneseList")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	data := rthk.GetYCantonese()
	util.CompressData(res, data)
	//fmt.Fprint(res, rthk.GetYCantonese())
}

func RTHKRadioList(res http.ResponseWriter, req *http.Request) {
	vars := getVars(req)
	saveAccessInfo(res, req, "RTHKRadioList "+vars["cate"])
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	data := rthk.GetRthkRadio(vars["cate"])
	util.CompressData(res, data)

	//fmt.Fprint(res, rthk.GetRthkRadio(vars["cate"]))
}

func RTHKPodcastList(res http.ResponseWriter, req *http.Request) {
	vars := getVars(req)
	saveAccessInfo(res, req, "RTHKPodcastList "+vars["cate"])
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	data := rthk.GetRthkPodcast(vars["cate"])
	util.CompressData(res, data)

	//fmt.Fprint(res, rthk.GetRthkPodcast(vars["cate"]))
}

func RTHKPodcastDesc(res http.ResponseWriter, req *http.Request) {
	vars := getVars(req)
	saveAccessInfo(res, req, "RTHKPodcastDesc "+vars["seq"])
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(res, rthk.GetRTHKPodcasDesc(vars["seq"]))
}

func UpdateRthkRadio(res http.ResponseWriter, req *http.Request) {
	vars := getVars(req)
	saveAccessInfo(res, req, "UpdateRthkRadio "+vars["cate"])
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	rthk.UpdateRthkRadio(vars["cate"])
	fmt.Fprint(res, "Update finished for "+vars["cate"])
}

func RTHKRadioProxy(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "RTHKRadioProxy")
	vars := getVars(req)

	audio_path := "http://stmw.rthk.hk/aod/_definst_/radio/archive/%s/%s/mp3/mp3:%s.mp3/%s"
	urlPath := fmt.Sprintf(audio_path, vars["channel"], vars["cate"], vars["date"], vars["suffix"])
	//Info.Println("m3u urlPath", urlPath)

	body := util.GetLink(string(urlPath))
	defer body.Close()

	var buffer [8192]byte

	complete := false
	for complete == false {
		n, err := body.Read(buffer[:])

		res.Write(buffer[:n])

		if err == io.EOF {
			body.Close()
			complete = true
			//Info.Println("RTHKRadioProxy Fetch completed")
		}
	}
}

// MyPocket ...
func MyPocket(res http.ResponseWriter, req *http.Request) {
	vars := getVars(req)
	saveAccessInfo(res, req, "MyPocket "+vars["tag"])
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	data := news.MyPocketItems(vars["tag"])
	util.CompressData(res, data)
	//fmt.Fprint(res, news.MyPocketItems(vars["tag"]))
}

// GetQuoteWithoutFormat ...
func GetQuoteWithoutFormat(res http.ResponseWriter, req *http.Request) {
	vars := getVars(req)
	saveAccessInfo(res, req, "GetQuoteWithoutFormat "+vars["cate"])
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	data := securities.GetQuoteWithoutFormat(vars["cate"])
	util.CompressData(res, []byte(data))
	//fmt.Fprint(res, securities.GetQuoteWithoutFormat(vars["cate"]))
}

func GetRthkNews(res http.ResponseWriter, req *http.Request) {
	vars := getVars(req)
	saveAccessInfo(res, req, "GetRthkNews "+vars["cate"])
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	data := rthk.GetRthkNews(vars["cate"])
	util.CompressData(res, data)
	//fmt.Fprint(res, rthk.GetRthkNews(vars["cate"]))
}

/*func GetRthkNewsDetail(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "GetRthkNewsDetail")
	vars := getVars(req)
	fmt.Fprint(res, rthk.GetRthkNewsDetail(vars["suffix"], vars["lang"]))
}*/

func Oschina(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "Oschina")
	result := ""
	code := req.URL.Query().Get("code")
	util.Info.Println("code = ", code)

	if code != "" {
		result = news.GetAccessToken(code)
	}

	fmt.Fprint(res, result)
}

func GetFile(res http.ResponseWriter, req *http.Request) {
	vars := getVars(req)
	file := util.GetExecutePath() + "/" + vars["file"]
	saveAccessInfo(res, req, "GetFile "+file)

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		util.Error.Println("fail to open file", file, err)
	}
	res.Write(bytes)
}

func OschinaAction(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "OschinaAction")
	vars := getVars(req)
	action := vars["action"]
	token := req.URL.Query().Get("token")
	id := req.URL.Query().Get("id")
	fmt.Fprint(res, news.OschinaApi(action, token, id))
}

// ProxyAccess ...
func ProxyAccess(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "ProxyAccess")

	vars := getVars(req)
	urlstr := vars["url"]
	urlstr, _ = url.QueryUnescape(urlstr)

	urlstr2, err := base64.URLEncoding.DecodeString(urlstr)
	if err != nil {
		util.Info.Printf("Fail to decode url string %s", urlstr)
		fmt.Fprint(res, "Fail to decode url")
	}

	body := util.GetLink(string(urlstr2))

	var buffer [8192]byte

	complete := false
	for complete == false {
		n, err := body.Read(buffer[:])

		//fmt.Println("buffer ", string(buffer[:n]))

		if err != nil && err != io.EOF {
			util.Info.Println("Error Access", err)
			body.Close()
			complete = true
		}

		if _, werr := res.Write(buffer[:n]); werr != nil {
			util.Info.Println("Error Write Data", werr)
			body.Close()
			complete = true
		}

		if err == io.EOF {
			util.Info.Println("Fetch completed")
			body.Close()
			complete = true
		}
	}
}

func SqlConsole(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "SqlConsole method - "+req.Method)

	if req.Method == "GET" {
		//dat, _ := ioutil.ReadFile("sqlconsole.html")
		fmt.Fprint(res, util.SqlConsole)
	} else if req.Method == "POST" {
		stmt := util.ReadReqBody(req)["command"]
		out := util.RunStmt(stmt)
		fmt.Fprint(res, out)
	}
}

func ScdCallback(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "ScdCallback")

	/*dump, _ := httputil.DumpRequest(req, true)
	  util.Info.Printf("dump %q\n", dump)

	  for a, b := range req.URL.Query(){
	      util.Info.Println("req.URL.Query()" , a, ",", b)
	  }*/

	accessToken := req.URL.Query().Get("access_token")
	expiresIn := req.URL.Query().Get("expires_in")

	fmt.Fprint(res, accessToken+"|"+expiresIn)
}

func ScdLogin(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "ScdLogin")

	s := soundcloud.Soundcloud{}
	path := s.Login()

	fmt.Fprint(res, string(path))
}

func GetPlaylists(res http.ResponseWriter, req *http.Request) {
	oauth_token := req.URL.Query().Get("oauth_token")
	saveAccessInfo(res, req, "GetPlaylists, oauth_token : "+oauth_token)

	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	s := soundcloud.Soundcloud{oauth_token}
	util.CompressData(res, s.GetPlaylists())
}

func GetPlaylist(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "GetPlaylist")

	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	s := soundcloud.Soundcloud{}
	vars := getVars(req)
	util.CompressData(res, s.GetPlaylist(vars["playlist"]))
}

func GetSubscriptions(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "GetSubscriptions")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	f := news.Feedly{}

	util.CompressData(res, f.GetSubscriptions())
}

func GetFeedList(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "GetFeedList")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	f := news.Feedly{}
	id := ""
	if cate, ok := req.URL.Query()["id"]; ok {
		id = cate[0]
	}

	util.CompressData(res, f.GetListByFeed(id))
}

func GetEntryContent(res http.ResponseWriter, req *http.Request) {
	saveAccessInfo(res, req, "GetEntryContent")
	res.Header().Set("Content-Type", "text/html; charset=utf-8")

	f := news.Feedly{}
	id := ""
	alternateHref := ""
	if param, ok := req.URL.Query()["id"]; ok {
		id = param[0]
	}

	if param, ok := req.URL.Query()["alternateHref"]; ok {
		alternateHref = param[0]
	}

	util.CompressData(res, f.GetEntryContent(id, alternateHref))
}

func saveAccessInfo(res http.ResponseWriter, req *http.Request, funcName string) {
	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		util.Info.Printf("userip: %q is not IP:port\n", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		//return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
		util.Info.Printf("userip: %q is not IP:port\n", req.RemoteAddr)
	}

	// This will only be defined when site is accessed via non-anonymous proxy
	// and takes precedence over RemoteAddr
	// Header.Get is case-insensitive
	forward := req.Header.Get("X-Forwarded-For")

	util.Info.Printf("<<<<< Client IP: %s, Port %s, Forwarded for: %s, Function called [%s] >>>>>> \n", ip, port, forward, funcName)
}

var (
	ReqVars  map[string]string
	Handlers []*Handler
)

type Handler struct {
	Path    string
	Fn      func(w http.ResponseWriter, req *http.Request)
	Method  string
	PathReg *regexp.Regexp
}

func getVars(req *http.Request) map[string]string {
	if ReqVars != nil {
		return ReqVars
	}
	return mux.Vars(req)
}

func init() {
	Handlers = append(Handlers, &Handler{Path: "/index", Fn: Index})
	Handlers = append(Handlers, &Handler{Path: "/mobile", Fn: Mobile})
	Handlers = append(Handlers, &Handler{Path: "/stock/{cate}", Fn: GetQuoteWithoutFormat})

	Handlers = append(Handlers, &Handler{Path: "/loginAndShowBookInfo", Fn: LoginAndShowBookInfo})
	Handlers = append(Handlers, &Handler{Path: "/showBookInfo", Fn: ShowBookInfo})
	Handlers = append(Handlers, &Handler{Path: "/renewBooks", Fn: RenewBooks})
	Handlers = append(Handlers, &Handler{Path: "/bookCheckoutHist", Fn: BookCheckoutHist})
	Handlers = append(Handlers, &Handler{Path: "/bookSearchHist", Fn: BookSearchHist})
	Handlers = append(Handlers, &Handler{Path: "/searchBooks/{queryText}/{location}/{pageNum}", Fn: SearchBooks})
	Handlers = append(Handlers, &Handler{Path: "/bookDetail/{chamoId}", Fn: GetBookDetail})

	Handlers = append(Handlers, &Handler{Path: "/newsList/query/{query}", Fn: NewsQuery})
	Handlers = append(Handlers, &Handler{Path: "/newsList/{offset}", Fn: NewsList})

	Handlers = append(Handlers, &Handler{Path: "/news/{favor}/{id:[0-9]+}", Fn: NewsFavor, Method: "POST"})
	Handlers = append(Handlers, &Handler{Path: "/news/{id:[0-9]+}", Fn: NewsContent})

	Handlers = append(Handlers, &Handler{Path: "/newsBookmarkList", Fn: BookmarkNews})
	Handlers = append(Handlers, &Handler{Path: "/newsAudioList", Fn: NewsAudioList})
	Handlers = append(Handlers, &Handler{Path: "/news/translated/{id:[0-9]+}/{word}", Fn: TranslateWord, Method: "POST"})

	Handlers = append(Handlers, &Handler{Path: "/citySnapShot/{cate}", Fn: CitySnapShotList})
	Handlers = append(Handlers, &Handler{Path: "/yCantonese", Fn: YCantoneseList})
	Handlers = append(Handlers, &Handler{Path: "/radioList/{cate}", Fn: RTHKRadioList})
	Handlers = append(Handlers, &Handler{Path: "/radioProxy/{channel}/{cate}/{date}/{suffix}", Fn: RTHKRadioProxy})
	Handlers = append(Handlers, &Handler{Path: "/radioUpdate/{cate}", Fn: UpdateRthkRadio, Method: "POST"})
	Handlers = append(Handlers, &Handler{Path: "/podcastList/{cate}", Fn: RTHKPodcastList})
	Handlers = append(Handlers, &Handler{Path: "/podcastDesc/{seq}", Fn: RTHKPodcastDesc})

	Handlers = append(Handlers, &Handler{Path: "/myPocket/{tag}", Fn: MyPocket})
	Handlers = append(Handlers, &Handler{Path: "/myPocket", Fn: MyPocket})

	Handlers = append(Handlers, &Handler{Path: "/oschina/{action}", Fn: OschinaAction})
	Handlers = append(Handlers, &Handler{Path: "/oschina", Fn: Oschina})

	//Handlers = append(Handlers, &Handler{Path: "/hkNewsDetail/{lang}/{suffix}", Fn: GetRthkNewsDetail})
	Handlers = append(Handlers, &Handler{Path: "/hkNews/{cate}", Fn: GetRthkNews})

	Handlers = append(Handlers, &Handler{Path: "/download/{file}", Fn: GetFile})

	Handlers = append(Handlers, &Handler{Path: "/proxy/{url}", Fn: ProxyAccess})

	Handlers = append(Handlers, &Handler{Path: "/sqlConsole", Fn: SqlConsole})


	Handlers = append(Handlers, &Handler{Path: "/scd/callback", Fn: ScdCallback})
	Handlers = append(Handlers, &Handler{Path: "/scd/login", Fn: ScdLogin})
	Handlers = append(Handlers, &Handler{Path: "/scd/playlists", Fn: GetPlaylists})
	Handlers = append(Handlers, &Handler{Path: "/scd/playlist/{playlist}", Fn: GetPlaylist})

	Handlers = append(Handlers, &Handler{Path: "/feedly/getSubscriptions", Fn: GetSubscriptions})
	Handlers = append(Handlers, &Handler{Path: "/feedly/getFeedList", Fn: GetFeedList})
	Handlers = append(Handlers, &Handler{Path: "/feedly/getEntryContent", Fn: GetEntryContent})

	Handlers = append(Handlers, &Handler{Path: "/", Fn: Home})
}
