package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	hkplUrl   = "https://www.hkpl.gov.hk/"
	webcatUrl = "https://webcat.hkpl.gov.hk/"
	wicketUrl = webcatUrl + "wicket/bookmarkable/"
)

// S2T: false for cgi
var S2T = true

func LoginAndShowBookInfo() []byte {
	cookieJar, _ := cookiejar.New(nil)
	client := http.Client{Jar: cookieJar}

	info := Decrypt("cxTR606psLy8HnUONNT6lQ==")
	i := strings.Split(string(info), "X")

	//loginUrl := webcatUrl + "auth/"
	//resp, _ := client.Get(loginUrl + "login")
	//resp, _ := client.Get(hkplUrl  + "tc/login.html")
	//defer resp.Body.Close()
	//body, _ := ioutil.ReadAll(resp.Body)
	//ioutil.WriteFile("login.html", body, 0644)

	//r := regexp.MustCompile(`<form id="login" .+action="(.+?)"`) // get login Form
	//login_form_action := r.FindStringSubmatch(string(body))[1]

	urlData := url.Values{}
	urlData.Set("USER", i[0])
	urlData.Set("PASSWORD", i[1])
	urlData.Set("target", "/auth/login?target=https://www.hkpl.gov.hk/mobile/en/index.html")

	resp, _ := client.PostForm(hkplUrl+"siteminderagent/forms/login.fcc", urlData)

	//Info.Println(loginUrl + login_form_action)
	//resp, _ = client.PostForm(loginUrl+login_form_action, urlData)
	defer resp.Body.Close()

	/*Info.Println("------")
	for _, c := range resp.Cookies() {
		Info.Println(c.Name + ", " + c.Value)
	}*/

	body, _ := ioutil.ReadAll(resp.Body)
	//ioutil.WriteFile("login-result.html", body, 0644)

	resp, _ = client.Get(wicketUrl + "com.vtls.chamo.webapp.component.patron.PatronAccountPage")

	body, _ = ioutil.ReadAll(resp.Body)
	//ioutil.WriteFile("MyAccount.html", body, 0644)

	var _cookies []string
	cookiesMap := make(map[string]string)

	_cookies, cookiesMap = retrieveCookies(*cookieJar, _cookies, cookiesMap, hkplUrl)
	_cookies, cookiesMap = retrieveCookies(*cookieJar, _cookies, cookiesMap, webcatUrl)

	/*for _, c := range _cookies {
		Info.Println(c)
	}*/

	form_action, bookInfos := parseMyAccount(string(body))

	myList := &MyList{form_action, bookInfos, _cookies}

	_myList, _ := json.Marshal(myList)
	return _myList
}

func ShowBookInfo(_cookies []byte) []byte {
	cookieList := strings.Split(string(_cookies), ";")

	var cookies []*http.Cookie
	for _, cookie := range cookieList {
		a := strings.Split(cookie, "=")
		cookies = append(cookies, &http.Cookie{Name: a[0], Value: a[1]})
	}

	cookieJar, _ := cookiejar.New(nil)
	u, _ := url.Parse(webcatUrl)
	cookieJar.SetCookies(u, cookies)

	client := http.Client{Jar: cookieJar}
	resp, _ := client.Get(wicketUrl + "com.vtls.chamo.webapp.component.patron.PatronAccountPage")
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	//ioutil.WriteFile("afterlogin-MyAccount.html", body, 0644)

	form_action, bookInfos := parseMyAccount(string(body))

	_myList, _ := json.Marshal(&MyList{form_action, bookInfos, cookieList})

	return _myList
}

func retrieveCookies(cookieJar cookiejar.Jar, _cookies []string, cookiesMap map[string]string, path string) ([]string, map[string]string) {
	u, _ := url.Parse(path)
	cookies := cookieJar.Cookies(u)

	for _, c := range cookies {
		if _, exist := cookiesMap[c.Name]; !exist {
			_cookies = append(_cookies, c.Name+"="+c.Value)
			cookiesMap[c.Name] = c.Value
		}
	}
	return _cookies, cookiesMap
}

