package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/deepzz0/goblog/helper"
	"github.com/deepzz0/goblog/models"
	// "github.com/deepzz0/go-common/log"
)

const (
	DUOSHUO_COMMENT_NUM_URL = "http://api.duoshuo.com/threads/counts.json?short_name=%s&threads=%d"
)

type TopicController struct {
	Common
}

func (this *TopicController) Get() {
	this.Layout = "homelayout.html"
	this.TplName = "topicTemplate.html"
	this.Leftbar("")
	this.Topic()
}

func (this *TopicController) Post() {
	resp := helper.NewResponse()
	defer resp.WriteJson(this.Ctx.ResponseWriter)
	resp.Data = "文章索引错误."
	id := this.Ctx.Input.Param(":id")
	ID, err := strconv.Atoi(id)
	if err == nil {
		topic := models.TMgr.GetTopic(int32(ID))
		if topic != nil {
			resp.Data = string(topic.Content)
		}
	}
}

func (this *TopicController) Topic() {
	id := this.Ctx.Input.Param(":id")
	ID, err := strconv.Atoi(id)
	if err != nil {
		this.Data["IsFalse"] = true
		return
	}
	topic := models.TMgr.GetTopic(int32(ID))
	if topic == nil {
		this.Data["IsFalse"] = true
		return
	}
	this.Data["IsFalse"] = false
	this.Data["Title"] = topic.Title + " - " + models.Blogger.BlogName
	this.Data["Topic"] = topic
	this.Data["Domain"] = this.domain
	this.Data["Description"] = fmt.Sprintf("%s,%s,,%s,blog", topic.Title, models.Blogger.Introduce, models.Blogger.UserName)
	this.Data["KeyWords"] = fmt.Sprintf("%s,%s,%s,%s,%s", topic.Title, topic.CategoryID, strings.Join(topic.TagIDs, ","), models.Blogger.Introduce, models.Blogger.UserName)
}
