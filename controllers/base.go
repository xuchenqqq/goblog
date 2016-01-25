package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/smalltree0/beego_goblog/RS"
)

var sessionname = beego.AppConfig.String("sessionname")

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Prapare() {

}

func (this *BaseController) isLogin() bool {
	val := this.GetSession(sessionname)
	if val == nil || val.(string) == "" {
		return false
	}
	return true
}

// -------------------------- response --------------------------
const (
	WARNING = "warning"
	SUCCESS = "success"
	ALERT   = "alert"
	INFO    = "info"
)

type Response struct {
	Status int
	Data   interface{}
	Err    Error
}
type Error struct {
	Level string
	Msg   string
}

func NewResponse() *Response {
	return &Response{Status: RS.RS_success}
}

func (resp *Response) Tips(level string, rs int) {
	resp.Err = Error{level, "code=" + fmt.Sprint(rs) + "|" + RS.Desc(rs)}
}

func (resp *Response) WriteJson(w http.ResponseWriter) {
	b, err := json.Marshal(resp)
	if err != nil {
		w.Write([]byte(`{Status:-1,Err:Error{Level:"alert",Msg:"code=-1|序列化失败！"}}`))
	} else {
		w.Write(b)
	}
}
