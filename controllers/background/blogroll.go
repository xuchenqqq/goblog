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
	blogrollT := beego.BeeTemplates["manage/blogroll.html"]
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

func (this *BlogrollController) getBlogrolls(resp *helper.Response) {
	var html string
	for _, br := range models.Blogger.Blogrolls {
		html += "<tr>"
		html += "<th scope='row'><input id='" + br.ID + "'' type='checkbox'></th>"
		html += "<td>" + fmt.Sprint(br.SortID) + "</td>"
		html += "<td>" + br.ID + "</td>"
		html += "<td>" + br.Node.Children[0].Extra + "</td>"
		html += "<td>" + br.Node.Children[0].Text + "</td>"
		html += "<td>" + br.CreateTime.Format(helper.Layout_y_m_d_time) + "</td>"
		html += `<td><button type="button" data-toggle="modal" data-target="#gridSystemModal" class="btn btn-info btn-xs modify">修改</button><button type="button" class="btn btn-warning btn-xs delete-blogroll">删除</button></td>`
		html += "</tr>"
	}
	var script string
	script = `<script>$('.modify').on('click',function(){
			var id = $(this).parent().parent().find('th input').attr('id');
			if (id==""){pushMessage('info', '对不起|系统错误。');}
			var resp = get('post', location.pathname, {flag:'modify', id:id}, false);
			if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
			$('.modal-body .container-fluid').html(resp.Data);
		});
		$('.delete-blogroll').on('click', function(){
			var id = $(this).parent().parent().find('th input').attr('id');
			if (id==""){pushMessage('info', '对不起|系统错误。');}
			if (!confirm('确定要删除该链接吗？')){return;}
			var resp = get('post', location.pathname, {flag:'deleteblogroll', id:id}, false);
			if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
			location.reload();
		});	
		</script>`
	resp.Data = html + script
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
	var html string = `<div class="row"><textarea class="form-control" id="blogroll-content" rows="20">`
	id := this.GetString("id")
	if id != "" {
		for _, br := range models.Blogger.Blogrolls {
			if br.ID == id {
				b, _ := json.Marshal(br)
				html += string(b)
			}
		}
	} else {
		resp.Data = ""
	}
	html += "</textarea></div>"
	script := `<script>
		$('#gridModalLabel').text('修改分类');
		$('#blogroll-content').text(JSONFormat($('#blogroll-content').val()));
		$('.modal-footer').html('<button type="button" class="btn btn-default" data-dismiss="modal">Close</button><button type="button" id="modifyblogroll" class="btn btn-primary">Save changes</button>');
		$('#modifyblogroll').on('click',function(){
			var content = $('#blogroll-content').val();
			if (content==""){pushMessage('warning', '错误|请填写完整。');return;}
			var resp = get('post', location.pathname, {flag:'domodify',json:content},false);
			console.log(resp)
			if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
			location.reload();
		});</script>`
	resp.Data = html + script
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
