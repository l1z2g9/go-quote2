package rthk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/eduncan911/podcast"
)

func ExportFeedForCitySnapShot() []byte {
	now := time.Now()

	// instantiate a new Podcast
	p := podcast.New(
		"RTHK City Snapshot",
		"http://www.rthk.hk/radio/radio1/programme/City_Snapshot",
		"我們邀請旅居世界各地的名人為節目撰稿及以廣東話聲音演繹，以感性角度去分析他們身處的國家時事。",
		&now, &now,
	)

	// add some channel properties
	p.AddAuthor("Radio Television Hong Kong", "webmaster@rthk.hk")
	//p.AddAtomLink("http://eduncan911.com/feed.rss")
	//p.AddImage("http://janedoe.com/i.jpg")
	p.AddSummary("我們邀請旅居世界各地的名人為節目撰稿及以廣東話聲音演繹，以感性角度去分析他們身處的國家時事。")
	p.IExplicit = "no"

	urlTmpl := "http://www.rthk.hk/radio/catchUp?c=radio1&p=City_Snapshot&page=%d&m="
	for i := 1; i < 27; i++ {
		url := fmt.Sprintf(urlTmpl, i)
		fmt.Println(url)
		resp, _ := http.Get(url)
		defer resp.Body.Close()

		dat, _ := ioutil.ReadAll(resp.Body)

		var audioRow AudioRow
		if err := json.Unmarshal(dat, &audioRow); err != nil {
			panic(err)
		}

		for _, i := range audioRow.Content {
			// text := fmt.Sprintf("%s %s %s\n", i.ID, i.Title, i.Date)
			// fmt.Println(text)

			t, err := time.Parse("02/01/2006", i.Date)
			if err != nil {
				log.Fatal(err)
			}

			// create an Item
			url := fmt.Sprintf("http://stmw3.rthk.hk/aod/_definst_/radio/archive/radio1/City_Snapshot/mp3/mp3:%s.mp3/playlist.m3u8", t.Format("20060102"))
			item := podcast.Item{
				Title:       i.Title,
				Link:        url,
				Description: i.Title,
				PubDate:     &t,
			}

			item.AddEnclosure(url, podcast.MP3, 1000)
			if _, err := p.AddItem(item); err != nil {
				log.Fatal(item.Title, ": error", err.Error())
			}
		}
	}

	return p.Bytes()
}

type AudioRow struct {
	Status  string `json:"status"`
	Content []struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Date   string `json:"date"`
		Photos struct {
			Photo []interface{} `json:"photo"`
		} `json:"photos"`
		Part []string `json:"part"`
	} `json:"content"`
	NextPage int `json:"nextPage"`
}

