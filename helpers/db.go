package helpers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//DbConfig 配置本项目数据库连接
func DbConfig() {
	var db_host = beego.AppConfig.DefaultString("db_host", "127.0.0.1")
	var db_port = beego.AppConfig.DefaultInt("db_port", 3306)
	var db_name = beego.AppConfig.DefaultString("db_name", "weather_kid")
	var db_user = beego.AppConfig.DefaultString("db_user", "root")
	var db_pwd = beego.AppConfig.DefaultString("db_pwd", "password")

	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=Local", db_user, db_pwd, db_host, db_port, db_name), 30)

	orm.Debug = beego.AppConfig.DefaultBool("orm_debug", true)

	if orm.Debug == true {
		sqllog := logs.NewLogger()
		sqllog.SetLogger(logs.AdapterFile, `{"filename":"logs/sql.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
		orm.DebugLog = orm.NewLog(sqllog)
	}
	//l, _ := time.LoadLocation("Asia/Shanghai")
	//orm.SetDataBaseTZ("default", l)
}
