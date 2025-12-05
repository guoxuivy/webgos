package utils

import "strconv"

// SliceAtoi 将字符串切片转换为整数切片
func SliceAtoi(slice []string) []int {
	result := make([]int, len(slice))
	for i, v := range slice {
		// 忽略转换错误，使用默认值0
		result[i], _ = strconv.Atoi(v)
	}
	return result
}

// 这个文件保留用于其他通用工具函数