package background

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/astaxie/beego"
	"github.com/smalltree0/beego_goblog/RS"
	"github.com/smalltree0/beego_goblog/helper"
	"github.com/smalltree0/beego_goblog/models"
	"github.com/smalltree0/com/log"
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
	socialT := beego.BeeTemplates["manage/social.html"]
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
func (this *SocialController) getSocials(resp *helper.Response) {
	var html string
	for _, social := range models.Blogger.Socials {
		html += "<tr>"
		html += "<th scope='row'><input id='" + social.ID + "'' type='checkbox'></th>"
		html += "<td>" + fmt.Sprint(social.SortID) + "</td>"
		html += "<td>" + social.ID + "</td>"
		html += "<td>" + social.Node.Children[0].Extra + "</td>"
		html += "<td><i class='" + social.Node.Children[0].Children[0].Class + "'</i></td>"
		html += "<td>" + social.CreateTime.Format(helper.Layout_y_m_d_time) + "</td>"
		html += `<td><button type="button" data-toggle="modal" data-target="#gridSystemModal" class="btn btn-info btn-xs modify">修改</button><button type="button" class="btn btn-warning btn-xs delete-social">删除</button></td>`
		html += "</tr>"
	}
	var script string
	script = `<script>$('.modify').on('click',function(){
			var id = $(this).parent().parent().find('th input').attr('id');
			if (id==""){pushMessage('warning', '对不起|系统错误。');}
			var resp = get('post', location.pathname, {flag:'modify', id:id}, false);
			if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
			$('.modal-body .container-fluid').html(resp.Data);
		});
		$('.delete-social').on('click', function(){
			var id = $(this).parent().parent().find('th input').attr('id');
			if (id==""){pushMessage('warning', '对不起|系统错误。');}
			if (!confirm('确定要删除该工具吗？')){return;}
			var resp = get('post', location.pathname, {flag:'deletesocial', id:id}, false);
			if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
			location.reload();
		});	
		</script>`
	resp.Data = html + script
}
func (this *SocialController) getSocial(resp *helper.Response) {
	var html string = `<div class="row"><textarea class="form-control" id="social-content" rows="20">`
	id := this.GetString("id")
	if id != "" {
		for _, social := range models.Blogger.Socials {
			if social.ID == id {
				b, _ := json.Marshal(social)
				html += string(b)
			}
		}
	} else {
		resp.Data = ""
	}
	html += "</textarea></div>"
	script := `<script>
		$('#gridModalLabel').text('修改分类');
		$('#social-content').text(JSONFormat($('#social-content').val()));
		$('.modal-footer').html('<button type="button" class="btn btn-default" data-dismiss="modal">Close</button><button type="button" id="modifysocial" class="btn btn-primary">Save changes</button>');
		$('#modifysocial').on('click',function(){
			var content = $('#social-content').val();
			if (content==""){pushMessage('warning', '错误|请填写完整。');return;}
			var resp = get('post', location.pathname, {flag:'domodify',json:content},false);
			console.log(resp)
			if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
			location.reload();
		});</script>`
	resp.Data = html + script
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
