package controllers

import "github.com/astaxie/beego"

type BaseController struct {
	beego.Controller
}

//返回成功
func (c *BaseController) AjaxSucc(data interface{}, msg string) {
	c.Data[`json`] = map[string]interface{} {
		`status`:0,
		`msg`:msg,
		`data`:data,
	}
	//beego.Debug(fmt.Sprintf("ajaxSucc:%+v", c.Data[`json`]))
	c.ServeJSON()
}

//返回失败
func (c *BaseController) AjaxFail(data interface{}, msg string) {
	c.Data[`json`] = map[string]interface{} {
		`status`:1,
		`msg`:msg,
		`data`:data,
	}
	//beego.Debug(fmt.Sprintf("ajaxSucc:%+v", c.Data[`json`]))
	c.ServeJSON()
}