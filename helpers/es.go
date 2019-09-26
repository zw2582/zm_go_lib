package helpers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego"
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

func (this EsHelper) SendPut(basepath string, data interface{}) HArr {
	jd, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return this.SendPutRaw(basepath, jd)
}

func (this EsHelper) SendPost(basepath string, data interface{}) HArr {
	jd, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return this.SendPostRaw(basepath, jd)
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
