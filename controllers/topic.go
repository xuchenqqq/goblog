package controllers

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/deepzz/beego_goblog/helper"
	"github.com/deepzz/beego_goblog/models"
	// "github.com/deepzz/com/log"
)

const (
	DUOSHUO_COMMENT_NUM_URL = "http://api.duoshuo.com/threads/counts.json?short_name=%s&threads=%d"
)

type TopicController struct {
	BaseController
}

func (this *TopicController) Get() {
	this.TplName = "home.html"
	this.Leftbar("")
	this.Topic()
}

func (this *TopicController) Post() {
	resp := helper.NewResponse()
	resp.Data = "文章索引错误."
	id := this.Ctx.Input.Param(":id")
	ID, err := strconv.Atoi(id)
	if err == nil {
		topic := models.TMgr.GetTopic(int32(ID))
		if topic != nil {
			resp.Data = string(topic.Content)
		}
	}
	resp.WriteJson(this.Ctx.ResponseWriter)
}

func (this *TopicController) Topic() {
	id := this.Ctx.Input.Param(":id")
	ID, err := strconv.Atoi(id)
	if err != nil {
		this.Data["Content"] = ""
		return
	}
	Map := make(map[string]string)
	topic := models.TMgr.GetTopic(int32(ID))
	if topic == nil {
		this.Data["Content"] = "文章索引错误."
		return
	}
	this.Data["Title"] = topic.Title + " - " + models.Blogger.BlogName
	Map["Url"] = fmt.Sprintf("%s/%s/%d.html", this.domain, topic.CreateTime.Format(helper.Layout_y_m_d), topic.ID)
	Map["Title"] = topic.Title
	Map["Time"] = topic.CreateTime.Format(helper.Layout_y_m_d2)
	Map["Category"] = "<a " + topic.PCategory.Node.Children[0].Extra + "' rel='category tag'>" + topic.PCategory.Node.Children[0].Text + "</a>"
	Map["Tags"] = ""
	for i, tag := range topic.PTags {
		if i == 0 {
			Map["Tags"] += "<a " + tag.Node.Extra + " rel='tag'>" + tag.Node.Text + "</a>"
		} else {
			Map["Tags"] += ",<a " + tag.Node.Extra + " rel='tag'>" + tag.Node.Text + "</a>"
		}
	}
	Map["ID"] = fmt.Sprint(topic.ID)
	topicT := beego.BeeTemplates["topicTemplate.html"]
	var buff bytes.Buffer
	topicT.Execute(&buff, Map)
	this.Data["Content"] = buff.String()

}
