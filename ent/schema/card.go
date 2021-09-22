package schema

import (
	"context"
	"fmt"

	gen "grpc_gateway/ent"
	"grpc_gateway/ent/hook"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/astaxie/beego/logs"
)

// Card holds the schema definition for the Card entity.
type Card struct {
	ent.Schema
}

// Hooks of the Card.
func (Card) Hooks() []ent.Hook {

	hookA := func(next ent.Mutator) ent.Mutator {
		return hook.CardFunc(func(ctx context.Context, m *gen.CardMutation) (ent.Value, error) {
			if num, ok := m.Number(); ok && len(num) < 10 {
				logs.Error("card number is too short")
				return nil, fmt.Errorf("card number is too short")
			}
			return next.Mutate(ctx, m)
		})
		// return nil
	}

	hookB := func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if s, ok := m.(interface{ SetName(string) }); ok {
				s.SetName("Boring")
			}
			return next.Mutate(ctx, m)
		})
	}

	hookC := func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if s, ok := m.(interface {
				SetName(string)
			}); ok {
				s.SetName("Visa")
			}
			return next.Mutate(ctx, m)
		})
	}

	return []ent.Hook{
		// First hook. 限制卡号长度
		hook.On(
			hookA,
			// 只在插入更新时执行钩子
			ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne,
		),

		// Second hook.统一名字
		hook.Unless(
			hookB,
			//不在插入执行这个钩子
			ent.OpCreate,
		),
		// third hook.名字变化时修改名字并清掉owner_id字段
		hook.If(
			hookC,
			hook.And(hook.HasFields("name"), hook.HasClearedFields("owner_id")),
		),
	}
}

// Fields of the Card.
func (Card) Fields() []ent.Field {
	return []ent.Field{
		field.String("number").
			Optional(),
		field.String("name").
			Optional(),
		field.Int("owner_id").
			Optional(),
	}
}

// Edges of the Card.
func (Card) Edges() []ent.Edge {
	return []ent.Edge{
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
