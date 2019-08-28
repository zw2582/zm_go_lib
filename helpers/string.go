package helpers

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	rand2 "math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//截取字符串 start 起点下标 length 需要截取的长度
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//截取字符串 start 起点下标 end 终点下标(不包括)
func Substr2(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

//Md5encode md5编码
func Md5encode(src string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(src)))
}

//生成Guid字串
func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5encode(base64.URLEncoding.EncodeToString(b))
}

func OrderNo(prefix string) string {
	r := rand2.New(rand2.NewSource(time.Now().Unix()))
	rn := r.Intn(8999)+1000
	nowstr := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s%s%s",prefix, nowstr, rn)
}

//解析gbk
func DecodeGBK(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func ValidMobile(mobileNum string) bool {
	const regular = `^((\+?86)|(\(\+86\)))?((((13[^4]{1})|(14[5-9]{1})|147|(15[^4]{1})|166|(17\d{1})|(18\d{1})|(19[89]{1}))\d{8})|((134[^9]{1}|1410|1440)\d{7}))$`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func ValidEmail(email string) bool {
	const regular = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(regular)
	return reg.MatchString(email)
}

//ValidContainChinese 包含中文检测
func ValidContainChinese(str string) bool {
	const regular = `[^\x00-\x80]+`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(str)
}

//ValidChineName 验证中文姓名
func ValidChineName(str string) bool {
	const regular = `^[\u4E00-\u9FA5]{2,10}$`
	reg := regexp.MustCompile(regular)
	return reg.MatchString(str)
}

func Sha1Encode(raw string) string {
	b := sha1.Sum([]byte(raw))
	return base64.StdEncoding.EncodeToString(b[:])
}

func InetAtoN(ip string) int64 {
	bits := strings.Split(ip, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}