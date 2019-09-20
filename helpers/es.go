package helpers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"io"
	"io/ioutil"
	"net/http"
)

type EsHelper struct {
	EsHost string
	EsPort string
}

type HArr map[string]interface{}

func (this HArr) String() string {
	if b, err := json.Marshal(this); err != nil {
		panic(err)
	} else {
		return string(b)
	}
}

type EsSync interface {
	GetId() int
}

func DefaultEsHelper() EsHelper {
	return EsHelper{
		beego.AppConfig.DefaultString("es_host", "http://127.0.0.1"),
		beego.AppConfig.DefaultString("es_port", "9200"),
	}
}

//IndexExist 判断索引是否存在
func (this EsHelper) IndexExist(idx string) bool {
	es_path := this.EsHost + ":" + this.EsPort
	res, err := http.Get(es_path + "/_cat/indices")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	r := bufio.NewReader(res.Body)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		tindex := 0
		sc := bufio.NewScanner(bytes.NewReader(line))
		sc.Split(bufio.ScanWords)
		for sc.Scan() {
			tindex++
			if tindex == 3 && sc.Text() == idx {
				beego.Debug("index is existed")
				return true
			}
		}
	}
	return false
}

//测试分词器
func (this EsHelper) TestToken(idxname, word, tokenDriver string) (tags []string) {
	form, _ := json.Marshal(HArr{
		"analyzer": tokenDriver,
		"text":     word,
	})
	basepath := "/_analyze"
	if idxname != "" {
		basepath = "/" + idxname + basepath
	}
	res := this.SendPostRaw(basepath, form)
	if tokens, ok := res["tokens"].([]interface{}); ok {
		for _, token := range tokens {
			if val, ok := token.(map[string]interface{}); ok {
				if v, ok := val["token"].(string); ok {
					tags = append(tags, v)
				}
			}
		}
	}
	return
}

func (this EsHelper) SyncBeegoModel(indexName string, modelType EsSync) {
	o := orm.NewOrm()

	page := 0
	size := 1000
	for {
		//查询novels
		novels := make([]interface{}, 0)
		_, err := o.QueryTable(modelType).Offset(page*size).Limit(size).Filter("sync", 1).All(&novels)
		page++
		if err != nil {
			if err == orm.ErrNoRows {
				break
			}
			panic(err)
		}
		if len(novels) == 0 {
			beego.Info("no data to sync")
			break
		}
		//修改novel同步状态为2
		var novelIds []interface{}
		for _, val := range novels {
			novelIds = append(novelIds, val.(EsSync).GetId())
		}
		if _, err := o.QueryTable(modelType).Filter("id__in", novelIds...).Update(orm.Params{"sync": 2}); err != nil {
			panic(err)
		}
		//整理同步数据
		bulkTxt := ""
		for _, val := range novels {
			if v, err := json.Marshal(val); err != nil {
				beego.Error(err)
				continue
			} else {
				bulkTxt += fmt.Sprintf("{\"index\":{\"_id\":\"%d\"}}\n", val.(EsSync).GetId())
				bulkTxt += string(v) + "\n"
			}
		}
		//发送请求
		resj := this.SendPostRaw(fmt.Sprintf("/%s/_doc/_bulk", indexName), []byte(bulkTxt))
		if e, ok := resj["error"].(bool); ok && e {
			panic(resj.String())
		}
		for _, v := range resj["items"].([]interface{}) {
			index := v.(map[string]interface{})["index"].(map[string]interface{})
			if int(index["status"].(float64)/10) == 20 {
				o.QueryTable(modelType).Filter("id", index["_id"]).Update(orm.Params{"sync": 3})
				continue
			}

			beego.Error("没有同步成功:", v)
		}
		if len(novels) < size {
			beego.Info("同步结束")
			break
		}
	}
}

func (this EsHelper) SendPutRaw(basepath string, raw []byte) HArr {
	es_path := this.EsHost + ":" + this.EsPort
	req, err := http.NewRequest("PUT", es_path+basepath, bytes.NewReader(raw))
	if err != nil {
		panic(err)
	}
	req.Header.Set("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	resbyte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	resj := HArr{}
	if err := json.Unmarshal(resbyte, &resj); err != nil {
		panic(err)
	}
	return resj
}

func (this EsHelper) SendPostRaw(basepath string, raw []byte) HArr {
	fmt.Println(this.EsHost)
	es_path := this.EsHost + ":" + this.EsPort
	link := es_path + basepath
	beego.Info("eshelper:SendPostRaw:", link, "\r\n", string(raw))
	req, err := http.NewRequest("POST", link, bytes.NewReader(raw))
	if err != nil {
		panic(err)
	}
	req.Header.Set("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	resbyte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	beego.Info("eshelper:sendPostRaw:result:", string(resbyte))
	resj := HArr{}
	if err := json.Unmarshal(resbyte, &resj); err != nil {
		panic(err)
	}
	return resj
}
