package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

var setting Settings

var picMap map[string]string

func main() {
	initLog()
	fmt.Println("程序开始...")
	log.Info("程序开始...")
	setting = loadConfig()
	picMap = make(map[string]string, 1000)
	log.Info("请求地址：" + setting.BingUrl)
	loop := make(chan bool)
	go downloadloop(loop)

	ret := <-loop

	if !ret {
		fmt.Println("程序退出！")
		log.Info("程序退出！")
		close(loop)
		os.Exit(0)
	}
}
func downloadloop(loop chan bool) {
	var t time.Duration = time.Minute * time.Duration(setting.IntervalTime)
	var url = setting.BingUrl
	fmt.Println("循环开始！时间间隔：" + t.String())
	log.Info("循环开始！时间间隔：" + t.String())
	for {
		defer func() {
			if err := recover(); err != nil {
				log.Info(err)
				log.Info("Stack trace:\n" + string(debug.Stack()))
				log.Info("Got panic in goroutine, will start a new one...")
				go downloadloop(loop)
			}
		}()
		data := getDataFromUrl(url)
		s := analyzeHtml(string(data))
		fmt.Println(s)
		log.Info("获得图片地址：" + s)
		getPic(s, setting.SaveDir)

		time.Sleep(t)
	}
	log.Info("循环结束！")
	loop <- false
}
func getPic(jpgUrl string, savePath string) string {

	data := getDataFromUrl(jpgUrl)
	sary := strings.Split(jpgUrl, "/")
	name := sary[len(sary)-1]

	path := savePath + name
	p, exist := picMap[name]
	if exist {
		fmt.Println("图片[" + name + "]已下载过,保存目录[" + p + "]")
		log.Info("图片[" + name + "]已下载过,保存目录[" + p + "]")
		return p
	}
	if len(data) > 10 {
		fmt.Println("[" + name + "]")
		file, err := os.Create(savePath + name)
		if err != nil {
			panic(err)
		}
		size, err := file.Write(data)
		defer file.Close()
		if err != nil {
			panic(err)
		}
		fmt.Println("下载完成,保存为:[" + path + "],文件大小:" + strconv.Itoa(size))
		log.Info("下载完成,保存为:[" + path + "],文件大小:" + strconv.Itoa(size))
	}
	picMap[name] = path
	return path
}
func analyzeHtml(htmlStr string) string {

	reg, _ := regexp.Compile("g_img={url:'(\\S+?)',")
	match := reg.FindAllStringSubmatch(htmlStr, -1)
	for _, v := range match {
		tmp := string(v[1])
		if strings.HasPrefix(tmp, "http://") && strings.HasSuffix(tmp, ".jpg") {
			return tmp
		}
	}
	return ""
}

func getDataFromUrl(url string) []byte {

	rsp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		panic(err)
	}
	return data
}
