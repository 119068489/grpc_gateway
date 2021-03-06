package easygo

import (
	"grpc_gateway/easygo/base"
	"math"
	"math/rand"
)

const Underline = 0

//-----------------------------

var None = struct{}{}

type Pair = [2]string

type RpcMap = map[string]Pair

type DB_NAME = string

//-----------------------------

// 经常需要构造 Fail 消息，封装一下来得快
func NewFailMsg(reason string, codes ...string) *base.Fail {
	msg := &base.Fail{Reason: NewString(reason)}
	if len(codes) != 0 {
		msg.Code = NewString(codes[0])
	}
	return msg
}

//=================================
func Keys(dict interface{}) interface{} {
	// 怎么做?
	return nil
}

func Values(dict interface{}) interface{} {
	// 怎么做?
	return nil
}

//=================================

type KWAT map[string]interface{} //Key Word Arg Type

var KWAO = KWAT{} //避免每次都创建新的 //Key Word Arg Object
//=================================
func (kwatSelf KWAT) Add(key string, val interface{}) {
	kwatSelf[key] = val
}
func (kwatSelf KWAT) Del(key string) {
	delete(kwatSelf, key)
}

func (kwatSelf KWAT) GetString(key string) string {
	val, ok := kwatSelf[key]
	if ok {
		return AnytoA(val)
	}
	return ""
}
func (kwatSelf KWAT) GetInt(key string) int {
	val, ok := kwatSelf[key].(int)
	if ok {
		return val
	}
	return Atoi(kwatSelf.GetString(key))
}
func (kwatSelf KWAT) GetInt32(key string) int32 {
	val, ok := kwatSelf[key].(int32)
	if ok {
		return val
	}
	return AtoInt32(kwatSelf.GetString(key))
}
func (kwatSelf KWAT) GetInt64(key string) int64 {
	val, ok := kwatSelf[key].(int64)
	if ok {
		return val
	}
	return AtoInt64(kwatSelf.GetString(key))
}
func (kwatSelf KWAT) GetFloat32(key string) float32 {
	val, ok := kwatSelf[key].(float32)
	if ok {
		return val
	}
	return 0
}
func (kwatSelf KWAT) GetFloat64(key string) float64 {
	val, ok := kwatSelf[key].(float64)
	if ok {
		return val
	}
	return 0
}
func (kwatSelf KWAT) GetBool(key string) bool {
	val, ok := kwatSelf[key].(bool)
	if ok {
		return val
	}
	return false
}

// 把多个 rpc map 拼成一个
func CombineRpcMap(args ...map[string]Pair) map[string]Pair {
	result := make(map[string]Pair, len(args))
	for _, m := range args {
		for method, pair := range m {
			result[method] = pair
		}
	}
	return result
}

//=================================//
//随机制定字符串
func RandString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//随机制定字符串
func RandStringInt(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

const FLOAT64_EPSINON = 0.000000000000001

func IsEqual(f1, f2 float64) bool {
	return math.Abs(f1-f2) < FLOAT64_EPSINON
}

func IsEqualZero(x float64) bool {
	return math.Abs(x) < FLOAT64_EPSINON
}

//----------------------------------
type NoCopy struct{}

func (*NoCopy) Lock()   {}
func (*NoCopy) Unlock() {}

//---------------------------------------------------------------

// 回滚类
type ScopeGuard struct {
	Function  func()
	Dismissed bool
}

func NewScopeGuard(f func()) *ScopeGuard {
	p := &ScopeGuard{}
	p.Init(f)
	return p
}

func (kwatSelf *ScopeGuard) Init(f func()) {
	kwatSelf.Function = f
}

func (kwatSelf *ScopeGuard) Rollback() {
	if !kwatSelf.Dismissed {
		kwatSelf.Function()
	}
}

func (kwatSelf *ScopeGuard) Dismiss() {
	kwatSelf.Dismissed = true
}

/* 用法举例
func ScopeGuardUsage() {
	guard := NewScopeGuard(func() {这里写要回滚的逻辑})
	defer guard.Rollback()

	// do something here
	// ...
	// 中间要是 return 了，或是 panic 了，就会执行 “回滚逻辑”

	guard.Dismiss() // 成功到这一句了，就不会回滚，表示成功了
}
*/