func parseMyAccount(content string) (string, []*BookInfo) {
	//var output []string
	var bookInfos []*BookInfo

	r := regexp.MustCompile(`<form id=.+action="(.+?renewalForm)">`) // get Form
	form_action := r.FindStringSubmatch(content)[1]

	r = regexp.MustCompile(`(?s)<table id="checkout" class="table">(.+?)</table>`) // get table

	table := r.FindStringSubmatch(content)[1]

	r = regexp.MustCompile(`(?s)<tbody>(.+?)</tbody>`) // get tbody

	tbody := r.FindStringSubmatch(table)[1]

	r = regexp.MustCompile(`(?s)<tr.+?>(.+?)</tr>`)

	reBookname := regexp.MustCompile(`<a href=".+dir="ltr">(.+?)</a>`) // get book name

	reDiv := regexp.MustCompile(`<div>(.*?)</div>`)
	reRenewGroup := regexp.MustCompile(`<input type="checkbox".+name="(.+?)".+value="(.+?)" class=.+>`) // get renewable checkbox

	for _, x := range r.FindAllStringSubmatch(tbody, -1) {
		tr := x[1]
		bookName := reBookname.FindStringSubmatch(tr)
		renewCheckbox := reRenewGroup.FindStringSubmatch(tr)
		Info.Println("renewCheckbox ", renewCheckbox)

		var name string
		if len(bookName) > 0 {
			name = bookName[1]

			div := reDiv.FindAllStringSubmatch(tr, -1)
			var date string
			var unit string
			var timesRenewed string
			var renewName string
			var renewValue string

			if len(div) > 0 {
				unit = string(div[0][1])
				date = string(div[2][1])
				timesRenewed = string(div[3][1])
			}
			if len(renewCheckbox) > 0 {
				renewName = renewCheckbox[1]
				renewValue = renewCheckbox[2]
			}

			bookInfo := &BookInfo{name, date, unit, timesRenewed, renewName, renewValue}
			//item, _ := json.Marshal(bookInfo)
			//output = append(output, string(item))
			bookInfos = append(bookInfos, bookInfo)
		}
	}

	return form_action, bookInfos
}

func RenewBooks(myListJson []byte) string {
	var myList MyList
	if err := json.Unmarshal(myListJson, &myList); err != nil {
		panic(err)
	}

	submitUrl := wicketUrl + myList.FormAction
	//Info.Println(submitUrl)

	urlData := url.Values{}
	for _, bookInfo := range myList.BookInfo {
		if bookInfo.RenewValue != "" {
			urlData.Set(bookInfo.RenewName, bookInfo.RenewValue)
			Info.Println("Renew book param", bookInfo.RenewName, bookInfo.RenewValue)
		}
	}

	var cookies []*http.Cookie
	for _, cookie := range myList.Cookies {
		a := strings.Split(cookie, "=")
		cookies = append(cookies, &http.Cookie{Name: a[0], Value: a[1]})
	}

	cookieJar, _ := cookiejar.New(nil)
	u, _ := url.Parse(webcatUrl)
	cookieJar.SetCookies(u, cookies)

	client := http.Client{Jar: cookieJar}
	resp, _ := client.PostForm(submitUrl, urlData)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	//ioutil.WriteFile("renew-result.html", body, 0644)

	r := regexp.MustCompile(`<form id=.+action="(.+?group-confirmationForm)">`) // get confirm form
	form_action := r.FindStringSubmatch(string(body))
	if len(form_action) > 0 {
		Info.Println("Renew a group of items")
		resp, _ = client.PostForm(webcatUrl+"wicket/"+form_action[1], url.Values{})
		defer resp.Body.Close()

		body, _ = ioutil.ReadAll(resp.Body)
		//ioutil.WriteFile("renew-confirm-result.html", body, 0644)
		return retrieveResult(string(body))
	} else {
		return retrieveResult(string(body))
	}

	return "Fail to Renew"
}

func retrieveResult(content string) string {
	r := regexp.MustCompile(`<h2>(.+?)</h2>`)
	renew_result := r.FindStringSubmatch(content)
	if len(renew_result) > 0 {
		return strings.Replace(renew_result[1], "&quot;", "\"", -1)
	}

	Info.Println("retrieve renew result ", content)
	return "Fail to find Result, please check log"
}

