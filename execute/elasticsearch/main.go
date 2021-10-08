package main

import (
	"fmt"
	"grpc_gateway/easygo"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
)

// Elasticsearch demo
var (
	Elastic *easygo.ElasticManager
	host    string = "http://127.0.0.1:9200/"
	index   string = "megacorp"
	wg      sync.WaitGroup
)

func init() {
	Elastic = easygo.NewElasticManager(host)
	initializer := easygo.NewInitializer()
	initializer.Execute(nil)
}

type Employee struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Age       int      `json:"age"`
	About     string   `json:"about"`
	Interests []string `json:"interests"`
}

func main() {

	fmt.Scan(&index)
	// Elastic.Create("megacorp", "1", Employee{"Jane", "Smith", 32, "I like to collect rock albums", []string{"music"}})ss
	// Elastic.Create("megacorp", "2", `{"first_name":"John","last_name":"Smith","age":25,"about":"I love to go rock climbing","interests":["sports","music"]}`)
	// Elastic.Create("megacorp", "3", `{"first_name":"Douglas","last_name":"Fir","age":35,"about":"I like to build cabinets","interests":["forestry"]}`)
	// Elastic.Update("megacorp", "1", map[string]interface{}{"age": 88})
	// Elastic.Gets(index, "1")
	// Elastic.Delete("megacorp", "oYuyDHwBmHZlmkWxHr7L")

	// for i := 0; i < 10; i++ {
	// 	wg.Add(1)
	// 	easygo.Spawn(GetLock, i)
	// }
	// wg.Wait()

}

func GetLock(i int) {
	logs.Info(i)
	defer wg.Done()
	a := easygo.RandInt(1, 10)
	lockKey := "test"
	value, errLock := easygo.RedisMgr.GetC().DoRedisLockWithRetry(lockKey, 10, int32(a))
	// defer easygo.RedisMgr.GetC().DoRedisUnlock(lockKey, value)
	if errLock != nil {
		// logs.Error(errLock)
	} else {
		logs.Debug(i, value, a)
	}
	time.Sleep(time.Duration(easygo.RandInt(1, 5)) * time.Second)
}
