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
	topicT := beego.BeeTemplates["manage/topicsmanage.html"]
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

func (this *TopicsController) getTopics(resp *helper.Response) {
	cat := this.GetString("cat")
	page, err := this.GetInt("page")
	if err != nil {
		log.Debug(err)
		page = 1
	}
	var html string
	var topics []*models.Topic
	var remainpage int
	if cat == "" || cat == "0" {
		topics, remainpage = models.TMgr.GetTopicsByPage(page)
		if remainpage == -1 {
			resp.Data = "<script>$('#previous').addClass('disabled');$('#next').addClass('disabled');pushMessage('info','莫得数据得嘛|是不是页码错了。')</script>"
			return
		}
	} else {
		topics, remainpage = models.TMgr.GetTopicsByCatgory(cat, page)
		if remainpage == -1 {
			resp.Data = "<script>$('#previous').addClass('disabled');$('#next').addClass('disabled');pushMessage('info','莫得数据哦|我看是分类或者页码错误。')</script>"
			return
		}
	}
	var script string = "<script>"
	if page == 1 {
		script += "$('#previous').addClass('disabled');"
	} else {
		script += `$('#previous').removeClass('disabled');$('#previous').one('click', function(){
			var resp = get('post', location.pathname, {flag:'topics',cat:'` + cat + `',page:'` + fmt.Sprint(page-1) + `'},false);
			if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
			$('#topics-tbody').html(resp.Data);
		});`
	}
	if remainpage == 0 {
		script += "$('#next').addClass('disabled');"
	} else {
		script += `$('#next').removeClass('disabled');$('#next').one('click', function(){
			var resp = get('post', location.pathname, {flag:'topics', cat:'` + cat + `',page:'` + fmt.Sprint(page+1) + `'},false);
			if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
			$('#topics-tbody').html(resp.Data);
		});`
	}
	for _, topic := range topics {
		html += `<tr><th scope="row"><input type="checkbox" id="` + fmt.Sprint(topic.ID) + `"></th>`
		html += "<td>" + fmt.Sprint(topic.ID) + "</td><td>" + topic.Title + "</td><td>" + topic.CategoryID + "</td><td>" + fmt.Sprint(topic.TagIDs) + "</td><td>" + topic.Author + "</td><td>" + topic.CreateTime.Format(helper.Layout_y_m_d_time) + "</td><td>" + topic.EditTime.Format(helper.Layout_y_m_d_time) + "</td>"
		html += `<td><button type="button" data-toggle="modal" data-target="#gridSystemModal" class="btn btn-info btn-xs modify">修改</button><button type="button" class="btn btn-warning btn-xs delete-topic">删除</button></td>`
		html += "</tr>"
	}
	script += `$('.modify').on('click',function(){
			var id = $(this).parent().parent().find('th input').attr('id');
			if (id==""){pushMessage('warning', '对不起|系统错误。');}
			var resp = get('post', location.pathname, {flag:'modify', id:id}, false);
			if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
			$('#gridModalLabel').text('修改文章');
			$('.modal-body .container-fluid').html(resp.Data);
		});
		$('.delete-topic').on('click', function(){
			var id = $(this).parent().parent().find('th input').attr('id');
			if (id==""){pushMessage('warning', '对不起|系统错误。');return;}
			if (!confirm('确定要删除该文章吗？')){return;}
			var resp = get('post', location.pathname, {flag:'deletetopic', id:id}, false);
			if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
			location.reload();
		});	
		</script>`
	resp.Data = html + script
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
func (this *TopicsController) getTopic(resp *helper.Response) {
	// var html string = `<div class="row"><textarea class="form-control" id="topic-content" rows="20">`
	// id := this.GetString("id")
	// if id != "" {
	// 	for _, cat := range models.Blogger.Categories {
	// 		if cat.IsCategory && cat.ID == id {
	// 			b, _ := json.Marshal(cat)
	// 			html += string(b)
	// 		}
	// 	}
	// } else {
	// 	resp.Data = ""
	// }
	// html += "</textarea></div>"
	// script := `<script>
	// 	$('#gridModalLabel').text('修改分类');
	// 	$('#cat-content').text(JSONFormat($('#cat-content').val()));
	// 	$('#change').on('click',function(){
	// 		var content = $('#cat-content').val();
	// 		if (content==""){pushMessage('warning', '错误|请填写完整。');return;}
	// 		var resp = get('post', location.pathname, {flag:'domodify',json:content},false);
	// 		console.log(resp)
	// 		if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
	// 		location.reload();
	// 	});</script>`
	// resp.Data = html + script
}
func (this *TopicsController) doModify(resp *helper.Response) {

}
func (this *TopicsController) category(resp *helper.Response) {
	var html string
	html += `
	<h6><div id='newtopic-tag'>
	</div></h6>
	<div class="btn-group choose-tag">
	  <button type="button" id="choose-tag" class="btn btn-default btn-success" aria-haspopup="true" aria-expanded="false">
	    选择TAG <span class="caret"></span>
	  </button>
	  <div class="dropdown-menu" id="dropdown-tag" style="width:240px;">
	  </div>
	</div>
	<select id="selecet-addtopic" class="form-control" style="width:auto">
	`
	for _, cat := range models.Blogger.Categories {
		if cat.IsCategory {
			html += "<option value='" + cat.ID + "'>" + cat.ID + "</option>"
		}
	}
	html += `</select>
	<script>
	  	$("#choose-tag").on('click', function(){
	  		if($("#dropdown-tag").is(":hidden")){
	  			$('#dropdown-tag').fadeIn();
	  			var resp = get('post', location.pathname, {flag:'tag'}, false);
	  			if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
	  			$("#dropdown-tag").html(resp.Data);
	  		}else{
	  			$('#dropdown-tag').fadeOut();
	  		}
		});
		var tags = new Array();
		function removetag(e){
			var parent = e.parentElement;
			var tag = parent.innerText.substring(0,parent.innerText.length-1);
			console.log(tag);
			tags.remove(tag);
			parent.remove();
		};
	 </script>`
	resp.Data = html
}
func (this *TopicsController) tag(resp *helper.Response) {
	var html string = `
		<div class="input-group input-group-sm" style="margin-bottom:2px;">
	     		<input type="text" class="form-control" id='tag-content' placeholder="Add a tag">
	     		<span class="input-group-btn">
	       			<button class="btn btn-default" id='add-tag' type="button">Add</button>
	     		</span>
	   	</div><!-- /input-group -->`
	for _, tag := range models.Blogger.Tags {
		html += tag.TagStyle()
	}
	script := `<script>
		$('#dropdown-tag span.label').each(function(i){
			$(this).on('click', function(){
				var tag = $(this).text();
				console.log(tag);
				if (tags.indexOf(tag) != -1){
					return;
				}
				tags.push(tag);
				var classname = $(this).attr('class')
				var html = '<span class="'+ classname +'">'+ tag +'<span aria-hidden="true" class="remove-tag" onclick="removetag(this);">×</span></span>';
				$('#newtopic-tag').append(html);
			});
		});
		$('#add-tag').on('click', function(){
			var tag = $('#tag-content').val();
			if (tag == ''){
				pushMessage('info', '错误|你需要填写新tag。');
				return;
			}
			tags.push(tag);
			$('#newtopic-tag').append('<span class="label label-primary">'+ tag +'<span aria-hidden="true" class="remove-tag" onclick="removetag(this);">×</span></span>');
		});
	</script>`
	resp.Data = html + script
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
