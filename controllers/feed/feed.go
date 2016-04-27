package feed

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/deepzz0/go-common/log"
	"github.com/deepzz0/goblog/helper"
	"github.com/deepzz0/goblog/models"
)

func init() {
	go scheduler()
}

func Feed(ctx *context.Context) {
	http.ServeFile(ctx.ResponseWriter, ctx.Request, feedFile)
}

func scheduler() {
	t := time.NewTicker(time.Hour)
	for {
		select {
		case <-t.C:
			doFeed()
		}
	}
}

const (
	version = "0.0.1"
	year    = "2016"

	templateFile = "./static/feedTemplate.xml"
	feedFile     = "./static/feed.xml"
)

type Topic struct {
	Title    string
	URL      string
	PubDate  string
	Author   string
	Category string
	Desc     string
}

var buildDate time.Time

func doFeed() {
	temp, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Error(err)
		return
	}
	domain := beego.AppConfig.String("mydomain")
	ts, _ := models.TMgr.GetTopicsByPage(1)
	var Topics []*Topic
	for i, v := range ts {
		if i == 0 && v.CreateTime == buildDate {
			return
		}
		t := &Topic{}
		t.Title = v.Title
		t.URL = fmt.Sprintf("%s/%s/%d.html", domain, v.CreateTime.Format(helper.Layout_y_m_d), v.ID)
		t.PubDate = v.CreateTime.Format(time.RFC1123Z)
		t.Author = v.Author
		t.Category = v.CategoryID
		t.Desc = v.Preview
		Topics = append(Topics, t)
	}
	buildDate = time.Now()
	params := make(map[string]interface{})
	params["Title"] = "Deepzz 的个人博客"
	params["Domain"] = domain
	params["Desc"] = "Golang爱好者，专注于游戏服务器开发，轻架构，高并发"
	params["PubDate"] = time.Now().Format(time.RFC1123Z)
	params["BuildDate"] = buildDate.Format(time.RFC1123Z)
	params["Year"] = year
	params["Version"] = version
	params["Author"] = "deepzz"
	params["Topics"] = Topics

	_, err = os.Stat(feedFile)
	if err != nil && !strings.Contains(err.Error(), "no such file") {
		log.Error(err)
		return
	} else {
		os.Remove(feedFile)
	}
	f, err := os.Create(feedFile)
	if err != nil {
		log.Error(err)
		return
	}
	defer f.Close()
	err = temp.Execute(f, params)
	if err != nil {
		log.Error(err)
		return
	}
}
