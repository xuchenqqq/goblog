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
	// db "github.com/deepzz0/go-common/mongo"
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
	catT := beego.BeeTemplates["manage/category/categoryTemplate.html"]
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
	id := this.GetString("id")
	if id != "" {
		modifycategoryT := beego.BeeTemplates["manage/category/modifycategory.html"]
		var buffer bytes.Buffer
		if cat := models.Blogger.GetCategoryByID(id); cat != nil {
			b, _ := json.Marshal(cat)
			modifycategoryT.Execute(&buffer, map[string]string{"Content": string(b)})
			resp.Data = buffer.String()
			return
		}
	}
	resp.Data = ""
}

type category struct {
	ID     string
	SortID int
	Extra  string
	Name   string
	Count  int
	Time   string
}

func (this *CategoryController) getCategories(resp *helper.Response) {
	categoriesT := beego.BeeTemplates["manage/category/categories.html"]
	var buffer bytes.Buffer
	var categories []*category
	for _, cat := range models.Blogger.Categories {
		if len(cat.Node.Children) > 0 && cat.IsCategory {
			temp := &category{}
			temp.ID = cat.ID
			temp.SortID = cat.SortID
			temp.Extra = cat.Node.Children[0].Extra
			temp.Name = cat.Node.Children[0].Text
			temp.Count = cat.Count
			temp.Time = cat.CreateTime.Format(helper.Layout_y_m_d_time)
			categories = append(categories, temp)
		}
	}
	categoriesT.Execute(&buffer, map[string]interface{}{"Categories": categories})
	resp.Data = buffer.String()
}

type tag struct {
	ID    string
	Extra string
	Name  string
	Count int
}

func (this *CategoryController) getTag(resp *helper.Response) {
	tagsT := beego.BeeTemplates["manage/category/tags.html"]
	var buffer bytes.Buffer
	var tags []*tag
	for _, t := range models.Blogger.Tags {
		temp := &tag{}
		temp.ID = t.ID
		temp.Extra = t.Node.Extra
		temp.Name = t.Node.Text
		temp.Count = t.Count
		tags = append(tags, temp)
	}
	tagsT.Execute(&buffer, map[string]interface{}{"Tags": tags})
	resp.Data = buffer.String()
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
	if id == "" || id == "default" {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "哦噢。。。|参数错误,default不能删除。"}
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
