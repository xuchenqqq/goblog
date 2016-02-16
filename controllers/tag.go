package controllers

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/smalltree0/beego_goblog/models"
	// "github.com/smalltree0/com/log"
	"github.com/smalltree0/beego_goblog/helper"
)

type TagController struct {
	BaseController
}

func (this *TagController) Get() {
	this.TplName = "home.html"
	this.Leftbar("")
	this.ListTopic()
}

func (this *TagController) ListTopic() {
	tagName := this.Ctx.Input.Param(":tag")
	groupT := beego.BeeTemplates["groupTemplate.html"]
	var buff bytes.Buffer
	Map := make(map[string]interface{})
	tag := models.Blogger.GetTagByID(tagName)
	Map["Name"] = "无效TAG"
	if tag != nil {
		Map["Name"] = tag.ID
		page := 1
		tagName := this.Ctx.Input.Param(":tag")
		pageStr := this.Ctx.Input.Param(":page")
		if temp, err := strconv.Atoi(pageStr); err == nil {
			page = temp
		}
		topics, remainpage := models.TMgr.GetTopicsByTag(tagName, page)
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
				Map["UrlOlder"] = this.domain + "/tag/" + tagName + fmt.Sprintf("/p/%d", page-1)
			}
			if remainpage == 0 {
				Map["ClassNewer"] = "disabled"
				Map["UrlNewer"] = "#"
			} else {
				Map["ClassNewer"] = ""
				Map["UrlNewer"] = this.domain + "/tag/" + tagName + fmt.Sprintf("/p/%d", page+1)
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
	}
	groupT.Execute(&buff, Map)
	this.Data["Title"] = tagName + " - " + models.Blogger.BlogName
	this.Data["Content"] = fmt.Sprintf("%s", buff.Bytes())
}
