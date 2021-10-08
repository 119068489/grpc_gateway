package easygo

import (
	"fmt"
	"io/ioutil"

	"github.com/astaxie/beego/logs"

	"gopkg.in/yaml.v2"
)

type IYamlConfig interface {
	ReLoadYaml()
	LoadYaml(yamlPath string) map[string]interface{}
	GetValueAsInt(key string, defaultVal ...int) int
	GetValueAsString(key string, defaultVal ...string) string
	GetValueAsBool(key string, defaultVal ...bool) bool
	GetValueAsArrayString(key string, defaultVal ...string) []string

	GetConfig() map[string]interface{}
	DealWithSameKeyItem(config map[string]interface{}, key string, newVal interface{}, oldVal interface{})
	SpecialParseAfterUnmarshal(config map[string]interface{})

	ParseMySQLDsn() map[DB_NAME]map[string]string
	GetMongoDBDsnMaster() map[DB_NAME]map[string]string
	GetMongoDBDsnSlave() map[DB_NAME]map[string]string
	GetMongoInfoMaster(name string) map[string]string
	GetMongoInfoSlave(name string) map[string]string
	GetSpecificInfoOrDefault(name string, defaultName string) map[string]string
	GetSpecificInfoOrDefaultSlave(name string, defaultName string) map[string]string
}

type YamlConfig struct {
	Me       IYamlConfig
	Config   map[string]interface{}
	YamlPath string
}

func NewYamlConfig(yamlPath string) *YamlConfig {
	p := &YamlConfig{}
	p.Init(p, yamlPath)
	return p
}

func (yself *YamlConfig) Init(me IYamlConfig, yamlPath string) {
	yself.Me = me
	yself.YamlPath = yamlPath
	yself.Me.ReLoadYaml()
}

func (yself *YamlConfig) ReLoadYaml() {
	yself.Config = yself.Me.LoadYaml(yself.YamlPath)
}

func (yself *YamlConfig) LoadYaml(yamlPath string) map[string]interface{} {
	logs.Info(yamlPath)
	bytes, err := ioutil.ReadFile(yamlPath)
	PanicError(err)

	var config map[string]interface{}
	err = yaml.Unmarshal(bytes, &config)
	PanicError(err)
	yself.Me.SpecialParseAfterUnmarshal(config)
	value, ok := config["INCLUDE"]
	if !ok {
		return config
	}

	// 如果 include 的多个 yaml 文件有相同的 key，则 INCLUDE 列表前面 yaml 生效
	// 如果当前 yaml 文件与被 INCLUDE 的 yaml 有相同的 key，则当前 yaml 生效
	for _, path1 := range value.([]interface{}) {
		path := path1.(string)
		config2 := yself.Me.LoadYaml(path)
		for k, oldVal := range config2 {
			if newVal, ok := config[k]; !ok {
				config[k] = oldVal // 继承祖业
			} else {
				yself.Me.DealWithSameKeyItem(config, k, newVal, oldVal)
			}
		}
	}
	return config
}

func (yself *YamlConfig) GetConfig() map[string]interface{} {
	return yself.Config
}

func (yself *YamlConfig) GetValueAsInt(key string, defaultVal ...int) int {
	value, ok := yself.Config[key]
	if !ok {
		if len(defaultVal) == 0 {
			//panic(fmt.Sprintf("yaml 配置文件中没有 %v 这一项", key))
			return 0
		} else {
			return defaultVal[0]
		}
	}
	return value.(int)
}

func (yself *YamlConfig) GetValueAsString(key string, defaultVal ...string) string {
	value, ok := yself.Config[key]
	if !ok {
		if len(defaultVal) == 0 {
			//panic(fmt.Sprintf("yaml 配置文件中没有 %v 这一项", key))
			return ""
		} else {
			return defaultVal[0]
		}
	}
	return value.(string)
}

func (yself *YamlConfig) GetValueAsBool(key string, defaultVal ...bool) bool {
	value, ok := yself.Config[key]
	if !ok {
		if len(defaultVal) == 0 {
			panic(fmt.Sprintf("yaml 配置文件中没有 %v 这一项", key))
		} else {
			return defaultVal[0]
		}
	}
	return value.(bool)
}

func (yself *YamlConfig) GetValueAsArrayString(key string, defaultVal ...string) []string {
	value, ok := yself.Config[key]
	if !ok {
		if len(defaultVal) == 0 {
			panic(fmt.Sprintf("yaml 配置文件中没有 %v 这一项", key))
		} else {
			return defaultVal
		}
	}

	list := make([]string, 0)
	for _, v := range value.([]interface{}) {
		list = append(list, v.(string))
	}

	return list
}

