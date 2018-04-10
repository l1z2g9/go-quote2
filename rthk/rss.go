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

const ycanList = `0B1tCMe1zn1gHdGxnUlI4eGlMZ00,唐詩七絕選賞14 劉長卿 送李判官之潤州行營,0B1tCMe1zn1gHQUZ3Q3o3cXhWdkE
0B1tCMe1zn1gHbEFRcWJLUTRlOTA,唐詩七絕選賞03 王翰 涼州詞,0B1tCMe1zn1gHb0Q5cWF4djFuUGM
0B1tCMe1zn1gHRlVWRTFhdnl5QUk,唐詩七絕選賞10 李白 早發白帝城,0B1tCMe1zn1gHcm9mel9FRjZ5UFE
0B1tCMe1zn1gHNU51d2s5TUJ4cXM,唐詩七絕選賞11 高適 別董大,0B1tCMe1zn1gHdXBueTdLamNCWWs
0B1tCMe1zn1gHNG04d2gwSW9XeHM,唐詩七絕選賞17 韓翃 寒食,0B1tCMe1zn1gHMXloeDc3WmJkelE
0B1tCMe1zn1gHSjhBLVJvOGFsd3M,唐詩七絕選賞09 王維 送元二使安西,0B1tCMe1zn1gHbDg3WUtoLXVfM1k
0B1tCMe1zn1gHLXRQSmNzeGctZG8,唐詩七絕選賞07 王昌齡 閨怨,0B1tCMe1zn1gHUkEzdFN5Z0tEUUk
0B1tCMe1zn1gHQmhheGZtUWswT1E,唐詩七絕選賞12 張繼 楓橋夜泊,0B1tCMe1zn1gHMXNsMU9uWXBZQlk
0B1tCMe1zn1gHMGpjdm5sRDB3Y0k,唐詩七絕選賞15 岑參 虢州後亭送李判官使赴晉絳,0B1tCMe1zn1gHRERqYVpVWFREa28
0B1tCMe1zn1gHOXZucF9ZSlR4TGc,唐詩七絕選賞18 王建 十五夜望月寄杜郎中,0B1tCMe1zn1gHNHpBNGkwbkVDeE0
0B1tCMe1zn1gHclBmckpxOWtCbUU,唐詩七絕選賞05 王之渙 涼州詞,0B1tCMe1zn1gHRUxGN3B3RWN2OE0
0B1tCMe1zn1gHdlYwU0tRcURiWlk,唐詩七絕選賞13 杜甫 江南逢李龜年,0B1tCMe1zn1gHMUJFMnh2UFl6ekE
0B1tCMe1zn1gHaDZ4Mk43WU5JMDA,唐詩七絕選賞04 張旭 桃花谿 山中留客,0B1tCMe1zn1gHVjBEYVVLZm9WOHM
0B1tCMe1zn1gHRmZpdnpNRnBOQ2s,唐詩七絕選賞06 王昌齡 出塞,0B1tCMe1zn1gHN3JHRUtjUk5sNjA
0B1tCMe1zn1gHencyRkJUYjNKa2s,唐詩七絕選賞02 賀知章 回鄉偶書,0B1tCMe1zn1gHOWFqazZfcmJvUWM
0B1tCMe1zn1gHRFlzTUtZbE5HdFU,唐詩七絕選賞01 杜審言 贈蘇綰書記,0B1tCMe1zn1gHYUl6ZDNZbm5vSW8
0B1tCMe1zn1gHV1BRU0JFeWFPQW8,唐詩七絕選賞08 王維 九月九日憶山東兄弟,0B1tCMe1zn1gHc1Z5TEp1NndNLU0
0B1tCMe1zn1gHTGJyUDFiQlFrM2c,唐詩七絕選賞16 錢起 暮春歸故山草堂,0B1tCMe1zn1gHZENsQ1dOZjFDV3c
0B1tCMe1zn1gHTzRaYnpJN3ZrNU0,唐詩七絕選賞19 顧況 宮詞 葉上題詩,0B1tCMe1zn1gHSmwxZjc2eVEtMFU
0B1tCMe1zn1gHS2VYbWdYRnc0Ukk,唐詩七絕選賞33 朱慶餘 閨意 近試上張水部,0B1tCMe1zn1gHN2gxVGxJeFRaQ2s
0B1tCMe1zn1gHbTJtMTRrSVk1R1U,唐詩七絕選賞21 張籍 秋思,0B1tCMe1zn1gHSFBpZk14MWY0ckE
0B1tCMe1zn1gHUTlLb2d0VGdYUzQ,唐詩七絕選賞38 杜牧 金谷園,0B1tCMe1zn1gHdnZpSFlBYlZndWs
0B1tCMe1zn1gHeDJCNWxSQ2dtdlk,唐詩七絕選賞39 杜牧 清明,0B1tCMe1zn1gHcC1fdVJoTUFiejA
0B1tCMe1zn1gHc3BYbEY1TzVWMmc,唐詩七絕選賞27 劉禹錫 烏衣巷,0B1tCMe1zn1gHMlg2eWowYjM4SGc
0B1tCMe1zn1gHZUlBNUlfY2hJeG8,唐詩七絕選賞30 李涉 井欄砂宿遇夜客,0B1tCMe1zn1gHNW1zeGR3REhXd2s
0B1tCMe1zn1gHb003MUdhb2M0U0U,唐詩七絕選賞31 韋應物 滁州西澗,0B1tCMe1zn1gHRjQyZDBCM3ZUSnM
0B1tCMe1zn1gHVzNTLXZCaVNGM1U,唐詩七絕選賞32 元稹 離思 五首之四,0B1tCMe1zn1gHY3NOSU5fcWF1Vms
0B1tCMe1zn1gHWTlBODVTNmJPdU0,唐詩七絕選賞29 張祜 贈內人,0B1tCMe1zn1gHYVpTTm1DQktXcGs
0B1tCMe1zn1gHQlAyamlRMExBRkE,唐詩七絕選賞23 柳中庸 征人怨,0B1tCMe1zn1gHR0o0TEtDU1VFbEE
0B1tCMe1zn1gHc2w1U19ERy1iVWc,唐詩七絕選賞22 韓愈 早春呈水部張十八員外,0B1tCMe1zn1gHQW1NVGdZSGIzVVE
0B1tCMe1zn1gHLTFiWjAySEg1SG8,唐詩七絕選賞28 劉禹錫 石頭城,0B1tCMe1zn1gHOFVKQUM1N1lSbnM
0B1tCMe1zn1gHczEyTENOcmtIcDg,唐詩七絕選賞40 溫庭筠 蔡中郎墳,0B1tCMe1zn1gHXzk5X0ZZS0NPUUU
0B1tCMe1zn1gHMFZTdGlVNld4SlU,唐詩七絕選賞26 柳宗元 與浩初上人同看山寄京華親故,0B1tCMe1zn1gHa0x3OXpSZTdPS00
0B1tCMe1zn1gHTllNZU5kOWF4Ylk,唐詩七絕選賞24 崔護 題都城南莊,0B1tCMe1zn1gHRXhqY0l3ZlJFQm8
0B1tCMe1zn1gHUmxLaWxJbXh3eFU,唐詩七絕選賞25 李益 夜上受降城聞笛,0B1tCMe1zn1gHSElaYUNFMkhWeW8
0B1tCMe1zn1gHOHlpSldfcGtCNm8,唐詩七絕選賞37 杜牧 赤壁,0B1tCMe1zn1gHU2laVGg4bE85U2M
0B1tCMe1zn1gHYzVac1VGaXlsMDA,唐詩七絕選賞34 李賀 南園十三首之六,0B1tCMe1zn1gHanJMbUNLTzRMRms
0B1tCMe1zn1gHWVRvek9LT2xaX2s,唐詩七絕選賞20 劉方平 月夜,0B1tCMe1zn1gHZTVITUdwZ2FsQXM
0B1tCMe1zn1gHd0FTWDkxbjVUU3c,唐詩七絕選賞36 杜牧 泊秦淮,0B1tCMe1zn1gHWmllUEdVOFZEM2s
0B1tCMe1zn1gHMzVsa0VUQUdVLWs,唐詩七絕選賞35 杜牧 題烏江亭,0B1tCMe1zn1gHZGRnYjhwZllyb00
0B1tCMe1zn1gHYXpJVDhkZzdaZ2s,唐詩七絕選賞43 李商隱 嫦娥,0B1tCMe1zn1gHS2VrMzUwY0VLVlU
0B1tCMe1zn1gHOWpKdDZJNmdHdkE,唐詩七絕選賞45 鄭畋 馬嵬坡,0B1tCMe1zn1gHWHFEQm96djZITlU
0B1tCMe1zn1gHUGdkcF9mV05Kazg,唐詩七絕選賞46 司空圖 河湟有感,0B1tCMe1zn1gHRzY4VTAxWHZUenM
0B1tCMe1zn1gHTUtzUks1eEJmWEE,唐詩七絕選賞49 羅隱 蜂,0B1tCMe1zn1gHbk1xN2JGYV9pTjg
0B1tCMe1zn1gHaUhUN0RaYXJkRVk,唐詩七絕選賞52 無名氏 雜詩,0B1tCMe1zn1gHSXVaY2E0TFBDU3c
0B1tCMe1zn1gHZGhRR1pfS3h6WEU,唐詩七絕選賞51 張泌 寄人,0B1tCMe1zn1gHY1lSR2x2UXR6Qjg
0B1tCMe1zn1gHdDBTNkhzX1lJNzg,唐詩七絕選賞54 花蕊夫人 答宋君,0B1tCMe1zn1gHcy03em9oS3BJSnc
0B1tCMe1zn1gHUVRsRkJSbHpXcDg,唐詩七絕選賞53 杜秋娘 金縷衣,0B1tCMe1zn1gHWVllWWlyRVlkc2M
0B1tCMe1zn1gHWDVHRmotaGZMaUk,唐詩七絕選賞41 薛濤 籌邊樓,0B1tCMe1zn1gHaERIZWdGbXdNYzA
0B1tCMe1zn1gHM0xEc05aeEhJSlk,唐詩七絕選賞44 陳陶 隴西行,0B1tCMe1zn1gHcU9mMDhfbjA1UTQ
0B1tCMe1zn1gHOGxiMVJMRDZfX1k,唐詩七絕選賞50 羅隱 偶題,0B1tCMe1zn1gHS3ZjV3JZZWh1OTQ
0B1tCMe1zn1gHZXU2ZVVhaFlJbkU,唐詩七絕選賞48 韓偓 已涼,0B1tCMe1zn1gHNnppR1RjbGxIVmc
0B1tCMe1zn1gHTnhJM1lSbmxfYWM,唐詩七絕選賞47 韋莊 臺城,0B1tCMe1zn1gHMlBYNnl3eDJxdzA
0B1tCMe1zn1gHMzM5TG5GN25PWEk,唐詩七絕選賞42 李商隱 夜雨寄北,0B1tCMe1zn1gHTThvT1ZLcGstU2M`

