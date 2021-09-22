package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/examples/privacytenant/ent/privacy"
	"entgo.io/ent/examples/privacytenant/rule"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// TimeMixin implements the ent.Mixin for sharing
// time fields with package schemas.
type CommMixin struct {
	// We embed the `mixin.Schema` to avoid
	// implementing the rest of the methods.
	mixin.Schema
}

func (CommMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("nickname").
			Optional(),
	}
}

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
			rule.AllowIfAdmin(),
			// Filter out entities that are not connected to the tenant.
			// If the viewer is admin, this policy rule is skipped above.
			rule.FilterTenantRule(),
		},
	}
}

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
