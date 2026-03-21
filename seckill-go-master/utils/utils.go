package utils

import (
	"strconv"
)

// StringToInt 将字符串转换为整数
func StringToInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

// StringToUint 将字符串转换为uint
func StringToUint(s string) uint {
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint(v)
}
