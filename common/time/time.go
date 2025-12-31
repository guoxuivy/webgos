package time

import (
	"time"
)

const (
	DateFormat         = "2006-01-02"
	TimeFormat         = "15:04:05"
	DateTimeFormat     = "2006-01-02 15:04:05"
	DateTimeZoneFormat = "2006-01-02 15:04:05 Z07:00"
)

// FormatDate 将 time.Time 格式化为 "YYYY-MM-DD" 格式的字符串
// 参数: t - 要格式化的时间
// 返回值: 格式化后的日期字符串
func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}

// FormatTime 将 time.Time 格式化为 "HH:MM:SS" 格式的字符串
// 参数: t - 要格式化的时间
// 返回值: 格式化后的时间字符串
func FormatTime(t time.Time) string {
	return t.Format(TimeFormat)
}

// FormatDateTime 将 time.Time 格式化为 "YYYY-MM-DD HH:MM:SS" 格式的字符串
// 参数: t - 要格式化的时间
// 返回值: 格式化后的日期时间字符串
func FormatDateTime(t time.Time) string {
	return t.Format(DateTimeFormat)
}

// FormatDateTimeZone 将 time.Time 格式化为带时区的 "YYYY-MM-DD HH:MM:SS ZZZ" 格式的字符串
// 参数: t - 要格式化的时间
// 返回值: 格式化后的带时区的日期时间字符串
func FormatDateTimeZone(t time.Time) string {
	return t.Format(DateTimeZoneFormat)
}

// FormatDateTimeForTimestamp 将时间戳格式化为 "YYYY-MM-DD HH:MM:SS" 格式的字符串
// 参数: timestamp - Unix时间戳
// 返回值: 格式化后的日期时间字符串
func FormatDateTimeForTimestamp(timestamp int64) string {
	return FormatDateTime(time.Unix(timestamp, 0))
}

// ParseDate 解析 "YYYY-MM-DD" 格式的日期字符串为 time.Time
// 参数: dateStr - 日期字符串
// 返回值: 解析后的时间对象和错误（如果解析失败）
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse(DateFormat, dateStr)
}

// ParseTime 解析 "HH:MM:SS" 格式的时间字符串为 time.Time
// 参数: timeStr - 时间字符串
// 返回值: 解析后的时间对象和错误（如果解析失败）
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(TimeFormat, timeStr)
}

// ParseDateTime 解析 "YYYY-MM-DD HH:MM:SS" 格式的日期时间字符串为 time.Time
// 参数: dateTimeStr - 日期时间字符串
// 返回值: 解析后的时间对象和错误（如果解析失败）
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	return time.Parse(DateTimeFormat, dateTimeStr)
}

// ParseDateTimeZone 解析带时区的 "YYYY-MM-DD HH:MM:SS ZZZ" 格式的日期时间字符串为 time.Time
// 参数: dateTimeZoneStr - 带时区的日期时间字符串
// 返回值: 解析后的时间对象和错误（如果解析失败）
func ParseDateTimeZone(dateTimeZoneStr string) (time.Time, error) {
	return time.Parse(DateTimeZoneFormat, dateTimeZoneStr)
}

// ParseLocalDate 解析 "YYYY-MM-DD" 格式的日期字符串为本地时间（上海时区）
// 参数: dateStr - 日期字符串
// 返回值: 解析后的时间对象和错误（如果解析失败）
func ParseLocalDate(dateStr string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(DateFormat, dateStr, loc)
}

// ParseLocalTime 解析 "HH:MM:SS" 格式的时间字符串为本地时间（上海时区）
// 参数: timeStr - 时间字符串
// 返回值: 解析后的时间对象和错误（如果解析失败）
func ParseLocalTime(timeStr string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(TimeFormat, timeStr, loc)
}

// ParseLocalDateTime 解析 "YYYY-MM-DD HH:MM:SS" 格式的日期时间字符串为本地时间（上海时区）
// 参数: dateTimeStr - 日期时间字符串
// 返回值: 解析后的时间对象和错误（如果解析失败）
func ParseLocalDateTime(dateTimeStr string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(DateTimeFormat, dateTimeStr, loc)
}

// ParseLocalDateTimeZone 解析带时区的 "YYYY-MM-DD HH:MM:SS ZZZ" 格式的日期时间字符串为本地时间（上海时区）
// 参数: dateTimeZoneStr - 带时区的日期时间字符串
// 返回值: 解析后的时间对象和错误（如果解析失败）
func ParseLocalDateTimeZone(dateTimeZoneStr string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(DateTimeZoneFormat, dateTimeZoneStr, loc)
}