func ExportFeedForYCantonese() []byte {
	now := time.Now()

	// instantiate a new Podcast
	p := podcast.New(
		"粵講越有趣 - 推廣粵文化 粵講越有趣",
		"http://ycantonese.org/",
		"陳耀南教授主講",
		&now, &now,
	)

	// add some channel properties
	p.AddAuthor("Radio Television Hong Kong", "webmaster@rthk.hk")

	p.AddSummary("本網站是為推廣粵文化而設。粵文化流傳的地區包括南中國、港、澳及海外華人僑居地。我們提供的視頻有風物誌及唐詩等。粵語是由中原南傳的唐代古語─雅言, 故用粵語讀唐詩及古文最傳神。敬希指導。")
	p.IExplicit = "no"

	for _, line := range strings.Split(ycanList, "\n") {
		part := strings.Split(line, ",")
		//fmt.Println("line ", part[0])

		// create an Item
		url := fmt.Sprintf("https://drive.google.com/uc?id=%s&authuser=0&export=download", part[0])
		item := podcast.Item{
			Title:       part[1],
			Link:        url,
			Description: part[1],
			PubDate:     &now,
		}

		item.AddEnclosure(url, podcast.MP3, 1000)
		if _, err := p.AddItem(item); err != nil {
			log.Fatal(item.Title, ": error", err.Error())
		}
	}

	return p.Bytes()
}
