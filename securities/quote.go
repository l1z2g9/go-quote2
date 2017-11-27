package securities

import (
	"github.com/l1z2g9/go-quote2/util"
	"database/sql"
	"encoding/json"
	"fmt"
	iconv "github.com/djimenez/iconv-go"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	sortField = "PercentChange" //Change, PercentDayRangeWithYesterday, PercentDayRangeWithTodayLow,
	//DaysRange, Money,  Volumn, Bid, YesterdayClosePrice
	date string
)

const (
//DB_PATH string = "/home/l1z2g9/Applications/gocode/src/go-quote/stock.db"
//DB_PATH string = "/var/lib/openshift/5524d4c0e0b8cd7ac5000109/app-root/runtime/myapp/go-quote/stock.db"
//DB_PATH string = "/var/lib/openshift/5524d4c0e0b8cd7ac5000109/app-root/runtime/repo/stock.db"
)

func SetSortField(sf string) {
	sortField = sf
}

func GetQuoteWithFormat(category string) string {
	t := time.Now()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t = t.In(loc)
	time := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	dataList := getRawQuote(category)

	var output []string
	// local index
	index := []string{"000001", "399001", "399005", "399006"}
	var sumMoney float64

	for _, i := range index {
		q := dataList.Composite_Index[i]
		output = append(output, fmt.Sprintf(
			"%s(%s) - Bid: %5.3f, Change %5.2f, PercentChange %5.2f%%, Volumn %3.1f亿手 , Money %5.2f亿",
			q.Symbol, q.Name, q.Bid, q.Change, q.PercentChange, q.Volumn, q.Money))

		sumMoney = sumMoney + q.Money
		if q.Symbol == "399001" {
			output = append(output, fmt.Sprintf("Total: %5.2f亿", sumMoney))
		}
	}
	output = append(output, " ")

	// index future
	index = []string{"IH" + date, "IF" + date, "IC" + date}
	for _, i := range index {
		q := dataList.Index_Future[i]
		output = append(output, fmt.Sprintf(
			"%-13s - Bid: %5.3f, Change %5.2f, PercentChange %5.2f%%, Low %-4.2f, High %-4.2f, Volumn %3.0f张 , Positions %5.0f",
			q.Name, q.Bid, q.Change, q.PercentChange, q.Low, q.High, q.Volumn, q.Position))
	}
	output = append(output, " ")

	// foreign index
	index = []string{"hangseng", "nikkei", "dji", "nasdaq"}
	for _, i := range index {
		q := dataList.Composite_Index[i]
		output = append(output, fmt.Sprintf(
			"%s(%s) - Bid: %5.3f, Change %5.2f, PercentChange %5.2f%%", q.Symbol, q.Name, q.Bid, q.Change, q.PercentChange))
	}
	output = append(output, " ")

	// quote
	for _, quote := range dataList.Quotes {
		output = append(output, fmt.Sprintf("%s(%-4s) Bid: %-5.2f Change %5.2f PercentChange %5.2f%% DaysRange %-10s / "+
			"%4.2f (%4.2f%% | %4.2f%%) Position [%-4.2f]  Low-High Stop %-5.2f - %-5.2f, Volumn %3.0f万手  Money %5.2f亿 %s", quote.Symbol,
			quote.Name, quote.Bid, quote.GetChange(), quote.GetPercentChange(), quote.GetDaysRange(),
			quote.High-quote.Low, quote.GetPercentDayRangeWithYesterday(),
			quote.GetPercentDayRangeWithTodayLow(), quote.Position, quote.YesterdayClosePrice*0.9, quote.YesterdayClosePrice*1.1, quote.Volumn, quote.Money,
			quote.Message))
	}
	output = append(output, " ")

	// goldprice
	output = append(output, "Gold Price From goldprice: "+dataList.GoldPrice)
	output = append(output, " ")

	// GoldPriceFromKITCO
	/*output = append(output, "Gold Price From Kitco: "+dataList.GoldPriceFromKITCO[0])
	for _, g := range dataList.GoldPriceFromKITCO[1:] {
		output = append(output, g)
	}*/

	goldAndFx := dataList.GoldPriceFromKITCO
	output = append(output, "Gold Price From Kitco: \n")
	output = append(output, goldAndFx["gold"])
	output = append(output, goldAndFx["nikkei"])
	output = append(output, goldAndFx["oil"])
	output = append(output, goldAndFx["usdIndex"])

	output = append(output, " ")

	// ExchangeRate
	rate := dataList.ExchangeRate
	output = append(output, "USDCNY: "+strconv.FormatFloat(rate["USDCNY"], 'g', 4, 32))
	output = append(output, "USDHKD: "+strconv.FormatFloat(rate["USDHKD"], 'g', 4, 32))
	output = append(output, "HKDCNY: "+strconv.FormatFloat(rate["HKDCNY"], 'g', 4, 32))
	output = append(output, "USDJPY: "+strconv.FormatFloat(rate["USDJPY"], 'g', 4, 32))
	output = append(output, " ")

	output = append(output, "Update at: "+time)

	return strings.Join(output, "\n")
}

