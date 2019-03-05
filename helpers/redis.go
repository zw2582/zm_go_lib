package helpers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
	"net"
)

var _client = registRedisPool()

//GetClient：获取redis连接
func GetRedisClient() *redis.Client {
	return _client
}

//RegistRedisPool:注册全局redis连接池
func registRedisPool() *redis.Client {
	host := beego.AppConfig.String("redis_host")
	password := beego.AppConfig.String("redis_password")
	port := beego.AppConfig.String("redis_port")

	beego.Debug(fmt.Sprintf("连接redis配置信息,host:%s,password:%s,port:%s", host, password, port))

	client := redis.NewClient(&redis.Options{
		Addr:net.JoinHostPort(host, port),
		Password:password,
		DB:0,
	})
	if err := client.Ping().Err(); err != nil {
		beego.Error("连接redis失败")
	}
	return client
}