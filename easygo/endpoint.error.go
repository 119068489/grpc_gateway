package easygo

import (
	"errors"
	"fmt"
	"grpc_gateway/easygo/base"
	"strings"
)

//---------------------------------------------------------------------------------------

type IRpcInterrupt interface {
	error
	__IRpcInterrupt__()
	Reason() string
	Code() string
	AddPrefix(prefix interface{})
	AddPostfix(postfix interface{})
}

type RpcInterrupt struct {
	error
	Me      IRpcInterrupt
	failMsg *base.Fail
}

/* 抽象类，不提供实例化方法
func NewRpcInterrupt(text string)*RpcInterrupt{}
*/

func (rSelf *RpcInterrupt) Init(me IRpcInterrupt, methodName string, failMsg *base.Fail) {
	rSelf.Me = me

	var list []string
	if methodName != "" {
		list = append(list, "method: "+methodName)
	}

	code := failMsg.GetCode()
	if code != "" {
		list = append(list, "code: "+code)
	}

	reason := failMsg.GetReason()
	if reason != "" {
		list = append(list, "reason: "+reason)
	}

	text := strings.Join(list, ";")
	rSelf.error = errors.New(text)
	rSelf.failMsg = failMsg
}

func (rSelf *RpcInterrupt) Reason() string {
	return rSelf.failMsg.GetReason() // 没有把 methodName 拼在这里
}

func (rSelf *RpcInterrupt) Code() string {
	return rSelf.failMsg.GetCode()
}

func (rSelf *RpcInterrupt) AddPrefix(prefix interface{}) {
	s := fmt.Sprintf("%v;%v", prefix, rSelf.error)
	rSelf.error = errors.New(s)
}

func (rSelf *RpcInterrupt) AddPostfix(postfix interface{}) {
	s := fmt.Sprintf("%v;%v", rSelf.error, postfix)
	rSelf.error = errors.New(s)
}

func (rSelf *RpcInterrupt) __IRpcInterrupt__() {}

//----------------------------------------------------------------------------------------

type IRpcFail interface {
	IRpcInterrupt
	__IRpcFail__()
}

type RpcFail struct {
	RpcInterrupt
}

func NewRpcFail(methodName string, failMsg *base.Fail) *RpcFail {
	p := &RpcFail{}
	p.Init(p, methodName, failMsg)
	return p
}

func (rSelf *RpcFail) __IRpcFail__() {}

//----------------------------------------------------------------------------------------

type IRpcTimeout interface {
	IRpcInterrupt
	__IRpcTimeout__()
}
type RpcTimeout struct {
	RpcInterrupt
}

func NewRpcTimeout(methodName string, failMsg *base.Fail) *RpcTimeout {
	p := &RpcTimeout{}
	p.Init(p, methodName, failMsg)
	return p
}

func (rSelf *RpcTimeout) __IRpcTimeout__() {}
