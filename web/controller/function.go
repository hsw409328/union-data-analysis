package controller

import (
	"errors"
	"github.com/hsw409328/gofunc/go_hlog"
	"gofunc"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	lg *go_hlog.Logger
)

func init() {
	lg = go_hlog.GetInstance("")
}

func SendMobile(mobileString string, authCode string) error {
	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		"https://api.netease.im/sms/sendcode.action",
		strings.NewReader("mobile="+mobileString+"&authCode="+authCode),
	)
	if err != nil {
		lg.Error(err)
	}
	currentTime := gofunc.InterfaceToString(time.Now().Unix())
	randStr := GetRandomSalt()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("AppKey", "0ce2da85ef52bcc4160414446104e49c")
	req.Header.Set("CurTime", currentTime)
	req.Header.Set("CheckSum", gofunc.Sha1Encrypt("b09cf0f8197f"+randStr+currentTime))
	req.Header.Set("Nonce", randStr)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		lg.Error(err)
	}
	r, err := gofunc.StringToMap(string(body))
	if err != nil {
		lg.Error(err)
	}
	if gofunc.InterfaceToString(r["code"]) == "200" {
		return nil
	}
	return errors.New(string(body))
}

// return len=8  salt
func GetRandomSalt() string {
	return GetRandomString(6)
}

//生成随机字符串
func GetRandomString(l int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
