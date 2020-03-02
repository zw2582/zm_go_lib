package helpers

import (
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
)

//redis消息队列
type MqRedisSer struct {
	QueueName string //消息队列名称
	QueueBackName string //防丢队列名称
	RedisCli *redis.Client	//redis连接客户端
}

//PushMessage 塞入消息
func (this *MqRedisSer) PushMessage(msg string) int64 {
	if msg == "" {
		return 0
	}

	logs.Debug("redis消息队列 PushMessage："+msg, this.QueueName)
	return this.RedisCli.LPush(this.QueueName, msg).Val()
}

//GetMessage 获取消息
func (this *MqRedisSer) GetMessage() string {
	logs.Debug("redis消息队列 GetMessage,key:", this.QueueName)
	res := this.RedisCli.RPopLPush(this.QueueName, this.QueueBackName)
	return res.Val()
}

//MessageAck 消息
func (this *MqRedisSer) MessageAck(msg string) int64 {
	if msg == "" {
		return 0
	}
	logs.Debug("redis消息队列 MessageAck："+msg, this.QueueName)
	return this.RedisCli.LRem(this.QueueBackName, -1, msg).Val()
}

//Repeat 回滚异常队列
func (this *MqRedisSer) Repeat() int {
	n := 0
	repeatKey := this.QueueName+"_repeat"
	for {
		res := this.RedisCli.RPop(this.QueueBackName).Val()
		if res == "" {
			break
		}

		md5res := Md5encode(res)
		cnt := this.RedisCli.HIncrBy(repeatKey, md5res, 1).Val()
		if cnt > 3 {
			this.RedisCli.HDel(repeatKey, md5res)
			continue
		}
		this.RedisCli.LPush(this.QueueName, res)
		n++
	}
	return n
}
