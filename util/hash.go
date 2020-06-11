package util

import (
	"crypto/md5"
	"fmt"
	"io"
)

func MD5(data string) string {
	hasher := md5.New()
	io.WriteString(hasher, data)
	res := hasher.Sum([]byte{})
	return fmt.Sprintf("%x", res)
}