func GetQuoteWithoutFormat(category string) string {
	dataList := getRawQuote(category)
	d, _ := json.Marshal(&dataList)
	return string(d)
}

func getRawQuote(category string) DataList {
	code, holds, observation := getFavorCodes(category)
	url, data := GetPriceFromSina(code)

	quotes, composite_index, index_future := parseData(url, data, holds, observation)

	dataList := DataList{quotes, composite_index, index_future, getGoldPrice(), getGoldPriceFromKITCO(), getExchangeRate()}

	return dataList
}

func getFavorCodes(category string) (string, string, map[string]float64) {
	util.Info.Println("category", category)
	var codes []string
	var holds []string
	var observations map[string]float64 = make(map[string]float64)

	db := util.GetDB()
	defer db.Close()

	var rows *sql.Rows

	if len(category) > 0 {
		if category == "favorite" {
			rows, _ = db.Query("select symbol, held, observation from profile where favor ='T'")
		} else if category == "held" {
			rows, _ = db.Query("select symbol, held, observation from profile where held = 'T'")
		} else if category == "nb" { //沪港通
			rows, _ = db.Query("select symbol, held, observation from profile where northbound = 'T'")
		} else if category == "hk" { //港股
			rows, _ = db.Query("select symbol, held, observation from profile where hkShare = 'T'")
		} else if category == "ah" { // A + H 股
			rows, _ = db.Query("select symbol, held, observation from profile where industry_1 = 'A+H'")
		} else {
			rows, _ = db.Query("select symbol, held, observation from profile where Industry_4 = ?", category)
		}
	} else {
		rows, _ = db.Query("select symbol, held, observation from profile where favor ='T'")
	}

	defer rows.Close()

	for rows.Next() {
		var code string
		var held sql.NullString
		var observation sql.NullFloat64

		err := rows.Scan(&code, &held, &observation)
		if err != nil {
			panic(err)
		}

		codes = append(codes, code)

		if held.Valid && strings.Contains(held.String, "T") {
			holds = append(holds, code)
		}

		if observation.Valid {
			observations[code] = observation.Float64
		}
	}

	return strings.Join(codes, ","), strings.Join(holds, ","), observations
}

func GetPriceFromSina(codes string) (string, string) {
	var codes_ []string

	for _, code := range strings.Split(codes, ",") {
		if strings.HasPrefix(code, "60") {
			codes_ = append(codes_, "sh"+code)
		} else if len(code) == 5 {
			codes_ = append(codes_, "hk"+code)
		} else {
			codes_ = append(codes_, "sz"+code)
		}
	}

	url := "http://hq.sinajs.cn/list=" + strings.Join(codes_, ",") + ",s_sh000001,s_sz399001,s_sz399005,s_sz399006"
	//url += ",CFF_RE_IC1507,CFF_RE_IC1508,CFF_RE_IC1509,CFF_RE_IC1512,CFF_RE_IF1507,CFF_RE_IF1508,CFF_RE_IF1509,CFF_RE_IF1512,CFF_RE_IH1507,CFF_RE_IH1508,CFF_RE_IH1509,CFF_RE_IH1512"

	t := time.Now()

	if t.After(getBalanceSheetDate()) {
		if int(t.Month()) == 12 {
			date = strconv.Itoa(t.Year() + 1)[2:4] + "01"
		} else {
			date = strconv.Itoa(t.Year())[2:4] + leftPad(strconv.Itoa(int(t.Month())+1))
		}
	} else {
		date = strconv.Itoa(t.Year())[2:4] + leftPad(strconv.Itoa(int(t.Month())))
	}

	url += ",CFF_RE_IH" + date + ",CFF_RE_IF" + date + ",CFF_RE_IC" + date
	url += ",int_hangseng,int_nikkei,int_dji,int_nasdaq" //hk00001长和, hk00700腾讯控股
	util.Info.Println("url = " + url)
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	content, _ := iconv.ConvertString(string(body), "GBK", "utf-8")

	return url, content
}

