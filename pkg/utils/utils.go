package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type RandDict string

const (
	LatterDict = "abcdefghijklmnopqrstuvwxyz"
	NumDict    = "0123456789"
	AllDict    = LatterDict + NumDict
)

// RandStr 生成随机字符串
func RandStr(num int, dict RandDict) (str string) {
	if num <= 0 {
		num = 1
	}
	//s := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	for i := 0; i < num; i++ {
		str += string(dict[rand.Intn(len(dict))])
	}
	return
}

// RandNum 随机数
func RandNum(n int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(n)
}

// RandStrPrefix 随机字符串固定前缀
func RandStrPrefix(num int, prefix string, dict RandDict) (str string) {
	return prefix + RandStr(num, dict)
}

// TrinaryOperation 三目运算
func TrinaryOperation[T interface{}](condition bool, resultA, resultB T) T {
	if condition {
		return resultA
	}
	return resultB
}

// MD5 MD5加密
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// MD5AndSalt MD5+Salt加密
func MD5AndSalt(str, salt string) (s string) {
	for i := 0; i < 2; i++ {
		s += MD5(str + salt)
	}
	return
}

// Result 结构体转json字符串
func Result(data any) string {
	marshal, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(marshal)
}

// GetVerErr 处理错误信息
func GetVerErr(err error) string {
	er, ok := err.(validator.ValidationErrors)
	if ok {
		field := er[0].Field()
		if field == er[0].StructField() {
			field = strings.ToLower(field[0:1]) + field[1:]
		}
		switch er[0].Tag() {
		case "required":
			return field + "不能为空"
		case "min":
			if er[0].Type().String() == "string" {
				return field + "不能小于" + er[0].Param() + "位"
			}
			return field + "不能小于" + er[0].Param()
		}
		return field + "错误"
	} else {
		return "参数格式错误"
	}
}

func InArray[T string | int](data T, list []T) bool {
	for _, item := range list {
		if data == item {
			return true
		}
	}
	return false
}

func Map[T any, F any](list []T, handler func(T) (F, bool)) (result []F) {
	for _, item := range list {
		value, ok := handler(item)
		if ok {
			result = append(result, value)
		}
	}
	return result
}

// 验证内网IP
func ValidateLocalIP(ip string) bool {
	{
		// 验证10.0.0.0 - 10.255.255.255
		compile := regexp.MustCompile("^10\\.(25[0-5]|2[0-4]\\d|1\\d{2}|\\d{1,2})\\.(25[0-5]|2[0-4]\\d|1\\d{2}|\\d{1,2})\\.(25[0-5]|2[0-4]\\d|1\\d{2}|\\d{1,2})$")
		if compile.MatchString(ip) {
			return true
		}
	}
	{
		// 验证172.16.0.0 - 172.31.255.255
		compile := regexp.MustCompile("^172\\.(1[6-9]|2\\d|3[01])\\.(25[0-5]|2[0-4]\\d|1\\d{2}|\\d{1,2})\\.(25[0-5]|2[0-4]\\d|1\\d{2}|\\d{1,2})$")
		if compile.MatchString(ip) {
			return true
		}
	}
	{
		// 验证192.168.0.0 - 192.168.255.255
		compile := regexp.MustCompile("^192\\.168\\.(25[0-5]|2[0-4]\\d|1\\d{2}|\\d{1,2})\\.(25[0-5]|2[0-4]\\d|1\\d{2}|\\d{1,2})$")
		if compile.MatchString(ip) {
			return true
		}
	}
	{
		// 验证127.0.0.0 - 127.255.255.255
		compile := regexp.MustCompile("^127\\.(25[0-5]|2[0-4]\\d|1\\d{2}|\\d{1,2})\\.(25[0-5]|2[0-4]\\d|1\\d{2}|\\d{1,2})\\.(25[0-5]|2[0-4]\\d|1\\d{2}|\\d{1,2})$")
		if compile.MatchString(ip) {
			return true
		}
	}
	return false
}

// 验证端口
func ValidatePort(port string) bool {
	compile := regexp.MustCompile("^([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$")
	return compile.MatchString(port)
}

// 验证域名
func ValidateDomain(domain string) bool {
	compile := regexp.MustCompile(`^(?:[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?\.)+[A-Za-z]{2,63}$`)
	return compile.MatchString(domain)
}

func FormatTimes(layout string, times ...string) (list []time.Time, ok bool) {
	for _, item := range times {
		t, err := time.ParseInLocation(layout, item, time.Local)
		if err != nil {
			return nil, false
		}
		list = append(list, t)
	}
	return list, true
}

func MustString2Int(data string) int {
	num, _ := strconv.Atoi(data)
	return num
}
