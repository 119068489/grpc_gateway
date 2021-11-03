package easygo

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"runtime"
	"time"

	"github.com/garyburd/redigo/redis"
)

func ProtectRun(entry func()) {
	// 延迟处理的函数
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		switch err.(type) {
		case runtime.Error: // 运行时错误
			fmt.Println("runtime error:", err)
		default: // 非运行时错误
			if err != nil {
				fmt.Println("error:", err)
			}
		}
	}()
	entry()
}

func GetMillSecond() int64 {
	return time.Now().UnixNano() / 1e6
}

//两个相似字段的结构体，相同字段值数据互转
func StructToOtherStruct(src interface{}, dest interface{}) {
	js, err := json.Marshal(src)
	PanicError(err)
	err = json.Unmarshal(js, dest)
	PanicError(err)
}
func StructToMap(src interface{}, dest interface{}) {
	js, err := json.Marshal(src)
	PanicError(err)
	err = json.Unmarshal(js, dest)
	PanicError(err)
}

//[]uint8数组转int64
func InterfersToInt64s(src []interface{}, dest *[]int64) {
	for i := range src {
		v := string(src[i].([]uint8))
		*dest = append(*dest, AtoInt64(v))
	}
}

func InterfersToInt64(src []interface{}) []int64 {
	dest := make([]int64, 0, len(src))
	for _, val := range src {
		if v, ok := val.(int64); ok {
			dest = append(dest, v)
		}
	}
	return dest
}

//[]uint8数组转int32
func InterfersToInt32s(src []interface{}, dest *[]int32) {
	for i := range src {
		v := string(src[i].([]uint8))
		*dest = append(*dest, AtoInt32(v))
	}
}

//[]uint8数组转[]string
func InterfersToStrings(src []interface{}, dest *[]string) {
	for i := range src {
		v := string(src[i].([]uint8))
		*dest = append(*dest, v)
	}
}

func Int64StringMap(result interface{}, err error) (map[int64]string, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: StringMap expects even number of values result")
	}
	m := make(map[int64]string, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, okKey := values[i].([]byte)
		value, okValue := values[i+1].([]byte)
		if !okKey || !okValue {
			return nil, errors.New("redigo: StringMap key not a bulk string value")
		}
		k := string(key)
		k1 := AtoInt64(k)
		m[k1] = string(value)
	}
	return m, nil
}
func ObjListExistStrKey(result interface{}, err error, ikey string) (bool, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return false, err
	}
	if len(values)%2 != 0 {
		return false, errors.New("redigo: StringMap expects even number of values result")
	}

	for i := 0; i < len(values); i += 2 {
		key, okKey := values[i].([]byte)
		_, okValue := values[i+1].([]byte)
		if !okKey || !okValue {
			return false, errors.New("redigo: StringMap key not a bulk string value")
		}
		k := string(key)
		if k == ikey {
			return true, err
		}
	}
	return false, err
}
func ObjListToStrKeyList(result interface{}, err error) ([]string, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: StringMap expects even number of values result")
	}
	lst := []string{}
	for i := 0; i < len(values); i += 2 {
		key, okKey := values[i].([]byte)
		_, okValue := values[i+1].([]byte)
		if !okKey || !okValue {
			return nil, errors.New("redigo: StringMap key not a bulk string value")
		}
		k := string(key)
		lst = append(lst, k)
	}
	return lst, err
}

func StrkeyStringMap(result interface{}, err error) (map[string]string, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: StringMap expects even number of values result")
	}
	m := make(map[string]string, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, okKey := values[i].([]byte)
		value, okValue := values[i+1].([]byte)
		if !okKey || !okValue {
			return nil, errors.New("redigo: StringMap key not a bulk string value")
		}
		k := string(key)
		m[k] = string(value)
	}
	return m, nil
}

//redis obj 转int64 数组
func ObjInt64List(result interface{}, err error) ([]int64, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: StringMap expects even number of values result")
	}
	list := []int64{}
	for i := 0; i < len(values); i += 2 {
		key, okKey := values[i].([]byte)
		_, okValue := values[i+1].([]byte)
		if !okKey || !okValue {
			return nil, errors.New("redigo: StringMap key not a bulk string value")
		}
		k := string(key)
		k1 := AtoInt64(k)
		//m[k1] = string(value)
		list = append(list, k1)
	}
	return list, nil
}

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(RandInt(65, 90))
	}
	return string(bytes)
}

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// Fib 是一个计算第n个斐波那契数的函数
func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

//判断2个字符串切片是否相等
func StringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

//高效拼接字符串
func BuilderConcat(str ...string) string {
	var long int
	for _, v := range str {
		long += len(v)
	}
	var builder strings.Builder
	builder.Grow(long)
	for _, v := range str {
		builder.WriteString(v)
	}

	return builder.String()
}

//模拟打字效果输出
func PrintMsg(f interface{}, v ...interface{}) {
	str := formatLog(f, v)
	for i := range str {
		<-time.NewTicker(time.Second / 100).C
		fmt.Print(string(str[i]))
	}
	fmt.Println()
}

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch d := f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(d)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}

//获取hmac签名
func GetHmacToken(key string, data interface{}) string {
	dd, err := json.Marshal(data)
	if err != nil {
		fmt.Print(err)
	}
	hmac := hmac.New(md5.New, []byte(key))
	hmac.Write(dd)

	return hex.EncodeToString(hmac.Sum([]byte(key)))
}