const ycanList = `0B1tCMe1zn1gHdGxnUlI4eGlMZ00,唐詩七絕選賞14 劉長卿 送李判官之潤州行營,2015-11-28
0B1tCMe1zn1gHbEFRcWJLUTRlOTA,唐詩七絕選賞03 王翰 涼州詞,2015-11-13
0B1tCMe1zn1gHRlVWRTFhdnl5QUk,唐詩七絕選賞10 李白 早發白帝城,2015-11-24
0B1tCMe1zn1gHNU51d2s5TUJ4cXM,唐詩七絕選賞11 高適 別董大,2015-11-25
0B1tCMe1zn1gHNG04d2gwSW9XeHM,唐詩七絕選賞17 韓翃 寒食,2015-12-03
0B1tCMe1zn1gHSjhBLVJvOGFsd3M,唐詩七絕選賞09 王維 送元二使安西,2015-11-20
0B1tCMe1zn1gHLXRQSmNzeGctZG8,唐詩七絕選賞07 王昌齡 閨怨,2015-11-19
0B1tCMe1zn1gHQmhheGZtUWswT1E,唐詩七絕選賞12 張繼 楓橋夜泊,2015-11-26
0B1tCMe1zn1gHMGpjdm5sRDB3Y0k,唐詩七絕選賞15 岑參 虢州後亭送李判官使赴晉絳,2015-12-01
0B1tCMe1zn1gHOXZucF9ZSlR4TGc,唐詩七絕選賞18 王建 十五夜望月寄杜郎中,2015-12-04
0B1tCMe1zn1gHclBmckpxOWtCbUU,唐詩七絕選賞05 王之渙 涼州詞,2015-11-16
0B1tCMe1zn1gHdlYwU0tRcURiWlk,唐詩七絕選賞13 杜甫 江南逢李龜年,2015-11-27
0B1tCMe1zn1gHaDZ4Mk43WU5JMDA,唐詩七絕選賞04 張旭 桃花谿 山中留客,2015-11-13
0B1tCMe1zn1gHRmZpdnpNRnBOQ2s,唐詩七絕選賞06 王昌齡 出塞,2015-11-17
0B1tCMe1zn1gHencyRkJUYjNKa2s,唐詩七絕選賞02 賀知章 回鄉偶書,2015-11-11
0B1tCMe1zn1gHRFlzTUtZbE5HdFU,唐詩七絕選賞01 杜審言 贈蘇綰書記,2015-11-10
0B1tCMe1zn1gHV1BRU0JFeWFPQW8,唐詩七絕選賞08 王維 九月九日憶山東兄弟,2015-11-19
0B1tCMe1zn1gHTGJyUDFiQlFrM2c,唐詩七絕選賞16 錢起 暮春歸故山草堂,2015-12-01
0B1tCMe1zn1gHTzRaYnpJN3ZrNU0,唐詩七絕選賞19 顧況 宮詞 葉上題詩,2015-12-04
0B1tCMe1zn1gHS2VYbWdYRnc0Ukk,唐詩七絕選賞33 朱慶餘 閨意 近試上張水部,2015-12-24
0B1tCMe1zn1gHbTJtMTRrSVk1R1U,唐詩七絕選賞21 張籍 秋思,2015-12-08
0B1tCMe1zn1gHUTlLb2d0VGdYUzQ,唐詩七絕選賞38 杜牧 金谷園,2015-12-31
0B1tCMe1zn1gHeDJCNWxSQ2dtdlk,唐詩七絕選賞39 杜牧 清明,2016-01-01
0B1tCMe1zn1gHc3BYbEY1TzVWMmc,唐詩七絕選賞27 劉禹錫 烏衣巷,2015-12-17
0B1tCMe1zn1gHZUlBNUlfY2hJeG8,唐詩七絕選賞30 李涉 井欄砂宿遇夜客,2015-12-21
0B1tCMe1zn1gHb003MUdhb2M0U0U,唐詩七絕選賞31 韋應物 滁州西澗,2015-12-23
0B1tCMe1zn1gHVzNTLXZCaVNGM1U,唐詩七絕選賞32 元稹 離思 五首之四,2015-12-23
0B1tCMe1zn1gHWTlBODVTNmJPdU0,唐詩七絕選賞29 張祜 贈內人,2015-12-18
0B1tCMe1zn1gHQlAyamlRMExBRkE,唐詩七絕選賞23 柳中庸 征人怨,2015-12-10
0B1tCMe1zn1gHc2w1U19ERy1iVWc,唐詩七絕選賞22 韓愈 早春呈水部張十八員外,2015-12-09
0B1tCMe1zn1gHLTFiWjAySEg1SG8,唐詩七絕選賞28 劉禹錫 石頭城,2015-12-18
0B1tCMe1zn1gHczEyTENOcmtIcDg,唐詩七絕選賞40 溫庭筠 蔡中郎墳,2016-01-04
0B1tCMe1zn1gHMFZTdGlVNld4SlU,唐詩七絕選賞26 柳宗元 與浩初上人同看山寄京華親故,2015-12-16
0B1tCMe1zn1gHTllNZU5kOWF4Ylk,唐詩七絕選賞24 崔護 題都城南莊,2015-12-11
0B1tCMe1zn1gHUmxLaWxJbXh3eFU,唐詩七絕選賞25 李益 夜上受降城聞笛,2015-12-14
0B1tCMe1zn1gHOHlpSldfcGtCNm8,唐詩七絕選賞37 杜牧 赤壁,2015-12-30
0B1tCMe1zn1gHYzVac1VGaXlsMDA,唐詩七絕選賞34 李賀 南園十三首之六,2015-12-25
0B1tCMe1zn1gHWVRvek9LT2xaX2s,唐詩七絕選賞20 劉方平 月夜,2015-12-07
0B1tCMe1zn1gHd0FTWDkxbjVUU3c,唐詩七絕選賞36 杜牧 泊秦淮,2015-12-30
0B1tCMe1zn1gHMzVsa0VUQUdVLWs,唐詩七絕選賞35 杜牧 題烏江亭,2015-12-29
0B1tCMe1zn1gHYXpJVDhkZzdaZ2s,唐詩七絕選賞43 李商隱 嫦娥,2016-01-07
0B1tCMe1zn1gHOWpKdDZJNmdHdkE,唐詩七絕選賞45 鄭畋 馬嵬坡,2016-01-11
0B1tCMe1zn1gHUGdkcF9mV05Kazg,唐詩七絕選賞46 司空圖 河湟有感,2016-01-12
0B1tCMe1zn1gHTUtzUks1eEJmWEE,唐詩七絕選賞49 羅隱 蜂,2016-01-15
0B1tCMe1zn1gHaUhUN0RaYXJkRVk,唐詩七絕選賞52 無名氏 雜詩,2016-01-20
0B1tCMe1zn1gHZGhRR1pfS3h6WEU,唐詩七絕選賞51 張泌 寄人,2016-01-19
0B1tCMe1zn1gHdDBTNkhzX1lJNzg,唐詩七絕選賞54 花蕊夫人 答宋君,2016-01-22
0B1tCMe1zn1gHUVRsRkJSbHpXcDg,唐詩七絕選賞53 杜秋娘 金縷衣,2016-01-21
0B1tCMe1zn1gHWDVHRmotaGZMaUk,唐詩七絕選賞41 薛濤 籌邊樓,2016-01-05
0B1tCMe1zn1gHM0xEc05aeEhJSlk,唐詩七絕選賞44 陳陶 隴西行,2016-01-09
0B1tCMe1zn1gHOGxiMVJMRDZfX1k,唐詩七絕選賞50 羅隱 偶題,2016-01-18
0B1tCMe1zn1gHZXU2ZVVhaFlJbkU,唐詩七絕選賞48 韓偓 已涼,2016-01-14
0B1tCMe1zn1gHTnhJM1lSbmxfYWM,唐詩七絕選賞47 韋莊 臺城,2016-01-13
0B1tCMe1zn1gHMzM5TG5GN25PWEk,唐詩七絕選賞42 李商隱 夜雨寄北,2016-01-06`

