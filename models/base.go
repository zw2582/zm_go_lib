package models

import "time"

type BaseModel struct {
	Id int `json:"id" orm:"pk;auto"`
	Created time.Time `json:"created" orm:"auto_now_add;type(datetime)" `
	Updated time.Time `json:"updated" orm:"auto_now;type(datetime)"`
	Status int `json:"status" orm:"default(1);description:(状态 1.有效 0.无效)"`		//状态 1.有效 0.无效
}
