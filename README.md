## beego辅助类库

#### 1.根据mysql数据库表结构生成model文件

入口文件main.go中加入代码，改代码必须在数据库注册完成之后加入
```
commands.ModelGeneratorCommand()
```

.\main.exe model_generator -force -table=t_account -db

-force : 非必选，是否强制生成，默认false，可以不传

-table : 非必选，指定生成某个表，不传则生成所有表

-db : 非必选，指定数据库

#### 2.用户jwt认证
初始化示例：
```
    //初始化authUser
	this.AuthUser = helpers.AuthUser{
		JwtSecret:jwt_secret,
		BlackKey:"auth_shilian_black",
		UserModel:full_user.TUser{},
	}
	//分析token
	tok := this.Ctx.Request.Header.Get("Authorization")
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		tok = tok[7:]
	}
	if tok == "" {
		tok = this.Ctx.GetCookie("auth_token")
	}
```

使用示例：
``` 
    //获取用户uid
	uid, err := this.AuthUser.Uid()
	if err != nil {
		//todo登录失败
		this.AjaxFail(nil, err.Error())
		return
	}

    //用户登录,获取jwttoken
    user := models.User{Id:1,Name:"test"}
    token, err := this.AuthUser.JwtLogin(user)
```