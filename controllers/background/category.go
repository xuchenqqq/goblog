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
	// db "github.com/smalltree0/com/mongo"
)

type CategoryController struct {
	BackgroundController
}

func (this *CategoryController) Get() {
	this.TplName = "manage/adminTemplate.html"
	this.Data["Title"] = "分类管理 - " + models.Blogger.BlogName
	this.LeftBar("category")
	this.Content()
}
func (this *CategoryController) Content() {
	catT := beego.BeeTemplates["manage/categoryTag.html"]
	var buffer bytes.Buffer
	catT.Execute(&buffer, "")
	this.Data["Content"] = fmt.Sprintf("%s", string(buffer.Bytes()))
}

func (this *CategoryController) Post() {
	resp := helper.NewResponse()
	flag := this.GetString("flag")
	switch flag {
	case "cat":
		this.getCategories(resp)
	case "tag":
		this.getTag(resp)
	case "addcategory":
		this.addCategory(resp)
	case "modify":
		this.getCategory(resp)
	case "domodify":
		this.doModify(resp)
	case "deletecat":
		this.doDeleteCat(resp)
	case "deletetag":
		this.doDeleteTag(resp)
	default:
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "参数错误|未知的flag标志。"}
	}
	resp.WriteJson(this.Ctx.ResponseWriter)
}
func (this *CategoryController) getCategory(resp *helper.Response) {
	var html string = `<div class="row"><textarea class="form-control" id="cat-content" rows="20">`
	id := this.GetString("id")
	if id != "" {
		for _, cat := range models.Blogger.Categories {
			if cat.IsCategory && cat.ID == id {
				b, _ := json.Marshal(cat)
				html += string(b)
			}
		}
	} else {
		resp.Data = ""
	}
	html += "</textarea></div>"
	script := `<script>
		$('#gridModalLabel').text('修改分类');
		$('#cat-content').text(JSONFormat($('#cat-content').val()));
		$('.modal-footer').html('<button type="button" class="btn btn-default" data-dismiss="modal">Close</button><button type="button" id="modifycat" class="btn btn-primary">Save changes</button>');
		$('#modifycat').on('click',function(){
			var content = $('#cat-content').val();
			if (content==""){pushMessage('info', '错误|请填写完整。');return;}
			var resp = get('post', location.pathname, {flag:'domodify',json:content},false);
			if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
			location.reload();
		});</script>`
	resp.Data = html + script
}
func (this *CategoryController) getCategories(resp *helper.Response) {
	var html string
	for _, cat := range models.Blogger.Categories {
		if len(cat.Node.Children) > 0 && cat.IsCategory {
			html += "<tr>"
			html += "<th scope='row'><input id='" + cat.ID + "'' type='checkbox'></th>"
			html += "<td>" + fmt.Sprint(cat.SortID) + "</td>"
			html += "<td>" + cat.ID + "</td>"
			html += "<td>" + cat.Node.Children[0].Extra + "</td>"
			html += "<td>" + cat.Node.Children[0].Text + "</td>"
			html += "<td>" + fmt.Sprint(cat.Count) + "</td>"
			html += "<td>" + cat.CreateTime.Format(helper.Layout_y_m_d_time) + "</td>"
			html += `<td><button type="button" data-toggle="modal" data-target="#gridSystemModal" class="btn btn-info btn-xs modify">修改</button><button type="button" class="btn btn-warning btn-xs delete-cat">删除</button></td>`
			html += "</tr>"
		}
	}
	var script string
	script = `<script>$('.modify').on('click',function(){
			var id = $(this).parent().parent().find('th input').attr('id');
			if (id==""){pushMessage('info', '对不起|系统错误。');}
			var resp = get('post', location.pathname, {flag:'modify', id:id}, false);
			if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
			$('.modal-body .container-fluid').html(resp.Data);
		});
		$('.delete-cat').on('click', function(){
			var id = $(this).parent().parent().find('th input').attr('id');
			if (id==""){pushMessage('info', '对不起|系统错误。');}
			if (!confirm('确定要删除该分类吗？')){return;}
			var resp = get('post', location.pathname, {flag:'deletecat', id:id}, false);
			if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
			location.reload();
		});	
		</script>`
	resp.Data = html + script
}
func (this *CategoryController) getTag(resp *helper.Response) {
	var html string
	for _, tag := range models.Blogger.Tags {
		html += "<tr>"
		html += "<th scope='row'><input id='" + tag.ID + "'' type='checkbox'></th>"
		html += "<td>" + tag.ID + "</td>"
		html += "<td>" + tag.Node.Extra + "</td>"
		html += "<td>" + tag.Node.Text + "</td>"
		html += "<td>" + fmt.Sprint(tag.Count) + "</td>"
		html += `<td><button type="button" class="btn btn-warning btn-xs delete-tag">删除</button></td>`
		html += "</tr>"
	}
	script := `<script>$('.delete-tag').on('click', function(){
			var id = $(this).parent().parent().find('th input').attr('id');
			if (id==""){pushMessage('info', '对不起|系统错误。');}
			if (!confirm('确定要删除该标签吗？')){return;}
			var resp = get('post', location.pathname, {flag:'deletetag', id:id}, false);
			if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
			location.reload();
		});</script>`
	resp.Data = html + script
}

func (this *CategoryController) doModify(resp *helper.Response) {
	content := this.GetString("json")
	if content != "" {
		cat := models.Category{}
		err := json.Unmarshal([]byte(content), &cat)
		if err != nil {
			resp.Status = RS.RS_failed
			resp.Err = helper.Error{Level: helper.WARNING, Msg: "内容错误|反序列化失败。"}
			return
		}
		category := models.Blogger.GetCategoryByID(cat.ID)
		if category != nil {
			*category = cat
			sort.Sort(models.Blogger.Categories)
		} else {
			resp.Status = RS.RS_failed
			resp.Err = helper.Error{Level: helper.WARNING, Msg: "索引错误|没有找到该分类。"}
			return
		}
	} else {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "获取内容失败|你都干了什么。"}
		return
	}
	resp.Success()
}
func (this *CategoryController) addCategory(resp *helper.Response) {
	content := this.GetString("json")
	cat := models.NewCategory()
	err := json.Unmarshal([]byte(content), &cat)
	if err != nil {
		log.Error(err)
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "内容错误|要仔细检查哦。"}
		return
	}
	if models.Blogger.AddCategory(cat) != RS.RS_success {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "分类已存在|你确定要添加分类吗？"}
		return
	}
}
func (this *CategoryController) doDeleteCat(resp *helper.Response) {
	id := this.GetString("id")
	if id == "" {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "哦噢。。。|参数错误。"}
		return
	}
	if code := models.Blogger.DelCatgoryByID(id); code != RS.RS_success {
		resp.Status = code
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "抱歉|系统没有找到该分类。"}
		return
	}
}
func (this *CategoryController) doDeleteTag(resp *helper.Response) {
	id := this.GetString("id")
	if id == "" {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "哦噢。。。|参数错误。"}
		return
	}
	if code := models.Blogger.DelBlogrollByID(id); code != RS.RS_success {
		resp.Status = code
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "抱歉|系统没有找到该标签。"}
		return
	}
}
