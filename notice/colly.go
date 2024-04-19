package notice

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"strings"
)

var dzsUrl string = "https://dzs.qq.com"

type Notice struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Author      string `json:"author"`
	Context     string `json:"context"`
	Url         string `json:"url"`
	ReleaseTime string `json:"releaseTime"`
	CreateTime  string `json:"createTime"`
}

var queue = make(chan Notice, 10000)

// 根据页面地址保存内容
func Colly(c *gin.Context) {
	r, e := db.Query("SELECT url FROM t_pages")
	if e != nil {
		log.Error(e)
		return
	}
	for r.Next() {
		var url string
		_ = r.Scan(&url)

		log.Info("爬取url:", url)
		c := colly.NewCollector()
		c.DetectCharset = true

		// Find and visit all links
		c.OnHTML(".list_con", func(e *colly.HTMLElement) {
			//下一页url
			e.DOM.Each(func(i int, selection *goquery.Selection) {
				selection.Find(".list_ul").Each(func(i int, s1 *goquery.Selection) {
					s1.Find("li").Each(func(i int, s2 *goquery.Selection) {
						var notice Notice

						createTime := s2.Find(".time").Text()
						createTime = strings.Replace(createTime, " ", "", -1)

						notice.CreateTime = createTime

						title := s2.Find("font ").Text()
						title = "[公告]" + title
						notice.Subtitle = title

						//获得跳转url
						href, _ := s2.Find("a:nth-child(2)").Attr("href")
						notice.Url = dzsUrl + href
						queue <- notice
					})
				})
			})
		})
		_ = c.Visit(url)
	}
}

func SaveData() {
	for {
		notice := <-queue
		if notice == (Notice{}) {
			continue
		}

		sub := colly.NewCollector()
		sub.DetectCharset = true

		sub.OnHTML("#detailCon", func(element *colly.HTMLElement) {
			context, e := element.DOM.Html()
			if e != nil {
				log.Error("get detail error,", e)
			}
			notice.Context = context
		})

		sub.OnHTML("#detailTitle", func(element *colly.HTMLElement) {
			notice.Title = element.DOM.Text()
		})

		sub.OnHTML(".detail_info", func(element *colly.HTMLElement) {
			top := element.DOM.Find("p:nth-child(1)").Text()

			a := strings.Split(top, "  ")
			notice.Author = strings.Split(a[0], "：")[1]
			sendTime := strings.Split(a[1], "：")[1]

			notice.ReleaseTime = sendTime
		})

		sub.OnHTML("#detailCon", func(element *colly.HTMLElement) {
			context, e := element.DOM.Html()
			if e != nil {
				log.Error("get detail error,", e)
			}
			notice.Context = context
		})
		//根据url获取详情
		_ = sub.Visit(notice.Url)

		insertNotice(notice)
	}

}
