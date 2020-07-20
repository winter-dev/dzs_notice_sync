package notice

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"time"
)

var db *sql.DB

func InitDB() {
	log.Info("init database start...")
	var e error
	db, e = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/notice?charset=utf8&loc=Local")
	if e != nil {
		log.Error("open database error,", e)
		panic(e)
	}

	db.SetConnMaxLifetime(100 * time.Second) //最大连接周期，超过时间的连接就close
	db.SetMaxOpenConns(20)                   //设置最大连接数
	db.SetMaxIdleConns(5)                    //设置闲置连接数

	e = db.Ping()
	if e != nil {
		log.Error("ping DB error", e)
		panic(e)
	}

	r, e := db.Query("SELECT VERSION()")
	if e != nil {
		panic(e)
	}
	defer r.Close()
	var version string
	for r.Next() {
		e = r.Scan(&version)
		if e != nil {
			panic(e)
		}
		log.Info("current Mysql version:", version)
	}

	log.Info("connect database success!")
}

func DBClose() {
	if db != nil {
		_ = db.Close()
	}
}

func insertNotice(notice Notice) {
	stmt, e := db.Prepare("INSERT INTO t_notice(title,subtitle,author,context,url,release_time,create_time) VALUES (?,?,?,?,?,?,?)")
	if e != nil {
		log.Error(e)
		panic(e)
	}

	ret, e := stmt.Exec(notice.Title, notice.Subtitle, notice.Author, notice.Context, notice.Url, notice.ReleaseTime, notice.CreateTime)
	if e != nil {
		log.Error(e)
		panic(e)
	}
	lastId, _ := ret.LastInsertId()
	refCnt, _ := ret.RowsAffected()
	log.Info("this insert insert count:%d,lastId:%d", refCnt, lastId)
}

func insertPage(sli []Pages) {
	log.Info(sli)
	var count int64 = 0
	for _, page := range sli {
		stmt, e := db.Prepare("INSERT INTO t_pages(url,idx) VALUES (?,?)")
		if e != nil {
			log.Error(e)
			panic(e)
		}

		ret, e := stmt.Exec(page.url, page.idx)
		if e != nil {
			log.Error(e)
			panic(e)
		}
		tc, _ := ret.RowsAffected()
		count = count + tc
	}
	log.Info("this insert count:", count)
}

func GetNotice() ([]Notice, error) {
	sql := "SELECT id,title,subtitle,author,context,url,release_time,create_time FROM t_notice ORDER BY id"

	r, e := db.Query(sql)
	if e != nil {
		log.Error(e)
		return nil, e
	}
	var ns []Notice
	for r.Next() {
		var notice Notice
		_ = r.Scan(&notice.Id, &notice.Title, &notice.Subtitle, &notice.Author, &notice.Context, &notice.Url, &notice.ReleaseTime, &notice.CreateTime)
		ns = append(ns, notice)
	}
	return ns, nil
}