func SearchBooks(query string, location string, pageNum string) string {
	urlData := url.Values{}
	//urlData.Set("query", query)

	cookieJar, _ := cookiejar.New(nil)
	client := http.Client{Jar: cookieJar}

	var path string

	Info.Println("keyword", query)
	if S2T {
		query = convertS2t(query)
		Info.Println("keyword converted (s2t)", query)
	}

	if strings.Contains(query, "+") {
		var queryArray []string
		// MUST:包括所有的字, PHASE:完整詞組, SHOULD:至少包括一個字, NOT 不包括, just use MUST
		for i, q := range strings.Split(query, "+") {
			i = i + 1
			queryArray = append(queryArray, fmt.Sprintf("match_%d=%s&field_%d=text&term_%d=%s", i, "PHRASE", i, i, url.QueryEscape(q)))
		}

		path = strings.Join(queryArray, "&")
	} else {
		path = "term_1=" + url.QueryEscape(query)
	}

	if location != "all" {
		if strings.Contains(location, "|") {
			var loc []string
			for _, l := range strings.Split(location, "|") {
				loc = append(loc, "facet_loc="+l)
			}
			path = path + "&" + strings.Join(loc, "&")
		} else {
			path = path + "&facet_loc=" + location
		}
	}

	path = fmt.Sprintf("%ssearch/query?%s&theme=mobile&pageNumber=%s", webcatUrl, path, pageNum)
	Info.Println("search path", path)

	resp, _ := client.PostForm(path, urlData)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	//ioutil.WriteFile("search-result.html", body, 0644)

	//body, _ := ioutil.ReadFile("search-result.html")
	var searchRecords []SearchRecord

	r := regexp.MustCompile(`(?s)<div class="resultCount">.+<span>(結果總數：(.+?)項.+?)</span>.+?</div>`)
	resultDesc := r.FindStringSubmatch(string(body))
	if len(resultDesc) == 0 {
		return ""
	}

	resultTitle := resultDesc[1]
	resultCount := resultDesc[2]

	numExp := regexp.MustCompile(`(?s)<div class="recordNumber">.+?(\d{1,2}).+?</div>`)
	nums := numExp.FindAllStringSubmatch(string(body), -1)
	for _, a := range nums {
		searchRecord := SearchRecord{Title: a[1] + ". "}
		searchRecords = append(searchRecords, searchRecord)
		//fmt.Println("BBB ", a[1])
	}

	titleExp := regexp.MustCompile(`<a href=".+id=chamo:(\d+).+" class="title" dir="ltr">(.+?)</a>`)
	titles := titleExp.FindAllStringSubmatch(string(body), -1)

	for i, a := range titles {
		//fmt.Println("AAA ", a[1], a[2])
		searchRecords[i].Id = a[1]
		searchRecords[i].Title = searchRecords[i].Title + a[2]
	}

	labelExp := regexp.MustCompile(`<td class="label" dir="ltr">(.+?)</td>`)
	label2Exp := regexp.MustCompile(`<td><div><span dir="ltr">(.+?)</span></div></td>`)

	itemSearchResultFieldsExp := regexp.MustCompile(`(?s)<div class="itemSearchResultFields">(.+?)    </div>`)
	itemSearchResultFields := itemSearchResultFieldsExp.FindAllStringSubmatch(string(body), -1)

	for i, a := range itemSearchResultFields {
		//fmt.Println("EE", a[1])
		labels := labelExp.FindAllStringSubmatch(a[1], -1)
		label2s := label2Exp.FindAllStringSubmatch(a[1], -1)

		publication := ""
		edition := ""
		callNumber := ""

		for j, label := range labels {
			name := label[1]
			//fmt.Println("SS ", name)
			if name == "出版項" {
				if len(label2s) > 0 {
					publication = label2s[j][1]
				}
			} else if name == "版次" {
				if len(label2s) > 0 {
					edition = label2s[j][1]
				}
			} else if name == "索書號" {
				if len(label2s) > 0 {
					callNumber = label2s[j][1]
				}
			}

			searchRecords[i].Publication = publication
			searchRecords[i].Edition = edition
			searchRecords[i].CallNumber = callNumber
		}
	}

	locationAvailabilityExp := regexp.MustCompile(`(?s)<ul class="locationAvailability" .+?>(.+?)      </ul><script type="text/javascript">`)
	locationAvailability := locationAvailabilityExp.FindAllStringSubmatch(string(body), -1)

	availabilityLocationExp := regexp.MustCompile(`<span class="availabilityLocation" dir="ltr">(.+?)</span>`)
	//availabilityCountExp := regexp.MustCompile(`<span class="availabilityCount" dir="ltr"(.+?)</span>`)

	for i, a := range locationAvailability {

		if len(availabilityLocationExp.FindAllStringSubmatch(a[1], -1)) > 0 {
			var shelfs []string
			for _, b := range availabilityLocationExp.FindAllStringSubmatch(a[1], -1) {
				shelfs = append(shelfs, b[1])
			}

			/*for _, b := range availabilityCountExp.FindAllStringSubmatch(a[1], -1) {
				fmt.Println("b2", b[1])
			}*/

			searchRecords[i].AvailabeShelfs = shelfs
		}
	}

	searchRecordResult := SearchRecordResult{searchRecords, resultTitle, resultCount}

	_searchRecordResult, _ := json.Marshal(searchRecordResult)

	pn, _ := strconv.Atoi(pageNum)
	if pn == 1 {
		keyword, _ := url.QueryUnescape(query)
		saveKeyword(keyword)
	}

	return string(_searchRecordResult)
}

