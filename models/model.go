package models

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/deepzz/beego_goblog/RS"
	// "github.com/deepzz/com/log"
	"github.com/deepzz/com/monitor"
)

const (
	DB         = "newblog" // database数据库
	C_USER     = "user"    // collections表
	C_TOPIC    = "topic"
	C_TOPIC_ID = "topic_id" // 文章ID计数
)

var Blogger *User

func init() {
	monitor.RegistExitFunc("flushdata", flushdata)
	// 以下三句保证调用顺序
	UMgr.loadUsers()
	Blogger = UMgr.Get("deepzz")
	if Blogger == nil {
		path, _ := os.Getwd()
		f, err := os.Open(path + "/conf/backup/user.json")
		if err != nil {
			panic(err)
		}
		user := User{}
		b, _ := ioutil.ReadAll(f)
		err = json.Unmarshal(b, &user)
		if err != nil {
			panic(err)
		}
		UMgr.RegisterUser(&user)
		code := UMgr.UpdateUsers()
		if code != RS.RS_success {
			panic("更新用户数据失败。")
		}
		Blogger = UMgr.Get("deepzz")
	}
	TMgr.loadTopics()
}

func flushdata() {
	UMgr.UpdateUsers()
	TMgr.UpdateTopics()
}
