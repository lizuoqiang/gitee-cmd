package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"path/filepath"
	"strconv"
	"strings"
)

// StrPad
//
//	@Description: 填充字符串
//	@param str
//	@param length
//	@param padStr
//	@return string
func StrPad(str string, length int, padStr string) string {
	if len(str) >= length {
		return str
	}
	return strings.Repeat(padStr, length-len(str)) + str
}

// Interface2String
//
//	@Description: interface转string
//	@param obj
//	@return string
func Interface2String(obj interface{}) string {
	switch v := obj.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.Itoa(int(v))
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return ""
	}
}

// GenerateRandomString
//
//	@Description: 生成随机字符串
//	@param length
//	@return string
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result string

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			fmt.Println("Error generating random number:", err)
			return ""
		}
		result += string(charset[randomIndex.Int64()])
	}

	return result
}

// StrExplode
//
//	@Description: 字符串切割为切片
//	@param str
//	@param sep
//	@return []string
func StrExplode(str, sep string) []string {
	temp := strings.Split(str, sep)
	result := make([]string, 0)
	for _, v := range temp {
		if v != "" {
			result = append(result, v)
		}
	}

	return result
}

// StringToInt
//
//	@Description: 字符串转整型
//	@param value
//	@return int
func StringToInt(value string) int {
	if value == "" {
		return 0
	}

	num, _ := strconv.Atoi(value)
	return num
}

// GetFileExtension
//
//	@Description: 获取文件后缀
//	@param filename
//	@return string
func GetFileExtension(filename string) string {
	return filepath.Ext(filename)
}

// GetFileNameWithoutExtension
//
//	@Description: 获取文件名(不包含后缀)
//	@param filename
//	@return string
func GetFileNameWithoutExtension(filename string) string {
	// 获取文件名（包含后缀）
	baseName := filepath.Base(filename)
	// 获取文件后缀
	ext := filepath.Ext(baseName)
	// 去掉后缀
	return baseName[:len(baseName)-len(ext)]
}
