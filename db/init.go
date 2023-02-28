package db

import (
	"array/model"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"sync"
	"time"
)

var mainDB *gorm.DB
var _once sync.Once
var Conn redis.Conn

func Connect(showSQL bool) {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("connect redis failed,", err)
		return
	}
	Conn = conn
	_once.Do(func() {
		//db, err := gorm.Open(postgres.New(postgres.Config{
		//	DSN:                  config.DatabaseStr(),
		//	PreferSimpleProtocol: true, // disables implicit prepared statement usage
		//}), &gorm.Config{})
		//db, err := gorm.Open(mysql.Open(config.DatabaseStr()), &gorm.Config{})
		array := "root:root@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := gorm.Open(mysql.Open(array), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Panic(err)
		}
		tmp, _ := db.DB()
		tmp.SetMaxIdleConns(10)
		tmp.SetConnMaxIdleTime(5 * time.Minute)
		tmp.SetMaxOpenConns(40)

		if showSQL {
			db = db.Debug()
		}
		mainDB = db
		// connectToRedis()
		if err := db.AutoMigrate(
			&model.User{},
			&model.UserProperty{},
			&model.UserLimit{},
		); err != nil {
			panic(err)
		}

	})
}

func MainDB() *gorm.DB {
	return mainDB
}

func Core() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token")
		c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
		c.Header("Access-Control-Allow-Credentials", "True")
		//放行索引options
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		//处理请求
		c.Next()
	}
}
