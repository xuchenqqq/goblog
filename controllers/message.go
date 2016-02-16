package controllers

import (
	"bytes"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/smalltree0/beego_goblog/models"
)

type MessageController struct {
	BaseController
}

func (this *MessageController) Get() {
	this.TplName = "home.html"
	this.Leftbar("message")
	this.Content()
	this.Data["Title"] = "给我留言 - " + models.Blogger.BlogName
}

func (this *MessageController) Content() {
	messageT := beego.BeeTemplates["message.html"]
	var buffer bytes.Buffer
	content := `
	<div class="post-content">
	<p><a href="javascript:;"><img class＝"img-responsive img-rounded" src="/static/message.jpg"><img></a></p>
	<p>非常感谢你关注我的博客，如果你想联系我，可以通过下面的联系方式。</p>
	<p><a class="btn btn-sm btn-primary" href="mailto:chenqijing2@qq.com" target="_blank"><i class="fa fa-qq"></i>QQ邮箱</a><a class="btn btn-sm btn-primary" href="mailto:chenqijing2@163.com" target="_blank"><i class="fa fa-envelope-o"></i>网易邮箱</a></p>
	<p>当然，如果你有新浪微博或者腾讯微博的话，也可以在上面给我留言。</p>
	<p><a class="btn btn-sm btn-primary" href="http://weibo.com/52dxs" target="_blank"><i class="fa fa-weibo"></i>新浪微博</a><a class="btn btn-sm btn-primary" href="http://t.qq.com/chenqijing2" target="_blank"><i class="fa fa-tencent-weibo"></i>腾讯微博</a></p>
	<p>或者，直接在本页留言也可以，不过不确定什么时候会看到，:D</p></div>`
	Map := make(map[string]string)
	Map["Title"] = "给我留言"
	Map["ID"] = "99999"
	Map["Url"] = this.domain + "/message"
	Map["Content"] = content
	messageT.Execute(&buffer, Map)
	this.Data["Content"] = fmt.Sprintf("%s", buffer.Bytes())
}
