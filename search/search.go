package search

import (
	"context"
	"dzs_notice_sync/notice"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

var client *elastic.Client

var host = "http://127.0.0.1:9200/"

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

func InitEs() {
	var e error
	client, e = elastic.NewClient(elastic.SetURL(host))
	if e != nil {
		panic(e)
	}
	r, i, e := client.Ping(host).Do(context.Background())
	if e != nil {
		panic(e)
	}
	log.Info("ping elasticsearch code:", i, " , elasticsearch version:", r.Version.Number, " ,lucene version:", r.Version.LuceneVersion)
}

func CloseES() {
	if client != nil {
		client.Stop()
	}
}

// 导入数据到es
func Import2ES(c *gin.Context) {
	ns, e := notice.GetNotice()
	if e != nil {
		log.Error(e)
		return
	}

	bulkRequest := client.Bulk()
	for _, v := range ns {
		esReq := elastic.NewBulkIndexRequest().Index("dzs_notice").Type("_doc").Id(fmt.Sprintf("%d", int64(v.Id))).Doc(v)
		bulkRequest = bulkRequest.Add(esReq)
	}
	_, e = bulkRequest.Do(context.Background())
	if e != nil {
		log.Error("add index error,", e)
	}
	log.Info("add index success!")
}

func Search(c *gin.Context) {
	kw := c.Query("q")

	query := elastic.NewBoolQuery()
	query = query.Must(elastic.NewMatchQuery("title", kw))
	//query = query.Must(elastic.NewMatchQuery("title", kw))

	resp, e := client.Search().Index("dzs_notice").Type("_doc").Query(query).
		Sort("id", true).From(0).Size(100).Do(context.Background())
	if e != nil {
		log.Error("search occured an error,", e)
		c.JSON(200, e.Error())
		return
	}

	var docs []Notice
	for _, v := range resp.Hits.Hits {
		var doc Notice
		e = json.Unmarshal(*v.Source, &doc)
		if e != nil {
			log.Error("Unmarshal occured an error ,", e)
			continue
		}
		docs = append(docs, doc)
	}
	c.JSON(200, docs)
}
