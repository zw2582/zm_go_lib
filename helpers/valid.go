package helpers

import (
	"github.com/astaxie/beego/validation"
	"github.com/pkg/errors"
)

//参数验证
func Valid(data interface{}) error {
	//验证数据
	valid := validation.Validation{}
	b, err := valid.Valid(data)
	if err != nil {
		panic(err)
	}
	if !b {
		//验证失败，报告用户提交的数据有误
		err := valid.Errors[0]
		return errors.New(err.Key+`:`+err.Message)
	}
	return nil
}
