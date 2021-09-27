package easygo

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 使用go自带的testing单元测试工具==========================================================================
//go test -v -run=TestInterfersToInt64
func TestInterfersToInt64(t *testing.T) { // 测试函数名必须以Test开头，必须接收一个*testing.T类型参数
	got := InterfersToInt64([]interface{}{int64(1), int64(5), int64(6)}) // 程序输出的结果
	want := []int64{1, 5, 6}                                             // 期望的结果
	if !reflect.DeepEqual(want, got) {                                   // 因为slice不能比较直接，借助反射包中的方法比较
		t.Errorf("excepted:%v, got:%v", want, got) // 测试失败输出错误提示
	}
}

////go test -v -run=TestInterfersToInt64s 子测试
func TestInterfersToInt64s(t *testing.T) {
	// 定义一个测试用例类型
	type test struct {
		input []interface{}
		want  []int64
	}
	// 定义一个存储测试用例的切片
	tests := map[string]test{
		"int64": {input: []interface{}{int64(1), int64(5), int64(6)}, want: []int64{1, 5, 6}},
		"int":   {input: []interface{}{int64(4), int64(5), int64(6)}, want: []int64{4, 5, 6}},
		// "string": {input: []interface{}{"1", "5", "6"}, want: []int64{1, 5, 6}},
	}
	// 遍历切片，逐一执行测试用例
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) { // 使用t.Run()执行子测试
			got := InterfersToInt64(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("name:%s,excepted:%v, got:%v", name, tc.want, got)
			}
		})

	}
}

//go test -bench=BenchmarkInterfersToInt64 -benchmem
func BenchmarkInterfersToInt64(b *testing.B) {
	// b.StopTimer() //调用该函数停止压力测试的时间计数

	// //做一些初始化的工作,例如读取文件数据,数据库连接之类的,
	// //这样这些时间不影响我们测试函数本身的性能

	// b.StartTimer() //重新开始时间

	for i := 0; i < b.N; i++ {
		InterfersToInt64([]interface{}{int64(1), int64(5), int64(6)})
	}
}

//并行计算 go test -bench=BenchmarkInterfersToInt64Parallel -benchmem
func BenchmarkInterfersToInt64Parallel(b *testing.B) {
	// b.SetParallelism(4) // 设置使用的CPU数
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			InterfersToInt64([]interface{}{int64(1), int64(5), int64(6)})
		}
	})
}

func benchmarkFib(b *testing.B, n int) {
	for i := 0; i < b.N; i++ {
		Fib(n)
	}
}

//go test -bench=BenchmarkFib
func BenchmarkFib1(b *testing.B)  { benchmarkFib(b, 1) }
func BenchmarkFib2(b *testing.B)  { benchmarkFib(b, 2) }
func BenchmarkFib3(b *testing.B)  { benchmarkFib(b, 3) }
func BenchmarkFib10(b *testing.B) { benchmarkFib(b, 10) }
func BenchmarkFib20(b *testing.B) { benchmarkFib(b, 20) }
func BenchmarkFib40(b *testing.B) { benchmarkFib(b, 40) }

//go test -run ExampleInterfersToInt64
func ExampleInterfersToInt64() {
	fmt.Println(InterfersToInt64([]interface{}{int64(1), int64(5), int64(6)}))
	// Output:
	// [1 5 6]
}

// 使用GoConvey单元测试框架==========================================================================
// go test -v -run=TestConvey
func TestConveyInterfersToInt64(t *testing.T) {
	Convey("InterfersToInt64 should return [1 5 6]", t, func() {
		input := []interface{}{int64(1), int64(5), int64(6)}
		want := []int64{1, 5, 6}

		//自定义匿名assertion函数 判断结果
		So(InterfersToInt64(input), func(actual interface{}, expected ...interface{}) string {
			if !reflect.DeepEqual(expected[0], actual) { // 因为slice不能比较直接，借助反射包中的方法比较
				return fmt.Sprintf("excepted:%v, got:%v", expected[0], actual) // 测试失败输出错误提示
			}
			return ""
		}, want)
	})
}

/* 自定义assertion函数 判断结果
func shouldEqual(actual interface{}, expected ...interface{}) string {
	if !reflect.DeepEqual(expected[0], actual) { // 因为slice不能比较直接，借助反射包中的方法比较
		return fmt.Sprintf("excepted:%v, got:%v", expected[0], actual) // 测试失败输出错误提示
	}
	return ""
}
*/

// 多个测试用例
// go test -v -run=TestConveyFib
func TestConveyFib(t *testing.T) {
	Convey("TestConveyFib", t, func() {

		in := 6
		want := 8
		So(Fib(in), ShouldEqual, want)
	})

	Convey("TestConveyFib", t, func() {

		in := 1
		want := 1
		So(Fib(in), ShouldEqual, want)
	})

	Convey("TestConveyFib", t, func() {
		Convey("should return number", func() {
			in := 1
			want := 1
			So(Fib(in), ShouldEqual, want)
		})

		Convey("should return number 2", func() {
			in := 6
			want := 8
			So(Fib(in), ShouldEqual, want)
		})

	})
}

// go test -v -run=TestStringSliceEqual
func TestStringSliceEqual(t *testing.T) {
	Convey("TestStringSliceEqual should return true when a != nil  && b != nil", t, func() {
		a := []string{"hello", "goconvey"}
		b := []string{"hello", "goconvey"}
		So(StringSliceEqual(a, b), ShouldBeTrue)
	})
}

// go test -v -run=TestStringSliceEquals
func TestStringSliceEquals(t *testing.T) {
	Convey("TestStringSliceEqual", t, func() {
		Convey("should return true when a != nil  && b != nil", func() {
			a := []string{"hello", "goconvey"}
			b := []string{"hello", "goconvey"}
			So(StringSliceEqual(a, b), ShouldBeTrue)
		})

		Convey("should return true when a ＝= nil  && b ＝= nil", func() {
			So(StringSliceEqual(nil, nil), ShouldBeTrue)
		})

		Convey("should return false when a ＝= nil  && b != nil", func() {
			a := []string(nil)
			b := []string{}
			So(StringSliceEqual(a, b), ShouldBeFalse)
		})

		Convey("should return false when a != nil  && b != nil", func() {
			a := []string{"hello", "world"}
			b := []string{"hello", "goconvey"}
			So(StringSliceEqual(a, b), ShouldBeFalse)
		})
	})
}

//go test -bench=BenchmarkStringSliceEqual -benchmem
func BenchmarkStringSliceEqual(b *testing.B) {
	// b.StopTimer() //调用该函数停止压力测试的时间计数

	// //做一些初始化的工作,例如读取文件数据,数据库连接之类的,
	// //这样这些时间不影响我们测试函数本身的性能

	// b.StartTimer() //重新开始时间

	for i := 0; i < b.N; i++ {
		StringSliceEqual([]string{"hello", "goconvey"}, []string{"hello", "goconvey"})
	}
}

//并行计算 go test -bench=BenchmarkInterfersToInt64Parallel -benchmem
func BenchmarkStringSliceEqualParallel(b *testing.B) {
	// b.SetParallelism(4) // 设置使用的CPU数
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			StringSliceEqual([]string{"hello", "goconvey"}, []string{"hello", "goconvey"})
		}
	})
}
