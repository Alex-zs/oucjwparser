package util

import (
	"bytes"
	"encoding/base64"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"strconv"
	"strings"
)

// GBK编码字节数组转换为UTF8编码字符串
func GBK2UTF8(src []byte) []byte {
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

// 加密教务系统请求参数
func EncParamStr(param, key, time string ) map[string]string {
	token := TripleMD5(param, time)
	param = Base64Encoding([]byte(encodeDes(param, key)))
	paramMap := map[string]string{
		"params": param,
		"token": token,
		"timestamp": time,
	}
	return paramMap
}

func encodeDes(str, key string) string {
	lenStr := len(str)
	var encData bytes.Buffer
	var keyBt = make([][]int, 0)
	var length int
	if key != "" {
		keyBt = getKeyBytes(key)
		length = len(keyBt)
		iterator := lenStr / 4
		remainder := lenStr % 4
		for i := 0; i < iterator; i++ {
			tempData := str[i * 4 : i * 4 + 4]
			tempByte := strToBt(tempData)
			var encByte []int
			tempBt := tempByte
			for x := 0; x < length; x++ {
				tempBt = enc(tempBt, keyBt[x])
			}
			encByte = tempBt
			encData.WriteString(bt64ToHex(encByte))
		}
		if remainder > 0 {
			remainderData := str[iterator * 4 : lenStr]
			tempByte := strToBt(remainderData)
			var encByte []int
			var tempBt = tempByte
			for x := 0; x < length; x++ {
				tempBt = enc(tempBt, keyBt[x])
			}
			encByte = tempBt
			encData.WriteString(bt64ToHex(encByte))
		}
	}
	return encData.String()
}

func getKeyBytes(key string) [][]int {
	keyBytes := make([][]int, 0)
	len := len(key)
	iterator := len / 4
	remainder := len % 4
	i := 0
	for ; i < iterator; i++ {
		keyBytes = append(keyBytes, strToBt(key[i * 4 : i * 4 + 4]))
	}
	if remainder  > 0 {
		keyBytes = append(keyBytes, strToBt(key[i * 4: len]) )
	}

	return keyBytes
}

func strToBt(str string) []int {
	len := len(str)
	var bt = make([]int, 64)
	i := 0
	for i < len && i < 4 {
		k := int(str[i])
		fullBit(bt, i, k)
		i++
	}

	for p := len; p <= 3; p++ {
		k := 0
		fullBit(bt, p, k)
	}
	return bt
}

func fullBit(bt []int, p, k int)  {
	for q := 0; q <= 15; q++ {
		pow := 1
		for m := 15; m >= q+1; m-- {
			pow *= 2
		}
		bt[16 * p + q] = k / pow % 2
	}
}

func enc(dataByte , keyByte []int) []int {
	var index = make([]int, 16)
	for i := 0; i < 16; i++ {
		index[i] = i
	}
	return transferData(dataByte, keyByte, index)
}

func transferData(dataByte, keyByte, index []int) []int{
	keys := generateKeys(keyByte)
	ipByte := initPermute(dataByte)
	var ipLeft = make([]int, 32)
	var ipRight = make([]int, 32)
	var tempLeft = make([]int, 32)
	for k := 0; k < 32; k++ {
		ipLeft[k] = ipByte[k]
		ipRight[k] = ipByte[32 + k]
	}
	for _, i := range index {
		for j := 0; j < 32; j++ {
			tempLeft[j] = ipLeft[j]
			ipLeft[j] = ipRight[j]
		}
		var key = make([]int, 48)
		for index := 0; index < 48; index++ {
			key[index] = keys[i][index]
		}
		tempRight :=  xor(pPermute(sBoxPermute(xor(expandPermute(ipRight), key))), tempLeft)
		for i := 0; i < 32; i++ {
			ipRight[i] = tempRight[i]
		}
	}

	var finalData = make([]int, 64)
	for i := 0; i < 32; i++ {
		finalData[i] = ipRight[i]
		finalData[32 + i] = ipLeft[i]
	}
	return finallyPermute(finalData)
}

func generateKeys(keyByte []int) [][]int  {
	var e [56]int
	keys := make([][]int, 16)
	for i := 0; i < 16; i++ {
		keys[i] = make([]int, 48)
	}

	loop := []int{1, 1, 2, 2, 2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 1}
	for i := 0; i < 7; i++ {
		j := 0
		k := 7
		for j < 8 {
			e[i * 8 + j] = keyByte[8 * k + i]
			j++
			k--
		}
	}

	for i := 0; i < 16; i++ {
		var tempLeft int
		var tempRight int
		for j := 0; j < loop[i]; j++ {
			tempLeft = e[0]
			tempRight = e[28]
			for k := 0; k < 27; k++ {
				e[k] = e[k + 1]
				e[28 + k] = e[29 + k]
			}
			e[27] = tempLeft
			e[55] = tempRight
		}
		t := []int{e[13], e[16], e[10], e[23], e[0], e[4], e[2], e[27], e[14], e[5], e[20], e[9], e[22], e[18], e[11], e[3], e[25], e[7], e[15], e[6], e[26], e[19], e[12], e[1], e[40], e[51], e[30], e[36], e[46], e[54], e[29], e[39], e[50], e[44], e[32], e[47], e[43], e[48], e[38], e[55], e[33], e[52], e[45], e[41], e[49], e[35], e[28], e[31]}
		for index := 0; index < 48; index++ {
			keys[i][index] = t[index]
		}
	}
	return keys
}

func initPermute(originalData []int)[64]int  {
	var ipByte [64]int
	i := 0
	m := 1
	n := 0
	var j int
	var k int
	for i < 4 {
		j = 7
		k = 0
		for j >= 0 {
			ipByte[i * 8 + k] = originalData[j * 8 + m]
			ipByte[i * 8 + k + 32] = originalData[j * 8 + n]
			j--
			k++
		}
		i++
		m += 2
		n += 2
	}
	return ipByte
}

func expandPermute(rightData []int)[]int  {
	var epByte = make([]int, 48)
	for i := 0; i < 8; i++ {
		if i == 0 {
			epByte[i * 6] = rightData[31]
		}else {
			epByte[i * 6] = rightData[i * 4 - 1]
		}
		epByte[i * 6 + 1] = rightData[i * 4]
		epByte[i * 6 + 2] = rightData[i * 4 + 1]
		epByte[i * 6 + 3] = rightData[i * 4 + 2]
		epByte[i * 6 + 4] = rightData[i * 4 + 3]
		if i == 7 {
			epByte[i * 6 + 5] = rightData[0]
		}else {
			epByte[i * 6 + 5] = rightData[i * 4 + 4]
		}
	}
	return epByte
}

func xor(byteOne, byteTwo []int)[]int  {
	var xorByte = make([]int, len(byteOne))
	for i := range byteOne {
		xorByte[i] = byteOne[i] ^ byteTwo[i]
	}
	return xorByte
}

func sBoxPermute(expandByte []int) []int {
	var sBoxByte = make([]int, 32)
	binary := ""
	s1 := [][]int{{14, 4, 13, 1, 2, 15, 11, 8, 3, 10, 6, 12, 5, 9, 0, 7}, {0, 15, 7, 4, 14, 2, 13, 1, 10, 6, 12, 11, 9, 5, 3, 8}, {4, 1, 14, 8, 13, 6, 2, 11, 15, 12, 9, 7, 3, 10, 5, 0}, {15, 12, 8, 2, 4, 9, 1, 7, 5, 11, 3, 14, 10, 0, 6, 13}}
	s2 := [][]int{{15, 1, 8, 14, 6, 11, 3, 4, 9, 7, 2, 13, 12, 0, 5, 10}, {3, 13, 4, 7, 15, 2, 8, 14, 12, 0, 1, 10, 6, 9, 11, 5}, {0, 14, 7, 11, 10, 4, 13, 1, 5, 8, 12, 6, 9, 3, 2, 15}, {13, 8, 10, 1, 3, 15, 4, 2, 11, 6, 7, 12, 0, 5, 14, 9}}
	s3 := [][]int{{10, 0, 9, 14, 6, 3, 15, 5, 1, 13, 12, 7, 11, 4, 2, 8}, {13, 7, 0, 9, 3, 4, 6, 10, 2, 8, 5, 14, 12, 11, 15, 1}, {13, 6, 4, 9, 8, 15, 3, 0, 11, 1, 2, 12, 5, 10, 14, 7}, {1, 10, 13, 0, 6, 9, 8, 7, 4, 15, 14, 3, 11, 5, 2, 12}}
	s4 := [][]int{{7, 13, 14, 3, 0, 6, 9, 10, 1, 2, 8, 5, 11, 12, 4, 15}, {13, 8, 11, 5, 6, 15, 0, 3, 4, 7, 2, 12, 1, 10, 14, 9}, {10, 6, 9, 0, 12, 11, 7, 13, 15, 1, 3, 14, 5, 2, 8, 4}, {3, 15, 0, 6, 10, 1, 13, 8, 9, 4, 5, 11, 12, 7, 2, 14}}
	s5 := [][]int{{2, 12, 4, 1, 7, 10, 11, 6, 8, 5, 3, 15, 13, 0, 14, 9}, {14, 11, 2, 12, 4, 7, 13, 1, 5, 0, 15, 10, 3, 9, 8, 6}, {4, 2, 1, 11, 10, 13, 7, 8, 15, 9, 12, 5, 6, 3, 0, 14}, {11, 8, 12, 7, 1, 14, 2, 13, 6, 15, 0, 9, 10, 4, 5, 3}}
	s6 := [][]int{{12, 1, 10, 15, 9, 2, 6, 8, 0, 13, 3, 4, 14, 7, 5, 11}, {10, 15, 4, 2, 7, 12, 9, 5, 6, 1, 13, 14, 0, 11, 3, 8}, {9, 14, 15, 5, 2, 8, 12, 3, 7, 0, 4, 10, 1, 13, 11, 6}, {4, 3, 2, 12, 9, 5, 15, 10, 11, 14, 1, 7, 6, 0, 8, 13}}
	s7 := [][]int{{4, 11, 2, 14, 15, 0, 8, 13, 3, 12, 9, 7, 5, 10, 6, 1}, {13, 0, 11, 7, 4, 9, 1, 10, 14, 3, 5, 12, 2, 15, 8, 6}, {1, 4, 11, 13, 12, 3, 7, 14, 10, 15, 6, 8, 0, 5, 9, 2}, {6, 11, 13, 8, 1, 4, 10, 7, 9, 5, 0, 15, 14, 2, 3, 12}}
	s8 := [][]int{{13, 2, 8, 4, 6, 15, 11, 1, 10, 9, 3, 14, 5, 0, 12, 7}, {1, 15, 13, 8, 10, 3, 7, 4, 12, 5, 6, 11, 0, 14, 9, 2}, {7, 11, 4, 1, 9, 12, 14, 2, 0, 6, 10, 13, 15, 3, 5, 8}, {2, 1, 14, 7, 4, 10, 8, 13, 15, 12, 9, 0, 3, 5, 6, 11}}
	for m := 0; m < 8; m++ {
		i := expandByte[m * 6] * 2 + expandByte[m * 6 + 5]
		j := expandByte[m * 6 + 1] * 2 * 2 * 2 + expandByte[m * 6 + 2] * 2 * 2 + expandByte[m * 6 + 3] * 2 + expandByte[m * 6 + 4]
		switch m {
		case 0:binary = getBoxBinary(s1[i][j])
		case 1 : binary = getBoxBinary(s2[i][j])
		case 2 : binary = getBoxBinary(s3[i][j])
		case 3 : binary = getBoxBinary(s4[i][j])
		case 4 : binary = getBoxBinary(s5[i][j])
		case 5 : binary = getBoxBinary(s6[i][j])
		case 6 : binary = getBoxBinary(s7[i][j])
		case 7 : binary = getBoxBinary(s8[i][j])
		}
		sBoxByte[m * 4], _ = strconv.Atoi(binary[0:1])
		sBoxByte[m * 4 + 1], _ = strconv.Atoi(binary[1: 2])
		sBoxByte[m * 4 + 2], _ = strconv.Atoi(binary[2: 3])
		sBoxByte[m * 4 + 3], _ = strconv.Atoi(binary[3: 4])
	}
	return sBoxByte
}

var binaryArray = []string{"0000", "0001", "0010", "0011", "0100", "0101", "0110", "0111", "1000", "1001", "1010", "1011", "1100", "1101", "1110", "1111"}
func getBoxBinary(i int) string  {
	binary := ""
	if i > -1 && i < 16 {
		binary = binaryArray[i]
	}
	return binary
}

func pPermute(e []int)[]int {
	return []int{e[15], e[6], e[19], e[20], e[28], e[11], e[27], e[16], e[0], e[14], e[22], e[25], e[4], e[17], e[30], e[9], e[1], e[7], e[23], e[13], e[31], e[26], e[2], e[8], e[18], e[12], e[29], e[5], e[21], e[10], e[3], e[24]}
}

func finallyPermute(e []int)[]int {
	return []int{ e[39], e[7], e[47], e[15], e[55], e[23], e[63], e[31], e[38], e[6], e[46], e[14], e[54], e[22], e[62], e[30], e[37], e[5], e[45], e[13], e[53], e[21], e[61], e[29], e[36], e[4], e[44], e[12], e[52], e[20], e[60], e[28], e[35], e[3], e[43], e[11], e[51], e[19], e[59], e[27], e[34], e[2], e[42], e[10], e[50], e[18], e[58], e[26], e[33], e[1], e[41], e[9], e[49], e[17], e[57], e[25], e[32], e[0], e[40], e[8], e[48], e[16], e[56], e[24]}
}

func bt64ToHex(byteData []int) string {
	var hex bytes.Buffer
	for i := 0; i < 16; i++ {
		var bt bytes.Buffer
		for j := 0; j < 4; j++ {
			bt.WriteString(strconv.Itoa(byteData[i * 4 + j]))
		}
		hex.WriteString(bt4ToHex(bt.String()))
	}
	return hex.String()
}

func bt4ToHex(binary string) string {
	hex := ""
	i, _ := strconv.ParseInt(binary, 2, 64)
	if i > -1 && i < 16 {
		hex = strconv.FormatInt(i, 16)
	}
	return strings.ToUpper(hex)
}