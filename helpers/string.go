package helpers

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"regexp"
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

func validMobile(mobileNum string) bool {
	const regular = "^(13[0-9]|14[57]|15[0-35-9]|18[07-9]|166)\\\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func validEmail(email string) bool {
	const regular = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(regular)
	return reg.MatchString(email)
}