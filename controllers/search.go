package controllers

import (
	"fmt"

	"github.com/deepzz0/goblog/helper"
	"github.com/deepzz0/goblog/models"
)

type SearchController struct {
	Common
}

func (this *SearchController) Get() {
	this.Layout = "homelayout.html"
	this.TplName = "groupTemplate.html"
	this.Leftbar("")
	this.ListTopic()
}

func (this *SearchController) ListTopic() {
	search := this.GetString("title")
	this.Data["ClassOlder"] = "disabled"
	this.Data["UrlOlder"] = "#"
	this.Data["ClassNewer"] = "disabled"
	this.Data["UrlNewer"] = "#"
	this.Data["Name"] = "搜索：" + search
	this.Data["URL"] = fmt.Sprintf("%s/search?title=%s", this.domain, search)
	var ts []*listOfTopic
	topics := models.TMgr.GetTopicsSearch(search)
	for _, topic := range topics {
		t := &listOfTopic{}
		t.ID = topic.ID
		t.Time = topic.CreateTime.Format(helper.Layout_y_m_d2)
		t.URL = fmt.Sprintf("%s/%s/%d.html", this.domain, topic.CreateTime.Format(helper.Layout_y_m_d), topic.ID)
		t.Title = topic.Title
		t.Preview = topic.Preview
		t.PCategory = topic.PCategory
		t.PTags = topic.PTags
		ts = append(ts, t)
	}
	this.Data["ListTopics"] = ts
	this.Data["Title"] = fmt.Sprintf("搜索: %s - %s", search, models.Blogger.BlogName)
	this.Data["Description"] = fmt.Sprintf("搜索标题,%s,%s,blog", models.Blogger.Introduce, models.Blogger.UserName)
	this.Data["Keywords"] = fmt.Sprintf("搜索标题,find,%s,%s", models.Blogger.Introduce, models.Blogger.UserName)
}
