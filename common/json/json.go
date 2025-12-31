// Package json 使用开源第三方库json-iterator封装的json api。
// 与标准库api完全一模一样，只需将import路径由encoding/json改成going/json即可。
// 标准包使用了反射来实现，性能极低，使用json-iterator解码能提升5倍性能，编码也比标准包性能好，不过较不明显
package json

import (
	jsoniter "github.com/json-iterator/go"
)

var j = jsoniter.ConfigCompatibleWithStandardLibrary

// Marshal 利用json-iterator进行json编码
func Marshal(v any) ([]byte, error) {
	return j.Marshal(v)
}

// Unmarshal 利用json-iterator进行json解码
func Unmarshal(data []byte, v any) error {
	return j.Unmarshal(data, v)
}
