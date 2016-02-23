package background

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/astaxie/beego"
	"github.com/deepzz0/go-common/log"
	"github.com/deepzz0/goblog/RS"
	"github.com/deepzz0/goblog/helper"
	"github.com/deepzz0/goblog/models"
)

type BlogrollController struct {
	BackgroundController
}

func (this *BlogrollController) Get() {
	this.TplName = "manage/adminTemplate.html"
	this.Data["Title"] = "友情链接 - " + models.Blogger.BlogName
	this.LeftBar("blogroll")
	this.Content()
}

func (this *BlogrollController) Content() {
	blogrollT := beego.BeeTemplates["manage/blogroll/blogrollTemplate.html"]
	var buffer bytes.Buffer
	blogrollT.Execute(&buffer, "")
	this.Data["Content"] = fmt.Sprintf("%s", string(buffer.Bytes()))
}

func (this *BlogrollController) Post() {
	resp := helper.NewResponse()
	flag := this.GetString("flag")
	switch flag {
	case "blogroll":
		this.getBlogrolls(resp)
	case "addblogroll":
		this.addBlogroll(resp)
	case "modify":
		this.getBlogroll(resp)
	case "domodify":
		this.doModify(resp)
	case "deleteblogroll":
		this.doDelete(resp)
	default:
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "参数错误|未知的flag标志。"}
	}
	resp.WriteJson(this.Ctx.ResponseWriter)
}

type blogroll struct {
	ID     string
	SortID int
	Extra  string
	Name   string
	Time   string
}

func (this *BlogrollController) getBlogrolls(resp *helper.Response) {
	var blogrolls []*blogroll
	for _, br := range models.Blogger.Blogrolls {
		temp := &blogroll{}
		temp.ID = br.ID
		temp.SortID = br.SortID
		temp.Extra = br.Node.Children[0].Extra
		temp.Name = br.Node.Children[0].Text
		temp.Time = br.CreateTime.Format(helper.Layout_y_m_d_time)
		blogrolls = append(blogrolls, temp)
	}
	Map := map[string]interface{}{"Blogrolls": blogrolls}
	blogrollsT := beego.BeeTemplates["manage/blogroll/blogrolls.html"]
	var buffer bytes.Buffer
	blogrollsT.Execute(&buffer, Map)
	resp.Data = buffer.String()
}
func (this *BlogrollController) addBlogroll(resp *helper.Response) {
	content := this.GetString("json")
	blogroll := models.NewBlogroll()
	err := json.Unmarshal([]byte(content), &blogroll)
	if err != nil {
		log.Error(err)
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "内容错误|要仔细检查哦。"}
		return
	}
	if models.Blogger.AddBlogroll(blogroll) != RS.RS_success {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "分类已存在|你确定要添加分类吗？"}
		return
	}
}
func (this *BlogrollController) getBlogroll(resp *helper.Response) {
	id := this.GetString("id")
	if id != "" {
		modifyBlogrollT := beego.BeeTemplates["manage/blogroll/modifyblogroll.html"]
		var buffer bytes.Buffer
		if br := models.Blogger.GetBlogrollByID(id); br != nil {
			b, _ := json.Marshal(br)
			modifyBlogrollT.Execute(&buffer, map[string]string{"Content": string(b)})
			resp.Data = buffer.String()
			return
		}
	}
	resp.Data = ""
}
func (this *BlogrollController) doModify(resp *helper.Response) {
	content := this.GetString("json")
	if content != "" {
		br := models.Blogroll{}
		err := json.Unmarshal([]byte(content), &br)
		if err != nil {
			resp.Status = RS.RS_failed
			resp.Err = helper.Error{Level: helper.WARNING, Msg: "内容错误|反序列化失败。"}
			return
		}
		blogroll := models.Blogger.GetBlogrollByID(br.ID)
		if blogroll != nil {
			*blogroll = br
			sort.Sort(models.Blogger.Blogrolls)
		} else {
			resp.Status = RS.RS_failed
			resp.Err = helper.Error{Level: helper.WARNING, Msg: "索引错误|没有找到该友情链接。"}
			return
		}
	} else {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "获取内容失败|你都干了什么。"}
		return
	}
	resp.Success()
}
func (this *BlogrollController) doDelete(resp *helper.Response) {
	id := this.GetString("id")
	if id == "" {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "哦噢。。。|参数错误。"}
		return
	}
	if code := models.Blogger.DelBlogrollByID(id); code != RS.RS_success {
		resp.Status = code
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "抱歉|系统没有找到该友情链接。"}
		return
	}
}
