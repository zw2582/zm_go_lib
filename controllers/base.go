package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/zw2582/zm_go_lib/helpers"
	"strings"
)

type BaseController struct {
	beego.Controller
}

//返回成功
func (c *BaseController) AjaxSucc(data interface{}, msg ...string) {
	c.Data[`json`] = map[string]interface{}{
		`status`: 0,
		`msg`:    strings.Join(msg, ","),
		`data`:   data,
	}
	c.ServeJSON()
}

//返回失败
func (c *BaseController) AjaxFail(data interface{}, msg ...string) {
	c.Data[`json`] = map[string]interface{}{
		`status`: 1,
		`msg`:    strings.Join(msg, ","),
		`data`:   data,
	}
	c.ServeJSON()
}

//返回失败
func (c *BaseController) AjaxReturn(state int, data interface{}, msg ...string) {
	c.Data[`json`] = map[string]interface{}{
		`status`: state,
		`msg`:    strings.Join(msg, ","),
		`data`:   data,
	}
	c.ServeJSON()
}

//GetJsonData 获取json请求体的内容
func (this *BaseController) GetJsonData() helpers.HArr {
	data := make(helpers.HArr)
	json.Unmarshal(this.Ctx.Input.RequestBody, &data)
	return data
}

//GetJsonString 获取json请求体的内容中的一级字符串字段
func (this *BaseController) GetJsonString(field string, defaults ...string) string {
	data := this.GetJsonData()
	if data[field] != nil {
		return data[field].(string)
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

//GetJsonInt 获取json请求体的内容中的一级int字段
func (this *BaseController) GetJsonInt(field string, defaults ...int) int {
	data := this.GetJsonData()
	if data[field] != nil {
		return data[field].(int)
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

//GetJsonInt64 获取json请求体的内容中的一级int64字段
func (this *BaseController) GetJsonInt64(field string, defaults ...int64) int64 {
	data := this.GetJsonData()
	if data[field] != nil {
		return data[field].(int64)
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

//GetJsonFloat64 获取json请求体的内容中的一级float64字段
func (this *BaseController) GetJsonFloat64(field string, defaults ...float64) float64 {
	data := this.GetJsonData()
	if data[field] != nil {
		return data[field].(float64)
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}