func getBalanceSheetDate() time.Time {
	t := time.Now()

	beginningOfHour := time.Duration(-t.Hour()) * time.Hour

	beginningOfDay := t.Truncate(time.Hour).Add(beginningOfHour)

	beginningOfMonth := beginningOfDay.Add(time.Duration(-int(t.Day())+1) * 24 * time.Hour)

	fifthteen := beginningOfMonth.AddDate(0, 0, 14)

	balanceSheetDate := fifthteen.AddDate(0, 0, (5 - int(fifthteen.Weekday())))

	util.Info.Println("balanceSheetDate: ", balanceSheetDate.Day())

	return balanceSheetDate
}

func parseData(url string, content string, holds string, observations map[string]float64) (Quotes, map[string]Quote, map[string]Quote) {
	//var output []string
	//output = append(output, url+"\n")
	composite_index := make(map[string]Quote)
	index_future := make(map[string]Quote)

	var stocks = strings.Split(content, "var hq_str_")

	//var sumMoney float64
	var quotes Quotes
	r, _ := regexp.Compile("([a-z]+)(\\d+)=\"(.+)\"")

	for _, stock := range stocks {
		if strings.HasPrefix(stock, "int") {
			r, _ := regexp.Compile("int_(\\w+)=\"(.+)\"")
			result := r.FindStringSubmatch(stock)
			if len(result) > 0 {
				symbol := result[1]
				info := result[2]
				quote := strings.Split(info, ",")
				name := quote[0]

				bid_, _ := strconv.ParseFloat(quote[1], 32)
				change_, _ := strconv.ParseFloat(quote[2], 32)
				percentChange_, _ := strconv.ParseFloat(strings.Replace(quote[3], "%", "", -1), 32)
				//output = append(output, fmt.Sprintf(
				//  "%s(%s) - Bid: %5.3f, Change %5.2f, PercentChange %5.2f%%\n", symbol, name, bid_, change_, percentChange_))
				mQuote := Quote{Name: name, Bid: bid_, Symbol: symbol, Change: change_, PercentChange: percentChange_}

				composite_index[symbol] = mQuote
			}
			continue
		}

		if strings.HasPrefix(stock, "CFF_RE_") {
			r, _ := regexp.Compile("CFF_RE_(.+)=\"(.+)\"")
			result := r.FindStringSubmatch(stock)
			if len(result) > 0 {
				symbol := result[1]
				var name string
				if strings.HasPrefix(symbol, "IC") {
					name = symbol + "(中证500)"
				} else if strings.HasPrefix(symbol, "IF") {
					name = symbol + "(沪深300)"
				} else if strings.HasPrefix(symbol, "IH") {
					name = symbol + "(上证50)"
				}

				info := result[2]
				quote := strings.Split(info, ",")

				//start := quote[0]

				high_, _ := strconv.ParseFloat(quote[1], 32)
				low_, _ := strconv.ParseFloat(quote[2], 32)

				bid_, _ := strconv.ParseFloat(quote[3], 32)
				yesterdayClosePrice_, _ := strconv.ParseFloat(quote[14], 32)
				change_ := bid_ - yesterdayClosePrice_
				percentChange_ := change_ / yesterdayClosePrice_ * 100

				volumn_, _ := strconv.ParseFloat(quote[4], 32)

				positions_, _ := strconv.ParseFloat(quote[6], 32)

				/*output = append(output, fmt.Sprintf(
				      "%-13s - Bid: %5.3f, Change %5.2f, PercentChange %5.2f%%, Start: %s, Low %-4.2f, High %-4.2f, Volumn %3.0f张 , Positions %5.0f\n",
				      symbol, bid_, change_, percentChange_, start, low_, high_, volumn_, positions_))
				  if strings.Index(symbol, "中证500") != -1 {
				      output = append(output, "\n")
				  }*/

				mQuote := Quote{Bid: bid_, Symbol: symbol, Name: name, Change: change_, PercentChange: percentChange_,
					YesterdayClosePrice: yesterdayClosePrice_, High: high_, Low: low_, Volumn: volumn_, Position: positions_}
				index_future[symbol] = mQuote
			}
			continue
		}

		result := r.FindStringSubmatch(stock)
		if len(result) > 0 {
			kind := result[1]
			symbol := result[2]
			info := result[3]
			// util.Info.Println("XXXX kind " + kind)
			if strings.Contains(kind, "hk") {
				var infoprefix []string

				info2 := strings.Split(info, ",")
				infoprefix = append(infoprefix, (info2[1]))

				for i := 2; i < len(info2); i++ {
					infoprefix = append(infoprefix, info2[i])
				}

				info = strings.Join(infoprefix, ",")
			}

			quote := strings.Split(info, ",")
			name := quote[0]

			if strings.HasPrefix(symbol, "39900") || symbol == "000001" {
				bid_, _ := strconv.ParseFloat(quote[1], 32)
				change_, _ := strconv.ParseFloat(quote[2], 32)
				percentChange_, _ := strconv.ParseFloat(quote[3], 32)

				volumn_, _ := strconv.ParseFloat(quote[4], 32)
				volumn := volumn_ / 1000000

				if name == "深证成指" {
					volumn = volumn_ / 100000000
				}

				money_, _ := strconv.ParseFloat(quote[5], 32)

				money := money_ / 10000
				//sumMoney = sumMoney + money

				/*output = append(output, fmt.Sprintf(
				      "%s(%s) - Bid: %5.3f, Change %5.2f, PercentChange %5.2f%%, Volumn %3.1f亿手 , Money %5.2f亿\n",
				      symbol, name, bid_, change_, percentChange_, volumn, money))

				  if name == "深证成指" {
				      output = append(output, fmt.Sprintf("Total: %5.2f亿\n", sumMoney))
				  }
				  if strings.Index(symbol, "399006") != -1 {
				      output = append(output, "\n")
				  }*/

				mQuote := Quote{Name: name, Bid: bid_, Symbol: symbol, Change: change_, PercentChange: percentChange_, Volumn: volumn, Money: money}
				composite_index[symbol] = mQuote
			} else {
				yesterdayClosePrice_, _ := strconv.ParseFloat(quote[2], 32)
				bid_, _ := strconv.ParseFloat(quote[3], 32)
				high_, _ := strconv.ParseFloat(quote[4], 32)
				low_, _ := strconv.ParseFloat(quote[5], 32)

				if strings.Contains(kind, "hk") {
					bid_, _ = strconv.ParseFloat(quote[5], 32)
					high_, _ = strconv.ParseFloat(quote[3], 32)
					low_, _ = strconv.ParseFloat(quote[4], 32)
				}

				volumn_, _ := strconv.ParseFloat(quote[8], 32)
				money_, _ := strconv.ParseFloat(quote[9], 32)

				quote_ := Quote{Name: name, Symbol: symbol, Bid: bid_, YesterdayClosePrice: yesterdayClosePrice_}

				/*quote_.Name = name
				  quote_.Bid = bid_
				  quote_.Symbol = symbol
				  quote_.YesterdayClosePrice = yesterdayClosePrice_*/
				quote_.High = high_
				quote_.Low = low_
				if high_ != low_ {
					quote_.Position = (bid_ - low_) / (high_ - low_)
				} else {
					quote_.Position = 0
				}
				quote_.Volumn = volumn_ / (100 * 10000)
				quote_.Money = money_ / 100000000

				if _, ok := observations[symbol]; ok {
					msg := fmt.Sprintf("::: %5.2f -> %5.2f", observations[symbol], bid_)
					if quote_.Bid < observations[symbol] {
						msg += " ### Buy it ###"
					}
					quote_.Message = msg
				}

				if strings.Contains(holds, symbol) {
					//util.Info.Println("holds symbol = ", symbol, "("+holds+")")
					quote_.Message = quote_.Message + " <<<<<< held >>>>>>"
					quote_.Held = 1
				}

				quotes = append(quotes, quote_)
			}
		}
	}

	//output = append(output, fmt.Sprintf("\n"))

	//qs := QuoteSorter{quotes, sortField}
	sort.Sort(quotes)

	/*for _, quote := range quotes {
	      output = append(output, fmt.Sprintf("%s(%-4s) Bid: %-5.2f Change %5.2f PercentChange %5.2f%% DaysRange %-10s / "+
	          "%4.2f (%4.2f%% | %4.2f%%) Position [%-4.2f]  Low-High Stop %-5.2f - %-5.2f, Volumn %3.0f万手  Money %5.2f亿 %s\n", quote.Symbol,
	          quote.Name, quote.Bid, quote.GetChange(), quote.GetPercentChange(), quote.GetDaysRange(),
	          quote.High-quote.Low, quote.GetPercentDayRangeWithYesterday(),
	          quote.GetPercentDayRangeWithTodayLow(), (quote.Bid-quote.Low)/(quote.High-quote.Low), quote.YesterdayClosePrice*0.9, quote.YesterdayClosePrice*1.1, quote.Volumn, quote.Money,
	          quote.Message))
	  }

	  return strings.Join(output, "")*/

	//dataList := DataList{quotes, composite_index, index_future}
	return quotes, composite_index, index_future
}

