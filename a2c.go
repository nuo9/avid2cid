package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"log"
)

const (
	usageInfo string = "usage:\na2c 10492\na2c 10492 10493 10495\n"
	videoUrl  string = "http://www.bilibili.com/video/av%d"
	cidEquals string = "cid="
)

var (
	cidEqualsBytes       = []byte(cidEquals)
	cidEqualsBytesLength = len(cidEqualsBytes)
)

func main() {
	argLen := len(os.Args)
	if argLen < 1 {
		log.Fatal(usageInfo)
	}

	avids := make([]int32, 0, argLen)
	for _, v := range os.Args {
		iv, e := strconv.Atoi(v)
		if e != nil || int32(iv) <= 0 {
			log.Fatal(usageInfo)
		}
		avids = append(avids, int32(iv))
	}

	for _, v := range avids {
		cid, e := GetFromWeb(v)
		if e != nil {
			log.Println(e)
		} else {
			log.Println(cid)
		}
	}
}

func GetFromWeb(avid int32) (int32, error) {
	url := fmt.Sprintf(videoUrl, avid)
	r, e := http.Get(url)
	if e != nil {
		return 0, e
	}
	b, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return 0, e
	}
	r.Body.Close()

	cid, e := getCidFromHtml(&b)
	if e != nil {
		return 0, e
	}
	return cid, nil
}

func getCidFromHtml(htmlBytes *[]byte) (int32, error) {
	i := bytes.Index(*htmlBytes, cidEqualsBytes)
	if i == -1 {
		return 0, errors.New("not contains cid")
	}
	i += cidEqualsBytesLength

	numbers := make([]byte, 0, 10)
	tempC := (*htmlBytes)[i : i+1][0]
	for isNumber(tempC) {
		numbers = append(numbers, tempC)
		i++
		tempC = (*htmlBytes)[i : i+1][0]
	}
	i, e := strconv.Atoi(string(numbers))
	if e != nil {
		return 0, e
	}

	return int32(i), nil
}

func isNumber(c byte) bool {
	return c >= '0' && c <= '9'
}
