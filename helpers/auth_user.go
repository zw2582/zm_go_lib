package helpers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"mobange/models"
	"strconv"
	"time"
)

//认证user
type AuthUser struct {
	Token string	//认证token
	JwtSecret string //jwt加密秘钥
	BlackKey string	//黑名单redis key
	uid uint
	authed int //是否token认证过 1.已认证 0.未认证
	autherr error
}

var redisCli = GetRedisClient()

//Uid 用户id
func (this *AuthUser) Uid() (uint, error) {
	if this.authed == 1 {
		return this.uid, this.autherr
	}
	//判断token是否过期,5秒内仍有活性
	score := redisCli.ZScore(this.BlackKey, this.Token).Val()
	if score > 0 && float64(time.Now().Unix())-score > 5 {
		this.uid = 0
		this.authed = 1
		this.autherr = models.NewLogoutError("登录已失效")
		return this.uid, this.autherr
	}

	//校验token
	beego.Info("token:", this.Token)
	claim := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(this.Token, &claim, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(this.JwtSecret), nil
	})
	if err != nil {
		beego.Error(err)
		this.autherr = models.NewLogoutError("用户未登录")
		return this.uid, this.autherr
	}
	if !token.Valid {
		beego.Error(err)
		this.autherr = models.NewLogoutError("用户未登录")
		return this.uid, this.autherr
	}
	id,_ := strconv.Atoi(claim.Id)
	this.uid = uint(id)
	this.authed = 1
	return  this.uid, nil
}

//JwtRefresh 刷新
func (this *AuthUser) JwtRefresh() (string, error) {
	//将过去的token加入黑名单
	redisCli.ZAdd(this.BlackKey, redis.Z{float64(time.Now().Unix()), this.Token})

	// 清除redis中超过一个月的记录
	max := time.Now().AddDate(0, -1, 0)
	maxstr := strconv.Itoa(int(max.Unix()))
	redisCli.ZRemRangeByScore(this.BlackKey, "0", maxstr)

	//产生新的token
	uid, err := this.Uid()
	if err != nil {
		return "", err
	}
	user := models.UserGetById(uid)
	return this.JwtLogin(user)
}

//LogOut 退出
func (this *AuthUser) LogOut() {
	//判断是否认证成功
	_, err := this.Uid()
	if err != nil {
		return
	}
	//将过去的token加入黑名单
	redisCli.ZAdd(this.BlackKey, redis.Z{float64(time.Now().Unix()-10), this.Token})
}

//JwtLogin 登录
func (this *AuthUser) JwtLogin(user models.User) (tokenStr string, err error) {
	if user.Id == 0 {
		return "", fmt.Errorf("用户不存在")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour*24*5).Unix(),	//5天后失效
		Subject:"",
		Id:strconv.Itoa(int(user.Id)),
	})

	tokenStr, err = token.SignedString([]byte(this.JwtSecret))
	if err != nil {
		panic(err)
	}
	this.uid = user.Id
	this.autherr = nil
	this.authed = 1
	this.Token = tokenStr
	return
}