func getGoldPrice() string {
	/*resp, _ := http.Get("http://api2.goldprice.org/Service.svc/GetRaw/3")
	  defer resp.Body.Close()

	  body, _ := ioutil.ReadAll(resp.Body)
	  result := string(body)
	  price := strings.Split(result[3:len(result)-4], ",")
	  //result = fmt.Sprint("Gold Price From goldprice: ", price[0]) //, "\nUSD -> RMB: " , price[9])

	  /*for i, p := range price {
	      fmt.Println(i, p)
	  }*/
	//return price[0]
	return ""
}

func getGoldPriceFromKITCO() map[string]string {
	//var result []string
	// gold
	url := "http://wap.kitco.cn/"
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	page := string(body)
	r := regexp.MustCompile(`(?s)<p>(金.*?)</p>`)
	//result = append(result, fmt.Sprint("Gold Price From Kitco: ", r.FindStringSubmatch(page)[1]))

	result := make(map[string]string)
	result["gold"] = r.FindStringSubmatch(page)[1]
	//result = append(result, r.FindStringSubmatch(page)[1])

	// index
	url = "http://wap.kitco.cn/idxs.wml"
	resp, _ = http.Get(url)
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)

	page = string(body)

	//result = append(result, kitcoData("日经指数", page, "%-6s %-10s %s"))
	//result = append(result, kitcoData("原油", page, "%-8s %-10s %s"))
	//result = append(result, kitcoData("美元指数", page, "%-6s %-10s %s"))
	result["nikkei"] = kitcoData("日经指数", page, "%-6s %-10s %s")
	result["oil"] = kitcoData("原油", page, "%-8s %-10s %s")
	result["usdIndex"] = kitcoData("美元指数", page, "%-6s %-10s %s")
	//return strings.Join(result, "\n")
	return result
}