// 更深层次的转义，你需要自己去完成

//--------------------------------------------------------------------------

func (yself *YamlConfig) ParseMySQLDsn() map[DB_NAME]map[string]string {
	value, ok := yself.Config["MYSQL_DSN"]
	if !ok {
		panic("配置文件中没有 MYSQL_DSN 这一项")
	}
	dsn := map[DB_NAME]map[string]string{}

	for databaseName, val := range value.(map[interface{}]interface{}) {
		pair := map[string]string{}
		dsn[databaseName.(string)] = pair

		for k, v := range val.(map[interface{}]interface{}) {
			pair[k.(string)] = v.(string)
		}
	}
	return dsn
}

//
func (yself *YamlConfig) GetMongoDBDsnMaster() map[DB_NAME]map[string]string {
	value, ok := yself.Config["MONGODB_MASTER"]
	if !ok {
		panic("配置文件中没有 MONGODB_MASTER 这一项")
	}
	return value.(map[DB_NAME]map[string]string)
}
func (yself *YamlConfig) GetMongoDBDsnSlave() map[DB_NAME]map[string]string {
	value, ok := yself.Config["MONGODB_SLAVE"]
	if !ok {
		panic("配置文件中没有 MONGODB_SLAVE 这一项")
	}
	return value.(map[DB_NAME]map[string]string)
}

func (yself *YamlConfig) SpecialParseAfterUnmarshal(config map[string]interface{}) {
	value, ok := config["MONGODB_MASTER"]
	if ok {
		dsn := map[DB_NAME]map[string]string{}
		for databaseName, val := range value.(map[interface{}]interface{}) {
			pair := map[string]string{}
			dsn[databaseName.(string)] = pair

			for k, v := range val.(map[interface{}]interface{}) {
				pair[k.(string)] = v.(string)
			}
		}
		config["MONGODB_MASTER"] = dsn
	}
	value1, ok1 := config["MONGODB_SLAVE"]
	if ok1 {
		dsn := map[DB_NAME]map[string]string{}
		for databaseName, val := range value1.(map[interface{}]interface{}) {
			pair := map[string]string{}
			dsn[databaseName.(string)] = pair

			for k, v := range val.(map[interface{}]interface{}) {
				pair[k.(string)] = v.(string)
			}
		}
		config["MONGODB_SLAVE"] = dsn
	}
}

// 没有此项则返回 nil
func (yself *YamlConfig) GetMongoInfoMaster(dbName string) map[string]string {
	dict := yself.Me.GetMongoDBDsnMaster()
	return dict[dbName]
}

// 没有此项则返回 nil
func (yself *YamlConfig) GetMongoInfoSlave(dbName string) map[string]string {
	dict := yself.Me.GetMongoDBDsnSlave()
	return dict[dbName]
}

// 取指定 key 的信息，取不到则拿 default 的
func (yself *YamlConfig) GetSpecificInfoOrDefault(name string, defaultName string) map[string]string {
	dict := yself.Me.GetMongoInfoMaster(name)
	if dict != nil {
		return dict
	}
	return yself.Me.GetMongoInfoMaster(defaultName)
}

// 取指定 key 的信息，取不到则拿 default 的
func (yself *YamlConfig) GetSpecificInfoOrDefaultSlave(name string, defaultName string) map[string]string {
	dict := yself.Me.GetMongoInfoSlave(name)
	if dict != nil {
		return dict
	}
	return yself.Me.GetMongoInfoSlave(defaultName)
}

// 默认不需处理，config 已经是读到了最下层的 yaml 文件的值，一般情况下就是你想要的效果
func (yself *YamlConfig) DealWithSameKeyItem(config map[string]interface{}, key string, newVal interface{}, oldVal interface{}) {
	if key == "MONGODB_MASTER" || key == "MONGODB_SLAVE" {
		newVal2 := newVal.(map[DB_NAME]map[string]string)
		oldVal2 := oldVal.(map[DB_NAME]map[string]string)

		for site, oldDict := range oldVal2 {
			newDict, ok := newVal2[site]
			if ok {
				for k, v := range oldDict {
					if _, ok := newDict[k]; !ok { // 部分合并
						newDict[k] = v
					}
				}
			} else {
				newVal2[site] = oldDict //继承祖业
			}
		}
		// log.Println("after merge", newVal2)
	}
}
