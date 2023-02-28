package main

import (
	"array/db"
	"array/model"
	"array/server"
	"array/util"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {

	//创建服务
	ginServer := gin.Default()
	//跨域
	ginServer.Use(db.Core())

	//路由组
	userGroup := ginServer.Group("api/user")
	{
		//登录
		userGroup.GET("/login", func(context *gin.Context) {
			email := context.Query("email")       //获取登录用户名
			passWord := context.Query("passWord") //获取登录密码
			user, err := db.UserDao.GetUser(email)
			if err != nil {
				log.Println("user:", err)
			}
			if user.Email == "0" {
				// 无此用户(用户名不存在)
				context.JSON(200, gin.H{
					"success": false,
					"code":    400,
					"msg":     "无此用户",
				})
			} else {
				// 获取当前用户名密码
				// 获取登录用户名的密码 查询是否匹配
				// 密码是否匹配
				if passWord != user.PassWord { //密码不一致
					context.JSON(200, gin.H{
						"success": false,
						"code":    400,
						"msg":     "密码错误",
					})
				} else {
					claims := &server.JWTClaims{
						UserID:   user.Id,
						Email:    user.Email,
						Password: user.PassWord,
						UUID:     user.UUid,
					}
					claims.IssuedAt = time.Now().Unix()
					claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(server.ExpireTime)).Unix()
					signedToken, err := server.GetToken(claims)
					if err != nil {
						return
					}
					context.JSON(200, gin.H{ //用户名存在且密码匹配
						"token":   signedToken,
						"user":    user,
						"success": true,
						"code":    200,
						"msg":     "登录成功",
					})
				}
			}
		})
		//注册
		userGroup.GET("/register", func(context *gin.Context) {
			email := context.Query("email")       //获取注册用户名
			passWord := context.Query("passWord") //获取注册用户密码
			name := context.Query("name")
			dataBarth := context.Query("dataBarth")
			emailCode := context.Query("emailCode")
			status := server.EmailCode(email, emailCode)
			if status == true {
				_, err := db.Conn.Do("DEL", email)
				if err != nil {
					log.Println(err)
				}
			}

			user, err := db.UserDao.GetUser(email)
			if passWord == "" {
				context.JSON(200, gin.H{
					"code":    400,
					"success": false,
					"msg":     "注册失败，密码不能为空",
				})
				return
			}
			if status == false {
				context.JSON(200, gin.H{
					"code":    400,
					"success": false,
					"msg":     "验证码错误",
				})
				return
			}
			if err != nil {
				log.Println("user:", err)
			}
			if email == user.Email {
				fmt.Println("用户名已被注册")
				context.JSON(200, gin.H{
					"code":    400,
					"success": false,
					"msg":     "用户名已被注册",
				})
				return
			}
			if user.Email != email {
				//拼装参数
				user.Email = email
				user.PassWord = passWord
				user.Name = name
				user.DataBarth = dataBarth
				users, _ := db.UserDao.GetAll()
				uid := util.RandomString(8, 0)
				for _, u := range users {
					if uid == u.UUid {
						uid = util.RandomString(8, 0)
						continue
					}
				}
				user.UUid = uid
				//新建用户
				server.CreateUser(*user)
				us, _ := db.UserDao.GetUser(email)
				claims := &server.JWTClaims{
					UserID:   us.Id,
					Email:    user.Email,
					Password: user.PassWord,
					UUID:     user.UUid,
				}
				claims.IssuedAt = time.Now().Unix()
				claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(server.ExpireTime)).Unix()
				signedToken, err := server.GetToken(claims)
				if err != nil {
					log.Println(err)
				}
				context.JSON(200, gin.H{
					"token":   signedToken,
					"code":    200,
					"success": true,
					"msg":     "注册成功",
				})
				return
			}
		})
		//发送邮件
		userGroup.GET("/sendmail", func(context *gin.Context) {
			email := context.Query("email") //获取邮件用户名
			emailV := server.EmailRegister(email)
			context.JSON(200, gin.H{
				"data": emailV,
				"code": 200,
			})
		})
		//校验邮件验证码
		userGroup.POST("/emailVerify", func(context *gin.Context) {
			email := context.Query("email")         //获取邮件用户名
			emailCode := context.Query("emailCode") //获取验证码
			status := server.EmailCode(email, emailCode)
			if status == true {
				_, err := db.Conn.Do("DEL", email)
				if err != nil {
					log.Println(err)
				}
			}
			context.JSON(200, gin.H{
				"code": 200,
				"data": status,
			})
		})
		//重置密码
		userGroup.GET("/resetPassword", func(context *gin.Context) {
			email := context.Query("email")         //获取邮件用户名
			passWord := context.Query("passWord")   //获取重置密码
			emailCode := context.Query("emailCode") //获取重置密码
			status := server.EmailCode(email, emailCode)
			if status == true {
				user := model.User{}
				user.Email = email
				row := db.UserDao.UpdatePassWord(&user, passWord)
				_, err := db.Conn.Do("DEL", email)
				if err != nil {
					log.Println(err)
				}
				if row == 1 {
					context.JSON(200, gin.H{
						"code": 200,
						"data": "修改密码成功",
					})
				}
			}
			if status == false {
				context.JSON(200, gin.H{
					"code":    400,
					"success": false,
					"msg":     "邮箱验证失败",
				})
			}
			us, _ := db.UserDao.GetUser(email)
			if us.Email == "" {
				context.JSON(200, gin.H{
					"code":    400,
					"success": false,
					"msg":     "邮箱不存在",
				})
			}

		})
		//修改个人信息
		userGroup.POST("/updateUser", func(context *gin.Context) {
			strToken := context.Query("token")
			headPhoto := context.Query("headPhoto")
			name := context.Query("userName")
			claim, err := server.VerifyAction(strToken)
			if err == nil {
				claim.ExpiresAt = time.Now().Unix() + (claim.ExpiresAt - claim.IssuedAt)
				signedToken, err := server.GetToken(claim)
				if err != nil {
					log.Println(err)
				}
				user, _ := db.UserDao.GetUserID(claim.UserID)
				user.Name = name
				user.HeadPhoto = headPhoto
				rows := db.UserDao.UpdateUser(user)
				success := false
				if rows == 1 {
					success = true
				}
				context.JSON(200, gin.H{
					"token":   signedToken,
					"user":    user,
					"success": success,
					"code":    200,
					"msg":     "修改个人信息成功",
				})
			}
			if err != nil {
				context.JSON(200, gin.H{
					"success": false,
					"code":    400,
					"msg":     "请重新登录",
				})
			}

		})
		//根据token获取用户信息
		userGroup.GET("/getUser", func(context *gin.Context) {
			strToken := context.Query("token")
			claim, err := server.VerifyAction(strToken)
			claim.ExpiresAt = time.Now().Unix() + (claim.ExpiresAt - claim.IssuedAt)
			if err == nil {
				user, _ := db.UserDao.GetUserID(claim.UserID)
				success := false
				if user.Id > 1 {
					success = true
				}
				if success {
					userMess := model.UserMess{}
					userMess.UUID = user.UUid
					userMess.Email = user.Email
					userMess.HeadPhoto = user.HeadPhoto
					userMess.Name = user.Name
					context.JSON(200, gin.H{
						"user":    userMess,
						"success": success,
						"code":    200,
						"msg":     "返回个人信息",
					})
				}
				if !success {
					context.JSON(200, gin.H{
						"success": false,
						"code":    400,
						"msg":     "请重新登录",
					})
				}
			}

		})

	}
	server.RedisC()
	db.Connect(false)
	//服务器端口
	ginServer.Run(":80")
}

func init() {
	Conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("connect redis failed,", err)
		return
	}
	defer Conn.Close()
}
