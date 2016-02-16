package background

type DataAnalyseController struct {
	BackgroundController
}

func (this *DataAnalyseController) Get() {
	this.TplName = "manage/adminTemplate.html"
	this.LeftBar("data")
	this.Data["Content"] = ""
}

func (this *DataAnalyseController) Content() {

}
