package main

import (
	"flag"

	"github.com/astaxie/beego"
	_ "github.com/deepzz/beego_goblog/routers"
)

var runmode string

func init() {
	flag.StringVar(&runmode, "m", "dev", "runtime mode:should prod")
}
func main() {
	flag.Parse()
	if runmode != "" {
		beego.BConfig.RunMode = beego.PROD
	}
	beego.Run()
}
