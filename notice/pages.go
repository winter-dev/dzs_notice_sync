package notice

import (
	"github.com/gocolly/colly"
)

type Pages struct {
	id  int64  `json:"id"`
	url string `json:"url"`
	idx int64  `idx`
}

var basicUrl string = "https://dzs.qq.com"
var firstPage string = "https://dzs.qq.com/webplat/info/news_version3/394/3871/3872/3874/m2934/list_1.shtml"

func Crawling() {
	var sli []Pages
	var finished bool = false
	var idx int64 = 2
	var isFirst bool = true
	var currUrl string

	//初始化第一页
	var p Pages = Pages{0, firstPage, 1}
	sli = append(sli, p)
	for {
		c := colly.NewCollector()
		c.DetectCharset = true

		c.OnHTML("#page", func(e *colly.HTMLElement) {
			hrefs := e.ChildAttrs("a[class='page_next p2']", "href")
			if len(hrefs) == 0 || hrefs[0] == "" {
				finished = true
			}
			if len(hrefs) > 0 {
				href := basicUrl + hrefs[0]
				var page Pages
				page.idx = idx
				page.url = href
				sli = append(sli, page)
				currUrl = href
			}
		})
		if isFirst {
			c.Visit(firstPage)
		} else {
			c.Visit(currUrl)
		}
		isFirst = false

		if finished {
			break
		}
		idx++
	}
	//插入db
	insertPage(sli)
}
