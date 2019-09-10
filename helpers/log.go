package helpers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

//LogFileConfig 设置文件日志
func LogFileConfig() {
	loglevel := beego.AppConfig.DefaultInt("log_level", 7)
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/project.log","daily":true,"maxdays":7}`)
	beego.SetLevel(loglevel)
}

func LogMultiFileConfig() {
	loglevel := beego.AppConfig.DefaultInt("log_level", 7)
	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/project.log","daily":true,"maxdays":7,"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
	beego.SetLevel(loglevel)
}

func LogSmtpConfig() {
	user := beego.AppConfig.String("smtp_user")
	pwd := beego.AppConfig.String("smtp_pwd")
	host := beego.AppConfig.String("smtp_host")
	sendto := beego.AppConfig.String("smtp_log_sendto")
	if user == "" || pwd == "" || host == "" || sendto == "" {
		return
	}
	logs.SetLogger(logs.AdapterMail,
		fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\",\"host\":\"%s\",\"sendTos\":[\"%s\"],\"level\":4,\"subject\":\"%s\"}",
			user, pwd, host, sendto, "系统报错"))
}
