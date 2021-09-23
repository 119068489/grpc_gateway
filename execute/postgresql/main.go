package main

import (
	"context"
	"errors"
	"fmt"
	"grpc_gateway/easygo"
	"grpc_gateway/ent"
	"grpc_gateway/ent/car"
	"grpc_gateway/ent/card"
	"grpc_gateway/ent/group"
	"grpc_gateway/ent/node"
	"grpc_gateway/ent/post"
	"grpc_gateway/ent/predicate"
	"grpc_gateway/ent/user"
	"log"
	"time"

	_ "grpc_gateway/ent/runtime"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/examples/privacyadmin/viewer"
	"entgo.io/ent/privacy"
	"github.com/astaxie/beego/logs"
	_ "github.com/lib/pq"
)

var EntM *easygo.EntManager

func init() {
	EntM = easygo.NewEntManager()
}

func main() {
	psCfg := &easygo.PostgresCfg{
		Host:     "127.0.0.1",
		Port:     5432,
		User:     "postgres",
		PassWord: "123456",
		DbName:   "testdb",
	}
	client := EntM.Open(psCfg)
	defer client.Close()

	// Add a global hook that runs on all types and all operations.
	client.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			start := time.Now()
			defer func() {
				logs.Info("Op=%s\tType=%s\tTime=%s\tConcreteType=%T\n", m.Op(), m.Type(), time.Since(start), m)
			}()
			return next.Mutate(ctx, m)
		})
	})

	// CreateUser(context.Background(), client)
	// QueryUser(context.Background(), client)
	// UpdateUser(context.Background(), client.User.GetX(context.Background(), 1))
	// DeleteUser(context.Background(), client, client.User.GetX(context.Background(), 1))
	// CreateCars(context.Background(), client)
	// DeleteCar(context.Background(), client, client.Car.GetX(context.Background(), 2))
	// QueryCars(context.Background(), client.User.GetX(context.Background(), 4))
	// QueryCarUsers(context.Background(), client.User.GetX(context.Background(), 4))
	// CreateGraph(context.Background(), client)
	// QueryGithub(context.Background(), client)
	// QueryArielCars(context.Background(), client)
	// QueryGroupWithUsers(context.Background(), client)
	// QueryCarGroups(context.Background(), client.Car.GetX(context.Background(), 3))
	// NodeDo(context.Background(), client)
	// SpouseDo(context.Background(), client)
	// NodeTree(context.Background(), client)
	// FollowDo(context.Background(), client)
	// FriendsDo(context.Background(), client)
	// PostDo(context.Background(), client)
	// CityStreet(context.Background(), client)
	// CardDo(context.Background(), client)
	UpdateCard(context.Background(), client)
	// PrivacyDo(context.Background(), client)
	// ContentDo(context.Background(), client)
	// TenantDo(context.Background(), client)
	// GenTx(context.Background(), client)
	// TransactionDo(context.Background(), client)
	// GroupBys(context.Background(), client)
	// QueryUsers(context.Background(), client)
}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Create().
		SetAge(27).
		SetName("judth miss").
		Save(ctx)
	if err != nil {
		logs.Error(err)
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	logs.Info("user was created: ", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.Name("judth")).
		// `Only` 在 找不到用户 或 找到多于一个用户 时报错,
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("user returned: ", u.GoString())
	return u, nil
}

func UpdateUser(ctx context.Context, user *ent.User) (*ent.User, error) {
	u, err := user.
		Update().
		SetAge(27).
		SetName("judth").
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed update user: %w", err)
	}
	log.Println("user was updated: ", u)
	return u, nil
}

func DeleteUser(ctx context.Context, client *ent.Client, user *ent.User) error {
	logs.Debug("删除一个用户")
	err := client.User.
		DeleteOne(user).
		Exec(ctx)
	if err != nil {
		logs.Error(err)
		return fmt.Errorf("failed delete user: %w", err)
	}
	log.Println("user was deleted: ", user)
	return nil
}

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
		logs.Error(err)
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", a8m)
	return a8m, nil
}

func DeleteCar(ctx context.Context, client *ent.Client, car *ent.Car) error {
	logs.Debug("删除一辆车")
	err := client.Car.
		DeleteOne(car).
		Exec(ctx)
	if err != nil {
		logs.Error(err)
		return fmt.Errorf("failed delete user: %w", err)
	}
	log.Println("user was deleted: ", car)
	return nil
}

