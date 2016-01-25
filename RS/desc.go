package RS

// game server rpc status code and relevant description table
// Auto created by /Users/chen/gopath/src/github.com/smalltree0/beego_goblog/tool
//

var descDict = map[int]string{
	RS_failed  : "操作失败",
	RS_success : "操作成功",
	RS_user_exist       : "账号已存在",
	RS_user_inexistence : "账号不存在",
	RS_activate_failed  : "激活失败",
	RS_password_error   : "密码错误",
	RS_query_failed  : "查询失败",
	RS_update_failed : "更新失败",
	RS_create_failed : "创建失败",
	RS_delete_failed : "删除失败",
	RS_user_not_activate : "用户暂未激活",
	RS_user_not_login    : "用户没有登录",
	RS_params_error : "参数错误",
}

func Desc(code int) string {
    desc, found := descDict[code]
    if !found {
		return "未定义状态"
    }
    return desc 
}
