package main

import (
	"dzs_notice_sync/notice"
	"dzs_notice_sync/search"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Debug("run...")

	notice.InitDB()
	defer notice.DBClose()

	search.InitEs()
	defer search.CloseES()

	go notice.SaveData()

	g := gin.Default()
	g.StaticFile("/favicon.ico", "./resource/favicon.ico")
	g.GET("/", func(c *gin.Context) {
		c.String(200, "welcome use search")
	})
	g.GET("/search", search.Search)

	g.GET("/import2ES", search.Import2ES)
	g.GET("/crawling", notice.Crawling)

	g.GET("/colly", notice.Colly)
	g.Run(":80")
}
