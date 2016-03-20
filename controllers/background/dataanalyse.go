package background

type DataAnalyseController struct {
	Common
}

func (this *DataAnalyseController) Get() {
	this.TplName = "manage/adminlayout.html"
	this.LeftBar("data")
	this.Data["Content"] = ""
}

func (this *DataAnalyseController) Content() {

}
