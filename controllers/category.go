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

type CategoryController struct {
	BaseController
}

func (this *CategoryController) Get() {
	this.TplName = "home.html"
	this.ListTopic()
}

func (this *CategoryController) ListTopic() {
	cat := this.Ctx.Input.Param(":cat")
	this.Leftbar(cat)
	groupT := beego.BeeTemplates["groupTemplate.html"]
	var buff bytes.Buffer
	var Map = make(map[string]interface{})
	category := models.Blogger.GetCategoryByID(cat)
	var name string = "暂无该分类"
	if category != nil && len(category.Node.Children) > 0 {
		name = category.Node.Children[0].Text
	}
	Map["Name"] = name
	Map["Url"] = fmt.Sprintf("%s/cat/%s", this.domain, category.ID)
	pageStr := this.Ctx.Input.Param(":page")
	page := 1
	if temp, err := strconv.Atoi(pageStr); err == nil {
		page = temp
	}
	topics, remainpage := models.TMgr.GetTopicsByCatgory(cat, page)
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
			Map["UrlOlder"] = this.domain + "/cat/" + cat + fmt.Sprintf("/p/%d", page-1)
		}
		if remainpage == 0 {
			Map["ClassNewer"] = "disabled"
			Map["UrlNewer"] = "#"
		} else {
			Map["ClassNewer"] = ""
			Map["UrlNewer"] = this.domain + "/cat/" + cat + fmt.Sprintf("/p/%d", page+1)
			log.Debugf("%s", Map["UrlNewer"])
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
	groupT.Execute(&buff, Map)
	this.Data["Title"] = name + " - " + models.Blogger.BlogName
	this.Data["Content"] = buff.String()
}