func QueryCars(ctx context.Context, user *ent.User) error {
	cars, err := user.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
	}
	log.Println("returned cars:", cars)

	// What about filtering specific cars.
	ford, err := user.QueryCars().
		Where(car.Model("Ford")).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
	}
	log.Println(ford)
	return nil
}

func QueryCarUsers(ctx context.Context, user *ent.User) error {
	cars, err := user.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %w", err)
	}
	// Query the inverse edge.
	for _, ca := range cars {
		owner, err := ca.QueryOwner().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying car %q owner: %w", ca.Model, err)
		}
		log.Printf("car %q owner: %q\n", ca.Model, owner.Name)
	}
	return nil
}

func CreateGraph(ctx context.Context, client *ent.Client) error {
	// First, create the users.
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetState(user.StateOn).
		SetName("Ariel").
		Save(ctx)
	if err != nil {
		return err
	}
	neta, err := client.User.
		Create().
		SetAge(28).
		SetState(user.StateOn).
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

func QueryArielCars(ctx context.Context, client *ent.Client) error {
	// Get "Ariel" from previous steps.
	a8m := client.User.
		Query().
		Where(
			user.HasCars(),
			user.Name("Ariel"),
		).
		OnlyX(ctx)
	cars, err := a8m. // Get the groups, that a8m is connected to:
				QueryGroups(). // (Group(Name=GitHub), Group(Name=GitLab),)
				QueryUsers().  // (User(Name=Ariel, Age=30), User(Name=Neta, Age=28),)
				QueryCars().   //
				Where(         //
			car.Not( //  Get Neta and Ariel cars, but filter out
				car.Model("Mazda"), //  those who named "Mazda"
			), //
		). //
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting cars: %w", err)
	}
	log.Println("cars returned:", cars)
	// Output: (Car(Model=Tesla, RegisteredAt=<Time>), Car(Model=Ford, RegisteredAt=<Time>),)
	return nil
}

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

func NodeDo(ctx context.Context, client *ent.Client) error {
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

func SpouseDo(ctx context.Context, client *ent.Client) error {
	jay, err := client.User. //Query().Where(user.NameHasPrefix("jay")).First(ctx)
					Create().
					SetAge(38).
					SetName("jay").
					Save(ctx)
	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}
	judth, err := client.User. //Query().Where(user.NameHasPrefix("judth")).First(ctx)
					Create().
					SetAge(26).
					SetName("judth").
					SetSpouse(jay).
					Save(ctx)
	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	// Query the spouse edge.
	// Unlike `Only`, `OnlyX` panics if an error occurs.
	spouse := judth.QuerySpouse().OnlyX(ctx)
	fmt.Println(spouse.Name)
	// Output: jay

	spouse = jay.QuerySpouse().OnlyX(ctx)
	fmt.Println(spouse.Name)
	// Output: judth

	// Query how many users have a spouse.
	// Unlike `Count`, `CountX` panics if an error occurs.
	count := client.User.
		Query().
		Where(user.HasSpouse()).
		CountX(ctx)
	fmt.Println(count)
	// Output: 2

	// Get the user, that has a spouse with name="jay".
	spouse = client.User.
		Query().
		Where(user.HasSpouseWith(user.Name("jay"))).
		OnlyX(ctx)
	fmt.Println(spouse.Name)
	// Output: judth
	return nil
}

