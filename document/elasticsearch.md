# ElasticSearch

## 介绍
Elasticsearch（ES）是一个基于Lucene构建的开源、分布式、RESTful接口的全文搜索引擎。Elasticsearch还是一个分布式文档数据库，其中每个字段均可被索引，而且每个字段的数据均可被搜索，ES能够横向扩展至数以百计的服务器存储以及处理PB级的数据。可以在极短的时间内存储、搜索和分析大量的数据。通常作为具有复杂搜索场景情况下的核心发动机。

## 基本概念
ES概念	                                            关系型数据库
Index（索引）支持全文检索	                         Database（数据库）
Type（类型）	                                    Table（表）
Document（文档），不同文档可以有不同的字段集合         Row（数据行）
Field（字段）	                                    Column（数据列）
Mapping（映射）	                                    Schema（模式）

## 操作
引用第三方库来连接ES并进行操作。
`go get github.com/olivere/elastic/v7`

封装工具包文件 manager_elasticsearch.go
``` 
package easygo

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/astaxie/beego/logs"
	"github.com/olivere/elastic/v7"
)

// 基本概念对比关系型数据库
// ES概念	                                         关系型数据库
// Index（索引）支持全文检索	                      Database（数据库）
// Type（类型）	                                     Table（表）
// Document（文档），不同文档可以有不同的字段集合		Row（数据行）
// Field（字段）									 Column（数据列）
// Mapping（映射）	                                 Schema（模式）

//连接管理
type ElasticManager struct {
	Host   string // "http://127.0.0.1:9200/"
	Client *elastic.Client
}

func NewElasticManager(host string) *ElasticManager { // services map[string]interface{},
	p := &ElasticManager{}
	p.Init(host)
	return p
}

//初始化
func (c *ElasticManager) Init(host string) {
	c.Host = host
	errorlog := log.New(os.Stdout, "Elastic", log.LstdFlags)
	var err error
	c.Client, err = elastic.NewClient(elastic.SetErrorLog(errorlog), elastic.SetURL(host))
	if err != nil {
		logs.Error(err)
	}
}

func (c *ElasticManager) GetVersion(host string) {
	info, code, err := c.Client.Ping(host).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	esversion, err := c.Client.ElasticsearchVersion(host)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)
}

/*下面是简单的CURD*/
type Employee struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Age       int      `json:"age"`
	About     string   `json:"about"`
	Interests []string `json:"interests"`
}

//创建
func (c *ElasticManager) Create(index, id string, data interface{}) {
	logs.Info("Create")
	//使用结构体
	// e1 := Employee{"Jane", "Smith", 32, "I like to collect rock albums", []string{"music"}}
	//使用字符串
	// e2 := `{"first_name":"John","last_name":"Smith","age":25,"about":"I love to go rock climbing","interests":["sports","music"]}`
	put1, err := c.Client.Index().
		Index(index).
		Id(id).
		BodyJson(data).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put1.Id, put1.Index, put1.Type)
}

//删除
func (c *ElasticManager) Delete(index, id string) {
	res, err := c.Client.Delete().
		Index(index).
		Id(id).
		Do(context.Background())
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Printf("delete result %s\n", res.Result)
}

//修改
func (c *ElasticManager) Update(index, id string, data map[string]interface{}) {
	res, err := c.Client.Update().
		Index(index).
		Id(id).
		Doc(data).
		Do(context.Background())
	if err != nil {
		println(err.Error())
	}
	fmt.Printf("update age %s\n", res.Result)

}

//查找
func (c *ElasticManager) Gets(index, id string) {
	//通过id查找
	res, err := c.Client.Get().
		Index(index).
		Id(id).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	if res.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", res.Id, res.Version, res.Index, res.Type)
	}

	data, _ := res.Source.MarshalJSON()
	fmt.Println("result:", string(data))
}

//搜索
func (c *ElasticManager) Query(index, id string) {
	var res *elastic.SearchResult
	var err error
	//取所有
	res, err = c.Client.Search("megacorp").Do(context.Background())
	printEmployee(res, err)

	//字段相等
	q := elastic.NewQueryStringQuery("last_name:Smith")
	res, err = c.Client.Search("megacorp").Query(q).Do(context.Background())
	if err != nil {
		println(err.Error())
	}
	printEmployee(res, err)

	//条件查询
	//年龄大于30岁的
	boolQ := elastic.NewBoolQuery()
	boolQ.Must(elastic.NewMatchQuery("last_name", "smith"))
	boolQ.Filter(elastic.NewRangeQuery("age").Gt(30))
	res, err = c.Client.Search("megacorp").Query(q).Do(context.Background())
	printEmployee(res, err)

	//短语搜索 搜索about字段中有 rock climbing
	matchPhraseQuery := elastic.NewMatchPhraseQuery("about", "rock climbing")
	res, err = c.Client.Search("megacorp").Query(matchPhraseQuery).Do(context.Background())
	printEmployee(res, err)

	//分析 interests
	aggs := elastic.NewTermsAggregation().Field("interests")
	res, err = c.Client.Search("megacorp").Aggregation("all_interests", aggs).Do(context.Background())
	printEmployee(res, err)

}

//简单分页
func (c *ElasticManager) List(size, page int, index string) {
	if size < 0 || page < 1 {
		fmt.Printf("param error")
		return
	}
	res, err := c.Client.Search(index).
		Size(size).
		From((page - 1) * size).
		Do(context.Background())
	printEmployee(res, err)

}

//打印查询到的Employee
func printEmployee(res *elastic.SearchResult, err error) {
	if err != nil {
		print(err.Error())
		return
	}
	var typ Employee
	for _, item := range res.Each(reflect.TypeOf(typ)) { //从搜索结果中取数据的方法
		t := item.(Employee)
		fmt.Printf("%#v\n", t)
	}
}
```

客户端调用 main.go
```
package main

import (
	"grpc_gateway/easygo"
)

// Elasticsearch demo
var (
	Elastic *easygo.ElasticManager
	host    string = "http://127.0.0.1:9200/"
)

func init() {
	Elastic = easygo.NewElasticManager(host)
}

type Employee struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Age       int      `json:"age"`
	About     string   `json:"about"`
	Interests []string `json:"interests"`
}

func main() {
	Elastic.Create("megacorp", "1", Employee{"Jane", "Smith", 32, "I like to collect rock albums", []string{"music"}})
	Elastic.Create("megacorp", "2", `{"first_name":"John","last_name":"Smith","age":25,"about":"I love to go rock climbing","interests":["sports","music"]}`)
	Elastic.Create("megacorp", "3", `{"first_name":"Douglas","last_name":"Fir","age":35,"about":"I like to build cabinets","interests":["forestry"]}`)
	Elastic.Update("megacorp", "1", map[string]interface{}{"age": 88})
	Elastic.Gets("megacorp", "1")
	Elastic.Delete("megacorp", "oYuyDHwBmHZlmkWxHr7L")
}
```
