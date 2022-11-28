package main

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

const (
	magicVar = "dbbc3dd73364b4084c3a69346e0ce2b2"
)

var (
	sortArr = []string{
		"api_key",
		"app_version",
		"content",
		"created_at",
		"email",
		"file_ids",
		"password",
		"platform",
		"source",
		"timestamp",
		"tz",
		"webp",
	}
)

type BaseParm map[string]any

// 添加参数
func (bs BaseParm) Add(key string, val any) {
	bs[key] = val
}

func NewBaParm() BaseParm {
	bs := make(BaseParm)
	bs.Add("api_key", "flomo_web")
	bs.Add("app_version", "1.22.112")
	bs.Add("platform", "mac")
	bs.Add("webp", "1")
	bs.Add("timestamp", time.Now().Unix())
	return bs
}

// 解析
func (bs BaseParm) parse(raw any) string {
	switch raw.(type) {
	case string:
		return raw.(string)
	case int64:
		return strconv.FormatInt(raw.(int64), 10)
	default:
		return ""
	}
}

// 对所有参数进行排序，方便后面求sign
func (bs BaseParm) sortStr() string {
	builder := strings.Builder{}

	for i, v := range sortArr {
		if val, ok := bs[v]; ok {
			builder.WriteString(v)
			builder.WriteString("=")
			builder.WriteString(bs.parse(val))
			if i != len(sortArr)-1 {
				builder.WriteString("&")
			}
		}
	}

	return builder.String()
}

// GetTime 获取此parm中的时间
func (bs BaseParm) GetTime() int64 {
	i, ok := bs["timestamp"].(int64)
	if !ok {
		return 0
	}

	return i
}

// Encode 将parm编码为符合url中的query参数
func (bs BaseParm) Encode() string {
	if bs == nil {
		return ""
	}
	var buf strings.Builder
	i := 0
	for k, v := range bs {
		i++
		buf.Write([]byte(k))
		buf.WriteByte('=')
		buf.WriteString(bs.parse(v))
		if i != len(bs) {
			buf.WriteByte('&')
		}
	}

	return buf.String()
}

// GenSign 生成sign并添加到parm中
func (bs BaseParm) GenSign() BaseParm {
	h := md5.New()
	h.Write([]byte(bs.sortStr() + magicVar))
	str := hex.EncodeToString(h.Sum(nil))
	bs.Add("sign", str)

	return bs
}
