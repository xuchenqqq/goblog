package background

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/smalltree0/beego_goblog/RS"
	"github.com/smalltree0/beego_goblog/helper"
	"github.com/smalltree0/beego_goblog/models"
	"github.com/smalltree0/com/log"
)

type TopicsController struct {
	BackgroundController
}

func (this *TopicsController) Get() {
	this.TplName = "manage/adminTemplate.html"
	this.Data["Title"] = "博文管理 - " + models.Blogger.BlogName
	this.LeftBar("topics")
	this.Content()
}

func (this *TopicsController) Content() {
	topicT := beego.BeeTemplates["manage/topic/topicTemplate.html"]
	var buffer bytes.Buffer
	Map := make(map[string]string)
	var html string
	for _, cat := range models.Blogger.Categories {
		if cat.IsCategory {
			html += "<option value='" + cat.ID + "'>" + cat.ID + "</option>"
		}
	}
	Map["Category"] = html
	topicT.Execute(&buffer, Map)
	this.Data["Content"] = fmt.Sprintf("%s", string(buffer.Bytes()))
}

func (this *TopicsController) Post() {
	resp := helper.NewResponse()
	flag := this.GetString("flag")
	log.Debug(flag)
	switch flag {
	case "topics":
		this.getTopics(resp)
	case "addtopic":
		this.addTopic(resp)
	case "modify":
		this.getTopic(resp)
	case "domodify":
		this.doModify(resp)
	case "deletetopic":
		this.doDeleteTopic(resp)
	case "category":
		this.category(resp)
	case "tag":
		this.tag(resp)
	case "deletetopics":
		this.doDeleteTopics(resp)
	case "editor":
		this.editor(resp)
	default:
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "参数错误|未知的flag标志。"}
	}
	resp.WriteJson(this.Ctx.ResponseWriter)
}

type topic struct {
	ID         int32
	Title      string
	TagIDs     []string
	CategoryID string
	Author     string
	CreateTime string
	EditTime   string
}

func (this *TopicsController) getTopics(resp *helper.Response) {
	cat := this.GetString("cat")
	page, err := this.GetInt("page")
	if err != nil {
		page = 1
	}
	var pageTopics []*models.Topic
	var remainpage int
	if cat == "" || cat == "0" {
		pageTopics, remainpage = models.TMgr.GetTopicsByPage(page)
		log.Debugf("%s,%s", page, remainpage)
		if remainpage == -1 {
			resp.Data = "<script>$('#previous').addClass('disabled');$('#next').addClass('disabled');pushMessage('info','莫得数据得嘛|是不是页码错了。')</script>"
			return
		}
	} else {
		pageTopics, remainpage = models.TMgr.GetTopicsByCatgory(cat, page)
		if remainpage == -1 {
			resp.Data = "<script>$('#previous').addClass('disabled');$('#next').addClass('disabled');pushMessage('info','莫得数据哦|我看是分类或者页码错误。')</script>"
			return
		}
	}
	topicsT := beego.BeeTemplates["manage/topic/topics.html"]
	var buffer bytes.Buffer
	var topics []*topic
	Map := make(map[string]interface{})
	if page == 1 {
		Map["IsFirstPage"] = true
	} else {
		Map["IsFirstPage"] = false
		Map["CategoryID"] = cat
		Map["PrePage"] = page - 1
	}
	if remainpage == 0 {
		Map["IsLastPage"] = true
	} else {
		Map["IsLastPage"] = false
		Map["CategoryID"] = cat
		Map["NextPage"] = page + 1
	}
	for _, t := range pageTopics {
		temp := &topic{}
		temp.ID = t.ID
		temp.Title = t.Title
		temp.TagIDs = t.TagIDs
		temp.Author = t.Author
		temp.CategoryID = t.CategoryID
		temp.CreateTime = t.CreateTime.Format(helper.Layout_y_m_d_time)
		temp.EditTime = t.EditTime.Format(helper.Layout_y_m_d_time)
		topics = append(topics, temp)
	}
	Map["Topics"] = topics
	topicsT.Execute(&buffer, Map)
	resp.Data = buffer.String()
}
func (this *TopicsController) addTopic(resp *helper.Response) {
	title := this.GetString("title")
	content := this.GetString("content")
	cat := this.GetString("cat")
	tags := this.GetString("tags")
	log.Debugf("%s,%s,%s,%s", title, content, cat, tags)
	if title == "" || content == "" {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "错误|请检查是否填写完整。"}
		return
	}
	topic := models.NewTopic()
	topic.Title = title
	topic.Content = []rune(content)
	if category := models.Blogger.GetCategoryByID(cat); category == nil {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "错误|查找不到该分类。"}
		return
	} else {
		topic.CategoryID = category.ID
	}
	sliceTags := strings.Split(tags, ",")
	for _, tag := range sliceTags {
		if tag != "" {
			topic.TagIDs = append(topic.TagIDs, tag)
		}
	}
	if err := models.TMgr.AddTopic(topic, this.domain); err != nil {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "错误|" + err.Error()}
		return
	}
	resp.Success()
}

