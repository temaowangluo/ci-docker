package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type (
	Msg struct {
		Code     int    `json:"code"`
		Message  string `json:"message"`
		CodeDesc string `json:"codeDesc"`
		Data     Data   `json:"data,omitempty"`
	}
	Data struct {
		Count  int    `json:"count"`
		TaskId string `json:"task_id"`
	}
)

var (
	secretId = os.Getenv("COS_SECRET_ID")
	action = flag.String("action", "", "执行的动作")
	addr = flag.String("addr", "", "刷新地址")
	isDirectory = flag.Bool("dir", false, "是否是目录")
)

func main() {
	// 初始化参数
	flag.Parse()

	switch *action {
	case "refresh":
		Refresh()
	default:
		flag.PrintDefaults()
	}

}

func Refresh() {
	if *addr == "" {
		log.Fatal("请设置刷新地址")
	}

	addrResult := *addr
	nonce := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	// 获得当前时间戳
	timestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)
	params := url.Values{}

	if *isDirectory {
		params.Set("Action", "RefreshCdnDir")
		params.Set("dirs.0", addrResult)
	} else {
		params.Set("Action", "RefreshCdnUrl")
		params.Set("urls.0", addrResult)
	}

	params.Set("Timestamp", timestampStr)
	params.Set("Nonce", nonce)
	params.Set("SecretId", secretId)

	// 生成签名
	enCode := Signature(params)
	params.Set("Signature", enCode)
	resp, err := http.PostForm("https://cdn.api.qcloud.com/v2/index.php", params)
	defer func() {
		_ = resp.Body.Close()
	}()
	if err != nil {
		log.Fatal(err)
	}
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	result := Msg{}
	if err = json.Unmarshal(r, &result); err != nil {
		log.Fatal(err)
	}
	if result.Code != 0 {
		log.Fatal(result.Code)
	}
}

func Signature(params url.Values) string {
	secretKey := os.Getenv("COS_SECRET_KEY")
	// 利用key转换map中value为string
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	mapParams := map[string]string{}
	for _, v := range keys {
		mapParams[v] = params[v][0]
	}
	// 把map中的key value转为key=value格式
	var splicing []string
	for k, v := range mapParams {
		result := fmt.Sprintf("%s=%s", k, v)
		splicing = append(splicing, result)
	}
	//根据key进行排序
	sort.Strings(splicing)
	// string拼接&
	var buffer bytes.Buffer
	for _, area := range splicing {
		buffer.WriteString(area)
		buffer.WriteString("&")
	}
	str := strings.TrimRight(buffer.String(), "&")
	// 最终生成请求的URL
	resultUrl := fmt.Sprintf("%s?%s", "POSTcdn.api.qcloud.com/v2/index.php", str)
	// HMAC-SHA1   Base64 进行编码
	enCode := GetSignature(resultUrl, secretKey)
	return enCode
}

func GetSignature(input, key string) string {
	keyForSign := []byte(key)
	h := hmac.New(sha1.New, keyForSign)
	h.Write([]byte(input))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
