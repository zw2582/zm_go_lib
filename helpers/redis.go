package helpers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"net"
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