package controllers

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/deepzz0/go-common/log"
	"github.com/deepzz0/goblog/helper"
	"github.com/deepzz0/goblog/models"
)

type HomeController struct {
	BaseController
}

func (this *HomeController) Get() {
	this.TplName = "home.html"
	this.Data["Title"] = fmt.Sprintf("%s - %s", models.Blogger.Introduce, models.Blogger.BlogName)
	this.Leftbar("homepage")
	this.Home()
}
func (this *HomeController) Home() {
	homeT := beego.BeeTemplates["homeTemplate.html"]
	var buff bytes.Buffer
	var html string
	Map := make(map[string]interface{})
	for _, tag := range models.Blogger.Tags {
		html += tag.Node.String()
	}
	Map["TagCloud"] = html
	html = ""
	for _, br := range models.Blogger.Blogrolls {
		html += br.Node.String()
	}
	Map["Blogroll"] = html
	// 文章列表
	page := 1
	pageStr := this.Ctx.Input.Param(":page")
	if temp, err := strconv.Atoi(pageStr); err == nil {
		page = temp
	}
	topics, remainpage := models.TMgr.GetTopicsByPage(page)
	log.Debugf("page = %d，remainpage=%d	", page, remainpage)
	if remainpage == -1 {
		Map["ClassOlder"] = "disabled"
		Map["UrlOlder"] = "#"
		Map["ClassNewer"] = "disabled"
		Map["UrlNewer"] = "#"
	} else {
		if page == 1 {
			Map["ClassOlder"] = "disabled"
			Map["UrlOlder"] = "#"
		} else {
			Map["ClassOlder"] = ""
			Map["UrlOlder"] = this.domain + "/p/" + fmt.Sprint(page-1)
		}
		if remainpage == 0 {
			Map["ClassNewer"] = "disabled"
			Map["UrlNewer"] = "#"
		} else {
			Map["ClassNewer"] = ""
			Map["UrlNewer"] = this.domain + "/p/" + fmt.Sprint(page+1)
		}
		var ts []*listOfTopic
		for _, topic := range topics {
			t := &listOfTopic{}
			t.ID = topic.ID
			t.Time = topic.CreateTime.Format(helper.Layout_y_m_d2)
			t.Url = fmt.Sprintf("%s/%s/%d.html", this.domain, topic.CreateTime.Format(helper.Layout_y_m_d), topic.ID)
			t.Title = topic.Title
			if len(topic.Content) < 300 {
				t.Preview = string(topic.Content)
			} else {
				t.Preview = string(topic.Content[:300])
			}
			t.Category = "<a " + topic.PCategory.Node.Children[0].Extra + " rel='category tag'>" + topic.PCategory.Node.Children[0].Text + "</a>"
			for i, tag := range topic.PTags {
				if i == 0 {
					t.Tags += "<a " + tag.Node.Extra + " rel='tag'>" + tag.Node.Text + "</a>"
				} else {
					t.Tags += ",<a " + tag.Node.Extra + " rel='tag'>" + tag.Node.Text + "</a>"
				}
			}
			ts = append(ts, t)
		}
		Map["ListTopics"] = ts
	}
	homeT.Execute(&buff, Map)
	this.Data["Content"] = buff.String()
}
