package main

import (
	"fmt"
	"grpc_gateway/easygo"
	"reflect"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
)

type SpecialString interface{}
type SpecialInt interface{}

func main() {
	easygo.PrintMsg("asdddddddaaaaaaaasefefefefefeesdfafsafaf")
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

func PrintCat() {
	dd := []string{"cat", "dog", "pig"}
	s := len(dd)

	chMap := make([]chan int, s)
	for i := range dd {
		chMap[i] = make(chan int)
	}
	var sw sync.WaitGroup
	for i := range dd {
		go func(i int) {
			for {
				sw.Add(1)
				c := <-chMap[i]
				if c == 0 {
					if i < s-1 {
						close(chMap[i+1])
					} else {
						close(chMap[0])
					}
					sw.Done()
					return
				}

				logs.Info(c, dd[i])
				c--
				if i < s-1 {
					chMap[i+1] <- c
				} else {
					chMap[0] <- c
				}
				sw.Done()
			}
		}(i)
	}
	chMap[0] <- 100
	sw.Wait()

	logs.Info("%d秒后退出", s)
	for i := range dd {
		logs.Info(s - i)
		<-time.NewTicker(time.Second).C
	}
}

func do(taskCh chan int) {
	for {
		select {
		case t := <-taskCh:
			time.Sleep(time.Millisecond)
			fmt.Printf("task %d is done\n", t)
		default:
			return
		}
	}
}

func sendTasks() {
	taskCh := make(chan int, 10)
	go do(taskCh)
	for i := 0; i < 1000; i++ {
		taskCh <- i
	}
}
