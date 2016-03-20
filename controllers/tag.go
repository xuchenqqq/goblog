package controllers

import (
	"fmt"
	"strconv"

	"github.com/deepzz0/goblog/models"
	// "github.com/deepzz0/go-common/log"
	"github.com/deepzz0/goblog/helper"
)

type TagController struct {
	Common
}

func (this *TagController) Get() {
	this.Layout = "homelayout.html"
	this.TplName = "groupTemplate.html"
	this.Leftbar("")
	this.ListTopic()
}

func (this *TagController) ListTopic() {
	tagName := this.Ctx.Input.Param(":tag")
	tag := models.Blogger.GetTagByID(tagName)
	this.Data["Name"] = "无效TAG"
	if tag != nil {
		this.Data["Name"] = tag.ID
		page := 1
		tagName := this.Ctx.Input.Param(":tag")
		pageStr := this.Ctx.Input.Param(":page")
		if temp, err := strconv.Atoi(pageStr); err == nil {
			page = temp
		}
		topics, remainpage := models.TMgr.GetTopicsByTag(tagName, page)
		if remainpage == -1 {
			this.Data["ClassOlder"] = "disabled"
			this.Data["UrlOlder"] = "#"
			this.Data["ClassNewer"] = "disabled"
			this.Data["UrlNewer"] = "#"
		} else {
			if page == 1 {
				this.Data["ClassOlder"] = "disabled"
				this.Data["UrlOlder"] = "#"
			} else {
				this.Data["ClassOlder"] = ""
				this.Data["UrlOlder"] = this.domain + "/tag/" + tagName + fmt.Sprintf("/p/%d", page-1)
			}
			if remainpage == 0 {
				this.Data["ClassNewer"] = "disabled"
				this.Data["UrlNewer"] = "#"
			} else {
				this.Data["ClassNewer"] = ""
				this.Data["UrlNewer"] = this.domain + "/tag/" + tagName + fmt.Sprintf("/p/%d", page+1)
			}
			var ts []*listOfTopic
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
		}
	}
	this.Data["Title"] = tagName + " - " + models.Blogger.BlogName
}
