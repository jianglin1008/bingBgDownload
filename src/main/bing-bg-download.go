package main

import (
	"encoding/json"
	"fmt"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var setting Settings

var picMap map[string]string

func main() {
	initLog()
	fmt.Println("程序开始...")
	log.Info("程序开始...")
	setting = loadConfig()
	initContainer(setting.MaxDlCount) //初始化map,list
	getHashOfHpWall()                 //获取墙纸的hash

	fmt.Println("程序退出！")
	log.Info("程序退出！")

	os.Exit(0)
}
func expHandler() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
}
func initContainer(count int) {
	picMap = make(map[string]string, count)

}
func getHashOfHpWall() {
	var count int = 0
	for i := 1; i <= setting.MaxDlCount; i++ {
		reqObj := NewReqMsg(i)
		rspObj := getDataFromUrlAndData(setting.NextHpUrl + "?" + buildReqParams(reqObj))
		var rspMsg RspMsg
		json.Unmarshal(rspObj, &rspMsg)

		if len(rspMsg.Images) > 0 {

			log.Debug(rspMsg.Images[0].Hsh)
			url := setting.DlUrl + rspMsg.Images[0].Hsh
			count++
			go dl(url)
		}
	}
	log.Info("抓取数据数量:", count)
}
func buildReqParams(reqObj *ReqMsg) string {
	values := "format=" + reqObj.Format
	values = values + "&idx=" + strconv.Itoa(reqObj.Idx)
	values = values + "&n=" + strconv.Itoa(reqObj.N)

	values = values + "&nc=" + strconv.FormatInt(reqObj.Nc, 10)
	values = values + "&pid=" + reqObj.Pid
	return values
}
func parseStr(value interface{}) string {
	str, ok := value.(string)
	if ok {
		return str
	} else {
		panic("Parse error")
	}

}
func dl(url string) string {
	data, filename := getDataFromUrl(url)
	if filename == "fail" {
		return "fail"
	}
	path := setting.SaveDir + filename
	p, exist := picMap[filename]

	if exist || fileExist(path) {
		fmt.Println("图片[" + filename + "]已存在,位置[" + path + "]")
		log.Info("图片[" + filename + "]已存在,位置[" + path + "]")
		return p
	}
	if len(data) > 10 {
		fmt.Println("[" + filename + "]")
		file, err := os.Create(path)
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
	picMap[filename] = path
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

func getDataFromUrl(url string) ([]byte, string) {

	rsp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()
	filename := analyseContentDisposition(rsp.Header)
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		panic(err)
	}

	return data, filename
}

func analyseContentDisposition(header http.Header) string {
	log.Debug(header)
	disposition := header.Get("Content-Disposition")
	log.Debug(disposition)
	canSplit := func(c rune) bool { return c == ';' || c == '=' || c == ' ' }
	ary := strings.FieldsFunc(disposition, canSplit)
	if len(ary) == 3 && ary[0] == "attachment" && ary[1] == "filename" {
		return ary[2]
	}
	return "fail"
}
func getDataFromUrlAndData(url string) []byte {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Debug(err)
		panic(err)
	}
	client := new(http.Client)
	rsp, err := client.Do(req)
	if err != nil {
		log.Debug(err)
		panic(err)
	}
	defer rsp.Body.Close()
	rspData, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Debug(err)
		panic(err)
	}
	log.Debug("Request: " + url + " ,Rsp: " + string(rspData))
	return rspData
}
func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
