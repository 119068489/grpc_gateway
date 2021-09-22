// Code generated by entc, DO NOT EDIT.

package card

import (
	"grpc_gateway/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// Number applies equality check predicate on the "number" field. It's identical to NumberEQ.
func Number(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldNumber), v))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// OwnerID applies equality check predicate on the "owner_id" field. It's identical to OwnerIDEQ.
func OwnerID(v int) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldOwnerID), v))
	})
}

// NumberEQ applies the EQ predicate on the "number" field.
func NumberEQ(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldNumber), v))
	})
}

// NumberNEQ applies the NEQ predicate on the "number" field.
func NumberNEQ(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldNumber), v))
	})
}

// NumberIn applies the In predicate on the "number" field.
func NumberIn(vs ...string) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldNumber), v...))
	})
}

// NumberNotIn applies the NotIn predicate on the "number" field.
func NumberNotIn(vs ...string) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldNumber), v...))
	})
}

// NumberGT applies the GT predicate on the "number" field.
func NumberGT(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldNumber), v))
	})
}

// NumberGTE applies the GTE predicate on the "number" field.
func NumberGTE(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldNumber), v))
	})
}

// NumberLT applies the LT predicate on the "number" field.
func NumberLT(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldNumber), v))
	})
}

// NumberLTE applies the LTE predicate on the "number" field.
func NumberLTE(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldNumber), v))
	})
}

// NumberContains applies the Contains predicate on the "number" field.
func NumberContains(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldNumber), v))
	})
}

// NumberHasPrefix applies the HasPrefix predicate on the "number" field.
func NumberHasPrefix(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldNumber), v))
	})
}

// NumberHasSuffix applies the HasSuffix predicate on the "number" field.
func NumberHasSuffix(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldNumber), v))
	})
}

// NumberIsNil applies the IsNil predicate on the "number" field.
func NumberIsNil() predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldNumber)))
	})
}

// NumberNotNil applies the NotNil predicate on the "number" field.
func NumberNotNil() predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldNumber)))
	})
}

// NumberEqualFold applies the EqualFold predicate on the "number" field.
func NumberEqualFold(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldNumber), v))
	})
}

// NumberContainsFold applies the ContainsFold predicate on the "number" field.
func NumberContainsFold(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldNumber), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldName), v...))
	})
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldName), v...))
	})
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameIsNil applies the IsNil predicate on the "name" field.
func NameIsNil() predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldName)))
	})
}

// NameNotNil applies the NotNil predicate on the "name" field.
func NameNotNil() predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldName)))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// OwnerIDEQ applies the EQ predicate on the "owner_id" field.
func OwnerIDEQ(v int) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldOwnerID), v))
	})
}

// OwnerIDNEQ applies the NEQ predicate on the "owner_id" field.
func OwnerIDNEQ(v int) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldOwnerID), v))
	})
}

// OwnerIDIn applies the In predicate on the "owner_id" field.
func OwnerIDIn(vs ...int) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldOwnerID), v...))
	})
}

// OwnerIDNotIn applies the NotIn predicate on the "owner_id" field.
func OwnerIDNotIn(vs ...int) predicate.Card {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Card(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldOwnerID), v...))
	})
}

// OwnerIDIsNil applies the IsNil predicate on the "owner_id" field.
func OwnerIDIsNil() predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldOwnerID)))
	})
}

// OwnerIDNotNil applies the NotNil predicate on the "owner_id" field.
func OwnerIDNotNil() predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldOwnerID)))
	})
}

// HasOwner applies the HasEdge predicate on the "owner" edge.
func HasOwner() predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(OwnerTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, OwnerTable, OwnerColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasOwnerWith applies the HasEdge predicate on the "owner" edge with a given conditions (other predicates).
func HasOwnerWith(preds ...predicate.User) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(OwnerInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, OwnerTable, OwnerColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Card) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Card) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Card) predicate.Card {
	return predicate.Card(func(s *sql.Selector) {
		p(s.Not())
	})
}