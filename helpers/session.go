package helpers

import "github.com/astaxie/beego"
import _ "github.com/astaxie/beego/session/redis"

//InitSession 初始化session
func InitSession() {
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionProvider = "redis"
	beego.BConfig.WebConfig.Session.SessionName = beego.BConfig.AppName + "_ID"

	redisHost := beego.AppConfig.DefaultString("redis_host", "127.0.0.1")
	redisPort := beego.AppConfig.DefaultString("redis_port", "6379")
	beego.BConfig.WebConfig.Session.SessionProviderConfig = redisHost + ":" + redisPort
}