func kitcoData(index string, page string, format string) string {
	m := regexp.MustCompile("(?s)" + index + "<br/>(.*?)<br/>(.*?)<br/>").FindStringSubmatch(page)
	result := fmt.Sprintf(format, index, m[1][2:], m[2][2:])
	return result
}

func getExchangeRateOld() map[string]float64 {
	resp, _ := http.Get("http://tw.rter.info/capi.php")
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	result := string(body)

	var dat map[string]map[string]interface{}
	byt := []byte(result)
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}

	rate := make(map[string]float64)
	usdcny := dat["USDCNY"]["Exrate"].(float64)
	usdhkd := dat["USDHKD"]["Exrate"].(float64)
	usdjpy := dat["USDJPY"]["Exrate"].(float64)

	//hkdcny := strconv.FormatFloat(usdcny/usdhkd*100, 'f', 4, 64)
	hkdcny := usdcny / usdhkd * 100

	rate["USDCNY"] = usdcny
	rate["USDHKD"] = usdhkd
	rate["HKDCNY"] = hkdcny
	rate["USDJPY"] = usdjpy

	//result = fmt.Sprint("\nUSDCNY: ", usdcny, "\nUSDHKD: ", usdhkd, "\nHKDCNY: ", hkdcny, "\nUSDJPY: ", usdypy, "\n")

	return rate
}

