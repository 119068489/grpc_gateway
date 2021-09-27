# go 单元测试

单元测试分为基准测试和性能测试：

- 基准测试函数命名规则：TestName(t *testing.T)
- 性能测试函数命名规则：BenchmarkName(b *testing.B)

前缀必须这样命名，Name可以随便写。


## 基准测试

### go test工具
go test是go内置的测试工具，在进行单元测试的时候，必须在要测试的函数目录下创建XXX_test.go文件。

#### 基本用法
写一个判断2个字符串切片是否相等的函数,主逻辑如下：
- 两个字符串切片长度不相等时，返回false
- 两个字符串切片一个是nil，另一个不是nil时，返回false
- 遍历两个切片，比较对应索引的两个切片元素值，如果不相等，返回false
- 否则，返回true
  
``` 创建文件common.go,写下函数代码如下：
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
```

```在common.go文件所在的目录创建common_test.go文件，写测试代码如下：
func TestStringSliceEqual(t *testing.T) { // 测试函数名必须以Test开头，必须接收一个*testing.T类型参数
	a := []string{"hello", "test"}
	b := []string{"hello", "test"}
	got := StringSliceEqual(a, b) // 程序输出的结果
	want := true                  // 期望的结果
	if !got {                     // 因为slice不能比较直接，借助反射包中的方法比较
		t.Errorf("excepted:%v, got:%v", want, got) // 测试失败输出错误提示
	}
}
```

在文件目录下运行命令`go test -v -run=TestStringSliceEqual`，结果如下：
``` 表示测试通过
=== RUN   TestStringSliceEqualt
--- PASS: TestStringSliceEqualt (0.00s)
PASS
ok      grpc_gateway/easygo     0.092s
```

``` 表示测试不通过
=== RUN   TestStringSliceEqualt
    common_test.go:156: excepted:true, got:false
--- FAIL: TestStringSliceEqualt (0.00s)
FAIL
exit status 1
FAIL    grpc_gateway/easygo     0.126s
```


### GoConvey测试框架
GoConvey是一款针对Golang的测试框架，可以管理和运行测试用例，同时提供了丰富的断言函数，并支持很多 Web 界面特性。

#### 安装
`go get github.com/smartystreets/goconvey`

#### 基本使用方法

在函数所在的文件目录下创建一个XXX_test.go测试文件
``` 代码如下
import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

func TestStringSliceEqual(t *testing.T) {
    Convey("TestStringSliceEqual should return true when a != nil  && b != nil", t, func() {
        a := []string{"hello", "goconvey"}
        b := []string{"hello", "goconvey"}
        So(StringSliceEqual(a, b), ShouldBe)
    })
}
```

运行命令`go test -v -run=TestStringSliceEqual`,必须在文件目录下运行命令
``` 执行后的结果
=== RUN   TestStringSliceEqual

  TestStringSliceEqual should return true when a != nil  && b != nil .


1 total assertion

--- PASS: TestStringSliceEqual (0.00s)
PASS
ok      grpc_gateway/easygo     0.095s
```

#### Convey语句的嵌套
``` 代码如下
import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

func TestStringSliceEqual(t *testing.T) {
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
```

执行结果如下：
```
=== RUN   TestStringSliceEquals

  TestStringSliceEqual
    should return true when a != nil  && b != nil .
    should return true when a ＝= nil  && b ＝= nil .
    should return false when a ＝= nil  && b != nil .
    should return false when a != nil  && b != nil .


4 total assertions

--- PASS: TestStringSliceEquals (0.01s)
PASS
ok      grpc_gateway/easygo     0.101s
```


## 性能测试
 与基准测试不一样的地方是测试函数命名，必须Benchmark开头参数类型*testing.B，代码如下：
```
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
```

测试结果如下：
```
goos: windows
goarch: amd64
pkg: grpc_gateway/easygo
cpu: Intel(R) Core(TM) i7-4710MQ CPU @ 2.50GHz
BenchmarkStringSliceEqual-8             92161651                12.10 ns/op            0 B/op          0 allocs/op
BenchmarkStringSliceEqualParallel-8     344429662                3.434 ns/op           0 B/op          0 allocs/op
PASS
ok      grpc_gateway/easygo     2.613s
```

- BenchmarkStringSliceEqual-8：-cpu参数指定，-8表示8个CPU线程执行
- 92161651：表示总共执行了92161651次
- 12.10 ns/op：表示每次执行耗时12.10纳秒
- 0 B/op:表示每次执行分配的内存（字节）
- 0 allocs/op：表示每次执行分配了多少次对象

