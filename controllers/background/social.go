package background

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/astaxie/beego"
	"github.com/deepzz/beego_goblog/RS"
	"github.com/deepzz/beego_goblog/helper"
	"github.com/deepzz/beego_goblog/models"
	"github.com/deepzz/com/log"
)

type SocialController struct {
	BackgroundController
}

func (this *SocialController) Get() {
	this.TplName = "manage/adminTemplate.html"
	this.Data["Title"] = "社交工具 - " + models.Blogger.BlogName
	this.LeftBar("social")
	this.Content()
}
func (this *SocialController) Content() {
	socialT := beego.BeeTemplates["manage/social/socialTemplate.html"]
	var buffer bytes.Buffer
	socialT.Execute(&buffer, "")
	this.Data["Content"] = fmt.Sprintf("%s", string(buffer.Bytes()))
}

func (this *SocialController) Post() {
	resp := helper.NewResponse()

	flag := this.GetString("flag")
	switch flag {
	case "social":
		this.getSocials(resp)
	case "addsocial":
		this.addSocial(resp)
	case "modify":
		this.getSocial(resp)
	case "domodify":
		this.doModify(resp)
	case "deletesocial":
		this.doDelete(resp)
	default:
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "参数错误|未知的flag标志。"}
	}
	resp.WriteJson(this.Ctx.ResponseWriter)
}
func (this *SocialController) addSocial(resp *helper.Response) {
	content := this.GetString("json")
	social := models.NewSocial()
	err := json.Unmarshal([]byte(content), &social)
	if err != nil {
		log.Error(err)
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "内容错误|要仔细检查哦。"}
		return
	}
	if models.Blogger.AddSocial(social) != RS.RS_success {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "工具已存在|你确定要添加工具吗？"}
		return
	}
}

type social struct {
	ID     string
	SortID int
	Extra  string
	Class  string
	Time   string
}

func (this *SocialController) getSocials(resp *helper.Response) {
	socialsT := beego.BeeTemplates["manage/social/socials.html"]
	var buffer bytes.Buffer
	var socials []*social
	for _, s := range models.Blogger.Socials {
		temp := &social{}
		temp.ID = s.ID
		temp.SortID = s.SortID
		temp.Extra = s.Node.Children[0].Extra
		temp.Class = s.Node.Children[0].Children[0].Class
		temp.Time = s.CreateTime.Format(helper.Layout_y_m_d_time)
		socials = append(socials, temp)
	}
	socialsT.Execute(&buffer, map[string]interface{}{"Socials": socials})
	resp.Data = buffer.String()
}
func (this *SocialController) getSocial(resp *helper.Response) {
	id := this.GetString("id")
	if id != "" {
		modifysocialT := beego.BeeTemplates["manage/social/modifysocial.html"]
		var buffer bytes.Buffer
		if social := models.Blogger.GetSocialByID(id); social != nil {
			b, _ := json.Marshal(social)
			modifysocialT.Execute(&buffer, map[string]string{"Content": string(b)})
			resp.Data = buffer.String()
			return
		}
	}
	resp.Data = ""
}
func (this *SocialController) doModify(resp *helper.Response) {
	content := this.GetString("json")
	if content != "" {
		temp := models.Social{}
		err := json.Unmarshal([]byte(content), &temp)
		if err != nil {
			resp.Status = RS.RS_failed
			resp.Err = helper.Error{Level: helper.WARNING, Msg: "内容错误|反序列化失败。"}
			return
		}
		social := models.Blogger.GetSocialByID(temp.ID)
		if social != nil {
			*social = temp
			sort.Sort(models.Blogger.Socials)
		} else {
			resp.Status = RS.RS_failed
			resp.Err = helper.Error{Level: helper.WARNING, Msg: "索引错误|没有找到该工具。"}
			return
		}
	} else {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "获取内容失败|你都干了什么。"}
		return
	}
	resp.Success()
}
func (this *SocialController) doDelete(resp *helper.Response) {
	id := this.GetString("id")
	if id == "" {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "哦噢。。。|参数错误。"}
		return
	}
	if code := models.Blogger.DelSocialByID(id); code != RS.RS_success {
		resp.Status = code
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "抱歉|系统没有找到该工具。"}
		return
	}
}
