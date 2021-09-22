package schema

import (
	"net/url"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// // Policy defines the privacy policy of the User.
// func (User) Policy() ent.Policy {
// 	return privacy.Policy{
// 		Mutation: privacy.MutationPolicy{
// 			// Deny if not set otherwise.
// 			rule.DenyIfNoViewer(),
// 			rule.AllowIfAdmin(),
// 			privacy.AlwaysDenyRule(),
// 		},
// 		Query: privacy.QueryPolicy{
// 			// Allow any viewer to read anything.
// 			privacy.AlwaysAllowRule(),
// 		},
// 	}
// }

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		CommMixin{},  //自定义mixin
		mixin.Time{}, //默认mixin
		// BaseMixin{},
		// TenantMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("age").
			Positive(), //允许最小值为1的验证器
		field.Float("rank").
			Optional(), //设置字段可选
		field.Bool("active").
			Default(false), //设置默认值为false
		field.String("name").
			Unique(), //设置字段值唯一
		field.Time("current_at").
			Default(time.Now). //设置默认值为当前时间
			Annotations(&entsql.Annotation{
				Default: "CURRENT_TIMESTAMP",
			}),
		field.JSON("url", &url.URL{}).
			Optional(),
		field.JSON("strings", []string{}).
			Optional(),
		field.Enum("state").
			Values("on", "off"). //设置值为枚举值on/off
			Optional(),
		field.UUID("uuid", uuid.UUID{}).
			Default(uuid.New), //设置默认值为UUID
		field.String("password").
			Optional().
			Sensitive(), //设置成敏感字段
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("cars", Car.Type),
		edge.From("groups", Group.Type).Ref("users"),
		edge.To("spouse", User.Type).Unique(),
		edge.To("following", User.Type).
			From("followers"),
		edge.To("friends", User.Type),
		edge.To("card", Card.Type).
			Unique(),
		edge.To("posts", Post.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.From("manage", Group.Type).
			Ref("admin"),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Annotations(entsql.Prefix(128)),
		index.Fields("age", "current_at").
			Annotations(
				entsql.PrefixColumn("age", 100),
				entsql.PrefixColumn("current_at", 200),
			),
	}
}
