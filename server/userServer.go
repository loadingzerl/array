package server

import (
	"array/db"
	"array/model"
	"array/util"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

type EmailVerify struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// EmailRegister 发送邮件
func EmailRegister(email string) EmailVerify {
	//CreateAT := fmt.Sprintf("%v", time.Now().Unix())
	emailV := EmailVerify{}
	code := util.RandomString(6, 3)
	//发送邮件
	util.EmailSign(email, code)
	_, err := Conn.Do("Set", email, code)
	if err != nil {
		fmt.Println("set expire error: ", err)
		return emailV
	}
	_, err = Conn.Do("expire", email, 600)

	return emailV
}

func EmailCode(email string, code string) bool {
	isKeyExit, err := redis.Bool(Conn.Do("EXISTS", email))
	if err != nil {
		log.Println(err)
	}
	if isKeyExit == false {
		return false
	}
	emailCode, _ := redis.String(db.Conn.Do("Get", email))

	if emailCode == code {

		return true
	}

	return false
}

func CreateUser(user model.User) {
	user.Time = fmt.Sprintf("%v", time.Now().Unix())
	db.UserDao.CreateUser(&user)
}

var Conn redis.Conn

func RedisC() {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("connect redis failed,", err)
		return
	}
	Conn = conn
}

const (
	SECRETKEY = "243223ffslsfsldfl412fdsfsdf" //私钥
)

type JWTClaims struct { // token里面添加用户信息，验证token后可能会用到用户信息
	jwt.StandardClaims
	UserID      int64    `json:"user_id"`
	Password    string   `json:"password"`
	Email       string   `json:"email"`
	UUID        string   `json:"uuid"`
	FullName    string   `json:"full_name"`
	Permissions []string `json:"permissions"`
}

var (
	Secret     = "dong_tech" // 加盐
	ExpireTime = 360000      // token有效期
)

const (
	ErrorReason_ServerBusy = "服务器繁忙"
	ErrorReason_ReLogin    = "请重新登陆"
)

func GetToken(claims *JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", errors.New(ErrorReason_ServerBusy)
	}
	return signedToken, nil
}

func VerifyAction(strToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		return nil, errors.New(ErrorReason_ServerBusy)
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New(ErrorReason_ReLogin)
	}
	if err := token.Claims.Valid(); err != nil {
		return nil, errors.New(ErrorReason_ReLogin)
	}
	fmt.Println("verify")
	return claims, nil
}

// NewToken 创建token
func NewToken(user model.User) string {
	maxAge := 60 * 60 * 48
	claims := CustomClaims{
		UserId: user.Id,
		Email:  user.Email,
		Exp:    time.Now().Add(time.Duration(maxAge) * time.Second).Unix(), // 过期时间，必须设置，
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SECRETKEY))
	if err != nil {
		fmt.Println(err)
	}
	//_, err = Conn.Do("Set", user.Email, token)
	//if err != nil {
	//	fmt.Println("set expire error: ", err)
	//	return ""
	//}
	//_, err = Conn.Do("expire", user.Email, 172800)
	return tokenString
}

// ParseToken 解析tongId
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(SECRETKEY), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

type CustomClaims struct {
	UserId int64
	Email  string
	Exp    int64
	jwt.StandardClaims
}
