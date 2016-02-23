package background

import (
	"bytes"
	// "encoding/json"
	"fmt"
	// "sort"

	"github.com/astaxie/beego"
	// "github.com/deepzz0/goblog/RS"
	"github.com/deepzz0/goblog/models"
	// "github.com/deepzz0/go-common/log"
)

type MessageController struct {
	BackgroundController
}

func (this *MessageController) Get() {
	this.TplName = "manage/adminTemplate.html"
	this.Data["Title"] = "留言管理 - " + models.Blogger.BlogName
	this.LeftBar("message")
	this.Content()
}

func (this *MessageController) Content() {
	catT := beego.BeeTemplates["manage/message.html"]
	var buffer bytes.Buffer
	catT.Execute(&buffer, map[string]string{"ID": "99999", "Url": this.domain + "/message"})
	this.Data["Content"] = fmt.Sprintf("%s", string(buffer.Bytes()))
}
