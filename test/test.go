package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	//path := "D:\\360MoveData\\Users\\BlockPulse\\Desktop\\1.png"
	//GetUrlImgBase64(path)
	http.HandleFunc("/entrance", Entrance)
	http.HandleFunc("/uploadImg", UploadImg)
	http.ListenAndServe(":8000", nil)
}

var Conn redis.Conn

func sleep() {
	// 休眠1秒
	// time.Millisecond    表示1毫秒
	// time.Microsecond    表示1微妙
	// time.Nanosecond    表示1纳秒
	// 休眠100毫秒
	time.Sleep(10000 * time.Millisecond)
}
func redisC() {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("connect redis failed,", err)
		return
	}
	Conn = conn
}

const UPLOAD_PATH string = "C:/Users/benben/Desktop/"

type Img struct {
	Id     bson.ObjectId `bson:"_id"`
	ImgUrl string        `bson:"imgUrl"`
}

func Entrance(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("uploadImg.html")
	t.Execute(w, nil)
}

func UploadImg(w http.ResponseWriter, r *http.Request) {
	var img Img
	img.Id = bson.NewObjectId()

	r.ParseMultipartForm(1024)
	imgFile, imgHead, imgErr := r.FormFile("img")
	if imgErr != nil {
		fmt.Println(imgErr)
		return
	}
	defer imgFile.Close()

	imgFormat := strings.Split(imgHead.Filename, ".")
	img.ImgUrl = img.Id.Hex() + "." + imgFormat[len(imgFormat)-1]

	image, err := os.Create(UPLOAD_PATH + img.ImgUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer image.Close()

	_, err = io.Copy(image, imgFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	session, err := mgo.Dial("")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	err = session.DB("images").C("image").Insert(img)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("success to upload img")
}
