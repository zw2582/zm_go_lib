package helpers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)


//LogFileConfig 设置文件日志
func LogFileConfig() {
	loglevel := beego.AppConfig.DefaultInt("log_level", 7)
	logs.SetLogger(logs.AdapterFile,`{"filename":"logs/project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	beego.SetLevel(loglevel)
}