命令行命令解释
`go test -v -bench=. -cpu=8 -benchtime="3s" -timeout="5s" -benchmem`
* benchmem：输出内存分配统计
* benchtime：指定测试时间
* cpu：指定GOMAXPROCS
* timeout：超市限制

## pprof
go tools继承了pprof，以便进行性能测试并找出瓶颈。

1. 命令行生成测试数据文件
   `go test -bench=BenchmarkStringSliceEqual -cpuprofile cpu.out`

2. 用命令行分析 cpu.out和mem.out cpu和内存
   `go tool pprof -text cpu.out`

    ``` 输出结果
      Type: cpu
      Time: Sep 27, 2021 at 2:12pm (CST)
      Duration: 2.84s, Total samples = 11.08s (390.28%)
      Showing nodes accounting for 11s, 99.28% of 11.08s total
      Dropped 24 nodes (cum <= 0.06s)
            flat  flat%   sum%        cum   cum%
           5.90s 53.25% 53.25%      7.42s 66.97%  grpc_gateway/easygo.StringSliceEqual
           2.58s 23.29% 76.53%      9.58s 86.46%  grpc_gateway/easygo.BenchmarkStringSliceEqualParallel.func1
           1.52s 13.72% 90.25%      1.52s 13.72%  runtime.memequal
           0.55s  4.96% 95.22%      0.55s  4.96%  testing.(*PB).Next (inline)
           0.39s  3.52% 98.74%      1.37s 12.36%  grpc_gateway/easygo.BenchmarkStringSliceEqual
           0.06s  0.54% 99.28%      0.06s  0.54%  runtime.stdcall3
               0     0% 99.28%      0.06s  0.54%  runtime.(*pageAlloc).scavenge
               0     0% 99.28%      0.06s  0.54%  runtime.(*pageAlloc).scavengeOne
               0     0% 99.28%      0.06s  0.54%  runtime.(*pageAlloc).scavengeRangeLocked
               0     0% 99.28%      0.06s  0.54%  runtime.bgscavenge.func2
               0     0% 99.28%      0.06s  0.54%  runtime.mstart
               0     0% 99.28%      0.06s  0.54%  runtime.sysUnused
               0     0% 99.28%      0.06s  0.54%  runtime.systemstack
               0     0% 99.28%      9.58s 86.46%  testing.(*B).RunParallel.func1
               0     0% 99.28%      1.37s 12.36%  testing.(*B).launch
               0     0% 99.28%      1.37s 12.36%  testing.(*B).runN
    ```

   列名	含义
   flat	    本函数的执行耗时
   flat%	flat 占 CPU 总时间的比例。程序总耗时 11.08s, StringSliceEqual 的 5.90s 占了 53.25%
   sum%	    前面每一行的 flat 占比总和
   cum	    累计量。指该函数加上该函数调用的函数总耗时
   cum%	    cum 占 CPU 总时间的比例

3. pprof交互模式分析
   `go tool pprof testTB.test cpu.out`
   交互命令：
   `top` 
   `list StringSliceEqual`(StringSliceEqual函数名) 输出函数消耗时间的代码行详情
   `web` (执行 web 需要安装 graphviz，pprof 能够借助 grapgviz 生成程序的调用图，会生成一个 svg 格式的文件)
   `traces` 可以列出函数的调用栈
   
4. pdf或者svg分析
   `go tool pprof -svg cpu.out > cpu.svg`
   `go tool pprof -pdf cpu.out > cpu.pdf`



## 补充
有时候我们在测试之前需要调用一些外部函数或者设置，比如连接数据库，或者连接外部服务器等等，

这个时候我们需要TestMain函数，相当于钩子，可以在测试之前和之后做一些事情。
```
var c gateway.GatewayClient

func GetClient() *grpc.ClientConn {
	conn, err := grpc.Dial("localhost:9192", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logs.Error("did not connect: %v", err)
	}
	c = gateway.NewGatewayClient(conn)
	return conn
}

func TestMain(m *testing.M) {
	fmt.Println("测试之前的做一些设置,连接rpc服务器")
	conn := GetClient()
	defer conn.Close()
	// 如果 TestMain 使用了 flags，这里应该加上flag.Parse()
	retCode := m.Run() // 执行测试
	fmt.Println("测试之后做一些拆卸工作,关闭rpc服务器连接")
	os.Exit(retCode) // 退出测试
}
```
