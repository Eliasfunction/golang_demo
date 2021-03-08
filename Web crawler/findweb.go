package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

//Regular Expression
var (
	reMail = `\w+@\w+\.\w+`
	relink = `href="(https?://[\s\S]+?)"`
	reimg  = `"https?://[^"]+?(\.((jpg)|(png)|(jpng)|(gif)|(bmp)))"`
)

//Getmail 找尋信箱
func Getmail(url string) {
	pageStr := GetPageStr(url)
	//過濾內容
	re := regexp.MustCompile(reMail)
	results := re.FindAllStringSubmatch(pageStr, -1)
	//fmt.Println(result)
	for _, result := range results {
		fmt.Println("email:", result)
	}

}

//Getlink 連結類
func Getlink(url string) {
	pageStr := GetPageStr(url)
	//過濾內容
	re := regexp.MustCompile(relink)
	results := re.FindAllStringSubmatch(pageStr, -1)

	for _, result := range results {
		fmt.Println(result[1])
	}
}

//Getimg 圖片類
func Getimg(url string) {
	pageStr := GetPageStr(url)
	//過濾內容
	re := regexp.MustCompile(reimg)
	results := re.FindAllStringSubmatch(pageStr, -1)

	for _, result := range results {
		fmt.Println(result[0])
	}
}

//GetPageStr 讀取網頁訊息
func GetPageStr(url string) (pageStr string) {
	resp, err := http.Get(url)
	HandleError(err, "http.Get url")
	defer resp.Body.Close()
	//讀取網頁內容
	pageBytes, err := ioutil.ReadAll(resp.Body)
	HandleError(err, "ioutil.ReadAll")
	//字節轉字串
	pageStr = string(pageBytes)
	return pageStr
}

//HandleError 處理異常
func HandleError(err error, why string) {
	if err != nil {
		fmt.Println(why, err)
	}
}

func main() {
	url := "https://maiimage.com/furniture-kao/"
	//Getmail(url)
	//Getlink(url)
	Getimg(url)
}
