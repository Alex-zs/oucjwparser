package util

import (
	"bytes"
	"encoding/base64"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

// GBK编码字节数组转换为UTF8编码字符串
func GBKBytes2UTF8(src []byte) []byte {
	data, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(src), simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return make([]byte, 0)
	}
	return data
}

// base64编码
func Base64Encoding(data []byte) string  {
	return base64.StdEncoding.EncodeToString(data)
}