type modifyTopic struct {
	Title    string
	Content  string
	Category string
	Tags     []string
}

func (this *TopicsController) getTopic(resp *helper.Response) {
	id, err := this.GetInt("id")
	if err != nil {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "错误|ID格式不正确。"}
		return
	}
	if topic := models.TMgr.GetTopic(int32(id)); topic == nil {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "错误|系统未查询到该文章。"}
		return
	} else {
		mt := &modifyTopic{}
		mt.Title = topic.Title
		mt.Content = string(topic.Content)
		mt.Tags = topic.TagIDs
		mt.Category = topic.CategoryID
		resp.Data = mt
	}
}
func (this *TopicsController) doModify(resp *helper.Response) {
	id, err := this.GetInt32("id")
	if err != nil {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "错误|ID格式不正确。"}
		return
	}
	title := this.GetString("title")
	content := this.GetString("content")
	categoryID := this.GetString("cat")
	tags := this.GetString("tags")
	if title == "" || content == "" || categoryID == "" || tags == "" {
		resp.Status = RS.RS_failed
		resp.Tips(helper.WARNING, RS.RS_params_error)
		return
	}
	if t := models.TMgr.GetTopic(id); t == nil {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "错误|系统查找不到该文章ID。"}
		return
	} else {
		t.Title = title
		t.Content = []rune(content)
		t.CategoryID = categoryID
		t.TagIDs = strings.Split(tags, ",")
	}
	resp.Success()
}

type selectCat struct {
	ID         string
	IsSelected bool
}

func (this *TopicsController) category(resp *helper.Response) {
	categoryT := beego.BeeTemplates["manage/topic/category.html"]
	var buffer bytes.Buffer
	var selectCats []*selectCat
	selected := this.GetString("selected")
	for _, cat := range models.Blogger.Categories {
		if cat.IsCategory {
			temp := &selectCat{}
			temp.ID = cat.ID
			if cat.ID != "" && cat.ID == selected {
				temp.IsSelected = true
			}
			selectCats = append(selectCats, temp)
		}
	}
	err := categoryT.Execute(&buffer, map[string][]*selectCat{"Categories": selectCats})
	log.Error(err)
	resp.Data = buffer.String()
}
func (this *TopicsController) tag(resp *helper.Response) {
	var html string
	for _, tag := range models.Blogger.Tags {
		html += tag.TagStyle()
	}
	tagsT := beego.BeeTemplates["manage/topic/tags.html"]
	var buffer bytes.Buffer
	tagsT.Execute(&buffer, map[string]string{"Content": html})
	resp.Data = buffer.String()
}
func (this *TopicsController) doDeleteTopic(resp *helper.Response) {
	id, err := this.GetInt("id")
	if err != nil {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "ID错误|走正常途径哦。"}
		return
	}
	err = models.TMgr.DelTopic(int32(id))
	if err != nil {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "删除失败|" + err.Error()}
		return
	}
}
func (this *TopicsController) doDeleteTopics(resp *helper.Response) {
	ids := this.GetString("ids")
	log.Debugf("%s", ids)
	if ids == "" {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "ID错误|走正常途径哦。"}
		return
	}
	sliceID := strings.Split(ids, ",")
	for _, v := range sliceID {
		id, err := strconv.Atoi(v)
		if err != nil {
			log.Error(err)
			resp.Status = RS.RS_failed
			resp.Err = helper.Error{Level: helper.WARNING, Msg: "ID错误|走正常途径哦。"}
			return
		}
		err = models.TMgr.DelTopic(int32(id))
		if err != nil {
			resp.Status = RS.RS_failed
			resp.Err = helper.Error{Level: helper.WARNING, Msg: "删除失败|" + err.Error()}
			return
		}
	}
}
func (this *TopicsController) editor(resp *helper.Response) {
	editorT := beego.BeeTemplates["manage/editor.html"]
	var buffer bytes.Buffer
	editorT.Execute(&buffer, "")
	resp.Data = fmt.Sprintf("%s", string(buffer.Bytes()))
}