func saveKeyword(keyword string) {
	db := GetDB()
	defer db.Close()

	var cnt int
	db.QueryRow("select count(word) from Library_Search_Keyword where word = ?", keyword).Scan(&cnt)
	if cnt > 0 {
		return
	}

	Info.Println("Save search keyword", keyword)

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into Library_Search_Keyword values (?, datetime('now', '+8 hour'))")
	stmt.Exec(keyword)
	tx.Commit()
	stmt.Close()
}

func GetBookDetail(chamoId string) []byte {
	Info.Println("Get book of chamoId", chamoId, "detail")
	url := fmt.Sprintf("%slib/item?id=chamo:%s&theme=mobile", webcatUrl, chamoId)
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	//ioutil.WriteFile("googleDrive.html", body, 0644)

	//body, _ := ioutil.ReadFile("googleDrive.html")

	var fields []ItemField

	itemFieldsExp := regexp.MustCompile(`(?s)<div class="itemFields">(.+?)</table>`)
	itemFields := itemFieldsExp.FindStringSubmatch(string(body))

	trExp := regexp.MustCompile(`(?s)<tr>(.+?)</tr>`)
	tdExp := regexp.MustCompile(`<td class="label" dir="ltr">(.+?)</td>`)
	tagExp := regexp.MustCompile(`<div><.+ dir="ltr">(.+?)</.*></div>`)

	//Info.Println(itemFields[0])

	trs := trExp.FindAllStringSubmatch(itemFields[0], -1)

	for _, tr := range trs {
		//Info.Println(tr)
		tds := tdExp.FindAllStringSubmatch(tr[0], -1)
		tag := tagExp.FindAllStringSubmatch(tr[0], -1)

		for i, td := range tds {
			name := td[1]
			value := ""

			if len(tag) > 0 {
				value = tag[i][1]
			}

			fields = append(fields, ItemField{name, value})
		}
	}

	var copies []Copy
	copies = getLibraries(string(body), copies)

	navigatorLabelExp := regexp.MustCompile(`<div class="navigatorLabel"><div>顯示(\d+) 項中第.+</div></div>`)
	navigatorItems := navigatorLabelExp.FindStringSubmatch(string(body))
	if len(navigatorItems) > 0 {
		items, _ := strconv.ParseFloat(navigatorItems[1], 32)

		for i := 1; i <= int(math.Floor(items/10)); i++ {
			url = fmt.Sprintf("%slib/item?id=chamo:%s&theme=mobile&copies-page=%d", webcatUrl, chamoId, i)
			Info.Println("page url", url)
			resp, _ = http.Get(url)
			defer resp.Body.Close()

			body, _ = ioutil.ReadAll(resp.Body)
			copies = getLibraries(string(body), copies)
		}
	}

	titleExp := regexp.MustCompile(`<h1 class="title">(.+)</h1>`)
	title := titleExp.FindStringSubmatch(string(body))

	detail := ItemDetail{title[1], fields, copies}

	_detail, _ := json.Marshal(detail)
	return _detail
}

