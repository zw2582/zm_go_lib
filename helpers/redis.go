package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"net"
	"reflect"
	"time"
)

var _client *redis.Client

//GetClient：获取redis连接
func GetRedisClient() *redis.Client {
	if _client == nil {
		registRedisPool()
	}
	return _client
}

//RegistRedisPool:注册全局redis连接池
func registRedisPool() {
	host := beego.AppConfig.DefaultString("redis_host", "127.0.0.1")
	password := beego.AppConfig.DefaultString("redis_password","")
	port := beego.AppConfig.DefaultString("redis_port", "6379")

	if host == "" || port == "" {
		panic(errors.New("请在conf/app.conf中redis参数：redis_host，redis_password，redis_port"))
	}

	beego.Debug(fmt.Sprintf("连接redis:host:%s,password:%s,port:%s", host, password, port))

	_client = redis.NewClient(&redis.Options{
		Addr:net.JoinHostPort(host, port),
		Password:password,
		DB:0,
	})
	fmt.Printf("%+v\n", _client)
	if err := _client.Ping().Err(); err != nil {
		beego.Error("连接redis失败")
	}
}

func RedisCache(cacheKey string, d interface{}, ex time.Duration, f func() (interface{},error)) error {
	var redis = GetRedisClient()
	//存在缓存查询缓存
	if v := redis.Get(cacheKey); v != nil && v.Val() != "" {
		if err := json.Unmarshal([]byte(v.Val()), d); err != nil {
			panic(err)
		}
		beego.Debug("read from cache "+cacheKey)
		return nil
	}
	//不存在缓存查询源数据,并保存缓存
	if dd,err := f();err != nil {
		return err
	} else {
		dv := reflect.ValueOf(d)
		if dv.Kind() != reflect.Ptr || dv.IsNil() {
			panic("the data not ptr or is nil")
		}
		if ddv := reflect.ValueOf(dd); ddv.Kind() == reflect.Ptr {
			dv.Elem().Set(ddv.Elem())
		} else {
			dv.Elem().Set(ddv)
		}
	}
	if b,err := json.Marshal(d); err != nil {
		panic(err)
	} else {
		redis.Set(cacheKey, b, ex)
	}
	beego.Debug("read from source and cached in "+cacheKey+" with "+ex.String())
	return nil
}