func NodeTree(ctx context.Context, client *ent.Client) error {
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

func FollowDo(ctx context.Context, client *ent.Client) error {
	// Unlike `Save`, `SaveX` panics if an error occurs.
	a8m := client.User.
		Create().
		SetAge(30).
		SetName("a8m1").
		SaveX(ctx)
	nati := client.User.
		Create().
		SetAge(28).
		SetName("nati1").
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

func FriendsDo(ctx context.Context, client *ent.Client) error {
	// Unlike `Save`, `SaveX` panics if an error occurs.
	a8m := client.User.
		Create().
		SetAge(30).
		SetName("a8m2").
		SaveX(ctx)
	nati := client.User.
		Create().
		SetAge(28).
		SetName("nati2").
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

func PostDo(ctx context.Context, client *ent.Client) error {
	ps, err := client.Post.
		Create().
		SetAuthorID(3).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	u := client.User.Query().Where(user.ID(*ps.AuthorID)).OnlyX(ctx)

	p := client.Post.Query().
		Where(post.AuthorID(u.ID)).
		OnlyX(ctx)

	fmt.Println(p) // Access the "author" foreign-key.
	CreateUser(ctx, client)
	return nil
}

func CityStreet(ctx context.Context, client *ent.Client) error {
	// 和 `Save`不同，当出现错误是 `SaveX` 抛出panic。
	wh := client.City.
		Create().
		SetName("武汉").
		SaveX(ctx)
	gz := client.City.
		Create().
		SetName("广州").
		SaveX(ctx)
	// Add a street "吉庆街" to "武汉".
	client.Street.
		Create().
		SetName("吉庆街").
		SetCity(wh).
		SaveX(ctx)
	// 这一步的操作将会失败 because "吉庆街"
	// 因为 "吉庆街" 已经创建于 "武汉" 之下
	_, err := client.Street.
		Create().
		SetName("吉庆街").
		SetCity(wh).
		Save(ctx)
	if err == nil {
		return fmt.Errorf("expecting creation to fail")
	}
	// 将街道 "吉庆街" 添加到 "广州"
	client.Street.
		Create().
		SetName("吉庆街").
		SetCity(gz).
		SaveX(ctx)
	return nil
}

func CardDo(ctx context.Context, client *ent.Client) error {
	// 和 `Save`不同，当出现错误是 `SaveX` 抛出panic。
	user, err := client.User.
		Create().
		SetAge(38).
		SetName("visa4").
		Save(ctx)
	if err != nil {
		logs.Error(err)
	}

	card := client.Card.
		Create().
		SetNumber("9529564521478").
		SetName("boingings").
		SetOwner(user).
		SaveX(ctx)

	logs.Info(card)
	return nil
}

func UpdateCard(ctx context.Context, client *ent.Client) error {
	car := client.Card.Update().
		SetName("ccc").
		Where(card.ID(6)).
		SaveX(ctx)
	logs.Info(car)
	return nil
}
func PrivacyDo(ctx context.Context, client *ent.Client) error {
	// Expect operation to fail, because viewer-context
	// is missing (first mutation rule check).
	if _, err := client.User.Create().Save(ctx); !errors.Is(err, privacy.Deny) {
		logs.Error("expect operation to fail, but got %w", err)
		return fmt.Errorf("expect operation to fail, but got %w", err)
	}
	// Apply the same operation with "Admin" role.
	admin := viewer.NewContext(ctx, viewer.UserViewer{Role: viewer.Admin})
	if _, err := client.User.Create().Save(admin); err != nil {
		logs.Error("expect operation to pass, but got %w", err)
		return fmt.Errorf("expect operation to pass, but got %w", err)
	}
	// Apply the same operation with "ViewOnly" role.
	viewOnly := viewer.NewContext(ctx, viewer.UserViewer{Role: viewer.View})
	if _, err := client.User.Create().Save(viewOnly); !errors.Is(err, privacy.Deny) {
		logs.Error("expect operation to fail, but got %w", err)
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

func ContentDo(ctx context.Context, client *ent.Client) error {
	// Bind a privacy decision to the context (bypass all other rules).
	// allow := privacy.DecisionContext(ctx, privacy.Allow)
	admin := viewer.NewContext(ctx, viewer.UserViewer{Role: viewer.Admin})
	if _, err := client.User.Create().Save(admin); err != nil {
		logs.Error("expect operation to pass, but got %w", err)
		return fmt.Errorf("expect operation to pass, but got %w", err)
	}
	return nil
}

/*
func TenantDo(ctx context.Context, client *ent.Client) error {
	// Expect operation to fail, because viewer-context
	// is missing (first mutation rule check).
	if _, err := client.Tenant.Create().Save(ctx); !errors.Is(err, privacy.Deny) {
		logs.Error("1", err)
		return fmt.Errorf("expect operation to fail, but got %w", err)
	}
	// Deny tenant creation if the viewer is not admin.
	viewCtx := viewer.NewContext(ctx, viewer.UserViewer{Role: viewer.View})
	if _, err := client.Tenant.Create().Save(viewCtx); !errors.Is(err, privacy.Deny) {
		logs.Error("2", err)
		return fmt.Errorf("expect operation to fail, but got %w", err)
	}
	// Apply the same operation with "Admin" role, expect it to pass.
	adminCtx := viewer.NewContext(ctx, viewer.UserViewer{Role: viewer.Admin})
	hub, err := client.Tenant.Create().SetName("GitHub").Save(adminCtx)
	if err != nil {
		logs.Error("3", err)
		return fmt.Errorf("expect operation to pass, but got %w", err)
	}
	fmt.Println(hub)
	lab, err := client.Tenant.Create().SetName("GitLab").Save(adminCtx)
	if err != nil {
		logs.Error("4", err)
		return fmt.Errorf("expect operation to pass, but got %w", err)
	}
	fmt.Println(lab)

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
*/

// GenTx 在一次事务中生成一系列实体。
func GenTx(ctx context.Context, client *ent.Client) (error, error) {
	tx, err := client.Tx(ctx)
	if err != nil {
		logs.Error(err)
		return fmt.Errorf("starting a transaction: %w", err), nil
	}
	hub, err := tx.Group.
		Create().
		SetName("Github").
		Save(ctx)
	if err != nil {
		logs.Error(err)
		return fmt.Errorf("failed creating the group: %w", err), tx.Rollback()
	}
	// Create the admin of the group.
	dan, err := tx.User.
		Create().
		SetAge(29).
		SetName("Dan").
		AddManage(hub).
		Save(ctx)
	if err != nil {
		logs.Error(err)
		return err, tx.Rollback()
	}
	// Create user "Ariel".
	a8m, err := tx.User.
		Create().
		SetAge(26).
		SetName("Ariel").
		AddGroups(hub).
		AddFriends(dan).
		Save(ctx)
	if err != nil {
		logs.Error(err)
		return err, tx.Rollback()
	}
	fmt.Println(a8m)
	// Output:
	// User(id=2, age=30, name=Ariel)

	// Commit the transaction.
	return tx.Commit(), nil
}

func TransactionDo(ctx context.Context, client *ent.Client) {
	// WithTx helper.
	if err := WithTx(ctx, client, func(tx *ent.Tx) error {
		return PostDo(ctx, tx.Client())
	}); err != nil {
		log.Fatal(err)
	}
}

func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	// Add a hook on Tx.Commit.
	tx.OnCommit(func(next ent.Committer) ent.Committer {
		return ent.CommitFunc(func(ctx context.Context, tx *ent.Tx) error {
			// Code before the actual commit.
			logs.Info("事务准备提交")
			err := next.Commit(ctx, tx)
			// Code after the transaction was committed.
			logs.Info("事务提交完成")
			return err
		})
	})
	// Add a hook on Tx.Rollback.
	tx.OnRollback(func(next ent.Rollbacker) ent.Rollbacker {
		return ent.RollbackFunc(func(ctx context.Context, tx *ent.Tx) error {
			// Code before the actual rollback.
			logs.Info("事务准备回滚")
			err := next.Rollback(ctx, tx)
			// Code after the transaction was rolled back.
			logs.Info("事务回滚完成")
			return err
		})
	})

	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("rolling back transaction: %v", rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %v", err)
	}
	return nil
}

func GroupBys(ctx context.Context, client *ent.Client) {
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
	if err != nil {
		logs.Error(err)
	}
	logs.Info(v)
}

func GroupBy(ctx context.Context, client *ent.Client) {
	var v []struct {
		Age   int `json:"age"`
		Sum   int `json:"sum"`
		Count int `json:"count"`
	}
	err := client.User.
		Query().
		GroupBy(user.FieldAge).
		Aggregate(ent.Count(), ent.Sum(user.FieldAge)).
		Scan(ctx, &v)
	if err != nil {
		logs.Error(err)
	}
	logs.Info(v)
}

func QueryUsers(ctx context.Context, client *ent.Client) {
	// users, _ := client.User.Query().
	// 	Select(user.FieldID, user.FieldName).
	// 	Limit(3).
	// 	Offset(5).
	// 	Order(ent.Asc(user.FieldID)).
	// 	All(ctx)

	// users, _ := client.Card.Query().
	// 	Order(ent.Asc(card.FieldID)).
	// 	QueryOwner().
	// 	Order(ent.Asc(user.FieldID)).
	// 	All(ctx)

	users, _ := client.Card.Query().
		Order(func(s *sql.Selector) {
			// 连接用户表，以通过 owner-name 和 pet-name 进行排序.
			t := sql.Table(user.Table)
			s.Join(t).On(s.C(card.OwnerColumn), t.C(user.FieldID))
			s.OrderBy(t.C(user.FieldID), s.C(card.FieldID))
		}).
		Select(card.FieldName).
		Strings(ctx)

	logs.Info(users)
}
