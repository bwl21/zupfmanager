// Code generated by ent, DO NOT EDIT.

package project

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/bwl21/zupfmanager/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Project {
	return predicate.Project(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Project {
	return predicate.Project(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Project {
	return predicate.Project(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Project {
	return predicate.Project(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Project {
	return predicate.Project(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Project {
	return predicate.Project(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Project {
	return predicate.Project(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Project {
	return predicate.Project(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Project {
	return predicate.Project(sql.FieldLTE(FieldID, id))
}

// Title applies equality check predicate on the "title" field. It's identical to TitleEQ.
func Title(v string) predicate.Project {
	return predicate.Project(sql.FieldEQ(FieldTitle, v))
}

// ShortName applies equality check predicate on the "short_name" field. It's identical to ShortNameEQ.
func ShortName(v string) predicate.Project {
	return predicate.Project(sql.FieldEQ(FieldShortName, v))
}

// TitleEQ applies the EQ predicate on the "title" field.
func TitleEQ(v string) predicate.Project {
	return predicate.Project(sql.FieldEQ(FieldTitle, v))
}

// TitleNEQ applies the NEQ predicate on the "title" field.
func TitleNEQ(v string) predicate.Project {
	return predicate.Project(sql.FieldNEQ(FieldTitle, v))
}

// TitleIn applies the In predicate on the "title" field.
func TitleIn(vs ...string) predicate.Project {
	return predicate.Project(sql.FieldIn(FieldTitle, vs...))
}

// TitleNotIn applies the NotIn predicate on the "title" field.
func TitleNotIn(vs ...string) predicate.Project {
	return predicate.Project(sql.FieldNotIn(FieldTitle, vs...))
}

// TitleGT applies the GT predicate on the "title" field.
func TitleGT(v string) predicate.Project {
	return predicate.Project(sql.FieldGT(FieldTitle, v))
}

// TitleGTE applies the GTE predicate on the "title" field.
func TitleGTE(v string) predicate.Project {
	return predicate.Project(sql.FieldGTE(FieldTitle, v))
}

// TitleLT applies the LT predicate on the "title" field.
func TitleLT(v string) predicate.Project {
	return predicate.Project(sql.FieldLT(FieldTitle, v))
}

// TitleLTE applies the LTE predicate on the "title" field.
func TitleLTE(v string) predicate.Project {
	return predicate.Project(sql.FieldLTE(FieldTitle, v))
}

// TitleContains applies the Contains predicate on the "title" field.
func TitleContains(v string) predicate.Project {
	return predicate.Project(sql.FieldContains(FieldTitle, v))
}

// TitleHasPrefix applies the HasPrefix predicate on the "title" field.
func TitleHasPrefix(v string) predicate.Project {
	return predicate.Project(sql.FieldHasPrefix(FieldTitle, v))
}

// TitleHasSuffix applies the HasSuffix predicate on the "title" field.
func TitleHasSuffix(v string) predicate.Project {
	return predicate.Project(sql.FieldHasSuffix(FieldTitle, v))
}

// TitleEqualFold applies the EqualFold predicate on the "title" field.
func TitleEqualFold(v string) predicate.Project {
	return predicate.Project(sql.FieldEqualFold(FieldTitle, v))
}

// TitleContainsFold applies the ContainsFold predicate on the "title" field.
func TitleContainsFold(v string) predicate.Project {
	return predicate.Project(sql.FieldContainsFold(FieldTitle, v))
}

// ShortNameEQ applies the EQ predicate on the "short_name" field.
func ShortNameEQ(v string) predicate.Project {
	return predicate.Project(sql.FieldEQ(FieldShortName, v))
}

// ShortNameNEQ applies the NEQ predicate on the "short_name" field.
func ShortNameNEQ(v string) predicate.Project {
	return predicate.Project(sql.FieldNEQ(FieldShortName, v))
}

// ShortNameIn applies the In predicate on the "short_name" field.
func ShortNameIn(vs ...string) predicate.Project {
	return predicate.Project(sql.FieldIn(FieldShortName, vs...))
}

// ShortNameNotIn applies the NotIn predicate on the "short_name" field.
func ShortNameNotIn(vs ...string) predicate.Project {
	return predicate.Project(sql.FieldNotIn(FieldShortName, vs...))
}

// ShortNameGT applies the GT predicate on the "short_name" field.
func ShortNameGT(v string) predicate.Project {
	return predicate.Project(sql.FieldGT(FieldShortName, v))
}

// ShortNameGTE applies the GTE predicate on the "short_name" field.
func ShortNameGTE(v string) predicate.Project {
	return predicate.Project(sql.FieldGTE(FieldShortName, v))
}

// ShortNameLT applies the LT predicate on the "short_name" field.
func ShortNameLT(v string) predicate.Project {
	return predicate.Project(sql.FieldLT(FieldShortName, v))
}

// ShortNameLTE applies the LTE predicate on the "short_name" field.
func ShortNameLTE(v string) predicate.Project {
	return predicate.Project(sql.FieldLTE(FieldShortName, v))
}

// ShortNameContains applies the Contains predicate on the "short_name" field.
func ShortNameContains(v string) predicate.Project {
	return predicate.Project(sql.FieldContains(FieldShortName, v))
}

// ShortNameHasPrefix applies the HasPrefix predicate on the "short_name" field.
func ShortNameHasPrefix(v string) predicate.Project {
	return predicate.Project(sql.FieldHasPrefix(FieldShortName, v))
}

// ShortNameHasSuffix applies the HasSuffix predicate on the "short_name" field.
func ShortNameHasSuffix(v string) predicate.Project {
	return predicate.Project(sql.FieldHasSuffix(FieldShortName, v))
}

// ShortNameEqualFold applies the EqualFold predicate on the "short_name" field.
func ShortNameEqualFold(v string) predicate.Project {
	return predicate.Project(sql.FieldEqualFold(FieldShortName, v))
}

// ShortNameContainsFold applies the ContainsFold predicate on the "short_name" field.
func ShortNameContainsFold(v string) predicate.Project {
	return predicate.Project(sql.FieldContainsFold(FieldShortName, v))
}

// ConfigIsNil applies the IsNil predicate on the "config" field.
func ConfigIsNil() predicate.Project {
	return predicate.Project(sql.FieldIsNull(FieldConfig))
}

// ConfigNotNil applies the NotNil predicate on the "config" field.
func ConfigNotNil() predicate.Project {
	return predicate.Project(sql.FieldNotNull(FieldConfig))
}

// HasProjectSongs applies the HasEdge predicate on the "project_songs" edge.
func HasProjectSongs() predicate.Project {
	return predicate.Project(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, ProjectSongsTable, ProjectSongsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProjectSongsWith applies the HasEdge predicate on the "project_songs" edge with a given conditions (other predicates).
func HasProjectSongsWith(preds ...predicate.ProjectSong) predicate.Project {
	return predicate.Project(func(s *sql.Selector) {
		step := newProjectSongsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Project) predicate.Project {
	return predicate.Project(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Project) predicate.Project {
	return predicate.Project(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Project) predicate.Project {
	return predicate.Project(sql.NotPredicates(p))
}
