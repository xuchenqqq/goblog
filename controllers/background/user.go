package background

import (
	"bytes"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/smalltree0/beego_goblog/RS"
	"github.com/smalltree0/beego_goblog/helper"
	"github.com/smalltree0/beego_goblog/models"
)

type UserController struct {
	BackgroundController
}

func (this *UserController) Post() {
	resp := helper.NewResponse()
	flag := this.GetString("flag")
	switch flag {
	case "info":
		this.userInfo(resp)
	case "modifyinfo":
		this.doModifyInfo(resp)
	case "modpasswd":
		this.modifyPasswd(resp)
	case "domodpasswd":
		this.doModifyPasswd(resp)
	default:
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "参数错误|未知的flag标志。"}
	}
	resp.WriteJson(this.Ctx.ResponseWriter)
}
func (this *UserController) userInfo(resp *helper.Response) {
	infoT := beego.BeeTemplates["manage/user.html"]
	Map := make(map[string]string)
	Map["BlogName"] = models.Blogger.BlogName
	Map["Icon"] = models.Blogger.HeadIcon
	Map["Introduce"] = models.Blogger.Introduce
	Map["Sex"] = models.Blogger.Sex
	Map["Username"] = models.Blogger.UserName
	Map["Email"] = models.Blogger.Email
	Map["Address"] = models.Blogger.Address
	Map["Education"] = models.Blogger.Education
	Map["LoginTime"] = models.Blogger.LoginTime.Format(helper.Layout_y_m_d_time)
	Map["IP"] = models.Blogger.LoginIp

	var buffer bytes.Buffer
	infoT.Execute(&buffer, Map)
	resp.Data = fmt.Sprintf("%s", string(buffer.Bytes()))
}
func (this *UserController) modifyPasswd(resp *helper.Response) {
	resp.Data = `
<form class="form-horizontal">
  <div class="form-group">
    <label for="inputPasswd" class="col-sm-2 control-label">原密码</label>
    <div class="col-sm-10">
      <input type="password" class="form-control" id="inputPasswd" placeholder="Orgin Password">
    </div>
  </div>
  <div class="form-group">
    <label for="inputNewPasswd" class="col-sm-2 control-label">新密码</label>
    <div class="col-sm-10">
      <input type="password" class="form-control" id="inputNewPasswd" placeholder="New password">
    </div>
  </div>
  <div class="form-group">
    <label for="inputNewPasswd2" class="col-sm-2 control-label">再次确认</label>
    <div class="col-sm-10">
      <input type="password" class="form-control" id="inputNewPasswd2" placeholder="Input Again">
    </div>
  </div>
</form>
<script>
  $('.modal-footer').html('<button type="button" class="btn btn-default" data-dismiss="modal">Close</button><button type="button" id="changepasswd" class="btn btn-primary">Save changes</button>');
  $('#changepasswd').on('click', function(){
  	var originP = $('#inputPasswd').val();
	var newP = $('#inputNewPasswd').val();
	var againP = $('#inputNewPasswd2').val();
	if (newP!=againP||originP==""||newP==""||againP==""){
		pushMessage("info", "错误|请检查输入是否正确。");
		return;
	}
	var resp = get('post', "/admin/user", {flag:'domodpasswd', old:originP, new:newP}, false);
    if (resp.Status != success){pushMessage(resp.Err.Level, resp.Err.Msg);return;}
    pushMessage(resp.Data.Level, resp.Data.Msg);
    $(".modal-footer [data-dismiss='modal']").trigger("click");
  });	
</script>
`
}
func (this *UserController) doModifyInfo(resp *helper.Response) {
	if blogname := this.GetString("blogname"); blogname != "" {
		models.Blogger.BlogName = blogname
	}
	if icon := this.GetString("icon"); icon != "" {
		models.Blogger.HeadIcon = icon
	}
	if introduce := this.GetString("introduce"); introduce != "" {
		models.Blogger.Introduce = introduce
	}
	if sex := this.GetString("sex"); sex != "" {
		models.Blogger.Sex = sex
	}
	if email := this.GetString("email"); email != "" {
		models.Blogger.Email = email
	}
	if address := this.GetString("address"); address != "" {
		models.Blogger.Address = address
	}
	if education := this.GetString("education"); education != "" {
		models.Blogger.Education = education
	}
	resp.Success()
}
func (this *UserController) doModifyPasswd(resp *helper.Response) {
	oldPasswd := this.GetString("old")
	newPasswd := this.GetString("new")
	if oldPasswd == "" || newPasswd == "" {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "错误|参数错误。"}
		return
	}
	if !helper.VerifyPasswd(models.Blogger.PassWord, models.Blogger.UserName, oldPasswd, models.Blogger.Salt) {
		resp.Status = RS.RS_failed
		resp.Err = helper.Error{Level: helper.WARNING, Msg: "错误|原密码错误。"}
		return
	}
	models.Blogger.ChangePassword(newPasswd)
	resp.Success()
}