func ExportFeedForYCantonese() []byte {
	pubDate, _ := time.Parse("02/01/2006", "31/01/2016")

	// instantiate a new Podcast
	p := podcast.New(
		"粵講越有趣 - 推廣粵文化 粵講越有趣",
		"http://ycantonese.org/",
		"陳耀南教授主講",
		&pubDate, &pubDate,
	)

	// add some channel properties
	p.AddAuthor("Radio Television Hong Kong", "webmaster@rthk.hk")

	p.AddSummary("本網站是為推廣粵文化而設。粵文化流傳的地區包括南中國、港、澳及海外華人僑居地。我們提供的視頻有風物誌及唐詩等。粵語是由中原南傳的唐代古語─雅言, 故用粵語讀唐詩及古文最傳神。敬希指導。")
	p.IExplicit = "no"

	for _, line := range strings.Split(ycanList, "\n") {
		part := strings.Split(line, ",")

		// create an Item
		url := fmt.Sprintf("https://drive.google.com/uc?id=%s&authuser=0&export=download", part[0])

		date, _ := time.Parse("2006-01-02", part[2])
		item := podcast.Item{
			Title:       part[1],
			Link:        url,
			Description: part[1],
			PubDate:     &date,
		}

		item.AddEnclosure(url, podcast.MP3, 1000)
		if _, err := p.AddItem(item); err != nil {
			log.Fatal(item.Title, ": error", err.Error())
		}
	}

	return p.Bytes()
}