func trimDiv(char string) string {
	return strings.Trim(strings.Replace(strings.Replace(strings.Replace(char, "<div>", "", -1), "</div>", "", -1), "\n", "", -1), " ")
}

func getLibraries(body string, copies []Copy) []Copy {
	tableExp := regexp.MustCompile(`(?s)<tbody>(.+?)</tbody>`)
	table := tableExp.FindStringSubmatch(body)

	trExp := regexp.MustCompile(`(?s)<tr.+?>(.+?)</tr>`)
	tdExp := regexp.MustCompile(`(?s)<td>(.+?)</td>`)

	trs := trExp.FindAllStringSubmatch(table[1], -1)
	for _, tr := range trs {
		tds := tdExp.FindAllStringSubmatch(tr[1], -1)
		if len(tds) > 0 {
			library := trimDiv(tds[0][1])
			callNumber := trimDiv(tds[1][1])
			units := trimDiv(tds[2][1])
			mediaNumber := trimDiv(tds[3][1])
			status := trimDiv(tds[4][1])
			collection := trimDiv(tds[5][1])
			copies = append(copies, Copy{library, callNumber, units, mediaNumber, status, collection})
		}
	}
	return copies
}

func convertS2t(chars string) string {
	var input []string

	//charLen := utf8.RuneCountInString(chars)

	for len(chars) > 0 {
		r, size := utf8.DecodeRuneInString(chars)
		input = append(input, fmt.Sprintf("%c", r))
		chars = chars[size:]
	}

	dir := GetExecutePath()
	table, err := ioutil.ReadFile(dir + "/../util/ts.tab")
	//table, err := ioutil.ReadFile("ts.tab")
	if err != nil {
		Error.Fatal(err)
	}

	var output []string

	for _, char := range input {
		//Info.Println("char", char, i)
		j := 1
		found := false
		for i := 0; i < len(table)-3; {
			i = i + 3
			if string(table[i:i+3]) == char {
				if j%2 == 0 {
					output = append(output, char)
				} else {
					output = append(output, string(table[i-3:i]))
				}
				found = true
				break
			}
			j += 1
		}
		if !found {
			output = append(output, char)
		}
	}

	return strings.Join(output, "")
}

type BookInfo struct {
	Name         string `json:"name"`
	Date         string `json:"date"`
	Unit         string `json:"unit"`
	TimesRenewed string `json:"timesRenewed"`
	RenewName    string `json:"renewName"`
	RenewValue   string `json:"renewValue"`
}

type MyList struct {
	FormAction string      `json:"formAction"`
	BookInfo   []*BookInfo `json:"bookInfo"`
	Cookies    []string    `json:"cookies"`
}

type SearchRecord struct {
	Id             string
	Title          string
	Author         string
	Publication    string
	CallNumber     string
	Edition        string
	AvailabeShelfs []string
}

type SearchRecordResult struct {
	SearchRecords []SearchRecord
	Result        string
	Total         string
}

type ItemDetail struct {
	Title      string
	ItemFields []ItemField
	Copies     []Copy
}

type ItemField struct {
	Name  string
	Value string
}

type Copy struct {
	Library     string
	CallNumber  string
	Units       string
	MediaNumber string
	Status      string
	Collection  string
}

func BookCheckoutHist() string {
	db := GetDB()
	defer db.Close()
	var nameList []string
	stmt := "select BookName from CheckedOut_Book_History order by seq"
	rows, _ := db.Query(stmt)
	for rows.Next() {
		var name string
		rows.Scan(&name)
		nameList = append(nameList, name)
	}
	return strings.Join(nameList, "|")
}

func BookSearchHist() string {
	db := GetDB()
	defer db.Close()
	var nameList []string
	stmt := "select word from Library_Search_Keyword order by Last_Update"
	rows, _ := db.Query(stmt)
	for rows.Next() {
		var name string
		rows.Scan(&name)
		nameList = append(nameList, name)
	}
	return strings.Join(nameList, "|")
}
