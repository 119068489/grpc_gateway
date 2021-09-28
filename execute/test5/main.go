package main

import (
	"fmt"
	"grpc_gateway/easygo"
	"reflect"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/codegangsta/inject"
)

type SpecialString interface{}
type SpecialInt interface{}

func Say(name string, gender SpecialString, age int, c SpecialInt) {
	fmt.Printf("My name is %s, gender is %s, age is %d,c is %d!\n", name, gender, age, c)
}

func main() {

	type A struct {
		Name string
		Age  int
	}

	a := &A{}

	logs.Debug(reflect.TypeOf(*a).Kind())

	inj := inject.New()
	inj.Map(a)
	inj.Invoke(SetStructValue)

	logs.Info(a)

}

func Ranges() {
	a := []string{"1", "3", "2"}

	for i := range a {
		_, _ = i, a[i]
	}
}

func RangeMap() {
	a := make(map[int]string)
	a[0] = "ok"
	for key, v := range a {
		logs.Debug(v)
		key++
		a[key] = v + easygo.AnytoA(key)
	}
}

func RangeChannle() {
	ch := make(chan int, 5)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			<-time.NewTicker(time.Second).C
			select {
			case j := <-ch:
				if j == 0 {
					wg.Done()
					return
				}
				logs.Info(j)
			}
		}
	}()

	i := 0
	for {
		i++
		ch <- i
		if len(ch) == 5 {
			close(ch)
			break
		}
	}
	wg.Wait()
}

func SetStructValue(A interface{}) {
	typeofA := reflect.TypeOf(A)
	valueofA := reflect.ValueOf(A)
	for i := 0; i < typeofA.Elem().NumField(); i++ {
		if typeofA.Elem().Field(i).Type.Kind() == reflect.Int {
			valueofA.Elem().Field(i).SetInt(30)
		} else {
			valueofA.Elem().Field(i).SetString("bobo")
		}
	}
}
