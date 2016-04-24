package models

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/deepzz0/goblog/RS"
	// "github.com/deepzz0/go-common/log"
	"github.com/deepzz0/go-common/monitor"
)

const (
	DB         = "newblog" // database数据库
	C_USER     = "user"    // collections表
	C_TOPIC    = "topic"
	C_TOPIC_ID = "topic_id" // 文章ID计数
)

var Blogger *User

func init() {
	monitor.HookOnExit("flushdata", flushdata)
	// 以下三句保证调用顺序
	UMgr.loadUsers()
	Blogger = UMgr.Get("deepzz")
	if Blogger == nil { // 从配置初始化用户
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