func getExchangeRate() map[string]float64 {
	resp, _ := http.Get("http://wap.kitco.cn/exch.wml")
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	page := string(body)

	rate := make(map[string]float64)

	r := regexp.MustCompile(`人民币 (.*?)<br/>`)
	usdcny, _ := strconv.ParseFloat(r.FindStringSubmatch(page)[1], 32)

	r = regexp.MustCompile(`港币 (.*?)<br/>`)
	usdhkd, _ := strconv.ParseFloat(r.FindStringSubmatch(page)[1], 32)

	hkdcny := usdcny / usdhkd * 100

	r = regexp.MustCompile(`日元 (.*?)<br/>`)
	usdjpy, _ := strconv.ParseFloat(r.FindStringSubmatch(page)[1], 32)

	rate["USDCNY"] = usdcny
	rate["USDHKD"] = usdhkd
	rate["HKDCNY"] = hkdcny
	rate["USDJPY"] = usdjpy

	return rate
}

func leftPad(s string) string {
	if len(s) == 1 {
		s = "0" + s
	}
	return s
}

type Quotes []Quote

func (qs Quotes) Len() int {
	return len(qs)
}

func (qs Quotes) Swap(i, j int) {
	qs[i], qs[j] = qs[j], qs[i]
}

func (qs Quotes) Less(i, j int) bool {
	//x := reflect.ValueOf(&qs.quotes[i]).MethodByName("get" + qs.sortField).Call([]reflect.Value{}).Float()
	//y := reflect.ValueOf(&qs.quotes[j]).MethodByName("get" + qs.sortField).Call([]reflect.Value{})

	//fmt.Println(x)
	//fmt.Println(qs[i])

	a := reflect.ValueOf(&qs[i]).MethodByName("Get" + sortField).Call([]reflect.Value{})[0].Float()
	b := reflect.ValueOf(&qs[j]).MethodByName("Get" + sortField).Call([]reflect.Value{})[0].Float()

	//return qs[i].GetPercentChangea() > qs[j].GetPercentChangea()
	return a > b
}

type Quote struct {
	Name                string
	Bid                 float64
	Symbol              string
	YesterdayClosePrice float64
	Change              float64
	PercentChange       float64
	Position            float64
	High                float64
	Low                 float64
	Volumn              float64
	Money               float64
	Message             string
	Held                int
}

func (this Quote) GetPercentChange() float64 {
	return (this.Bid - this.YesterdayClosePrice) / this.YesterdayClosePrice * 100
}

func (this Quote) GetChange() float64 {
	return this.Bid - this.YesterdayClosePrice
}

func (this Quote) GetPercentDayRangeWithYesterday() float64 {
	var dayRange float64
	if this.High != this.Low {
		dayRange = (this.High - this.Low) / this.YesterdayClosePrice * 100
	}
	return dayRange
}

func (this Quote) GetPercentDayRangeWithTodayLow() float64 {
	var dayRange float64
	if this.High != this.Low {
		dayRange = (this.High - this.Low) / this.Low * 100
	}
	return dayRange
}

func (this Quote) GetDaysRange() string {
	return strconv.FormatFloat(this.Low, 'g', -1, 32) + "-" + strconv.FormatFloat(this.High, 'g', -1, 32)
}

func (this Quote) GetMoney() float64 {
	return this.Money
}

func (this Quote) GetVolumn() float64 {
	return this.Volumn
}

func (this Quote) GetBid() float64 {
	return this.Bid
}

func (this Quote) GetYesterdayClosePrice() float64 {
	return this.YesterdayClosePrice
}

type DataList struct {
	Quotes             []Quote
	Composite_Index    map[string]Quote
	Index_Future       map[string]Quote
	GoldPrice          string
	GoldPriceFromKITCO map[string]string
	ExchangeRate       map[string]float64
}
