# entgo

* 重点目录
  * [快速指南](#快速指南)
  * [字段](#字段fields)
  * [边](#边(Edges))
  * [索引](#索引(Indexes))
  * [增删改查](#增删改查API)
  * [聚合](#聚合)
  * [分页和排序](#分页和排序)
  * [事务](#事务)
  * [迁移](#数据库迁移)


* 简介
  entgo是一个go的实体框架，简单但强大的 ORM 用于建模和查询数据。

  - 简单地使用数据库结构作为图结构。
  - 使用Go代码定义结构。
  - 基于代码生成的静态类型。
  - 容易地进行数据库查询和图遍历。
  - 容易地使用Go模板扩展和自定义。

  应用场景
  entgo非常适合处理各种复杂的关系，定义好实体和实体之间的关系，就可以快速得到各种想要的数据。

  核心概念
  Schema：描述一个实体的定义以及他与其他实体的关系

  Edges：实体与实体之间的关系称为edge(边)


## 快速指南
以下帮助我们快速的学会如何使用entgo，本文档以Postgres数据库为例。

[在线SQL转entgo工具](https://www.printlove.cn/tools/sql2ent/)

[官方文档](https://entgo.io/zh/docs/getting-started/)

_官方文档是以sqlite为例，所以有部分代码和功能会出现报错或执行无效果_

### Installation

`go get entgo.io/ent/cmd/ent`


### 创建你的第一个结构

`go run entgo.io/ent/cmd/ent init User`

- 此命令将生成结构 User 于 <project>/ent/schema/ 目录内:

  <project>/ent/schema/user.go  命令生成的user.go文件

  ``` 代码如下：
  package schema
  
  // User holds the schema definition for the User entity.
  type User struct {
      ent.Schema
  }
  
  // Fields of the User.
  func (User) Fields() []ent.Field {
      return nil
  }
  
  // Edges of the User.
  func (User) Edges() []ent.Edge {
      return nil
  }
  
  ```

- 将 2 个字段添加到 User 结构：
  
  <project>/ent/schema/user.go

  ```代码如下：
  package schema

  // Fields of the User.
  func (User) Fields() []ent.Field {
      return []ent.Field{
          field.Int("age").
              Positive(),//限制最小值为1
          field.String("name").
              Default("unknown"),//设置缺省默认值为"unknown"
      }
  }
  ```

- 在项目根目录中运行 go generate 如下:
  `go generate ./ent`

  这将创建以下文件：
  ```
  ent
  ├── client.go
  ├── config.go
  ├── context.go
  ├── ent.go
  ├── generate.go
  ├── mutation.go
  ... truncated
  ├── schema
  │   └── user.go
  ├── tx.go
  ├── user
  │   ├── user.go
  │   └── where.go
  ├── user.go
  ├── user_create.go
  ├── user_delete.go
  ├── user_query.go
  └── user_update.go
  ```

### 创建你的第一个实体
创建一个 ent.Client。 作为例子，我们将使用 PostgreSql 13

- 在项目下找个合适的地方，创建一个main.go文件，比如：<project>/execute/dbclient/main.go

  ```代码如下：
  package main
  
  import (
  	"context"
  	"fmt"
  	"grpc_gateway/ent"
  	"grpc_gateway/ent/user"
  	"log"
  
  	"entgo.io/ent/dialect"
  	_ "github.com/lib/pq"
  )
  
  func main() {
      //先去手工创建一个数据库testdb
  	host := "127.0.0.1"
  	port := 5432
  	user := "postgres"
  	password := "123456"
  	dbname := "testdb"
  	pdqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
  
  	client, err := ent.Open(dialect.Postgres, pdqlInfo)
  	if err != nil {
  		log.Fatalf("failed opening connection to Postgres: %v", err)
  	}
  	defer client.Close()
  	// Run the auto migration tool.
  	if err := client.Schema.Create(context.Background()); err != nil {
  		log.Fatalf("failed creating schema resources: %v", err)
  	}

    // CreateUser(context.Background(), client)
	// QueryUser(context.Background(), client)
  }
  ```

  _注：如果运行报错可能是没有postgres的驱动，自行go get github.com/lib/pq驱动包_

- 现在，我们做好了创建用户的准备。 让我们调用 CreateUser 函数，比如：
  在文件 <project>/execute/dbclient/main.go 下写个插入函数
  ```代码如下：
  func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
      u, err := client.User.
          Create().
          SetAge(30).
          SetName("a8m").
          Save(ctx)
      if err != nil {
          return nil, fmt.Errorf("failed creating user: %w", err)
      }
      log.Println("user was created: ", u)
      return u, nil
  }
  ```
  
- 现在，让我们查询实体
  继续在文件 <project>/execute/dbclient/main.go 下写个查询函数
  ```代码如下：
  func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
  	u, err := client.User.
  		Query().
  		Where(user.Name("a8m")).
  		// `Only` 在 找不到用户 或 找到多于一个用户 时报错,
  		Only(ctx)
  	if err != nil {
  		return nil, fmt.Errorf("failed querying user: %w", err)
  	}
  	log.Println("user returned: ", u)
  	return u, nil
  }
  ```

- 最后我们在main函数下调用插入和查询函数
  ```
  CreateUser(context.Background(), client)
  QueryUser(context.Background(), client)
  ```
  `go run .\main.go` 

  以下是正确的执行结果：
  2021/08/23 17:54:26 user was created:  User(id=1, age=30, name=a8m)
  2021/08/23 17:54:26 user returned:  User(id=1, age=30, name=a8m)


- 我们进到SQL Shell查看一下数据库的状况
  
  _SQL Shell是postgres的命令行工具_

  - `\l` 命令查看数据，testdb为我们事先创建好的库
    ```
    postgres=# \l
                                                            数据库列表
       名称    |  拥有者  | 字元编码 |            校对规则            |             Ctype              |       存取权限
    -----------+----------+----------+--------------------------------+--------------------------------+-----------------------
     postgres  | postgres | UTF8     | Chinese (Simplified)_China.936 | Chinese (Simplified)_China.936 |
     template0 | postgres | UTF8     | Chinese (Simplified)_China.936 | Chinese (Simplified)_China.936 | =c/postgres          +
               |          |          |                                |                                | postgres=CTc/postgres
     template1 | postgres | UTF8     | Chinese (Simplified)_China.936 | Chinese (Simplified)_China.936 | =c/postgres          +
               |          |          |                                |                                | postgres=CTc/postgres
     testdb    | postgres | UTF8     | Chinese (Simplified)_China.936 | Chinese (Simplified)_China.936 |
    (4 行记录)
    ```

  - `\c testdb` 命令连接上数据库testdb
    ```
    postgres=# \c testdb
    您现在已经连接到数据库 "testdb",用户 "postgres".
    ```
  
  - `\dt` 命令看到users表已经被创建
    ```
    testdb=# \dt;
                 关联列表
     架构模式 | 名称  |  类型  |  拥有者
    ----------+-------+--------+----------
     public   | users | 数据表 | postgres
    (1 行记录)
    ```
  
  - `\d users` 命令查看users表结构
    ```
    testdb=# \d users
                                   数据表 "public.users"
     栏位 |       类型        | 校对规则 |  可空的  |               预设
    ------+-------------------+----------+----------+----------------------------------
     id   | bigint            |          | not null | generated by default as identity
     age  | bigint            |          | not null |
     name | character varying |          | not null | 'unknown'::character varying
    索引：
        "users_pkey" PRIMARY KEY, btree (id)
    ```

### 添加你的第一条边(edge)

- 我们想声明一条到另一个实体的边(关系)。 让我们另外创建2个实体，分别为 Car 和 Group，并添加一些字段。 我们使用 ent CLI 生成初始的结构(schema)：

  `go run entgo.io/ent/cmd/ent init Car Group`

- 手工给Car和Group添加字段
  
  <project>/ent/schema/car.go
  ```
  // Fields of the Car.
  func (Car) Fields() []ent.Field {
      return []ent.Field{
          field.String("model"),
          field.Time("registered_at"),
      }
  }
  ```

  <project>/ent/schema/group.go
  ```
  // Fields of the Group.
  func (Group) Fields() []ent.Field {
      return []ent.Field{
          field.String("name").
              // Regexp validation for group name.
              Match(regexp.MustCompile("[a-zA-Z_]+$")),
      }
  }
  ```

- 让我们来定义第一个关系。 从 User 到 Car 的关系，定义了一个用户可以拥有1辆或多辆汽车，但每辆汽车只有一个车主（一对多关系）。
  
  * 让我们将 "cars" 关系添加到 User 结构中
    <project>/ent/schema/user.go
    ```
    // Edges of the User.
    func (User) Edges() []ent.Edge {
        return []ent.Edge{
            edge.To("cars", Car.Type),
        }
    }
    ```
  * go generate重新生成文件
    `go generate ./ent`

  * 作为示例，我们创建2辆汽车并将它们添加到某个用户。
    ```
    func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
    	// Create a new car with model "Tesla".
    	tesla, err := client.Car.
    		Create().
    		SetModel("Tesla").
    		SetRegisteredAt(time.Now()).
    		Save(ctx)
    	if err != nil {
    		return nil, fmt.Errorf("failed creating car: %w", err)
    	}
    	log.Println("car was created: ", tesla)
    
    	// Create a new car with model "Ford".
    	ford, err := client.Car.
    		Create().
    		SetModel("Ford").
    		SetRegisteredAt(time.Now()).
    		Save(ctx)
    	if err != nil {
    		return nil, fmt.Errorf("failed creating car: %w", err)
    	}
    	log.Println("car was created: ", ford)
    
    	// Create a new user, and add it the 2 cars.
    	a8m, err := client.User.
    		Create().
    		SetAge(30).
    		SetName("a8m").
    		AddCars(tesla, ford).
    		Save(ctx)
    	if err != nil {
    		return nil, fmt.Errorf("failed creating user: %w", err)
    	}
    	log.Println("user was created: ", a8m)
    	return a8m, nil
    }
    ```
  
  * 想要查询 cars 关系怎么办？ 请参考以下：
    ```
    func QueryCars(ctx context.Context, a8m *ent.User) error {
    	cars, err := a8m.QueryCars().All(ctx)
    	if err != nil {
    		return fmt.Errorf("failed querying user cars: %w", err)
    	}
    	log.Println("returned cars:", cars)
    
    	// What about filtering specific cars.
    	ford, err := a8m.QueryCars().
    		Where(car.Model("Ford")).
    		Only(ctx)
    	if err != nil {
    		return fmt.Errorf("failed querying user cars: %w", err)
    	}
    	log.Println(ford)
    	return nil
    }
    ```

- 添加您的第一条逆向边(反向引用)
  假定我们有一个 Car 对象，我们想要得到它的所有者；即这辆汽车所属的用户。 为此，我们有另一种“逆向”的边，通过 edge.From 函数定义。

  新建的边是隐性的，以强调我们不会在数据库中创建另一个关联。 它只是真正边(关系) 的回溯。

  让我们把一个名为 owner 的逆向边添加到 Car 的结构中, 在 User 结构中引用它到 cars 关系 然后运行 go generate ./ent

  <project>/ent/schema/car.go
  ```
  // Edges of the Car.
  func (Car) Edges() []ent.Edge {
      return []ent.Edge{
          // Create an inverse-edge called "owner" of type `User`
          // and reference it to the "cars" edge (in User schema)
          // explicitly using the `Ref` method.
          edge.From("owner", User.Type).
              Ref("cars").
              // setting the edge to unique, ensure
              // that a car can have only one owner.
              Unique(),
      }
  }
  ```

  我们继续使用用户和汽车，作为查询逆向边的例子。
  ```
  func QueryCarUsers(ctx context.Context, user *ent.User) error {
  	cars, err := user.QueryCars().All(ctx)
  	if err != nil {
  		return fmt.Errorf("failed querying user cars: %w", err)
  	}
  	// Query the inverse edge.
  	for _, ca := range cars {
  		owner, err := ca.QueryOwner().Only(ctx)
  		if err != nil {
  			return fmt.Errorf("failed querying car %q owner: %w", ca.Model,   err)
  		}
  		log.Printf("car %q owner: %q\n", ca.Model, owner.Name)
  	}
  	return nil
  }
  ```

### 创建你的第二条边
继续我们的例子，在用户和组之间创建一个M2M（多对多）关系。

每个组实体可以拥有许多用户，而一个用户可以关联到多个组。

Group结构是users关系的所有者， 而User实体对这个关系有一个名为groups的反向引用。 让我们在结构中定义这种关系。

- group中添加一条边
  <project>/ent/schema/group.go
  ```
  // Edges of the Group.
  func (Group) Edges() []ent.Edge {
     return []ent.Edge{
         edge.To("users", User.Type),
     }
  }
  ```

- user中添加一条反向引用
  <project>/ent/schema/user.go
  ```
  // Edges of the User.
  func (User) Edges() []ent.Edge {
     return []ent.Edge{
         edge.To("cars", Car.Type),
         // Create an inverse-edge called "groups" of type `Group`
         // and reference it to the "users" edge (in Group schema)
         // explicitly using the `Ref` method.
         edge.From("groups", Group.Type).
             Ref("users"),
     }
  }
  ```

- 我们在schema目录上运行ent来重新生成资源文件。
  `go generate ./ent`

### 进行第一次图遍历
为了进行我们的第一次图遍历，我们需要生成一些数据（节点和边，或者说，实体和关系）。让我们创建如下图所示的框架：

user1 'Ariel'在gruop1 'GitHub'和group2 'GitLab'中，user2 'Neta'在group2中 'GitLab'。

user1 'Ariel'拥有car3 'Tesla'和car4 'Mazda',user2 'Neta'拥有car5 'Ford'

```函数代码如下：
func CreateGraph(ctx context.Context, client *ent.Client) error {
    // First, create the users.
    a8m, err := client.User.
        Create().
        SetAge(30).
        SetName("Ariel").
        Save(ctx)
    if err != nil {
        return err
    }
    neta, err := client.User.
        Create().
        SetAge(28).
        SetName("Neta").
        Save(ctx)
    if err != nil {
        return err
    }
    // Then, create the cars, and attach them to the users in the creation.
    _, err = client.Car.
        Create().
        SetModel("Tesla").
        SetRegisteredAt(time.Now()). // ignore the time in the graph.
        SetOwner(a8m).               // attach this graph to Ariel.
        Save(ctx)
    if err != nil {
        return err
    }
    _, err = client.Car.
        Create().
        SetModel("Mazda").
        SetRegisteredAt(time.Now()). // ignore the time in the graph.
        SetOwner(a8m).               // attach this graph to Ariel.
        Save(ctx)
    if err != nil {
        return err
    }
    _, err = client.Car.
        Create().
        SetModel("Ford").
        SetRegisteredAt(time.Now()). // ignore the time in the graph.
        SetOwner(neta).              // attach this graph to Neta.
        Save(ctx)
    if err != nil {
        return err
    }
    // Create the groups, and add their users in the creation.
    _, err = client.Group.
        Create().
        SetName("GitLab").
        AddUsers(neta, a8m).
        Save(ctx)
    if err != nil {
        return err
    }
    _, err = client.Group.
        Create().
        SetName("GitHub").
        AddUsers(a8m).
        Save(ctx)
    if err != nil {
        return err
    }
    log.Println("The graph was created successfully")
    return nil
}
```

调用这个函数将会创建cars,users,groups,group_users 4个表，分别包含了我们写入的初始数据

现在我们有一个含数据的图，我们可以对它运行一些查询：
- 1、获取名为 "GitHub" 的群组内所有用户拥有的汽车
  ```函数代码如下：
  func QueryGithub(ctx context.Context, client *ent.Client) error {
      cars, err := client.Group.
          Query().
          Where(group.Name("GitHub")). // (Group(Name=GitHub),)
          QueryUsers().                // (User(Name=Ariel, Age=30),)
          QueryCars().                 // (Car(Model=Tesla,   RegisteredAt=<Time>), Car(Model=Mazda, RegisteredAt=<Time>),)
          All(ctx)
      if err != nil {
          return fmt.Errorf("failed getting cars: %w", err)
      }
      log.Println("cars returned:", cars)
      // Output: (Car(Model=Tesla, RegisteredAt=<Time>), Car(Model=Mazda,   RegisteredAt=<Time>),)
      return nil
  }
  ```
    执行结果：cars returned: [Car(id=1, model=Tesla, registered_at=Tue Aug 24 14:22:18 2021) Car(id=2, model=Mazda, registered_at=Tue Aug 24 14:22:18 2021)]

- 2、修改上面的查询，从用户 Ariel 开始遍历。
  ```
  func QueryArielCars(ctx context.Context, client *ent.Client) error {
      // Get "Ariel" from previous steps.
      a8m := client.User.
          Query().
          Where(
              user.HasCars(),
              user.Name("Ariel"),
          ).
          OnlyX(ctx)
      cars, err := a8m.                       // Get the groups, that a8m is connected to:
              QueryGroups().                  // (Group(Name=GitHub), Group(Name=GitLab),)
              QueryUsers().                   // (User(Name=Ariel, Age=30), User(Name=Neta, Age=28),)
              QueryCars().                    //
              Where(                          //
                  car.Not(                    //  Get Neta and Ariel cars, but filter out
                      car.Model("Mazda"),     //  those who named "Mazda"
                  ),                          //
              ).                              //
              All(ctx)
      if err != nil {
          return fmt.Errorf("failed getting cars: %w", err)
      }
      log.Println("cars returned:", cars)
      // Output: (Car(Model=Tesla, RegisteredAt=<Time>), Car(Model=Ford, RegisteredAt=<Time>),)
      return nil
  }
  ```
  执行结果：cars returned: [Car(id=1, model=Tesla, registered_at=Tue Aug 24 14:22:18 2021) Car(id=3, model=Ford, registered_at=Tue Aug 24 14:22:18 2021)]

- 3、获取所有拥有用户的群组 (通过额外 [look-aside] 条件查询)：
  ```
  func QueryGroupWithUsers(ctx context.Context, client *ent.Client) error {
  	groups, err := client.Group.
  		Query().
  		Where(group.HasUsers()).
  		All(ctx)
  	if err != nil {
  		return fmt.Errorf("failed getting groups: %w", err)
  	}
  	log.Println("groups returned:", groups)
  	// Output: (Group(Name=GitHub), Group(Name=GitLab),)
  	return nil
  }
  ```
  执行结果:groups returned: [Group(id=2, name=GitHub) Group(id=1, name=GitLab)]

- 4、根据汽车查询用户所在的群组
  ```
  func QueryCarGroups(ctx context.Context, car *ent.Car) error {
  	owner, err := car.QueryOwner().Only(ctx)
  	if err != nil {
  		return fmt.Errorf("failed getting user: %w", err)
  	}
  
  	groups, err := owner.
  		QueryGroups().
  		Where(predicate.Group(user.HasCarsWith())).
  		All(ctx)
  	if err != nil {
  		return fmt.Errorf("failed getting groups: %w", err)
  	}
  	log.Println("groups returned:", groups)
  	return nil
  }
  ```
  执行结果：groups returned: [Group(id=1, name=GitLab)]


## 结构(Schema)

### 引言
通过快速指南的实践，我们对Schema已经不陌生了，Schema描述了图中一个实体类型的定义，如 User 或 Group， 并可以包含以下配置：

- 实体的字段 (或属性)，如：User 的姓名或年龄。
- 实体的边 (或关系)。如：User 所属用户组，或 User 的朋友。
- 数据库相关的配置，如：索引或唯一索引。

以下代码是一个示例：
```
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/index"
)

type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Int("age"),
        field.String("name"),
        field.String("nickname").
            Unique(),
    }
}

func (User) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("groups", Group.Type),
        edge.To("friends", User.Type),
    }
}

func (User) Index() []ent.Index {
    return []ent.Index{
        index.Fields("age", "name").
            Unique(),
    }
}
```

实体 Schema 通常存储在你项目根目录的 ent/schema 目录内， 且可以通过以下的 entc 命令生成：

`go run entgo.io/ent/cmd/ent init User Group`

### 字段(Fields)
- 概述
  Schema 中的字段（或属性）是节点的属性。 例如：User 有 age， name, username 和 created_at 4个字段。
  用 schema 的 Fields 方法可返回这些字段。 如：
  ```
  package schema
  
  import (
      "time"
  
      "entgo.io/ent"
      "entgo.io/ent/schema/field"
  )
  
  // User schema.
  type User struct {
      ent.Schema
  }
  
  // Fields of the user.
  func (User) Fields() []ent.Field {
      return []ent.Field{
          field.Int("age"),
          field.String("name"),
          field.String("username").
              Unique(),
          field.Time("created_at").
              Default(time.Now),
      }
  }
  ```
  默认所有字段都是必需的，可使用 Optional 方法设置为可选字段。

- 类型
  目前框架支持以下数据类型：

  * Go中所有数值类型。 如 int，uint8，float64 等
  * bool 布尔型
  * string 字符串
  * time.Time 时间类型
  * UUID
  * []byte (SQL only).
  * JSON (SQL only).
  * Enum (SQL only).
  * 其它类型 （仅限SQL）

  ```
  // Fields of the User.
  func (User) Fields() []ent.Field {
  	return []ent.Field{
  		field.Int("age").
  			Positive(),               //允许最小值为1的验证器
  		field.Float("rank").
  			Optional(),               //设置字段可选
  		field.Bool("active").
  			Default(false),           //设置默认值为false
  		field.String("name").
  			Unique(),                 //设置字段值唯一
  		field.Time("created_at").
  			Default(time.Now),        //设置默认值为当前时间
  		field.JSON("url", &url.URL{}).
  			Optional(),
  		field.JSON("strings", []string{}).
  			Optional(),
  		field.Enum("state").
  			Values("on", "off").      //设置值为枚举值on/off
  			Optional(),
  		field.UUID("uuid", uuid.UUID{}).
  			Default(uuid.New),        //设置默认值为UUID
  	}
  }
  ```

- ID 字段
  id 字段内置于架构中，无需声明。 在基于 SQL 的数据库中，它的类型默认为 int （可以使用 代码生成设置 更改）并自动递增。

  如果想要配置 id 字段在所有表中唯一，可以在运行 schema 迁移时使用 WithGlobalUniqueID 选项实现。
  
  如果需要对 id 字段进行其他配置，或者要使用由应用程序在实体创建时提供的 id （例如UUID），可以覆盖内置 id 配置。 
  
  如果你需要设置一个自定义函数来生成 ID， 使用 DefaultFunc 方法来指定一个函数，每次ID将由此函数生成。如：
  ```
  // Fields of the Group.
  func (Group) Fields() []ent.Field {
      return []ent.Field{
          field.Int("id").
              StructTag(`json:"oid,omitempty"`),
      }
  }
  
  // Fields of the Blob.
  func (Blob) Fields() []ent.Field {
      return []ent.Field{
          field.UUID("id", uuid.UUID{}).
              Default(uuid.New).
              StorageKey("oid"),
      }
  }
  
  // Fields of the Pet.
  func (Pet) Fields() []ent.Field {
      return []ent.Field{
          field.String("id").
              MaxLen(25).
              NotEmpty().
              Unique().
              Immutable(),
      }
  }

  // Fields of the User.
  func (User) Fields() []ent.Field {
      return []ent.Field{
          field.Int64("id").
              DefaultFunc(func() int64 {
                  // An example of a dumb ID generator - use a production-ready alternative instead.
                  return time.Now().Unix() << 8 | atomic.AddInt64(&counter, 1) % 256
              }),
      }
  }
  ```

- 数据库字段类型
  每个数据库方言都有自己Go类型与数据库类型的映射。 例如，MySQL 方言将Go类型为 float64 的字段创建为 double 的数据库字段。 当然，我们也可以通过 SchemaType 方法来重写默认的类型映射。
  ```
  // Fields of the Card.
  func (Card) Fields() []ent.Field {
      return []ent.Field{
          field.Float("amount").
              SchemaType(map[string]string{
                  dialect.MySQL:    "decimal(6,2)",   // Override MySQL.
                  dialect.Postgres: "numeric",        // Override Postgres.
              }),
      }
  }
  ```

- Go 类型
  字段的默认类型是基本的 Go 类型。 例如，对于字符串字段，类型是 string, 对于时间字段，类型是 time.Time。 GoType 方法提供了以自定义类型覆盖 默认类型的选项。

  自定义类型必须是可转换为Go基本类型，或者实现 ValueScanner 接口的类型。

  ```
  // Fields of the Card.
  func (Card) Fields() []ent.Field {
      return []ent.Field{
          field.Float("amount").
              GoType(Amount(0)),
          field.String("name").
              Optional().
              // A ValueScanner type.
              GoType(&sql.NullString{}),
          field.Enum("role").
              // A convertible type to string.
              GoType(role.Role("")),
          field.Float("decimal").
              // A ValueScanner type mixed with SchemaType.
              GoType(decimal.Decimal{}).
              SchemaType(map[string]string{
                  dialect.MySQL:    "decimal(6,2)",
                  dialect.Postgres: "numeric",
              }),
      }
  }
  ```

- 其它字段
  Other 代表一个不适合任何标准字段类型的字段。 示例为 Postgres 中的 Rage 类型或 Geospatial 类型
  ```
  // Fields of the User.
  func (User) Fields() []ent.Field {
      return []ent.Field{
          field.Other("duration", &pgtype.Tstzrange{}).
              SchemaType(map[string]string{
                  dialect.Postgres: "tstzrange",
              }),
      }
  }
  ```

- 默认值
  非唯一 字段可使用 Default 和 UpdateDefault 方法为其设置默认值。 你也可以指定 DefaultFunc 方法来自定义默认值生成。
  ```
  // Fields of the User.
  func (User) Fields() []ent.Field {
      return []ent.Field{
          field.Time("created_at").
              Default(time.Now),
          field.Time("updated_at").
              Default(time.Now).
              UpdateDefault(time.Now),
          field.String("name").
              Default("unknown"),
          field.String("cuid").
              DefaultFunc(cuid.New),
      }
  }
  ```
  可以通过 entsql.Annotation 将像函数调用的SQL特定表达式添加到默认值配置中：
  ```
  // Fields of the User.
  func (User) Fields() []ent.Field {
      return []ent.Field{
          // Add a new field with CURRENT_TIMESTAMP
          // as a default value to all previous rows.
          field.Time("created_at").
              Default(time.Now).
              Annotations(&entsql.Annotation{
                  Default: "CURRENT_TIMESTAMP",
              }),
      }
  }
  ```
  为避免你指定的 DefaultFunc 方法也返回了一个错误，最好使用 schema-hooks 处理它。

- 校验器
  字段校验器是一个 func(T) error 类型的函数，定义在 schema 的 Validate 方法中，字段在创建或更新前会执行此方法。支持 string 类型和所有数值类型。
  ```
  // Fields of the group.
  func (Group) Fields() []ent.Field {
      return []ent.Field{
          field.String("name").
              Match(regexp.MustCompile("[a-zA-Z_]+$")).
              Validate(func(s string) error {
                  if strings.ToLower(s) == s {
                      return errors.New("group name must begin with uppercase")
                  }
                  return nil
              }),
      }
  }
  ```
  又如：编写一个可复用的校验器
  ```
  // MaxRuneCount validates the rune length of a string by using the unicode/utf8 package.
  func MaxRuneCount(maxLen int) func(s string) error {
      return func(s string) error {
          if utf8.RuneCountInString(s) > maxLen {
              return errors.New("value is more than the max length")
          }
          return nil
      }
  }
  
  field.String("name").
      Validate(MaxRuneCount(10))
  field.String("nickname").
      Validate(MaxRuneCount(20))
  ```

- 内置校验器
  框架为每个类型提供了几个内置的验证器：

  * 数值类型：
    Positive() - 验证给定最小值为1
    Negative() - 验证给定最大值为-1
    NonNegative() - 验证给定最小值为0
    Min(i) - 验证给定的值 > i。
    Max(i) - 验证给定的值 < i。
    Range(i, j) - 验证给定值在 [i, j] 之间。

  * string 字符串
    MinLen(i)
    MaxLen(i)
    Match(regexp.Regexp)
    NotEmpty()

- Optional 可选项
  可选字段为创建时非必须的字段，在数据库中被设置为 null。 和 edges 不同，字段默认都为必需字段，可通过 Optional 方法显示的设为可选字段。
  ```
  // Fields of the user.
  func (User) Fields() []ent.Field {
      return []ent.Field{
          field.String("required_name"),
          field.String("optional_name").
              Optional(),
      }
  }
  ```

- Nillable
  有时您希望能够区分字段的零值和 nil； 例如，如果数据库列包含 0 或 NULL。那么 Nillable 就派上用场了。

  如果你有一个类型为 T 的 可选字段，设置为 Nillable 后，将生成一个类型为 *T 的结构体字段。 因此，如果数据库返回 NULL 字段， 结构体字段将为 nil 值。否则，它将包含一个指向实际数据的指针。

  ```
  // Fields of the user.
  func (User) Fields() []ent.Field {
      return []ent.Field{
          field.String("required_name"),
          field.String("optional_name").
              Optional(),
          field.String("nillable_name").
              Optional().
              Nillable(),
      }
  }
  ```
  以上生成的结构体如下：
  ```
  // ent/user.go
  package ent
  
  // User entity.
  type User struct {
      RequiredName string `json:"required_name,omitempty"`
      OptionalName string `json:"optional_name,omitempty"`
      NillableName *string `json:"nillable_name,omitempty"`
  }
  ```

- Immutable 不可变的
  字段可以使用 Immutable 方法定义不可变字段只能在创建实体时设置。 即：不会为实体生成任何更新方法。

- 唯一键
  字段可以使用 Unique 方法定义为唯一字段。 注意：唯一字段不能有默认值。

- 存储键名
  可以使用 StorageKey 方法自定义数据库中的字段名称。 在 SQL 中为字段名，在 Gremlin 中为属性名称。

- 索引
  索引可以在多字段和某些类型的 edges 上定义. 注意：目前只有 SQL 类型的数据库支持此功能。

- 结构体标记（tags）
  
- 外部模版
  默认情况下，ent 使用在 schema.Fields 方法中配置的字段生成实体模型。 例如，给定此架构配置：
  ```
  // Fields of the user.
  func (User) Fields() []ent.Field {
      return []ent.Field{
          field.Int("age").
              Optional().
              Nillable(),
          field.String("name").
              StructTag(`gqlgen:"gql_name"`),
      }
  }
  ```
  生成的模版如下：
  ```
  // User is the model entity for the User schema.
  type User struct {
      // Age holds the value of the "age" field.
      Age  *int   `json:"age,omitempty"`
      // Name holds the value of the "name" field.
      Name string `json:"name,omitempty" gqlgen:"gql_name"`
  }
  ```
  为了向生成的结构中添加未存储在数据库中的其他字段，请使用外部模板。 例如：
  ```
  {{ define "model/fields/additional" }}
      {{- if eq $.Name "User" }}
          // StaticField defined by template.
          StaticField string `json:"static,omitempty"`
      {{- end }}
  {{ end }}
  ```
  生成的模版如下：
  ```
  // User is the model entity for the User schema.
  type User struct {
      // Age holds the value of the "age" field.
      Age  *int   `json:"age,omitempty"`
      // Name holds the value of the "name" field.
      Name string `json:"name,omitempty" gqlgen:"gql_name"`
      // StaticField defined by template.
      StaticField string `json:"static,omitempty"`
  }
  ```

- 敏感字段
  可以使用 Sensitive 方法将字符串字段定义为敏感字段。 不会打印敏感字段，编码时将省略它们。 请注意，敏感字段不能有结构标记。
  ```
  // Fields of the user.
  func (User) Fields() []ent.Field {
      return []ent.Field{
          field.String("name").
			Unique(), //设置字段值唯一
          field.String("password").
              Sensitive(),//设置成敏感字段
      }
  }
  ```
  生成的字段如下
  ```
  // User is the model entity for the User schema.
  type User struct {
    Name string `json:"name,omitempty"`
    Password string `json:"-"`
  }
  ```

### 边(Edges)
- Quick Summary
  1. cars / owner edges; user's cars and car's owner
     ```user.go
     // Edges of the user.
     func (User) Edges() []ent.Edge {
         return []ent.Edge{
             edge.To("cars", Car.Type),
         }
     }
     ```
     ```car.go
     // Edges of the car.
     func (Car) Edges() []ent.Edge {
         return []ent.Edge{
             edge.From("owner", User.Type).
                 Ref("cars").
                 Unique(),
         }
     }
     ```

     如您所见，一个 User 实体可以拥有多个car，但一个 car 实体只能拥有一个owner。

     在关系定义中，car边是O2M（一对多）关系，owner边是M2O（多对一）关系。
     
     User 模式拥有car/owner关系，因为它使用 edge.To，而 Car 模式只有对它的反向引用，使用 edge.From 和 Ref 方法声明。
     
     Ref 方法描述了我们引用的 User 模型的哪条边，因为从一个模型到另一个模型可以有多个引用。
     
     可以使用 Unique 方法控制边/关系的基数，下面将对其进行更广泛的解释。

  2. users / groups edges; group's users and user's groups
     ```group.go
     // Edges of the group.
     func (Group) Edges() []ent.Edge {
         return []ent.Edge{
             edge.To("users", User.Type),
         }
     }
     ```
     ```user.go
     // Edges of the user.
     func (User) Edges() []ent.Edge {
         return []ent.Edge{
             edge.From("groups", Group.Type).
                 Ref("users"),
             // "pets" declared in the example above.
             edge.To("pets", Pet.Type),
         }
     }
     ```
     如您所见，一个 Group 实体可以有多个用户，一个 User 实体可以有多个组。

     在关系定义中，用户边是M2M（多对多）关系，组边也是M2M（多对多）关系。

- To and From
  edge.To 和 edge.From 是用于创建边/关系的 2 个构建器。

  使用 edge.To 构建器定义边的模式拥有该关系，这与使用 edge.From 构建器仅提供关系的反向引用（具有不同名称）不同。

- 关系
  1. O2O Two Types (O2O 两种类型) `简单点说，就是2个表之间的单对单关系`
    ![关系图例](https://entgo.io/images/assets/er_user_card.png "example")

     在这个例子中，一个用户只有一张信用卡，一张卡只有一个所有者。User 架构定义了一个名为 card 的 edge.To 卡，Card 架构使用 edge.From 命名所有者定义了对此边的反向引用。
     ```
     // Edges of the user.
     func (User) Edges() []ent.Edge {
         return []ent.Edge{
             edge.To("card", Card.Type).
                 Unique(),
         }
     }
     ```
     ```
     // Edges of the Card.
     func (Card) Edges() []ent.Edge {
         return []ent.Edge{
             edge.From("owner", User.Type).
                 Ref("card").
                 Unique().
                 // We add the "Required" method to the builder
                 // to make this edge required on entity creation.
                 // i.e. Card cannot be created without its owner.
                 Required(),
         }
     }
     ```
     与这些边交互的 API 如下：
     ```
     func Do(ctx context.Context, client *ent.Client) error {
        a8m, err := client.User.
            Create().
            SetAge(30).
            SetName("Mashraki").
            Save(ctx)
        if err != nil {
            return fmt.Errorf("creating user: %w", err)
        }
        log.Println("user:", a8m)
        card1, err := client.Card.
            Create().
            SetOwner(a8m).
            SetNumber("1020").
            SetExpired(time.Now().Add(time.Minute)).
            Save(ctx)
        if err != nil {
            return fmt.Errorf("creating card: %w", err)
        }
        log.Println("card:", card1)
        // Only returns the card of the user,
        // and expects that there's only one.
        card2, err := a8m.QueryCard().Only(ctx)
        if err != nil {
            return fmt.Errorf("querying card: %w", err)
        }
        log.Println("card:", card2)
        // The Card entity is able to query its owner using
        // its back-reference.
        owner, err := card2.QueryOwner().Only(ctx)
        if err != nil {
            return fmt.Errorf("querying owner: %w", err)
        }
        log.Println("owner:", owner)
        return nil
     }
     ```
     完整例子请去[GitHub](https://github.com/ent/ent/tree/master/examples/o2o2types)

  2. O2O Same Type (O2O 同类型) `简单点说，就是1个表里的单对单关系`
    ![关系图例](https://entgo.io/images/assets/er_linked_list.png "linked-list example")

     我们有一个名为 next/prev 的递归关系。 列表中的每个节点只能有一个下一个节点。 如果节点 A 指向（使用 next）节点 B，则 B 可以使用 prev（后向引用边）获取其指针。
     ```
     // Edges of the Node.
     func (Node) Edges() []ent.Edge {
         return []ent.Edge{
             edge.To("next", Node.Type).
                 Unique().
                 From("prev").
                 Unique(),
 
             //edge.To("next", Node.Type).
             //    Unique(),
             //edge.From("prev", Node.Type).
             //    Ref("next).
             //    Unique(),
         }
     }
     ```
     如上所见，对于相同类型的关系，可以在同一个构建器中声明边及其引用。
     ```与此edge交互的 API 如下：
     func Do(ctx context.Context, client *ent.Client) error {
         head, err := client.Node.
             Create().
             SetValue(1).
             Save(ctx)
         if err != nil {
             return fmt.Errorf("creating the head: %w", err)
         }
         curr := head
         // Generate the following linked-list: 1<->2<->3<->4<->5.
         for i := 0; i < 4; i++ {
             curr, err = client.Node.
                 Create().
                 SetValue(curr.Value + 1).
                 SetPrev(curr).
                 Save(ctx)
             if err != nil {
                 return err
             }
         }
 
         // Loop over the list and print it. `FirstX` panics if an error occur.
         for curr = head; curr != nil; curr = curr.QueryNext().FirstX(ctx) {
             fmt.Printf("%d ", curr.Value)
         }
         // Output: 1 2 3 4 5
 
         // Make the linked-list circular:
         // The tail of the list, has no "next".
         tail, err := client.Node.
             Query().
             Where(node.Not(node.HasNext())).
             Only(ctx)
         if err != nil {
             return fmt.Errorf("getting the tail of the list: %v", tail)
         }
         tail, err = tail.Update().SetNext(head).Save(ctx)
         if err != nil {
             return err
         }
         // Check that the change actually applied:
         prev, err := head.QueryPrev().Only(ctx)
         if err != nil {
             return fmt.Errorf("getting head's prev: %w", err)
         }
         fmt.Printf("\n%v", prev.Value == tail.Value)
         // Output: true
         return nil
     }
     ```
     完整例子请去[GitHub](https://github.com/ent/ent/tree/master/examples/o2orecur)

  3. O2O Bidirectional (O2O 双向)  `简单点说，就是1个表里的双向单对单关系`
    ![关系图例](https://entgo.io/images/assets/er_user_spouse.png "example")
   
     在这个用户配偶示例中，我们有一个名为配偶的对称 O2O 关系。 每个用户只能有一个配偶。 如果用户 A 将其配偶（使用配偶）设置为 B，则 B 可以使用配偶边缘获取其配偶。请注意，在双向边的情况​​下没有所有者/逆项。

     ```这将在表中生成user_spouse字段保存关系id
     // Edges of the User.
     func (User) Edges() []ent.Edge {
         return []ent.Edge{
             edge.To("spouse", User.Type).
                 Unique(),
         }
     }
     ```
     与此edge交互的 API 如下：
     ```
     func Do(ctx context.Context, client *ent.Client) error {
         a8m, err := client.User.
             Create().
             SetAge(30).
             SetName("a8m").
             Save(ctx)
         if err != nil {
             return fmt.Errorf("creating user: %w", err)
         }
         nati, err := client.User.
             Create().
             SetAge(28).
             SetName("nati").
             SetSpouse(a8m).
             Save(ctx)
         if err != nil {
             return fmt.Errorf("creating user: %w", err)
         }
     
         // Query the spouse edge.
         // Unlike `Only`, `OnlyX` panics if an error occurs.
         spouse := nati.QuerySpouse().OnlyX(ctx)
         fmt.Println(spouse.Name)
         // Output: a8m
     
         spouse = a8m.QuerySpouse().OnlyX(ctx)
         fmt.Println(spouse.Name)
         // Output: nati
     
         // Query how many users have a spouse.
         // Unlike `Count`, `CountX` panics if an error occurs.
         count := client.User.
             Query().
             Where(user.HasSpouse()).
             CountX(ctx)
         fmt.Println(count)
         // Output: 2
     
         // Get the user, that has a spouse with name="a8m".
         spouse = client.User.
             Query().
             Where(user.HasSpouseWith(user.Name("a8m"))).
             OnlyX(ctx)
         fmt.Println(spouse.Name)
         // Output: nati
         return nil
     }
     ```
     请注意，可以使用 Edge Field 选项配置外键列并将其公开为实体字段，如下所示：这将在表中生成spouse_id字段保存关系id
     ```
     // Fields of the User.
     func (User) Fields() []ent.Field {
         return []ent.Field{
             field.Int("spouse_id").
                 Optional(),
         }
     }
     
     // Edges of the User.
     func (User) Edges() []ent.Edge {
         return []ent.Edge{
             edge.To("spouse", User.Type).
                 Unique().
                 Field("spouse_id"),
         }
     }
     ```
     完整例子请去[GitHub](https://github.com/ent/ent/tree/master/examples/o2obidi)
     
  4. O2M Two Types (O2M 两种类型)  `简单点说，就是2个表里的单对多关系`
   ![关系图例](https://entgo.io/images/assets/er_user_pets.png "example")
    
     在这个 user-pets 示例中，我们在用户与其宠物之间建立了 O2M 关系。 每个用户有很多宠物，一个宠物有一个主人。 如果用户 A 使用 pets 边添加了宠物 B，则 B 可以使用所有者边（反向引用边）获取其所有者。请注意，从 Pet 模式的角度来看，这种关系也是 M2O（多对一）。_快速指南中的汽车和拥有者的关系跟本例一样_

     ```
     // Edges of the User.
     func (User) Edges() []ent.Edge {
         return []ent.Edge{
             edge.To("pets", Pet.Type),
         }
     }
     ```
     ```
     // Edges of the Pet.
     func (Pet) Edges() []ent.Edge {
         return []ent.Edge{
             edge.From("owner", User.Type).
                 Ref("pets").
                 Unique(),
         }
     }
     ```
     与此edge交互的 API 如下：
     ```
     func Do(ctx context.Context, client *ent.Client) error {
         // Create the 2 pets.
         pedro, err := client.Pet.
             Create().
             SetName("pedro").
             Save(ctx)
         if err != nil {
             return fmt.Errorf("creating pet: %w", err)
         }
         lola, err := client.Pet.
             Create().
             SetName("lola").
             Save(ctx)
         if err != nil {
             return fmt.Errorf("creating pet: %w", err)
         }
         // Create the user, and add its pets on the creation.
         a8m, err := client.User.
             Create().
             SetAge(30).
             SetName("a8m").
             AddPets(pedro, lola).
             Save(ctx)
         if err != nil {
             return fmt.Errorf("creating user: %w", err)
         }
         fmt.Println("User created:", a8m)
         // Output: User(id=1, age=30, name=a8m)
     
         // Query the owner. Unlike `Only`, `OnlyX` panics if an error occurs.
         owner := pedro.QueryOwner().OnlyX(ctx)
         fmt.Println(owner.Name)
         // Output: a8m
     
         // Traverse the sub-graph. Unlike `Count`, `CountX` panics if an error occurs.
         count := pedro.
             QueryOwner(). // a8m
             QueryPets().  // pedro, lola
             CountX(ctx)   // count
         fmt.Println(count)
         // Output: 2
         return nil
     }
     ```
     完整例子请去[GitHub](https://github.com/ent/ent/tree/master/examples/o2m2types)

  5. O2M Same Type (O2M 同类型) `简单点说，就是1个表里的单对多关系`
    ![关系图例](https://entgo.io/images/assets/er_tree.png "example")

     在这个例子中，我们在树的节点和它们的子节点（或者它们的父节点）之间有一个递归的 O2M 关系。树中的每个节点都有许多子节点，并且有一个父节点。 如果节点 A 将 B 添加到其子节点，则 B 可以使用所有者边获取其所有者。
     ```
     // Edges of the Node.
     func (Node) Edges() []ent.Edge {
         return []ent.Edge{
             // 对于相同类型的关系，可以在同一个构建器中声明边及其引用。
             edge.To("children", Node.Type).
                 From("parent").
                 Unique(),
             //edge.To("children", Node.Type),
             //edge.From("parent", Node.Type).
             //    Ref("children").
             //    Unique(),
         }
     }
     ```
     与此Edge交互的API如下：
     ```
     func Do(ctx context.Context, client *ent.Client) error {
         root, err := client.Node.
             Create().
             SetValue(2).
             Save(ctx)
         if err != nil {
             return fmt.Errorf("creating the root: %w", err)
         }
         // Add additional nodes to the tree:
         //
         //       2
         //     /   \
         //    1     4
         //        /   \
         //       3     5
         //
         // Unlike `Save`, `SaveX` panics if an error occurs.
         n1 := client.Node.
             Create().
             SetValue(1).
             SetParent(root).
             SaveX(ctx)
         n4 := client.Node.
             Create().
             SetValue(4).
             SetParent(root).
             SaveX(ctx)
         n3 := client.Node.
             Create().
             SetValue(3).
             SetParent(n4).
             SaveX(ctx)
         n5 := client.Node.
             Create().
             SetValue(5).
             SetParent(n4).
             SaveX(ctx)
     
         fmt.Println("Tree leafs", []int{n1.Value, n3.Value, n5.Value})
         // Output: Tree leafs [1 3 5]
     
         // Get all leafs (nodes without children).
         // Unlike `Int`, `IntX` panics if an error occurs.
         ints := client.Node.
             Query().                             // All nodes.
             Where(node.Not(node.HasChildren())). // Only leafs.
             Order(ent.Asc(node.FieldValue)).     // Order by their `value` field.
             GroupBy(node.FieldValue).            // Extract only the `value` field.
             IntsX(ctx)
         fmt.Println(ints)
         // Output: [1 3 5]
     
         // Get orphan nodes (nodes without parent).
         // Unlike `Only`, `OnlyX` panics if an error occurs.
         orphan := client.Node.
             Query().
             Where(node.Not(node.HasParent())).
             OnlyX(ctx)
         fmt.Println(orphan)
         // Output: Node(id=1, value=2)
     
         return nil
     }
     ```
     请注意，可以使用 Edge Field 选项配置外键列并将其公开为实体字段，如下所示:
     ```
     // Fields of the Node.
     func (Node) Fields() []ent.Field {
         return []ent.Field{
             field.Int("parent_id").
                 Optional(),
         }
     }
     
     // Edges of the Node.
     func (Node) Edges() []ent.Edge {
         return []ent.Edge{
             edge.To("children", Node.Type).
                 From("parent").
                 Unique().
                 Field("parent_id"),
         }
     }
     ```

     完整例子请去[GitHub](https://github.com/ent/ent/tree/master/examples/o2mrecur)

  6. M2M Two Types (M2M 两种类型) `简单点说，就是2个表里的多对多关系`
    ![关系图例](https://entgo.io/images/assets/er_user_groups.png "example")

     在这个组-用户示例中，我们在组与其用户之间建立了 M2M 关系。 每个组有多个用户，每个用户可以加入多个组。
     ```
     // Edges of the Group.
     func (Group) Edges() []ent.Edge {
         return []ent.Edge{
             edge.To("users", User.Type),
         }
     }
     ```
     ```
     // Edges of the User.
     func (User) Edges() []ent.Edge {
         return []ent.Edge{
             edge.From("groups", Group.Type).
                 Ref("users"),
         }
     }
     ```
     与此edge交互的API如下：
     ```
     func Do(ctx context.Context, client *ent.Client) error {
         // Unlike `Save`, `SaveX` panics if an error occurs.
         hub := client.Group.
             Create().
             SetName("GitHub").
             SaveX(ctx)
         lab := client.Group.
             Create().
             SetName("GitLab").
             SaveX(ctx)
         a8m := client.User.
             Create().
             SetAge(30).
             SetName("a8m").
             AddGroups(hub, lab).
             SaveX(ctx)
         nati := client.User.
             Create().
             SetAge(28).
             SetName("nati").
             AddGroups(hub).
             SaveX(ctx)
     
         // Query the edges.
         groups, err := a8m.
             QueryGroups().
             All(ctx)
         if err != nil {
             return fmt.Errorf("querying a8m groups: %w", err)
         }
         fmt.Println(groups)
         // Output: [Group(id=1, name=GitHub) Group(id=2, name=GitLab)]
     
         groups, err = nati.
             QueryGroups().
             All(ctx)
         if err != nil {
             return fmt.Errorf("querying nati groups: %w", err)
         }
         fmt.Println(groups)
         // Output: [Group(id=1, name=GitHub)]
     
         // Traverse the graph.
         users, err := a8m.
             QueryGroups().                                           // [hub, lab]
             Where(group.Not(group.HasUsersWith(user.Name("nati")))). // [lab]
             QueryUsers().                                            // [a8m]
             QueryGroups().                                           // [hub, lab]
             QueryUsers().                                            // [a8m, nati]
             All(ctx)
         if err != nil {
             return fmt.Errorf("traversing the graph: %w", err)
         }
         fmt.Println(users)
         // Output: [User(id=1, age=30, name=a8m) User(id=2, age=28, name=nati)]
         return nil
     }
     ```

     这将生成一个新的表group_users来保存多对多关系

     完整例子请去[GitHub](https://github.com/ent/ent/tree/master/examples/m2m2types)

  7. M2M Same Type (M2M 同类型) `简单点说，就是1个表里的多对多关系`
    ![关系图例](https://entgo.io/images/assets/er_following_followers.png "example")

     在下面的关注者示例中，我们在用户与其关注者之间建立了 M2M 关系。 每个用户可以关注多个用户，并且可以拥有多个关注者。

     ```
     // Edges of the User.
     func (User) Edges() []ent.Edge {
       return []ent.Edge{
           // 对于相同类型的关系，可以在同一个构建器中声明边及其引用。
           edge.To("following", User.Type).
               From("followers"),
           //edge.To("following", User.Type),
           //edge.From("followers", User.Type).
           //    Ref("following"),
       }
     }
     ```
     与此edge交互的API如下：
     ```
     func Do(ctx context.Context, client *ent.Client) error {
         // Unlike `Save`, `SaveX` panics if an error occurs.
         a8m := client.User.
             Create().
             SetAge(30).
             SetName("a8m").
             SaveX(ctx)
         nati := client.User.
             Create().
             SetAge(28).
             SetName("nati").
             AddFollowers(a8m).
             SaveX(ctx)

         // Query following/followers:

         flw := a8m.QueryFollowing().AllX(ctx)
         fmt.Println(flw)
         // Output: [User(id=2, age=28, name=nati)]

         flr := a8m.QueryFollowers().AllX(ctx)
         fmt.Println(flr)
         // Output: []

         flw = nati.QueryFollowing().AllX(ctx)
         fmt.Println(flw)
         // Output: []

         flr = nati.QueryFollowers().AllX(ctx)
         fmt.Println(flr)
         // Output: [User(id=1, age=30, name=a8m)]

         // Traverse the graph:

         ages := nati.
             QueryFollowers().       // [a8m]
             QueryFollowing().       // [nati]
             GroupBy(user.FieldAge). // [28]
             IntsX(ctx)
         fmt.Println(ages)
         // Output: [28]

         names := client.User.
             Query().
             Where(user.Not(user.HasFollowers())).
             GroupBy(user.FieldName).
             StringsX(ctx)
         fmt.Println(names)
         // Output: [a8m]
         return nil
     }
     ```


     完整例子请去[GitHub](https://github.com/ent/ent/tree/master/examples/m2mrecur) 

  8. M2M Bidirectional (M2M 双向)  `简单点说，就是2个表里的双向多对多关系`
    ![关系图例](https://entgo.io/images/assets/er_user_friends.png "example")

     在这个用户-朋友示例中，我们有一个名为朋友的对称 M2M 关系。 每个用户可以有很多朋友。 如果用户 A 成为 B 的朋友，则 B 也是 A 的朋友。请注意，在双向边的情况​​下没有所有者/逆项。
     ```
     // Edges of the User.
     func (User) Edges() []ent.Edge {
         return []ent.Edge{
             edge.To("friends", User.Type),
         }
     }
     ```
     与此edge交互的API如下：
     ```
     func Do(ctx context.Context, client *ent.Client) error {
         // Unlike `Save`, `SaveX` panics if an error occurs.
         a8m := client.User.
             Create().
             SetAge(30).
             SetName("a8m").
             SaveX(ctx)
         nati := client.User.
             Create().
             SetAge(28).
             SetName("nati").
             AddFriends(a8m).
             SaveX(ctx)
     
         // Query friends. Unlike `All`, `AllX` panics if an error occurs.
         friends := nati.
             QueryFriends().
             AllX(ctx)
         fmt.Println(friends)
         // Output: [User(id=1, age=30, name=a8m)]
     
         friends = a8m.
             QueryFriends().
             AllX(ctx)
         fmt.Println(friends)
         // Output: [User(id=2, age=28, name=nati)]
     
         // Query the graph:
         friends = client.User.
             Query().
             Where(user.HasFriends()).
             AllX(ctx)
         fmt.Println(friends)
         // Output: [User(id=1, age=30, name=a8m) User(id=2, age=28, name=nati)]
         return nil
     }
     ```
     
     执行会创建新表user_friends保存好友关系

     完整例子请去[GitHub](https://github.com/ent/ent/tree/master/examples/m2mbidi)


- Edge Field
  
  Edge的字段选项允许用户将外键公开为schema上的常规字段。 请注意，只有持有外键 (edge-id) 的关系才允许使用此选项。
  ```
  // Fields of the Post.
  func (Post) Fields() []ent.Field {
      return []ent.Field{
          field.Int("author_id").
              Optional(),
      }
  }
  
  // Edges of the Post.
  func (Post) Edges() []ent.Edge {
      return []ent.Edge{
          edge.To("author", User.Type).
              // Bind the "author_id" field to this edge.
              Field("author_id").
              Unique(),
      }
  }
  ```
  与此edge交互的API如下：
  ```
  func Do(ctx context.Context, client *ent.Client) error {
    ps, err := client.Post.
		Create().
		SetAuthorID(8).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	u := client.User.Query().Where(user.ID(ps.AuthorID)).OnlyX(ctx)

	p := client.Post.Query().
		Where(post.AuthorID(u.ID)).
		OnlyX(ctx)

	fmt.Println(p) // Access the "author" foreign-key.
	return nil
  }
  ```
  
  - Migration To Edge Fields
    如 StorageKey 部分所述，Ent 通过 edge.To 配置边缘存储键（例如外键）。 因此，如果要向现有边（已作为列存在于数据库中）添加字段，则需要使用   StorageKey 选项进行设置，如下所示：
    ```
    // Fields of the Post.
    func (Post) Fields() []ent.Field {
        return []ent.Field{
           field.Int("author_id").
    +           StorageKey("post_author").
               Optional(),
        }
    }
    ```

- Required
  
  可以使用构建器上的 Required 方法根据需要在实体创建中定义边。
  ```
  // Edges of the Card.
  func (Card) Edges() []ent.Edge {
      return []ent.Edge{
          edge.From("owner", User.Type).
              Ref("card").
              Unique().
              Required(),
      }
  }
  ``` 如上例子，没有所有者就不能创建卡片实体。

- StorageKey
  默认情况下，Ent 由边缘所有者（保存 edge.To 的架构）配置边缘存储键，而不是通过反向引用（edge.From）。 这是因为反向引用是可选的，可以删除。

  为了对边缘使用自定义存储配置，请使用 StorageKey 方法，如下所示：

  ``` 请注意注释中的各种关系场景
  // Edges of the User.
  func (User) Edges() []ent.Edge {
      return []ent.Edge{
          edge.To("pets", Pet.Type).
              // Set the column name in the "pets" table for O2M   relationship.
              StorageKey(edge.Column("owner_id")),
          edge.To("cars", Car.Type).
              // Set the symbol of the foreign-key constraint for O2M   relationship.
              StorageKey(edge.Symbol("cars_owner_id")),
          edge.To("friends", User.Type).
              // Set the join-table, and the column names for a M2M   relationship.
              StorageKey(edge.Table("friends"), edge.Columns("user_id",   "friend_id")),
          edge.To("groups", Group.Type).
              // Set the join-table, its column names and the symbols
              // of the foreign-key constraints for M2M relationship.
              StorageKey(
                  edge.Table("groups"),
                  edge.Columns("user_id", "group_id"),
                  edge.Symbols("groups_id1", "groups_id2")
              ),
      }
  }
  ```

- Struct Tags
  可以使用 StructTag 方法将自定义结构标记添加到生成的实体中。 请注意，如果未提供此选项，或提供但不包含 json 标签，则默认 json 标签将添加字段名称。

  ```
  // Edges of the User.
  func (User) Edges() []ent.Edge {
      return []ent.Edge{
          edge.To("pets", Pet.Type).
              // Override the default json tag "pets" with "owner" for   O2M relationship.
              StructTag(`json:"owner"`),
      }
  }
  ```

- Indexes
  索引可以定义在多个字段和某些类型的边上。 但是，您应该注意，这是目前仅限 SQL 的功能。

- Annotations
  注释用于在代码生成中将任意元数据附加到边缘对象。 模板扩展可以检索此元数据并在其模板中使用它。

  请注意，元数据对象必须可序列化为 JSON 原始值（例如 struct、map 或 slice）。
  ```
  // Pet schema.
  type Pet struct {
      ent.Schema
  }
  
  // Edges of the Pet.
  func (Pet) Edges() []ent.Edge {
      return []ent.Field{
          edge.To("owner", User.Type).
              Ref("pets").
              Unique().
              Annotations(entgql.Annotation{
                  OrderField: "OWNER",
              }),
      }
  }
  ```

### 索引(Indexes)
- 多个字段
  
  索引可以在一个或多个字段上配置以提高数据检索速度，也可以定义其唯一性。
  ```
  package schema

  import (
      "entgo.io/ent"
      "entgo.io/ent/schema/index"
  )
  
  // User holds the schema definition for the User entity.
  type User struct {
      ent.Schema
  }
  
  func (User) Indexes() []ent.Index {
      return []ent.Index{
          // 非唯一约束索引
          index.Fields("field1", "field2"),
          // 唯一约束索引
          index.Fields("first_name", "last_name").
              Unique(),
      }
  }
  ```
  请注意，如果要为单个字段设置唯一约束，请在字段生成器上使用 Unique 方法，如下：
  ```
  func (User) Fields() []ent.Field {
      return []ent.Field{
          field.String("phone").
              Unique(),
      }
  }
  ```

- 边上的索引
  
  索引可以为字段和边的组合进行配置。 主要用法是在特定关系下设置字段的唯一性。 让我们来看一个例子：

  ![关系图例](https://entgo.io/images/assets/er_city_streets.png "example")

  在上面的图例中，我们有一个带了许多 Street 的 City，并且我们想设置每个城市的街道名称都是唯一的。
  ``` city.go
  // City holds the schema definition for the City entity.
  type City struct {
      ent.Schema
  }
  
  // Fields of the City.
  func (City) Fields() []ent.Field {
      return []ent.Field{
          field.String("name"),
      }
  }
  
  // Edges of the City.
  func (City) Edges() []ent.Edge {
      return []ent.Edge{
          edge.To("streets", Street.Type),
      }
  }
  ```
  ```street.go
  // Street holds the schema definition for the Street entity.
  type Street struct {
      ent.Schema
  }
  
  // Fields of the Street.
  func (Street) Fields() []ent.Field {
      return []ent.Field{
          field.String("name"),
      }
  }
  
  // Edges of the Street.
  func (Street) Edges() []ent.Edge {
      return []ent.Edge{
          edge.From("city", City.Type).
              Ref("streets").
              Unique(),
      }
  }
  
  // Indexes of the Street.
  func (Street) Indexes() []ent.Index {
      return []ent.Index{
          index.Fields("name").
              Edges("city").
              Unique(),
      }
  }
  ```
  ```交互api如下
  func Do(ctx context.Context, client *ent.Client) error {
      // 和 `Save`不同，当出现错误是 `SaveX` 抛出panic。
      tlv := client.City.
          Create().
          SetName("TLV").
          SaveX(ctx)
      nyc := client.City.
          Create().
          SetName("NYC").
          SaveX(ctx)
      // Add a street "ST" to "TLV".
      client.Street.
          Create().
          SetName("ST").
          SetCity(tlv).
          SaveX(ctx)
      // 这一步的操作将会失败 because "ST"
      // 因为 "ST" 已经创建于 "TLV" 之下
      _, err := client.Street.
          Create().
          SetName("ST").
          SetCity(tlv).
          Save(ctx)
      if err == nil {
          return fmt.Errorf("expecting creation to fail")
      }
      // 将街道 "ST" 添加到 "NYC"
      client.Street.
          Create().
          SetName("ST").
          SetCity(nyc).
          SaveX(ctx)
      return nil
  }
  ```

  完整例子请去[GitHub](https://github.com/ent/ent/tree/master/examples/edgeindex)

- 边字段上的索引
  
  目前边列总是添加在字段列之后。 但是，某些索引要求这些列排在第一位以实现特定的优化。 您可以通过使用边字段来解决此问题。
  
  ```
  // Card holds the schema definition for the Card entity.
  type Card struct {
      ent.Schema
  }
  // Fields of the Card.
  func (Card) Fields() []ent.Field {
      return []ent.Field{
          field.String("number").
              Optional(),
          field.Int("owner_id").
              Optional(),
      }
  }
  // Edges of the Card.
  func (Card) Edges() []ent.Edge {
      return []ent.Edge{
          //此处需要在user.go 增加edge.To构建
          edge.From("owner", User.Type).
              Ref("card").
              Field("owner_id").
              Unique(),
      }
  }
  // Indexes of the Card.
  func (Card) Indexes() []ent.Index {
      return []ent.Index{
          index.Fields("owner_id", "number"),
      }
  }
  ``` 

_此处所说的排列应该是指索引排列_

- Dialect Support
  索引目前仅支持 SQL 方言，不支持 Gremlin。 使用注释允许方言特定功能。 例如，为了在 MySQL 中使用索引前缀，请使用以下配置：

  ```
  // Indexes of the User.
  func (User) Indexes() []ent.Index {
      return []ent.Index{
          index.Fields("description").
              Annotations(entsql.Prefix(128)),
          index.Fields("c1", "c2", "c3").
              Annotation(
                  entsql.PrefixColumn("c1", 100),
                  entsql.PrefixColumn("c2", 200),
              )
      }
  }
  ```
  简单点说，就是支持SQL数据库的一些特有特性,比如索引前缀

### Mixin
Mixin 允许您创建可重用的 ent.Schema 代码片段，这些代码可以使用组合注入到其他Schema中。

ent.Mixin接口如下：
```
type Mixin interface {
    // Fields 数组的返回值会被添加到 schema 中。
    Fields() []Field
    // Edges 数组返回值会被添加到 schema 中。
    Edges() []Edge
    // Indexes 数组的返回值会被添加到 schema 中。
    Indexes() []Index
    // Hooks 数组的返回值会被添加到 schema 中。
    // 请注意，mixin 的钩子会在 schema 的钩子之前被执行。
    Hooks() []Hook
    // Policy 数组的返回值会被添加到 schema 中。
    //请注意，mixin（混入）的 policy（策略）会在 schema 的 policy 之前被执行。
    Policy() Policy
    // Annotations 方法返回要添加到 Schema 中的注解列表。
    Annotations() []schema.Annotation
}
```

- 示例
  Mixin 的一个常见用例是将公共字段列表混合到您的架构中。

  ``` mixin.go
  package schema
  
  import (
      "time"
  
      "entgo.io/ent"
      "entgo.io/ent/schema/field"
      "entgo.io/ent/schema/mixin"
  )
  
  // -------------------------------------------------
  // Mixin 实现接口，写自定义模版
  
  // TimeMixin implements the ent.Mixin for sharing
  // time fields with package schemas.
  type TimeMixin struct{
      // We embed the `mixin.Schema` to avoid
      // implementing the rest of the methods.
      mixin.Schema
  }
  
  func (TimeMixin) Fields() []ent.Field {
      return []ent.Field{
          field.Time("created_at").
              Immutable().
              Default(time.Now),
          field.Time("updated_at").
              Default(time.Now).
              UpdateDefault(time.Now),
      }
  }
  
  // DetailsMixin implements the ent.Mixin for sharing
  // entity details fields with package schemas.
  type DetailsMixin struct{
      // We embed the `mixin.Schema` to avoid
      // implementing the rest of the methods.
      mixin.Schema
  }
  
  func (DetailsMixin) Fields() []ent.Field {
      return []ent.Field{
          field.Int("age").
              Positive(),
          field.String("name").
              NotEmpty(),
      }
  }
  ```
  
  ``` user.go,pet.go
  // -------------------------------------------------
  // Schema definition
  
  // User schema mixed-in the TimeMixin and DetailsMixin fields and therefore
  // has 5 fields: `created_at`, `updated_at`, `age`, `name` and `nickname`.
  type User struct {
      ent.Schema
  }
  
  + //user.go 增加Mixin配置字段
  func (User) Mixin() []ent.Mixin {
      return []ent.Mixin{
          TimeMixin{},
          DetailsMixin{},
      }
  }
  
  func (User) Fields() []ent.Field {
      return []ent.Field{
          field.String("nickname").
              Unique(),
      }
  }

  // Pet schema mixed-in the DetailsMixin fields and therefore
  // has 3 fields: `age`, `name` and `weight`.
  type Pet struct {
      ent.Schema
  }
  
  + //user.go 增加Mixin配置字段
  func (Pet) Mixin() []ent.Mixin {
      return []ent.Mixin{
          DetailsMixin{},
      }
  }
  
  func (Pet) Fields() []ent.Field {
      return []ent.Field{
          field.Float("weight"),
      }
  }
  ```


- 内置Mixin
  mixin包提供了一些内置的mixin，它们可以用于在schema中添加create_time和update_time 字段。

  若要使用它们，请将 mixin.Time mixin 添加到您的schema，如下：
  ```
  func (Pet) Mixin() []ent.Mixin {
      return []ent.Mixin{
          mixin.Time{},
          // Or, mixin.CreateTime only for create_time
          // and mixin.UpdateTime only for update_time.
      }
  }
  ```


### 注解(Annotations)
结构注解(Schema annotations) 允许附加元数据到结构对象(例如字段和边) 上面，并且将元数据注入到外部模板中。 注解是一种Go类型，它能进行JSON序列化(例如 struct, map 或 slice)，并且需要实现Annotation接口。内置注解能够配置不同的存储驱动(例如 SQL)，控制代码生成输出。

- 自定义表名
  使用 entsql 注解的类型, 可以自定义表名，如下所示：
  ```
  package schema

  import (
      "entgo.io/ent"
      "entgo.io/ent/dialect/entsql"
      "entgo.io/ent/schema"
      "entgo.io/ent/schema/field"
  )
  
  // User类型持有用户实体的结构(schema)定义
  type User struct {
      ent.Schema
  }
  
  // 用户实体的注解
  func (User) Annotations() []schema.Annotation {
      return []schema.Annotation{
          entsql.Annotation{Table: "Users"},
      }
  }
  
  // 用户实体的字段
  func (User) Fields() []ent.Field {
      return []ent.Field{
          field.Int("age"),
          field.String("name"),
      }
  }
  ```

- 外键配置
  Ent允许对外键的创建进行定制，并且为ON DELETE子句提供referential action。
  ```
  package schema

  import (
      "entgo.io/ent"
      "entgo.io/ent/dialect/entsql"
      "entgo.io/ent/schema/edge"
      "entgo.io/ent/schema/field"
  )
  
  // User类型持有用户实体的结构(schema)定义
  type User struct {
      ent.Schema
  }
  
  // 用户实体的字段
  func (User) Fields() []ent.Field {
      return []ent.Field{
          field.String("name").
              Default("Unknown"),
      }
  }
  
  // 用户实体的关系
  func (User) Edges() []ent.Edge {
      return []ent.Edge{
          edge.To("posts", Post.Type).
              Annotations(entsql.Annotation{
                  OnDelete: entsql.Cascade,   //只删除关系ID，不会删除整条数据
              }),
      }
  }
  ```上面的示例配置了外键，将父表的删除操作关联到子表，对子表中匹配的数据也进行删除。

## 代码生成

### 引言
- 安装

  本项目有一个叫做 ent 的代码工具。 若要安装 ent 运行以下命令：
  `go get entgo.io/ent/cmd/ent`

- 初始化一个新的Schema

  生成一个或多个 schema 模板，运行 ent init 如下：
  `go run entgo.io/ent/cmd/ent init User Pet`
  
  init 将在 ent/schema 目录下创建 2个 schemas (user.go 和 pet.go)。 如果 ent 目录不存在，将自动创建。 一般约定将 ent 目录放在项目的根目录下。

- 生成资源文件

  每次添加或修改 fields 和 edges后, 你都需要生成新的实体. 在项目的根目录执行 ent generate或直接执行go generate命令重新生成资源文件:
  `go generate ./ent`

  generate将会按照schema模板生成生成以下资源:
  * Client 和 Tx 对象用于与graph的交互。
  * 每个schema对应的增删改查。
  * 每个schema的实体对象(Go结构体)。
  * 含常量和查询条件的包，用于与生成器交互。
  * 用于数据迁移的migrate包。
  * 用于中间件的hook包。

- entc和ent之间的版本兼容性
  在项目中使用 ent CLI时，需要确保CLI使用的版本与项目使用的 ent 版本相同。

  保证此的一个方法是通过 go generate 来使用 go.mod 内所定义的 ent CLI版本。如果您的项目没有使用 Go modules, 请设置一个：
  `go mod init <project>`

  然后执行以下命令，以便将 ent 添加到您的 go.mod 文件：
  `go get entgo.io/ent/cmd/ent`

  将 generate.go 文件添加到你项目的<project>/ent 目录中：

  最后，你可以执行 go generate ./ent ，以便在您的项目方案中执行 ent 代码生成。

- 代码生成选项
  关于代码生成选项的更多信息，执行 ent generate -h：
  ```
  generate go code for the schema directory

  Usage:
    ent generate [flags] path
  
  Examples:
    ent generate ./ent/schema
    ent generate github.com/a8m/x
  
  Flags:
        --feature strings                         extend codegen with additional features
        --header string                           override codegen header
    -h, --help                                    help for generate
        --idtype [int int64 uint uint64 string]   type of the id field (default int)
        --storage string                          storage driver to support in codegen (default "sql")
        --target string                           target directory for codegen
        --template strings                        external templates to execute
  
  ```

- Storage选项
  ent 可以为 SQL 和 Gremlin 生成资源。 默认是 SQL

- 外部模板
  ent 接受执行外部 Go 模板文件。 如果模板名称已由 ent定义，它将覆盖现有的模板。 否则，它将把执行后的输出写入到 与模板相同名称的文件。 支持参数 file, dir 和 glob 如下所示：
  `go run entgo.io/ent/cmd/ent generate --template <dir-path> --template glob="path/to/*.tmpl" ./ent/schema`

  更多信息和示例可在 [外部模板](#外部模版) 中找到

- 使用 entc
  运行 ent CLI 的另一个方式是将其作为一个包，如下所示：
  ```
  package main

  import (
      "log"
  
      "entgo.io/ent/entc"
      "entgo.io/ent/entc/gen"
      "entgo.io/ent/schema/field"
  )
  
  func main() {
      err := entc.Generate("./schema", &gen.Config{
          Header: "// Your Custom Header",
          IDType: &field.TypeInfo{Type: field.TypeInt},
      })
      if err != nil {
          log.Fatal("running ent codegen:", err)
      }
  }
  ```
  完整示例请参阅 [GitHub](https://github.com/ent/ent/tree/master/examples/entcpkg)

- Schema描述
  要获取您Schema的描述，请执行：
  `go run entgo.io/ent/cmd/ent describe ./ent/schema`

  ```
  Pet:
    +-------+---------+--------+----------+----------+---------+---------------+-----------+-----------------------+------------+
    | Field |  Type   | Unique | Optional | Nillable | Default | UpdateDefault | Immutable |       StructTag       | Validators |
    +-------+---------+--------+----------+----------+---------+---------------+-----------+-----------------------+------------+
    | id    | int     | false  | false    | false    | false   | false         | false     | json:"id,omitempty"   |          0 |
    | name  | string  | false  | false    | false    | false   | false         | false     | json:"name,omitempty" |          0 |
    +-------+---------+--------+----------+----------+---------+---------------+-----------+-----------------------+------------+
    +-------+------+---------+---------+----------+--------+----------+
    | Edge  | Type | Inverse | BackRef | Relation | Unique | Optional |
    +-------+------+---------+---------+----------+--------+----------+
    | owner | User | true    | pets    | M2O      | true   | true     |
    +-------+------+---------+---------+----------+--------+----------+

  User:
    +-------+---------+--------+----------+----------+---------+---------------+-----------+-----------------------+------------+
    | Field |  Type   | Unique | Optional | Nillable | Default | UpdateDefault | Immutable |       StructTag       | Validators |
    +-------+---------+--------+----------+----------+---------+---------------+-----------+-----------------------+------------+
    | id    | int     | false  | false    | false    | false   | false         | false     | json:"id,omitempty"   |          0 |
    | age   | int     | false  | false    | false    | false   | false         | false     | json:"age,omitempty"  |          0 |
    | name  | string  | false  | false    | false    | false   | false         | false     | json:"name,omitempty" |          0 |
    +-------+---------+--------+----------+----------+---------+---------------+-----------+-----------------------+------------+
    +------+------+---------+---------+----------+--------+----------+
    | Edge | Type | Inverse | BackRef | Relation | Unique | Optional |
    +------+------+---------+---------+----------+--------+----------+
    | pets | Pet  | false   |         | O2M      | false  | true     |
    +------+------+---------+---------+----------+--------+----------+
  ```

- 代码生成Hooks (钩子)

  entc 软件包提供了一个将钩子(中间件) 添加到代码生成阶段的方法。 这个方法适合在你需要对graph schema进行自定义验证或者生成附加资源时使用.
  ```
  // +build ignore

  package main
  
  import (
      "fmt"
      "log"
      "reflect"
  
      "entgo.io/ent/entc"
      "entgo.io/ent/entc/gen"
  )
  
  func main() {
      err := entc.Generate("./schema", &gen.Config{
          Hooks: []gen.Hook{
              EnsureStructTag("json"),
          },
      })
      if err != nil {
          log.Fatalf("running ent codegen: %v", err)
      }
  }
  
  // EnsureStructTag ensures all fields in the graph have a specific tag name.
  func EnsureStructTag(name string) gen.Hook {
      return func(next gen.Generator) gen.Generator {
          return gen.GenerateFunc(func(g *gen.Graph) error {
              for _, node := range g.Nodes {
                  for _, field := range node.Fields {
                      tag := reflect.StructTag(field.StructTag)
                      if _, ok := tag.Lookup(name); !ok {
                          return fmt.Errorf("struct tag %q is missing for field %s.%s", name, node.Name, field.Name)
                      }
                  }
              }
              return next.Generate(g)
          })
      }
  }
  ```
  [详情](#hooks钩子)

- 特性开关
  entc 软件包提供了一系列代码生成特性，可以自行选择使用。特性开关可以通过 CLI 标志或作为参数提供给 gen 包。

  用法：
  `go run entgo.io/ent/cmd/ent generate --feature privacy,entql ./ent/schema`

### 增删改查API

- 创建新的客户端
  ``` MySql
  package main

  import (
      "log"
  
      "<project>/ent"
  
      _ "github.com/go-sql-driver/mysql"
  )
  
  func main() {
      client, err := ent.Open("mysql", "<user>:<pass>@tcp(<host>:<port>)/<database>?parseTime=True")
      if err != nil {
          log.Fatal(err)
      }
      defer client.Close()
  }
  ```
  ```PostgreSQL
  package main

  import (
      "log"
  
      "<project>/ent"
  
      _ "github.com/lib/pq"
  )
  
  func main() {
      client, err := ent.Open("postgres","host=<host> port=<port> user=<user> dbname=<database> password=<pass>")
      if err != nil {
          log.Fatal(err)
      }
      defer client.Close()
  }
  ```
  ```Sqlite
  package main

  import (
      "log"
  
      "<project>/ent"
  
      _ "github.com/mattn/go-sqlite3"
  )
  
  func main() {
      client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
      if err != nil {
          log.Fatal(err)
      }
      defer client.Close()
  }
  ```
  ```Gremlin (AWS Neptune)
  package main

  import (
      "log"
  
      "<project>/ent"
  )
  
  func main() {
      client, err := ent.Open("gremlin", "http://localhost:8182")
      if err != nil {
          log.Fatal(err)
      }
  }
  ```

- 创建一个实体
  
  ``` 通过 Save 保存一个用户.
  a8m, err := client.User.   // UserClient.
    Create().               // 用户创建构造器
    SetName("a8m").         // 设置字段的值
    SetNillableAge(age).    // 忽略nil检查
    AddGroups(g1, g2).      // 添加多个边
    SetSpouse(nati).        // 设置单个边
    Save(ctx)               // 创建并返回
  ```
 
  ```通过 SaveX 保存一个宠物; 和 Save 不一样， SaveX 在出错时 panic。
   pedro := client.Pet.    // PetClient.
    Create().           // 宠物创建构造器
    SetName("pedro").   // 设置字段的值
    SetOwner(a8m).      // 设置主人 (唯一的边)
     SaveX(ctx)          // 创建并返回
  ```

- 批量创建
  
  ``` 通过 Save 批量保存宠物
   names := []string{"pedro", "xabi", "layla"}
   bulk := make([]*ent.PetCreate, len(names))
   for i, name := range names {
       bulk[i] = client.Pet.Create().SetName(name).SetOwner(a8m)
   }
   pets, err := client.Pet.CreateBulk(bulk...).Save(ctx)
  ```

- 更新单个实体
  
  ``` 更新一个数据库内的实体。
    a8m, err = a8m.Update().    // 用户更新构造器
        RemoveGroup(g2).        // 移除特定的边
        ClearCard().            // 清空唯一的边
        SetAge(30).             // 设置字段的值
        Save(ctx)               // 保存并返回
  ```

- 通过ID更新
  
  ```
    pedro, err := client.Pet.   // PetClient.
        UpdateOneID(id).        // 宠物更新构造器
        SetName("pedro").       // 设置名字字段
        SetOwnerID(owner).      // 通过ID设置唯一的边
        Save(ctx)               // 保存并返回
  ```

- 批量更新
  
  ``` 通过断言筛选
    n, err := client.User.          // UserClient.
        Update().                   // 宠物更新构造器
        Where(                      //
            user.Or(                // (age >= 30 OR name = "bar") 
                user.AgeGT(30),     //
                user.Name("bar"),   // AND
            ),                      //  
            user.HasFollowers(),    // UserHasFollowers()  
        ).                          //
        SetName("foo").             // 设置名字字段
        Save(ctx)                   // 执行并返回
  ```

  ``` 通过边上的断言筛选
      n, err := client.User.      // UserClient.
      Update().                   // 宠物更新构造器
      Where(                      // 
          user.HasFriendsWith(    // UserHasFriendsWith (
              user.Or(            //   age = 20
                  user.Age(20),   //      OR
                  user.Age(30),   //   age = 30
              )                   // )
          ),                      //
      ).                          //
      SetName("a8m").             // 设置名字字段
      Save(ctx)                   // 执行并返回
  ```

- Upsert One
  Ent 使用 sql/upsert 功能标志支持 upsert 记录。
  ```
  err := client.User.
      Create().
      SetAge(30).
      SetName("Ariel").
      OnConflict().
      // Use the new values that were set on create.
      UpdateNewValues().
      Exec(ctx)
  
  id, err := client.User.
      Create().
      SetAge(30).
      SetName("Ariel").
      OnConflict().
      // Use the "age" that was set on create.
      UpdateAge().
      // Set a different "name" in case of conflict.
      SetName("Mashraki").
      ID(ctx)
  
  // Customize the UPDATE clause.
  err := client.User.
      Create().
      SetAge(30).
      SetName("Ariel").
      OnConflict().
      UpdateNewValues().
      // Override some of the fields with a custom update.
      Update(func(u *ent.UserUpsert) {
          u.SetAddress("localhost")
      }).
      Exec(ctx)
  ```
  ``` In PostgreSQL, the conflict target is required:
    // Setting the column names using the fluent API.
    err := client.User.
        Create().
        SetName("Ariel").
        OnConflictColumns(user.FieldName).
        UpdateNewValues().
        Exec(ctx)
    
    // Setting the column names using the SQL API.
    err := client.User.
        Create().
        SetName("Ariel").
        OnConflict(
            sql.ConflictColumns(user.FieldName),    
        ).
        UpdateNewValues().
        Exec(ctx)
    
    // Setting the constraint name using the SQL API.
    err := client.User.
        Create().
        SetName("Ariel").
        OnConflict(
            sql.ConflictConstraint(constraint), 
        ).
        UpdateNewValues().
        Exec(ctx)
  ```
  ``` 自定义执行的语句，使用 SQL API：
    id, err := client.User.
        Create().
        OnConflict(
            sql.ConflictColumns(...),
            sql.ConflictWhere(...),
            sql.UpdateWhere(...),
        ).
        Update(func(u *ent.UserUpsert) {
            u.SetAge(30)
            u.UpadteName()
        }).
        ID(ctx)
    
    // INSERT INTO "users" (...) VALUES (...) ON CONFLICT WHERE ... DO UPDATE SET ... WHERE ...
  ```

- Upsert Many
  ```
    err := client.User.             // UserClient
      CreateBulk(builders...).    // User bulk create.
      OnConflict().               // User bulk upsert.
      UpdateNewValues().          // Use the values that were set on create in case of conflict.
      Exec(ctx)                   // Execute the statement.
  ```

- Query The Graph
  
  ``` 获取所有拥有关注者的用户
   users, err := client.User.      // UserClient.
    Query().                    // User query builder.
    Where(user.HasFollowers()). // filter only users with followers.
    All(ctx)                    // query and return.
  ```

  ``` 获取特定用户的所有关注者。
    users, err := a8m.
      QueryFollowers().
      All(ctx)
  ```

  ``` 获取用户的追随者的所有宠物。
   pets, err := a8m.
    QueryFollowers().
    QueryPets().
    All(ctx)
  ```

- Field Selection
  ``` 获取所有宠物的名字
   names, err := client.Pet.
    Query().
    Select(pet.FieldName).
    Strings(ctx)
  ```
  
  ``` 选择部分对象和部分关联 获取所有宠物及其主人，但仅选择并填写 ID 和名称字段。
   pets, err := client.Pet.
    Query().
    Select(pet.FieldName).
    WithOwner(func (q *ent.UserQuery) {
        q.Select(user.FieldName)
    }).
    All(ctx)
  ```
  
  ``` 将所有宠物的年龄和名字扫描到自定义结构
   var v []struct {
        Age  int    `json:"age"`
        Name string `json:"name"`
    }
    err := client.Pet.
        Query().
        Select(pet.FieldAge, pet.FieldName).
        Scan(ctx, &v)
    if err != nil {
        log.Fatal(err)
    }
  ```

  ``` 更新一个实体并返回它的一部分
   pedro, err := client.Pet.
    UpdateOneID(id).
    SetAge(9).
    SetName("pedro").
    // Select 允许选择返回实体的一个或多个字段（列）。
	// 默认是选择实体架构中定义的所有字段。
    Select(pet.FieldName).
    Save(ctx)
  ``` 测试发现select此处无效(postgres),返回任然是全部字段

- Delete One
  ```删除一个实体
   err := client.User.
    DeleteOne(a8m).
    Exec(ctx)
  ```

  ```用ID删除
   err := client.User.
    DeleteOneID(id).
    Exec(ctx)
  ```

- Delete Many
  ```
   _, err := client.File.
    Delete().
    Where(file.UpdatedAtLT(date)).
    Exec(ctx)
  ```

- Mutation
  _其实就是解决一些公共的代码方法_

  每个生成的节点类型都有自己的 Mutation。 例如，所有用户构建器共享同一个生成的 UserMutation 对象。 但是，所有构建器类型都实现了通用的 ent.Mutation 接口。

  例如，为了编写在 ent.UserCreate 和 ent.UserUpdate 上应用一组方法的通用代码，请使用 UserMutation 对象：

  ``` 创建和更新用到了相同的代码
   func Do() {
        creator := client.User.Create()
        SetAgeName(creator.Mutation())
        updater := client.User.UpdateOneID(id)
        SetAgeName(updater.Mutation())
    }

    // SetAgeName sets the age and the name for any mutation.
    func SetAgeName(m *ent.UserMutation) {
        m.SetAge(32)
        m.SetName("Ariel")
    }
  ```

  在某些情况下，您希望对多种实体类型应用一组方法。 对于这种情况，要么使用通用的 ent.Mutation 接口，要么创建自己的接口。
  ```
   func Do() {
        creator1 := client.User.Create()
        SetName(creator1.Mutation(), "a8m")

        creator2 := client.Pet.Create()
        SetName(creator2.Mutation(), "pedro")
    }

    // SetNamer wraps the 2 methods for getting
    // and setting the "name" field in mutations.
    type SetNamer interface {
        SetName(string)
        Name() (string, bool)
    }

    func SetName(m SetNamer, name string) {
        if _, exist := m.Name(); !exist {
            m.SetName(name)
        }
    }
  ```
   
   
### Graph Traversal(表遍历)

![关系图例](https://entgo.io/images/assets/er_traversal_graph.png "example")

1. 创建3个模型Pet User Group
`go run entgo.io/ent/cmd/ent init Pet User Group`

2. 给模型添加必要的字段和边(关系)
   ``` pet.go
    // Pet holds the schema definition for the Pet entity.
     type Pet struct {
         ent.Schema
     }

     // Fields of the Pet.
     func (Pet) Fields() []ent.Field {
         return []ent.Field{
             field.String("name"),
         }
     }

     // Edges of the Pet.
     func (Pet) Edges() []ent.Edge {
         return []ent.Edge{
             edge.To("friends", Pet.Type),
             edge.From("owner", User.Type).
                 Ref("pets").
                 Unique(),
         }
     }
   ```
   ``` user.go
    // User holds the schema definition for the User entity.
     type User struct {
         ent.Schema
     }

     // Fields of the User.
     func (User) Fields() []ent.Field {
         return []ent.Field{
             field.Int("age"),
             field.String("name"),
         }
     }

     // Edges of the User.
     func (User) Edges() []ent.Edge {
         return []ent.Edge{
             edge.To("pets", Pet.Type),
             edge.To("friends", User.Type),
             edge.From("groups", Group.Type).
                 Ref("users"),
             edge.From("manage", Group.Type).
                 Ref("admin"),
         }
     }
   ```
   ``` group.go
    // Group holds the schema definition for the Group entity.
     type Group struct {
         ent.Schema
     }

     // Fields of the Group.
     func (Group) Fields() []ent.Field {
         return []ent.Field{
             field.String("name"),
         }
     }

     // Edges of the Group.
     func (Group) Edges() []ent.Edge {
         return []ent.Edge{
             edge.To("users", User.Type),
             edge.To("admin", User.Type).
                 Unique(),
         }
     }
   ```

3. 填充必要的数据
   ```
    func Gen(ctx context.Context, client *ent.Client) error {
         hub, err := client.Group.
             Create().
             SetName("Github").
             Save(ctx)
         if err != nil {
             return fmt.Errorf("failed creating the group: %w", err)
         }
         // Create the admin of the group.
         // Unlike `Save`, `SaveX` panics if an error occurs.
         dan := client.User.
             Create().
             SetAge(29).
             SetName("Dan").
             AddManage(hub).
             SaveX(ctx)

         // Create "Ariel" and its pets.
         a8m := client.User.
             Create().
             SetAge(30).
             SetName("Ariel").
             AddGroups(hub).
             AddFriends(dan).
             SaveX(ctx)
         pedro := client.Pet.
             Create().
             SetName("Pedro").
             SetOwner(a8m).
             SaveX(ctx)
         xabi := client.Pet.
             Create().
             SetName("Xabi").
             SetOwner(a8m).
             SaveX(ctx)

         // Create "Alex" and its pets.
         alex := client.User.
             Create().
             SetAge(37).
             SetName("Alex").
             SaveX(ctx)
         coco := client.Pet.
             Create().
             SetName("Coco").
             SetOwner(alex).
             AddFriends(pedro).
             SaveX(ctx)

         fmt.Println("Pets created:", pedro, xabi, coco)
         // Output:
         // Pets created: Pet(id=1, name=Pedro) Pet(id=2, name=Xabi) Pet(id=3, name=Coco)
         return nil
     }
   ```

4. 根据需求查询数据
   ![关系图例](https://entgo.io/images/assets/er_traversal_graph_gopher.png "example")

   上图的遍历查询从一个Group实体开始，继续到它的管理员（edge），继续到它的朋友（edge），得到他们的宠物（edge），得到每个宠物的朋友（edge），最终请求他们的主人。
    ```
      func Traverse(ctx context.Context, client *ent.Client) error {
         owner, err := client.Group.         // GroupClient.
             Query().                        // Query builder.
             Where(group.Name("Github")).    // Filter only Github group (only 1).
             QueryAdmin().                   // Getting Dan.
             QueryFriends().                 // Getting Dan's friends: [Ariel].
             QueryPets().                    // Their pets: [Pedro, Xabi].
             QueryFriends().                 // Pedro's friends: [Coco], Xabi's friends: [].
             QueryOwner().                   // Coco's owner: Alex.
             Only(ctx)                       // Expect only one entity to return in the query.
         if err != nil {
             return fmt.Errorf("failed querying the owner: %w", err)
         }
         fmt.Println(owner)
         // Output:
         // User(id=3, age=37, name=Alex)
         return nil
     }
    ```

    ![关系图例](https://entgo.io/images/assets/er_traversal_graph_gopher_query.png "example")
  
    我们想要获取所有拥有所有者（edge）的宠物（实体），该所有者（edge）是某个组管理员（edge）的朋友（edge）。
  
    _这句话有点费劲，实际上就是查询所有宠物主人的朋友是群组管理员的宠物_
   
    ```
      func Traverse(ctx context.Context, client *ent.Client) error {
         pets, err := client.Pet.
             Query().
             Where(
                 pet.HasOwnerWith(
                     user.HasFriendsWith(
                         user.HasManage(),
                     ),
                 ),
             ).
             All(ctx)
         if err != nil {
             return fmt.Errorf("failed querying the pets: %w", err)
         }
         fmt.Println(pets)
         // Output:
         // [Pet(id=1, name=Pedro) Pet(id=2, name=Xabi)]
         return nil
      }
    ```
    
    完整例子请去[GitHub](https://github.com/ent/ent/tree/master/examples/traversal)
  
### 预加载
- 概览
  ent 支持查询关联的实体 (通过其边)。 关联的实体会填充到返回对象的 Edges 字段中。

  ``` 查询所有用户及其宠物：
  users, err := client.User.
      Query().
      WithPets().
      All(ctx)
  if err != nil {
      return err
  }
  // The returned users look as follows:
  //
  //  [
  //      User {
  //          ID:   1,
  //          Name: "a8m",
  //          Edges: {
  //              Pets: [Pet(...), ...]
  //              ...
  //          }
  //      },
  //      ...
  //  ]
  //
  for _, u := range users {
      for _, p := range u.Edges.Pets {
          fmt.Printf("User(%v) -> Pet(%v)\n", u.ID, p.ID)
          // Output:
          // User(...) -> Pet(...)
      }
  } 
  ```

  预加载允许同时查询多个关联 (包括嵌套)，也可以对其结果进行筛选、排序或限制数量。 例如：
  ```
  admins, err := client.User.
      Query().
      Where(user.Admin(true)).
      // 填充与 `admins` 相关联的 `pets`
      WithPets().
      // 填充与 `admins` 相关联的前5个 `groups`
      WithGroups(func(q *ent.GroupQuery) {
          q.Limit(5)              // 限量5个
          q.WithUsers()           // Populate the `users` of each `groups`.
      }).
      All(ctx)
  if err != nil {
      return err
  }

  // 返回的结果类似于:
  //
  //  [
  //      User {
  //          ID:   1,
  //          Name: "admin1",
  //          Edges: {
  //              Pets:   [Pet(...), ...]
  //              Groups: [
  //                  Group {
  //                      ID:   7,
  //                      Name: "GitHub",
  //                      Edges: {
  //                          Users: [User(...), ...]
  //                          ...
  //                      }
  //                  }
  //              ]
  //          }
  //      },
  //      ...
  //  ]
  //
  for _, admin := range admins {
      for _, p := range admin.Edges.Pets {
          fmt.Printf("Admin(%v) -> Pet(%v)\n", u.ID, p.ID)
          // Output:
          // Admin(...) -> Pet(...)
      }
      for _, g := range admin.Edges.Groups {
          for _, u := range g.Edges.Users {
              fmt.Printf("Admin(%v) -> Group(%v) -> User(%v)\n", u.ID, g.ID, u.ID)
              // Output:
              // Admin(...) -> Group(...) -> User(...)
          }
      }
  } 
  ```

- API
  每个查询构造器，都会为每一条边生成形如 With<E>(...func(<N>Query)) 的方法。 <E> 是边的名称 (如, WithGroups)，<N> 是边的类型 (如, GroupQuery).

  请注意，只有 SQL后端 支持此功能。

- 实现
  由于查询构造器可以加载多个关联，无法使用一个 JOIN 操作加载它们。 因此，ent 加载关联时会进行额外的查询。 M2O/O2M 和 O2O 的边会进行一次查询， M2M 的边会进行2次查询。

  关于此，我们预计在下一个版本的 ent 中改进。

### Hooks(钩子)

 Hooks 选项允许在改变表的操作之前和之后添加自定义逻辑。

- Mutation
  
  Mutation操作是对数据库进行变化的操作。 例如，向图中添加一个新节点，删除 2 个节点之间的边或删除多个节点。
  有5种类型的变化：
    * Create - 在图中创建节点。
    * UpdateOne - 更新图中的节点。 例如，增加其字段。
    * Update - 更新图中匹配条件的多个节点。
    * DeleteOne - 从图中删除一个节点。
    * Delete - 删除与条件匹配的所有节点。
  
  每个生成的节点类型都有自己的 mutation 。 例如，所有用户构建器共享同一个生成的 UserMutation 对象。 但是，所有构建器类型都实现了通用的 ent.Mutation 接口。

- Hooks
  Hooks是获取 ent.Mutator 并返回一个 mutator 的函数。 它们充当mutator之间的中间件。 它类似于流行的 HTTP 中间件模式。
  ```
   type (
        // Mutator is the interface that wraps the Mutate method.
        Mutator interface {
            // Mutate apply the given mutation on the graph.
            Mutate(context.Context, Mutation) (Value, error)
        }

        // Hook defines the "mutation middleware". A function that gets a Mutator
        // and returns a Mutator. For example:
        //
        //  hook := func(next ent.Mutator) ent.Mutator {
        //      return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
        //          fmt.Printf("Type: %s, Operation: %s, ConcreteType: %T\n", m.Type(), m.Op(), m)
        //          return next.Mutate(ctx, m)
        //      })
        //  }
        //
        Hook func(Mutator) Mutator
    )
  ```

  有两种类型的钩子——Schema hooks和runtime hooks。 Schema hooks 主要用于在 schema 中定义自定义变更逻辑，runtime hooks 用于添加日志、度量、跟踪等内容。

- Runtime hooks
  让我们从一个简短的例子开始，它记录了所有model表的所有变化操作
  ```
   func main() {
        client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
        if err != nil {
            log.Fatalf("failed opening connection to sqlite: %v", err)
        }
        defer client.Close()
        ctx := context.Background()
        // Run the auto migration tool.
        if err := client.Schema.Create(ctx); err != nil {
            log.Fatalf("failed creating schema resources: %v", err)
        }
        // Add a global hook that runs on all types and all operations.
        client.Use(func(next ent.Mutator) ent.Mutator {
            return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
                start := time.Now()
                defer func() {
                    log.Printf("Op=%s\tType=%s\tTime=%s\tConcreteType=%T\n", m.Op(), m.Type(), time.Since(start), m)
                }()
                return next.Mutate(ctx, m)
            })
        })
        client.User.Create().SetName("a8m").SaveX(ctx)
        // Output:
        // 2020/03/21 10:59:10 Op=Create    Type=Card   Time=46.23µs    ConcreteType=*ent.UserMutation
    }
  ```
  全局钩子对于添加跟踪、指标、日志等很有用。 但有时，用户想要更多的粒度：比如只针对User实体的操作钩子
  ```
  func main() {
        // <client was defined in the previous block>

        // Add a hook only on user mutations.
        client.User.Use(func(next ent.Mutator) ent.Mutator {
            // Use the "<project>/ent/hook" to get the concrete type of the mutation.
            return hook.UserFunc(func(ctx context.Context, m *ent.UserMutation) (ent.Value, error) {
                return next.Mutate(ctx, m)
            })
        })

        // Add a hook only on update operations.
        client.Use(hook.On(Logger(), ent.OpUpdate|ent.OpUpdateOne))

        // Reject delete operations.
        client.Use(hook.Reject(ent.OpDelete|ent.OpDeleteOne))
    }
  ```
  假设你想共享一个钩子，在多个表（例如组和用户）之间改变一个字段。 有~2种方法可以做到这一点
  ``` 这2个钩子很变态，把name统一改成了Ariel Mashraki，无论你想改成什么都不行。
    // Option 1: use type assertion.
    client.Use(func(next ent.Mutator) ent.Mutator {
        type NameSetter interface {
            SetName(value string)
        }
        return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
            // A schema with a "name" field must implement the NameSetter interface. 
            if ns, ok := m.(NameSetter); ok {
                ns.SetName("Ariel Mashraki")
            }
            return next.Mutate(ctx, m)
        })
    })

    // Option 2: use the generic ent.Mutation interface.
    client.Use(func(next ent.Mutator) ent.Mutator {
        return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
            if err := m.SetField("name", "Ariel Mashraki"); err != nil {
                // An error is returned, if the field is not defined in
                // the schema, or if the type mismatch the field type.
            }
            return next.Mutate(ctx, m)
        })
    })
  ```

  
- Schema hooks  使用这个模式请一定要看Hooks Registration，否则空指针报错
  模型挂钩在类型模式中定义，仅应用于与模式类型匹配的更改。 在模式中定义钩子的动机是将所有关于节点类型的逻辑收集在一个地方，这就是schema。
  ```
    package schema

    import (
        "context"
        "fmt"

        gen "<project>/ent"
        "<project>/ent/hook"

        "entgo.io/ent"
    )

    // Card holds the schema definition for the CreditCard entity.
    type Card struct {
        ent.Schema
    }

    // Hooks of the Card.
    func (Card) Hooks() []ent.Hook {
        return []ent.Hook{
            // First hook.
            hook.On(
                func(next ent.Mutator) ent.Mutator {
                    return hook.CardFunc(func(ctx context.Context, m *gen.CardMutation) (ent.Value, error) {
                        if num, ok := m.Number(); ok && len(num) < 10 {
                            return nil, fmt.Errorf("card number is too short")
                        }
                        return next.Mutate(ctx, m)
                    })
                },
                // Limit the hook only for these operations.
                ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne,
            ),
            // Second hook.
            func(next ent.Mutator) ent.Mutator {
                return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
                    if s, ok := m.(interface{ SetName(string) }); ok {
                        s.SetName("Boring")
                    }
                    return next.Mutate(ctx, m)
                })
            },
        }
    }
  ```

- Hooks Registration(钩子注册)
  使用模式挂钩时，模式包和生成的 ent 包之间可能会出现循环导入。 为了避免这种情况，ent 生成一个 ent/runtime 包，负责在运行时注册模式挂钩。

  用户必须导入 ENT/RUNTIME 才能注册 Schema Hooks。 该包可以导入到主包中（靠近导入数据库驱动程序的位置），或者导入到创建 ENT.CLIENT 的包中。简单点说就是使用对象表的地方。这一步是必须的，否则空指针错误

  `import _ "<project>/ent/runtime"`

- Evaluation order(执行顺序)
  
  钩子按照它们注册到客户端的顺序被调用。 因此，client.Use(f, g, h) 对变化执行 f(g(h(...))) 。

  还要注意，运行时钩子在模式钩子之前被调用。 也就是说，如果在模式中定义了 g 和 h，并且 f 是使用 client.Use(...) 注册的，它们将按如下方式执行：f(g(h(...)))。

- Hook helpers(钩子帮助器)
  生成的 hooks 包提供了几个帮助器，可以帮助您控制何时执行钩子。
  ```
   package schema

    import (
        "context"
        "fmt"

        "<project>/ent/hook"

        "entgo.io/ent"
        "entgo.io/ent/schema/mixin"
    )


    type SomeMixin struct {
        mixin.Schema
    }

    func (SomeMixin) Hooks() []ent.Hook {
        return []ent.Hook{
            // 仅对 UpdateOne 和 DeleteOne 操作执行“HookA”。
            hook.On(HookA(), ent.OpUpdateOne|ent.OpDeleteOne),
            // 不要在创建操作上执行“HookB”。
            hook.Unless(HookB(), ent.OpCreate),
            // 仅当 ent.Mutation 正在改变“status”字段时才执行“HookC”，
            // 并清除“dirty”字段
            hook.If(HookC(), hook.And(hook.HasFields("status"), hook.HasClearedFields("dirty"))),
        }
    }
  ```

- Transaction Hooks(事务钩子)
  钩子也可以在活动事务上注册，并将在 Tx.Commit 或 Tx.Rollback 上执行。 [查看更多](#事务)。

- Codegen Hooks(代码生成钩子)
  entc 包提供了一个选项，可以将挂钩（中间件）列表添加到代码生成阶段。 [查看更多](#代码生成)。

### Privacy(隐私策略)
schema中的 Policy 选项允许为数据库中实体的查询和更改配置隐私策略。

隐私层的主要优点是，一次（在架构中）编写隐私策略，并且始终对其进行评估。 无论在您的代码库中的何处执行查询和更改，它都将始终通过隐私层。

我们将首先介绍我们在框架中使用的基本术语，然后继续介绍为您的项目配置策略功能的部分，并以几个示例结束。

#### Basic Terms(基本规则)
- Policy(策略)
  ent.Policy 接口包含两个方法：EvalQuery 和 EvalMutation。 第一个定义读取策略，第二个定义写入策略。 一项政策包含零个或多个隐私规则（见下文）。 这些规则按照它们在schema中声明的相同顺序进行评估。
  
  如果所有规则都被评估而没有返回错误，则评估成功完成，并且执行的操作可以访问目标节点。

  但是，如果评估的规则之一返回错误或privacy.Deny 决定（见下文），则执行的操作将返回错误，并被取消。

- Privacy Rules(规则)
  每个策略（更改或查询）包括一个或多个隐私规则。 这些规则的函数签名如下：
  ```
   // EvalQuery defines the a read-policy rule.
   func(Policy) EvalQuery(context.Context, Query) error
   
   // EvalMutation defines the a write-policy rule.
   func(Policy) EvalMutation(context.Context, Mutation) error
  ```

- Privacy Decisions(决策)
  有三种类型的决策可以帮助您控制隐私规则评估。
  * privacy.Allow - 如果从隐私规则返回，则评估停止（将跳过下一个规则），并且执行的操作（查询或更改）可以访问目标节点。
  * privacy.Deny - 如果从隐私规则返回，则评估停止（将跳过下一个规则），并取消执行的操作。 这相当于返回任何错误。
  * privacy.Skip - 跳过当前规则，跳转到下一个隐私规则。 这相当于返回一个 nil 错误。

#### Configuration(配置)
为了在您的代码生成中启用隐私选项，请使用以下两个选项之一启用隐私功能：

1. 如果您使用默认的 go generate 配置，请在 ent/generate.go 文件中添加 --feature 隐私选项，如下所示：

   ```
    package ent

    //go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature privacy ./schema
   ```
   建议添加 schema/snapshot 特性标志以及隐私以增强开发体验（例如--featureprivacy,schema/snapshot）

2. 如果您使用的是 GraphQL 文档中的配置，请按如下方式添加功能标志
   ```
   // Copyright 2019-present Facebook Inc. All rights reserved.
   // This source code is licensed under the Apache 2.0 license found
   // in the LICENSE file in the root directory of this source tree.

   // +build ignore

   package main


   import (
       "log"

       "entgo.io/ent/entc"
       "entgo.io/ent/entc/gen"
       "entgo.io/contrib/entgql"
   )

   func main() {
       opts := []entc.Option{
           entc.FeatureNames("privacy"),
       }
       err := entc.Generate("./schema", &gen.Config{
           Templates: entgql.AllTemplates,
       }, opts...)
       if err != nil {
           log.Fatalf("running ent codegen: %v", err)
       }
   }
   ```

   您应该注意，类似于 SCHEMA Hooks，如果您在您的 SCHEMA 中使用策略选项，您必须在主包中添加以下导入，因为在 SCHEMA 包和生成的包之间可能会出现循环导入：
   ```import _ "<project>/ent/runtime"```

#### Examples(例子)
- Admin Only
  我们从一个简单的应用程序示例开始，该应用程序允许任何用户读取任何数据，并且仅接受来自具有管理员角色的用户的更改。 我们将为示例创建 2 个额外的包：
  
  rule - 用于在我们的架构中保存不同的隐私规则。

  viewer - 用于获取和设置执行操作的用户/查看器。 在这个简单的例子中，它可以是普通用户或管理员。

  运行代码生成（带有隐私功能标志）后，我们添加带有 2 个生成策略规则的 Policy 方法。

  ``` examples/privacyadmin/ent/schema/user.go
  package schema

  import (
      "entgo.io/ent"
      "entgo.io/ent/examples/privacyadmin/ent/privacy"
  )

  // User holds the schema definition for the User entity.
  type User struct {
      ent.Schema
  }

  // Policy defines the privacy policy of the User.
  func (User) Policy() ent.Policy {
      return privacy.Policy{
          Mutation: privacy.MutationPolicy{
              // Deny if not set otherwise. 
              privacy.AlwaysDenyRule(),
          },
          Query: privacy.QueryPolicy{
              // Allow any viewer to read anything.
              privacy.AlwaysAllowRule(),
          },
      }
  }
  ```

  我们定义了一个拒绝任何突变并接受任何查询的策略。 但是，如上所述，在此示例中，我们仅接受来自具有管理员角色的查看者的更改。 让我们创建 2 条隐私规则来强制执行此操作：
  ``` examples/privacyadmin/rule/rule.go
  package rule

  import (
      "context"

      "entgo.io/ent/examples/privacyadmin/ent/privacy"
      "entgo.io/ent/examples/privacyadmin/viewer"
  )

  // DenyIfNoViewer is a rule that returns Deny decision if the viewer is
  // missing in the context.
  func DenyIfNoViewer() privacy.QueryMutationRule {
      return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
          view := viewer.FromContext(ctx)
          if view == nil {
              return privacy.Denyf("viewer-context is missing")
          }
          // Skip to the next privacy rule (equivalent to returning nil).
          return privacy.Skip
      })
  }

  // AllowIfAdmin is a rule that returns Allow decision if the viewer is admin.
  func AllowIfAdmin() privacy.QueryMutationRule {
      return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
          view := viewer.FromContext(ctx)
          if view.Admin() {
              return privacy.Allow
          }
          // Skip to the next privacy rule (equivalent to returning nil).
          return privacy.Skip
      })
  }
  ```
  如您所见，第一条规则 DenyIfNoViewer 确保每个操作在其上下文中都有一个查看器，否则操作将被拒绝。 第二个规则 AllowIfAdmin，接受来自具有管理员角色的查看者的任何操作。 让我们将它们添加到架构中，并运行代码生成

  ``` examples/privacyadmin/ent/schema/user.go
  // Policy defines the privacy policy of the User.
  func (User) Policy() ent.Policy {
      return privacy.Policy{
          Mutation: privacy.MutationPolicy{
              rule.DenyIfNoViewer(),
              rule.AllowIfAdmin(),
              privacy.AlwaysDenyRule(),
          },
          Query: privacy.QueryPolicy{
              privacy.AlwaysAllowRule(),
          },
      }
  }
  ```
  由于我们首先定义了 DenyIfNoViewer，它将在所有其他规则之前执行，并且在 AllowIfAdmin 规则中访问 viewer.Viewer 对象是安全的。

  添加上述规则并运行代码生成后，我们希望将隐私层逻辑应用于 ent.Client 操作。
  ``` examples/privacyadmin/example_test.go
  func Do(ctx context.Context, client *ent.Client) error {
      // Expect operation to fail, because viewer-context
      // is missing (first mutation rule check).
      if _, err := client.User.Create().Save(ctx); !errors.Is(err, privacy.Deny) {
          return fmt.Errorf("expect operation to fail, but got %w", err)
      }
      // Apply the same operation with "Admin" role.
      admin := viewer.NewContext(ctx, viewer.UserViewer{Role: viewer.Admin})
      if _, err := client.User.Create().Save(admin); err != nil {
          return fmt.Errorf("expect operation to pass, but got %w", err)
      }
      // Apply the same operation with "ViewOnly" role.
      viewOnly := viewer.NewContext(ctx, viewer.UserViewer{Role: viewer.View})
      if _, err := client.User.Create().Save(viewOnly); !errors.Is(err, privacy.Deny) {
          return fmt.Errorf("expect operation to fail, but got %w", err)
      }
      // Allow all viewers to query users.
      for _, ctx := range []context.Context{ctx, viewOnly, admin} {
          // Operation should pass for all viewers.
          count := client.User.Query().CountX(ctx)
          fmt.Println(count)
      }
      return nil
  }
  ```

- Decision Context
  有时，我们希望将特定的隐私决策绑定到 context.Context。 在这种情况下，我们可以使用 privacy.DecisionContext 函数创建一个附加隐私决策的新上下文。
  ``` examples/privacyadmin/example_test.go
  func Do(ctx context.Context, client *ent.Client) error {
      // Bind a privacy decision to the context (bypass all other rules).
      allow := privacy.DecisionContext(ctx, privacy.Allow)
      if _, err := client.User.Create().Save(allow); err != nil {
          return fmt.Errorf("expect operation to pass, but got %w", err)
      }
      return nil
  }
  ```

- Multi Tenancy
  在此示例中，我们将创建一个具有 3 种实体类型的架构 - 租户、用户和组。 本示例中还存在帮助程序包查看器和规则（如上所述）以帮助我们构建应用程序。
  ![关系图例](https://entgo.io/images/assets/tenant_medium.png "example")

  让我们一点一点地开始构建这个应用程序。 我们首先创建 3 个不同的模型(请参阅此处的[完整代码](https://github.com/ent/ent/tree/master/examples/privacytenant/ent/schema))，并且由于我们希望在它们之间共享一些逻辑，因此我们创建了另一个混合模型并将其添加到所有其他模式中，如下所示：

  ``` examples/privacytenant/ent/schema/mixin.go
  // BaseMixin for all schemas in the graph.
  type BaseMixin struct {
      mixin.Schema
  }

  // Policy defines the privacy policy of the BaseMixin.
  func (BaseMixin) Policy() ent.Policy {
      return privacy.Policy{
          Mutation: privacy.MutationPolicy{
              rule.DenyIfNoViewer(),
          },
          Query: privacy.QueryPolicy{
              rule.DenyIfNoViewer(),
          },
      }
  }
  ```
  ``` examples/privacytenant/ent/schema/tenant.go
  // Mixin of the Tenant schema.
  func (Tenant) Mixin() []ent.Mixin {
      return []ent.Mixin{
          BaseMixin{},
      }
  }
  ```
  如第一个示例中所述，DenyIfNoViewer 隐私规则在 context.Context 不包含 viewer.Viewer 信息时拒绝该操作。

  与前面的示例类似，我们希望添加一个约束，只有管理员用户才能创建租户（否则拒绝）。 我们通过从上面复制 AllowIfAdmin 规则，并将其添加到租户模式的策略来实现：

  ``` examples/privacytenant/ent/schema/tenant.go
  // Policy defines the privacy policy of the User.
  func (Tenant) Policy() ent.Policy {
      return privacy.Policy{
          Mutation: privacy.MutationPolicy{
              // For Tenant type, we only allow admin users to mutate
              // the tenant information and deny otherwise.
              rule.AllowIfAdmin(),
              privacy.AlwaysDenyRule(),
          },
      }
  }
  ```
  然后，我们期望以下代码能够成功运行：
  ```
  func Do(ctx context.Context, client *ent.Client) error {
      // Expect operation to fail, because viewer-context
      // is missing (first mutation rule check).
      if _, err := client.Tenant.Create().Save(ctx); !errors.Is(err, privacy.Deny) {
          return fmt.Errorf("expect operation to fail, but got %w", err)
      }
      // Deny tenant creation if the viewer is not admin.
      viewCtx := viewer.NewContext(ctx, viewer.UserViewer{Role: viewer.View})
      if _, err := client.Tenant.Create().Save(viewCtx); !errors.Is(err, privacy.Deny) {
          return fmt.Errorf("expect operation to fail, but got %w", err)
      }
      // Apply the same operation with "Admin" role, expect it to pass.
      adminCtx := viewer.NewContext(ctx, viewer.UserViewer{Role: viewer.Admin})
      hub, err := client.Tenant.Create().SetName("GitHub").Save(adminCtx)
      if err != nil {
          return fmt.Errorf("expect operation to pass, but got %w", err)
      }
      fmt.Println(hub)
      lab, err := client.Tenant.Create().SetName("GitLab").Save(adminCtx)
      if err != nil {
          return fmt.Errorf("expect operation to pass, but got %w", err)
      }
      fmt.Println(lab)
      return nil
  }
  ```

  我们继续在我们的数据模型中添加其余的边（见上图），由于 User 和 Group 都有一个 Tenant 模式的边，我们为此创建了一个名为 TenantMixin 的共享混合模型：
  ```examples/privacytenant/ent/schema/mixin.go
  // TenantMixin for embedding the tenant info in different schemas.
  type TenantMixin struct {
      mixin.Schema
  }

  // Edges for all schemas that embed TenantMixin.
  func (TenantMixin) Edges() []ent.Edge {
      return []ent.Edge{
          edge.To("tenant", Tenant.Type).
              Unique().
              Required(),
      }
  }
  ```

  接下来，我们可能想要强制执行一个规则，将查看者限制为仅查询连接到他们所属租户的组和用户。 对于这样的用例，Ent 有一种名为 Filter 的附加隐私规则类型。 我们可以使用过滤规则根据查看者的身份过滤掉实体。 与我们之前讨论的规则不同，过滤规则除了返回隐私决策之外，还可以限制查看者可以进行的查询范围。

  _请注意，需要使用 entql 功能标志启用隐私过滤选项（请参阅[上面](#configuration配置)的说明）。_

   ``` examples/privacytenant/rule/rule.go
    // FilterTenantRule is a query/mutation rule that filters out entities that are not in the tenant.
    func FilterTenantRule() privacy.QueryMutationRule {
        // TenantsFilter is an interface to wrap WhereHasTenantWith()
        // predicate that is used by both `Group` and `User` schemas.
        type TenantsFilter interface {
            WhereHasTenantWith(...predicate.Tenant)
        }
        return privacy.FilterFunc(func(ctx context.Context, f privacy.Filter) error {
            view := viewer.FromContext(ctx)
            if view.Tenant() == "" {
                return privacy.Denyf("missing tenant information in viewer")
            }
            tf, ok := f.(TenantsFilter)
            if !ok {
                return privacy.Denyf("unexpected filter type %T", f)
            }
            // Make sure that a tenant reads only entities that has an edge to it.
            tf.WhereHasTenantWith(tenant.Name(view.Tenant()))
            // Skip to the next privacy rule (equivalent to return nil).
            return privacy.Skip
        })
    }
   ```
   创建 FilterTenantRule 隐私规则后，我们将其添加到 TenantMixin 以确保使用此 mixin 的所有模式也将具有此隐私规则。
   ```examples/privacytenant/ent/schema/mixin.go
   // Policy for all schemas that embed TenantMixin.
    func (TenantMixin) Policy() ent.Policy {
        return privacy.Policy{
            Query: privacy.QueryPolicy{
                rule.AllowIfAdmin(),
                // Filter out entities that are not connected to the tenant.
                // If the viewer is admin, this policy rule is skipped above.
                rule.FilterTenantRule(),
            },
        }
    }
   ```
   然后，在运行代码生成之后，我们期望隐私规则对客户端操作生效。
   ```examples/privacytenant/example_test.go
   func Do(ctx context.Context, client *ent.Client) error {
        // A continuation of the code-block above.

        // Create 2 users connected to the 2 tenants we created above
        hubUser := client.User.Create().SetName("a8m").SetTenant(hub).SaveX(adminCtx)
        labUser := client.User.Create().SetName("nati").SetTenant(lab).SaveX(adminCtx)

        hubView := viewer.NewContext(ctx, viewer.UserViewer{T: hub})
        out := client.User.Query().OnlyX(hubView)
        // Expect that "GitHub" tenant to read only its users (i.e. a8m).
        if out.ID != hubUser.ID {
            return fmt.Errorf("expect result for user query, got %v", out)
        }
        fmt.Println(out)

        labView := viewer.NewContext(ctx, viewer.UserViewer{T: lab})
        out = client.User.Query().OnlyX(labView)
        // Expect that "GitLab" tenant to read only its users (i.e. nati).
        if out.ID != labUser.ID {
            return fmt.Errorf("expect result for user query, got %v", out)
        }
        fmt.Println(out)
        return nil
    }
   ```

### 事务
- 启动一个事务
  ```
  // GenTx 在一次事务中生成一系列实体。
  func GenTx(ctx context.Context, client *ent.Client) error {
      tx, err := client.Tx(ctx)
      if err != nil {
          return fmt.Errorf("starting a transaction: %w", err)
      }
      hub, err := tx.Group.
          Create().
          SetName("Github").
          Save(ctx)
      if err != nil {
          return rollback(tx, fmt.Errorf("failed creating the group: %w", err))
      }
      // Create the admin of the group.
      dan, err := tx.User.
          Create().
          SetAge(29).
          SetName("Dan").
          AddManage(hub).
          Save(ctx)
      if err != nil {
          return rollback(tx, err)
      }
      // Create user "Ariel".
      a8m, err := tx.User.
          Create().
          SetAge(30).
          SetName("Ariel").
          AddGroups(hub).
          AddFriends(dan).
          Save(ctx)
      if err != nil {
          return rollback(tx, err)
      }
      fmt.Println(a8m)
      // Output:
      // User(id=2, age=30, name=Ariel)

      // Commit the transaction.
      return tx.Commit()
  }

  // rollback calls to tx.Rollback and wraps the given error
  // with the rollback error if occurred.
  func rollback(tx *ent.Tx, err error) error {
      if rerr := tx.Rollback(); rerr != nil {
          err = fmt.Errorf("%w: %v", err, rerr)
      }
      return err
  }
  ```
  完整示例可参阅 [GitHub](https://github.com/ent/ent/tree/master/examples/traversal).

- 事务化客户端
  你可能已经有代码是用了 *ent.Client 的，并且你想将它修改成（或封装成）在事务中执行的。 对于这种情况，你可以编写一个事务化客户端。 你可以从已有事务中获取一个 *ent.Client
  ```
   // WrapGen wraps the existing "Gen" function in a transaction.
  func WrapGen(ctx context.Context, client *ent.Client) error {
      tx, err := client.Tx(ctx)
      if err != nil {
          return err
      }
      txClient := tx.Client()
      // Use the "Gen" below, but give it the transactional client; no code changes to "Gen".
      if err := Gen(ctx, txClient); err != nil {
          return rollback(tx, err)
      }
      return tx.Commit()
  }

  // Gen generates a group of entities.
  func Gen(ctx context.Context, client *ent.Client) error {
      // ...
      return nil
  }
  ```

- 最佳实践
  在事务中使用回调函数来实现代码复用：
  ```
   func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
      tx, err := client.Tx(ctx)
      if err != nil {
          return err
      }
      defer func() {
          if v := recover(); v != nil {
              tx.Rollback()
              panic(v)
          }
      }()
      if err := fn(tx); err != nil {
          if rerr := tx.Rollback(); rerr != nil {
              err = errors.Wrapf(err, "rolling back transaction: %v", rerr)
          }
          return err
      }
      if err := tx.Commit(); err != nil {
          return errors.Wrapf(err, "committing transaction: %v", err)
      }
      return nil
  }
  ```
  用法：
  ```
   func Do(ctx context.Context, client *ent.Client) {
      // WithTx helper.
      if err := WithTx(ctx, client, func(tx *ent.Tx) error {
          return Gen(ctx, tx.Client())
      }); err != nil {
          log.Fatal(err)
      }
  }
  ```

- 钩子
  与结构钩子和运行时钩子一样，钩子也可以注册在活跃的事务中，将会在Tx.Commit或者是Tx.Rollback时执行：
  ```
   func Do(ctx context.Context, client *ent.Client) error {
      tx, err := client.Tx(ctx)
      if err != nil {
          return err
      }
      // Add a hook on Tx.Commit.
      tx.OnCommit(func(next ent.Committer) ent.Committer {
          return ent.CommitFunc(func(ctx context.Context, tx *ent.Tx) error {
              // Code before the actual commit.
              err := next.Commit(ctx, tx)
              // Code after the transaction was committed.
              return err
          })
      })
      // Add a hook on Tx.Rollback.
      tx.OnRollback(func(next ent.Rollbacker) ent.Rollbacker {
          return ent.RollbackFunc(func(ctx context.Context, tx *ent.Tx) error {
              // Code before the actual rollback.
              err := next.Rollback(ctx, tx)
              // Code after the transaction was rolled back.
              return err
          })
      })
      //
      // <Code goes here>
      //
      return err
  }
  ```

### 断言
- 字段断言
  
  - 布尔类型：=, !=
  - 数字类型：=, !=, >, <, >=, <=,IN, NOT IN
  - 时间类型：=, !=, >, <, >=, <=,IN, NOT IN
  - 字符类型：
    * =, !=, >, <, >=, <=
    * IN, NOT IN,Contains, HasPrefix, HasSuffix,
    * ContainsFold, EqualFold （只能用于 SQL 语句）
  - JSON类型： 
    * =, !=
    * =, !=, >, <, >=, <= on nested values (JSON path).
    * Contains （只能用于包含嵌套的 JSON path 表达式）
    * HasKey, Len<P>
  - 可选字段:IsNil, NotNil

- Edge 断言
  * HasEdge. 例如，对于一个 Pet 类型的 edge owner 来说 ，可以使用:
   ```
     client.Pet.
      Query().
      Where(pet.HasOwner()).
      All(ctx)
   ```
  * HasEdgeWith. 也可以将断言表示为集合的形式
   ```
     client.Pet.
      Query().
      Where(pet.HasOwnerWith(user.Name("a8m"))).
      All(ctx)
   ```
- 否定 (NOT)
   ```
    client.Pet.
    Query().
    Where(pet.Not(pet.NameHasPrefix("Ari"))).
    All(ctx)
   ```
- 析取 (OR)
   ```
    client.Pet.
    Query().
    Where(
        pet.Or(
            pet.HasOwner(),
            pet.Not(pet.HasFriends()),
        )
    ).
    All(ctx)
   ```
- 合取 (AND)
   ```
    client.Pet.
    Query().
    Where(
        pet.And(
            pet.HasOwner(),
            pet.Not(pet.HasFriends()),
        )
    ).
    All(ctx)
   ```
- 自定义断言
  如果您想编写自己的特定于方言的逻辑或控制执行的查询，自定义谓词会很有用。
  获得用户 1、2、3 的所有宠物
  ```
   pets := client.Pet.
    Query().
    Where(func(s *sql.Selector) {
        s.Where(sql.InInts(pet.FieldOwnerID, 1, 2, 3))
    }).
    AllX(ctx)
  ```
  上面的代码将产生以下 SQL 查询：
  ```
   SELECT DISTINCT `pets`.`id`, `pets`.`owner_id` FROM `pets` WHERE `owner_id` IN (1, 2, 3)
  ```

  统计名为 URL 的 JSON 字段包含 Scheme 键的用户数
  ```
   count := client.User.
    Query().
    Where(func(s *sql.Selector) {
        s.Where(sqljson.HasKey(user.FieldURL, sqljson.Path("Scheme")))
    }).
    CountX(ctx)
  ```
  上面的代码将产生以下 SQL 查询：
  ```
   -- PostgreSQL
    SELECT COUNT(DISTINCT "users"."id") FROM "users" WHERE "url"->'Scheme' IS NOT NULL

    -- SQLite and MySQL
    SELECT COUNT(DISTINCT `users`.`id`) FROM `users` WHERE JSON_EXTRACT(`url`, "$.Scheme") IS NOT NULL
  ```
  获取所有拥有"Tesla"汽车的用户
  ```
   users := client.User.Query().
    Where(user.HasCarWith(car.Model("Tesla"))).
    AllX(ctx)
  ```
  此查询可以改写为 3 种不同的形式：IN、EXISTS 和 JOIN。
  ```
   // `IN` version.
    users := client.User.Query().
        Where(func(s *sql.Selector) {
            t := sql.Table(car.Table)
            s.Where(
                sql.In(
                    s.C(user.FieldID),
                    sql.Select(t.C(user.FieldID)).From(t).Where(sql.EQ(t.C(car.FieldModel), "Tesla")),
                ),
            )
        }).
        AllX(ctx)

    // `JOIN` version.
    users := client.User.Query().
        Where(func(s *sql.Selector) {
            t := sql.Table(car.Table)
            s.Join(t).On(s.C(user.FieldID), t.C(car.FieldOwnerID))
            s.Where(sql.EQ(t.C(car.FieldModel), "Tesla"))
        }).
        AllX(ctx)

    // `EXISTS` version.
    users := client.User.Query().
        Where(func(s *sql.Selector) {
            t := sql.Table(car.Table)
            p := sql.And(
                sql.EQ(t.C(car.FieldModel), "Tesla"),
                sql.ColumnsEQ(s.C(user.FieldID), t.C(car.FieldOwnerID)),
            )
            s.Where(sql.Exists(sql.Select().From(t).Where(p)))
        }).
        AllX(ctx)
  ```
  上面的代码将产生以下 SQL 查询：
  ```
   -- `IN` version.
  SELECT DISTINCT `users`.`id`, `users`.`age`, `users`.`name` FROM `users` WHERE `users`.`id` IN (SELECT `cars`.`id` FROM `cars` WHERE `cars`.`model` = 'Tesla')

  -- `JOIN` version.
  SELECT DISTINCT `users`.`id`, `users`.`age`, `users`.`name` FROM `users` JOIN `cars` ON `users`.`id` = `cars`.`owner_id` WHERE `cars`.`model` = 'Tesla'

  -- `EXISTS` version.
  SELECT DISTINCT `users`.`id`, `users`.`age`, `users`.`name` FROM `users` WHERE EXISTS (SELECT * FROM `cars` WHERE `cars`.`model` = 'Tesla' AND `users`.`id` = `cars`.`owner_id`)
  ```

### 聚合
- 分组
  对 users 按 name 和 age 字段分组，并计算 age 的总和。
  ```
   package main

  import (
      "context"

      "<project>/ent"
      "<project>/ent/user"
  )

  func Do(ctx context.Context, client *ent.Client) {
      var v []struct {
          Name  string `json:"name"`
          Age   int    `json:"age"`
          Sum   int    `json:"sum"`
          Count int    `json:"count"`
      }
      err := client.User.Query().
          GroupBy(user.FieldName, user.FieldAge).
          Aggregate(ent.Count(), ent.Sum(user.FieldAge)).
          Scan(ctx, &v)
  }
  ```
  按单个字段分组.
  ```
   package main

  import (
      "context"

      "<project>/ent"
      "<project>/ent/user"
  )

  func Do(ctx context.Context, client *ent.Client) {
      names, err := client.User.
          Query().
          GroupBy(user.FieldName).
          Strings(ctx)
  }
  ```
- 根据边进行分组
  如果您想按照自己的逻辑进行聚合，可以使用自定义聚合函数。

  下面展示了：如何根据用户的 id 和 name 进行分组，并计算其宠物的平均 age。
  ```
   package main

  import (
      "context"
      "log"

      "<project>/ent"
      "<project>/ent/pet"
      "<project>/ent/user"
  )

  func Do(ctx context.Context, client *ent.Client) {
      var users []struct {
          ID      int
          Name    string
          Average float64
      }
      err := client.User.Query().
          GroupBy(user.FieldID, user.FieldName).
          Aggregate(func(s *sql.Selector) string {
              t := sql.Table(pet.Table)
              s.Join(t).On(s.C(user.FieldID), t.C(pet.OwnerColumn))
              return sql.As(sql.Avg(t.C(pet.FieldAge)), "average")
          }).
          Scan(ctx, &users)
  }
  ```

### 分页和排序

- 限量
  Limit 限制查询只返回 n 个结果.
  ```
   users, err := client.User.
    Query().
    Limit(n).
    All(ctx)
  ```
  
- 偏移 
  Offset 设置第一个返回结果在所有结果中的位置.就是跳过多少条记录
  ```
   users, err := client.User.
    Query().
    Offset(10).
    All(ctx)
  ```
  
- 排序
   Order 设置按照一个或多个字段值来对返回结果排序。 注意，如果给定的字段不是有效的列或外键，将返回错误。
   ```
    users, err := client.User.Query().
    Order(ent.Asc(user.FieldName)).
    All(ctx)
   ```

- 按边排序 (测试无效,必须对user在做一个order)
  为了按边 (关系) 的字段排序，从此边开始遍历 (你想要进行排序的边) 并应用排序，然后转到目标类型。

  下面展示了 如何根据用户 "pets" 的 "name" 来对用户进行排序。
  ```
   users, err := client.Pet.Query().
      Order(ent.Asc(pet.FieldName)).
      QueryOwner().
      All(ctx)
  ```

- 自定义排序
  如果您想写自己的特定逻辑，可以使用自定义排序函数。

  下面展示了 如何根据 宠物的名字 和 主人的名字 来对宠物进行升序排序。
  ```
   names, err := client.Pet.Query().
    Order(func(s *sql.Selector) {
        // 连接用户表，以通过 owner-name 和 pet-name 进行排序.
        t := sql.Table(user.Table)
        s.Join(t).On(s.C(pet.OwnerColumn), t.C(user.FieldID))
        s.OrderBy(t.C(user.FieldName), s.C(pet.FieldName))
    }).
    Select(pet.FieldName).
    Strings(ctx)
  ```

## 迁移
### 数据库迁移
ent的迁移支持功能，可使数据库 schema 与你项目根目录下的 ent/migrate/schema.go 中定义的 schema 对象保持一致。

- 自动迁移
  在应用程序初始化过程中运行自动迁移逻辑：
  ```
   if err := client.Schema.Create(ctx); err != nil {
      log.Fatalf("failed creating schema resources: %v", err)
  }
  ```
  Create 创建你项目 ent 部分所需的的数据库资源 。 默认情况下，Create以"append-only"模式工作；这意味着，它只创建新表和索引，将列追加到表或扩展列类型。 例如，将int改为bigint。

  想要删除列或索引怎么办？

- 删除资源
  WithDropIndex 和 WithDropColumn 是用于删除表列和索引的两个选项。
  ```
   package main

  import (
      "context"
      "log"

      "<project>/ent"
      "<project>/ent/migrate"
  )

  func main() {
      client, err := ent.Open("mysql", "root:pass@tcp(localhost:3306)/test")
      if err != nil {
          log.Fatalf("failed connecting to mysql: %v", err)
      }
      defer client.Close()
      ctx := context.Background()
      // Run migration.
      err = client.Schema.Create(
          ctx, 
          migrate.WithDropIndex(true),
          migrate.WithDropColumn(true), 
      )
      if err != nil {
          log.Fatalf("failed creating schema resources: %v", err)
      }
  }
  ```
  为了在调试模式下运行迁移 (打印所有SQL查询)，请运行：
  ```
   err := client.Debug().Schema.Create(
      ctx, 
      migrate.WithDropIndex(true),
      migrate.WithDropColumn(true),
  )
  if err != nil {
      log.Fatalf("failed creating schema resources: %v", err)
  }
  ```

- 全局唯一ID
  默认情况下，每个表的SQL主键从1开始；这意味着不同类型的多个实体可以有相同的ID。 不像AWS Neptune，节点ID是UUID。

  如果您使用 GraphQL，这不会很好地工作，它要求对象 ID 是唯一的。

  要为您的项目启用 Universal-ID 支持，请将 WithGlobalUniqueID 选项传递给迁移。

  ```
   package main

    import (
        "context"
        "log"

        "<project>/ent"
        "<project>/ent/migrate"
    )

    func main() {
        client, err := ent.Open("mysql", "root:pass@tcp(localhost:3306)/test")
        if err != nil {
            log.Fatalf("failed connecting to mysql: %v", err)
        }
        defer client.Close()
        ctx := context.Background()
        // Run migration.
        if err := client.Schema.Create(ctx, migrate.WithGlobalUniqueID(true)); err != nil {
            log.Fatalf("failed creating schema resources: %v", err)
        }
    }
  ```

  它是如何工作的？ ent 迁移为每个实体（表）的 ID 分配一个 1<<32 的范围，并将此信息存储在名为 ent_types 的表中。 例如，类型 A 的 ID 范围为 [1,4294967296)，类型 B 的范围为 [4294967296,8589934592)，等等。请注意，如果启用此选项，则可能的最大表数为 65535。

- Offline Mode
  离线模式允许您在数据库上执行之前将架构更改写入 io.Writer。 这对于在数据库上执行 SQL 命令之前验证它们很有用，或者让 SQL 脚本手动运行。
  ```
   package main

  import (
      "context"
      "log"
      "os"

      "<project>/ent"
      "<project>/ent/migrate"
  )

  func main() {
      client, err := ent.Open("mysql", "root:pass@tcp(localhost:3306)/test")
      if err != nil {
          log.Fatalf("failed connecting to mysql: %v", err)
      }
      defer client.Close()
      ctx := context.Background()
      // Dump migration changes to stdout.
      if err := client.Schema.WriteTo(ctx, os.Stdout); err != nil {
          log.Fatalf("failed printing schema changes: %v", err)
      }
  }
  ```
  将更改写入文件
  ```
   package main

  import (
      "context"
      "log"
      "os"

      "<project>/ent"
      "<project>/ent/migrate"
  )

  func main() {
      client, err := ent.Open("mysql", "root:pass@tcp(localhost:3306)/test")
      if err != nil {
          log.Fatalf("failed connecting to mysql: %v", err)
      }
      defer client.Close()
      ctx := context.Background()
      // Dump migration changes to an SQL script.
      f, err := os.Create("migrate.sql")
      if err != nil {
          log.Fatalf("create migrate file: %v", err)
      }
      defer f.Close()
      if err := client.Schema.WriteTo(ctx, f); err != nil {
          log.Fatalf("failed printing schema changes: %v", err)
      }
  }
  ```

- Foreign Keys(外键)
  默认情况下，ent 在定义关系（边）时使用外键来强制执行数据库端的正确性和一致性。

  但是，ent 还提供了使用 WithForeignKeys 选项禁用此功能的选项。 注意，将此选项设置为 false，将告诉迁移不在模式 DDL 中创建外键，并且边缘验证和清除必须由开发人员手动处理。

  我们期望在不久的将来提供一套钩子来实现应用层的外键约束。
  ```
   package main

  import (
      "context"
      "log"

      "<project>/ent"
      "<project>/ent/migrate"
  )

  func main() {
      client, err := ent.Open("mysql", "root:pass@tcp(localhost:3306)/test")
      if err != nil {
          log.Fatalf("failed connecting to mysql: %v", err)
      }
      defer client.Close()
      ctx := context.Background()
      // Run migration.
      err = client.Schema.Create(
          ctx,
          migrate.WithForeignKeys(false), // Disable foreign keys.
      )
      if err != nil {
          log.Fatalf("failed creating schema resources: %v", err)
      }
  }
  ```

- Migration Hooks(迁移钩子)
  该框架提供了向迁移阶段添加钩子（中间件）的选项。 此选项非常适合修改或过滤迁移正在处理的表，或用于在数据库中创建自定义资源。
  ```
   package main

  import (
      "context"
      "log"

      "<project>/ent"
      "<project>/ent/migrate"

      "entgo.io/ent/dialect/sql/schema"
  )

  func main() {
      client, err := ent.Open("mysql", "root:pass@tcp(localhost:3306)/test")
      if err != nil {
          log.Fatalf("failed connecting to mysql: %v", err)
      }
      defer client.Close()
      ctx := context.Background()
      // Run migration.
      err = client.Schema.Create(
          ctx,
          schema.WithHooks(func(next schema.Creator) schema.Creator {
              return schema.CreateFunc(func(ctx context.Context, tables ...*schema.Table) error {
                  // Run custom code here.
                  return next.Create(ctx, tables...)
              })
          }),
      )
      if err != nil {
          log.Fatalf("failed creating schema resources: %v", err)
      }
  }
  ```
  

### 支持的数据库后端
- MySQL
MySQL 支持 迁移 部分中提到的所有功能。 而且以下3个版本会持续地进行测试。5.6.35, 5.7.26 和 8。

- MariaDB
MariaDB 支持迁移部分中提到的所有功能，并且正在以下 3 个版本上不断测试：10.2、10.3 和最新版本。

- PostgreSQL
PostgreSQL 支持迁移部分中提到的所有功能，并且正在以下 4 个版本上不断测试：10、11、12 和 13。

- SQLite
SQLite支持 迁移 部分中提到的所有 "append-only" 功能。 然而，SQLite 默认不支持删除或修改资源，例如 drop-index，将来会使用 临时表 来添加。

- Gremlin
Gremlin 不支持迁移和索引，而且 对它的支持是实验性的。

## 额外的补充

### 外部模版
使用外部模本生成的命令：
`go run entgo.io/ent/cmd/ent generate --template glob="./ent/template/*.tmpl" ./ent/schema`

使用此命令要在ent文件夹下，创建template文件夹，并创建additional.tmpl模版文件
``` additional.tmpl
  {{/* 此模板用于给实体添加一些字段。 */}}
  {{ define "model/fields/additional" }}
      {{- /* 给 "Card" 实体添加静态字段。 */}}
      {{- if eq $.Name "Card" }}
          // 通过模板定义静态字段。
          StaticField string `json:"static_field,omitempty"`
      {{- end }}
  {{ end }}
```

运行结果在user.go生成的结构模板中增加字段

// StaticField defined by template.
StaticField string `json:"static,omitempty"`

#### Helper Templates
如上所述， ent 将每个模板的执行输出写入与模板同名的文件中。 例如，定义为 {{define "stringer" }} 的模板的输出将写入名为 ent/stringer.go 的文件。

默认情况下，ent 将每个用 {{define "<name>" }} 声明的模板写入文件。 然而，有时需要定义帮助器模板 - 不会直接调用而是由其他模板执行的模板。 为了方便这个用例，ent 支持两种命名格式，将模板指定为帮助程序。 格式是：
1. {{define "helper/.+" }} 用于全局助手模板。 例如：
   ```
    {{ define "helper/foo" }}
    {{/* Logic goes here. */}}
   {{ end }}

   {{ define "helper/bar/baz" }}
       {{/* Logic goes here. */}}
   {{ end }}
   ```
2. {{ 定义 "<root-template>/helper/.+" }} 用于本地助手模板。 如果模板的执行输出写入文件，则该模板被视为“根”。 例如：
  ```
   {{/* 在 `gen.Graph` 上执行的根模板将被写入名为：`ent/http.go` 的文件中。*/}}
   {{ define "http" }}
       {{ range $n := $.Nodes }}
           {{ template "http/helper/get" $n }}
           {{ template "http/helper/post" $n }}
       {{ end }}
   {{ end }}

   {{/* 在 `gen.Type` 上执行的辅助模板 */}}
   {{ define "http/helper/get" }}
       {{/* Logic goes here. */}}
   {{ end }}

   {{/* 在 `gen.Type` 上执行的辅助模板 */}}
   {{ define "http/helper/post" }}
       {{/* Logic goes here. */}}
   {{ end }}
  ```

#### Annotations
模型注释允许将元数据附加到字段和边缘并将它们注入外部模板。

注释必须是可序列化为 JSON 原始值（例如 struct、map 或 slice）并实现 Annotation 接口的 Go 类型。

这是注释及其在模型和模板中的用法的示例
1. 定义一个注解
   ```
    package entgql

    // Annotation annotates fields with metadata for templates.
    type Annotation struct {
        // OrderField is the ordering field as defined in graphql schema.
        OrderField string
    }

    // Name implements ent.Annotation interface.
    func (Annotation) Name() string {
        return "EntGQL"
    }
   ```
2. ent/schema 中的注解用法
   ```
    // User schema.
   type User struct {
       ent.Schema
   }

   // Fields of the user.
   func (User) Fields() []ent.Field {
       return []ent.Field{
           field.Time("creation_date").
               Annotations(entgql.Annotation{
                   OrderField: "CREATED_AT",
               }),
       }
   }
   ```
3. 注解在外部模板中的用法：
   ```
    {{ range $node := $.Nodes }}
       {{ range $f := $node.Fields }}
           {{/* Get the annotation by its name. See: Annotation.Name */}}
           {{ if $annotation := $f.Annotations.EntGQL }}
               {{/* Get the field from the annotation. */}}
               {{ $orderField := $annotation.OrderField }}
           {{ end }}
       {{ end }}
   {{ end }}
   ```

### 特性开关
此框架提供了一系列代码生成特性，可以自行选择使用。
#### 用法
特性开关可以通过 CLI 标志或作为参数提供给 gen 包。
`go run entgo.io/ent/cmd/ent generate --feature privacy,entql ./ent/schema`

```
   package main

   import (
       "log"
       "text/template"

       "entgo.io/ent/entc"
       "entgo.io/ent/entc/gen"
   )

   func main() {
       err := entc.Generate("./schema", &gen.Config{
           Features: []gen.Feature{
               gen.FeaturePrivacy,
               gen.FeatureEntQL,
           },
           Templates: []*gen.Template{
               gen.MustParse(gen.NewTemplate("static").
                   Funcs(template.FuncMap{"title": strings.ToTitle}).
                   ParseFiles("template/static.tmpl")),
           },
       })
       if err != nil {
           log.Fatalf("running ent codegen: %v", err)
       }
   }
```

#### 特性列表
- 隐私层 --feature privacy
隐私层允许为数据库中实体的查询和更变操作配置隐私政策。
可以使用 --feature privacy 将此选项添加到项目中，其完整文档在[隐私页面](#privacy隐私策略)

- EntQL过滤器 --feature entql
entql配置项在运行时为不同的查询构造器提供了通用的动态筛选功能。
可以使用 --feature entql 标志将此选项添加到项目中，其完整文档在[隐私页面](#privacy隐私策略)

- 自动解决合并冲突 --feature schema/snapshot
schema/snapshot配置项告诉entc (ent 代码生成) ，对最新结构 (schema) 生成一个快照，当用户的 结构 (schema) 不能构建时，自动使用生成的快照解决合并冲突。
可以使用 --feature schema/snapshot 标志将此选项添加到项目中，但请参阅 [ent/ent/issues/852](https://github.com/ent/ent/issues/852) 以获取有关它的更多信息。

- Schema配置 --feature sql/schemaconfig
  使用sql/schemaconfig 配置项，你能给关系数据库中的对象定义别名，并将其映射到模型上。 当你的模型并不都在一个数据库下，而是根据 Schema 而有所不同，这很有用。
  可以使用 --feature sql/schemaconfig 将此选项添加到项目中，在生成代码之后，你就可以使用新的配置项，比如：
  ```
   c, err := ent.Open(dialect, conn, ent.AlternateSchema(ent.SchemaConfig{
        User: "usersdb",
        Car: "carsdb",
    }))
    c.User.Query().All(ctx) // SELECT * FROM `usersdb`.`users`
    c.Car.Query().All(ctx)  // SELECT * FROM `carsdb`.`cars`
  ```

- 行级锁
  sql/lock 选项允许使用 SQL SELECT ... FOR {UPDATE | 配置行级锁定。 SHARE} 语法。
  可以使用 --feature sql/lock 标志将此选项添加到项目中。
  ```
   tx, err := client.Tx(ctx)
    if err != nil {
        log.Fatal(err)
    }

    tx.Pet.Query().
        Where(pet.Name(name)).
        ForUpdate().
        Only(ctx)

    tx.Pet.Query().
        Where(pet.ID(id)).
        ForShare(
            sql.WithLockTables(pet.Table),
            sql.WithLockAction(sql.NoWait),
        ).
        Only(ctx)
  ```

- 自定义SQL修饰符
  sql/modifier 选项允许将自定义 SQL 修饰符添加到构建器并在执行之前改变语句。
  可以使用 --feature sql/modifier 标志将此选项添加到项目中。
  
  ``` 例子1：
   client.Pet.
    Query().
    Modify(func(s *sql.Selector) {
        s.Select("SUM(LENGTH(name))")
    }).
    IntX(ctx)
  ```
  上面的代码将产生以下 SQL 查询：
  `SELECT SUM(LENGTH(name)) FROM pet`

  ```例子2
   var v []struct {
      Count     int       `json:"count"`
      Price     int       `json:"price"`
      CreatedAt time.Time `json:"created_at"`
  }

  client.User.
      Query().
      Where(
          user.CreatedAtGT(x),
          user.CreatedAtLT(y),
      ).
      Modify(func(s *sql.Selector) {
          s.Select(
              sql.As(sql.Count("*"), "count"),
              sql.As(sql.Sum("price"), "price"),
              sql.As("DATE(created_at)", "created_at"),
          ).
          GroupBy("DATE(created_at)").
          OrderBy(sql.Desc("DATE(created_at)"))
      }).
      ScanX(ctx, &v)
  ```
  上面的代码将产生以下 SQL 查询：
  ```
  SELECT
      COUNT(*) AS `count`,
      SUM(`price`) AS `price`,
      DATE(created_at) AS `created_at`
  FROM
      `users`
  WHERE
      `created_at` > x AND `created_at` < y
  GROUP BY
      DATE(created_at)
  ORDER BY
      DATE(created_at) DESC
  ```

- Upsert
  sql/upsert 选项允许使用 SQL ON CONFLICT / ON DUPLICATE KEY 语法配置 upsert 和bulk-upsert 逻辑。 如需完整文档，请访问 [Upsert](#增删改查API)。

  可以使用 --feature sql/upsert 标志将此选项添加到项目中。
  ```
   // Use the new values that were set on create.
    id, err := client.User.
        Create().
        SetAge(30).
        SetName("Ariel").
        OnConflict().
        UpdateNewValues().
        ID(ctx)

    // In PostgreSQL, the conflict target is required.
    err := client.User.
        Create().
        SetAge(30).
        SetName("Ariel").
        OnConflictColumns(user.FieldName).
        UpdateNewValues().
        Exec(ctx)

    // Bulk upsert is also supported.
    client.User.
        CreateBulk(builders...).
        OnConflict(
            sql.ConflictWhere(...),
            sql.UpdateWhere(...),
        ).
        UpdateNewValues().
        Exec(ctx)

    // INSERT INTO "users" (...) VALUES ... ON CONFLICT WHERE ... DO UPDATE SET ... WHERE ...
  ```
