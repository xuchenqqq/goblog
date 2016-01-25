package mgodb

import (
	"github.com/smalltree0/com/monitor"
)

const (
	DB      = "goblog" // database数据库
	C_USER  = "user"   // collections表
	C_TOPIC = "content"
	// 初始化管理员账号
	ADMIN_NAME  = "g"
	ADMIN_PASS  = "g"
	ADMIN_EMAIL = "chenqijing2@163.com"
)

var Deepzz *User

func init() {
	monitor.RegistExitFunc("flushdata", flushdata)
	UMgr.RegisterUser(ADMIN_NAME, ADMIN_PASS, ADMIN_EMAIL) // 注册账号

	Deepzz = UMgr.Get("deepzz")
}

func flushdata() {
	UMgr.UpdateUsers()
}
