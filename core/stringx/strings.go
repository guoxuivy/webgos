package stringx

import "strconv"

func S2Int(v string) int {
	i, _ := strconv.Atoi(v)
	return i
}
func S2Int64(v string) int64 {
	i, _ := strconv.ParseInt(v, 10, 64)
	return i
}
func S2Float64(v string) float64 {
	f, _ := strconv.ParseFloat(v, 64)
	return f
}

// SliceAtoi 将字符串切片转换为整数切片
func SliceAtoi(slice []string) []int {
	result := make([]int, len(slice))
	for i, v := range slice {
		// 忽略转换错误，使用默认值0
		result[i], _ = strconv.Atoi(v)
	}
	return result
}
