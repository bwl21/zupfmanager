// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/bwl21/zupfmanager/internal/ent/predicate"
	"github.com/bwl21/zupfmanager/internal/ent/project"
	"github.com/bwl21/zupfmanager/internal/ent/projectsong"
)

// ProjectUpdate is the builder for updating Project entities.
type ProjectUpdate struct {
	config
	hooks    []Hook
	mutation *ProjectMutation
}

// Where appends a list predicates to the ProjectUpdate builder.
func (pu *ProjectUpdate) Where(ps ...predicate.Project) *ProjectUpdate {
	pu.mutation.Where(ps...)
	return pu
}

// SetTitle sets the "title" field.
func (pu *ProjectUpdate) SetTitle(s string) *ProjectUpdate {
	pu.mutation.SetTitle(s)
	return pu
}

// SetNillableTitle sets the "title" field if the given value is not nil.
func (pu *ProjectUpdate) SetNillableTitle(s *string) *ProjectUpdate {
	if s != nil {
		pu.SetTitle(*s)
	}
	return pu
}

// SetShortName sets the "short_name" field.
func (pu *ProjectUpdate) SetShortName(s string) *ProjectUpdate {
	pu.mutation.SetShortName(s)
	return pu
}

// SetNillableShortName sets the "short_name" field if the given value is not nil.
func (pu *ProjectUpdate) SetNillableShortName(s *string) *ProjectUpdate {
	if s != nil {
		pu.SetShortName(*s)
	}
	return pu
}

// SetConfig sets the "config" field.
func (pu *ProjectUpdate) SetConfig(m map[string]interface{}) *ProjectUpdate {
	pu.mutation.SetConfig(m)
	return pu
}

// ClearConfig clears the value of the "config" field.
func (pu *ProjectUpdate) ClearConfig() *ProjectUpdate {
	pu.mutation.ClearConfig()
	return pu
}

// AddProjectSongIDs adds the "project_songs" edge to the ProjectSong entity by IDs.
func (pu *ProjectUpdate) AddProjectSongIDs(ids ...int) *ProjectUpdate {
	pu.mutation.AddProjectSongIDs(ids...)
	return pu
}

// AddProjectSongs adds the "project_songs" edges to the ProjectSong entity.
func (pu *ProjectUpdate) AddProjectSongs(p ...*ProjectSong) *ProjectUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pu.AddProjectSongIDs(ids...)
}

// Mutation returns the ProjectMutation object of the builder.
func (pu *ProjectUpdate) Mutation() *ProjectMutation {
	return pu.mutation
}

// ClearProjectSongs clears all "project_songs" edges to the ProjectSong entity.
func (pu *ProjectUpdate) ClearProjectSongs() *ProjectUpdate {
	pu.mutation.ClearProjectSongs()
	return pu
}

// RemoveProjectSongIDs removes the "project_songs" edge to ProjectSong entities by IDs.
func (pu *ProjectUpdate) RemoveProjectSongIDs(ids ...int) *ProjectUpdate {
	pu.mutation.RemoveProjectSongIDs(ids...)
	return pu
}

// RemoveProjectSongs removes "project_songs" edges to ProjectSong entities.
func (pu *ProjectUpdate) RemoveProjectSongs(p ...*ProjectSong) *ProjectUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return pu.RemoveProjectSongIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pu *ProjectUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, pu.sqlSave, pu.mutation, pu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pu *ProjectUpdate) SaveX(ctx context.Context) int {
	affected, err := pu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pu *ProjectUpdate) Exec(ctx context.Context) error {
	_, err := pu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pu *ProjectUpdate) ExecX(ctx context.Context) {
	if err := pu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pu *ProjectUpdate) check() error {
	if v, ok := pu.mutation.Title(); ok {
		if err := project.TitleValidator(v); err != nil {
			return &ValidationError{Name: "title", err: fmt.Errorf(`ent: validator failed for field "Project.title": %w`, err)}
		}
	}
	if v, ok := pu.mutation.ShortName(); ok {
		if err := project.ShortNameValidator(v); err != nil {
			return &ValidationError{Name: "short_name", err: fmt.Errorf(`ent: validator failed for field "Project.short_name": %w`, err)}
		}
	}
	return nil
}

func (pu *ProjectUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := pu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(project.Table, project.Columns, sqlgraph.NewFieldSpec(project.FieldID, field.TypeInt))
	if ps := pu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pu.mutation.Title(); ok {
		_spec.SetField(project.FieldTitle, field.TypeString, value)
	}
	if value, ok := pu.mutation.ShortName(); ok {
		_spec.SetField(project.FieldShortName, field.TypeString, value)
	}
	if value, ok := pu.mutation.Config(); ok {
		_spec.SetField(project.FieldConfig, field.TypeJSON, value)
	}
	if pu.mutation.ConfigCleared() {
		_spec.ClearField(project.FieldConfig, field.TypeJSON)
	}
	if pu.mutation.ProjectSongsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ProjectSongsTable,
			Columns: []string{project.ProjectSongsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(projectsong.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.RemovedProjectSongsIDs(); len(nodes) > 0 && !pu.mutation.ProjectSongsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ProjectSongsTable,
			Columns: []string{project.ProjectSongsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(projectsong.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.ProjectSongsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ProjectSongsTable,
			Columns: []string{project.ProjectSongsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(projectsong.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, pu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{project.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	pu.mutation.done = true
	return n, nil
}

// ProjectUpdateOne is the builder for updating a single Project entity.
type ProjectUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ProjectMutation
}

// SetTitle sets the "title" field.
func (puo *ProjectUpdateOne) SetTitle(s string) *ProjectUpdateOne {
	puo.mutation.SetTitle(s)
	return puo
}

// SetNillableTitle sets the "title" field if the given value is not nil.
func (puo *ProjectUpdateOne) SetNillableTitle(s *string) *ProjectUpdateOne {
	if s != nil {
		puo.SetTitle(*s)
	}
	return puo
}

// SetShortName sets the "short_name" field.
func (puo *ProjectUpdateOne) SetShortName(s string) *ProjectUpdateOne {
	puo.mutation.SetShortName(s)
	return puo
}

// SetNillableShortName sets the "short_name" field if the given value is not nil.
func (puo *ProjectUpdateOne) SetNillableShortName(s *string) *ProjectUpdateOne {
	if s != nil {
		puo.SetShortName(*s)
	}
	return puo
}

// SetConfig sets the "config" field.
func (puo *ProjectUpdateOne) SetConfig(m map[string]interface{}) *ProjectUpdateOne {
	puo.mutation.SetConfig(m)
	return puo
}

// ClearConfig clears the value of the "config" field.
func (puo *ProjectUpdateOne) ClearConfig() *ProjectUpdateOne {
	puo.mutation.ClearConfig()
	return puo
}

// AddProjectSongIDs adds the "project_songs" edge to the ProjectSong entity by IDs.
func (puo *ProjectUpdateOne) AddProjectSongIDs(ids ...int) *ProjectUpdateOne {
	puo.mutation.AddProjectSongIDs(ids...)
	return puo
}

// AddProjectSongs adds the "project_songs" edges to the ProjectSong entity.
func (puo *ProjectUpdateOne) AddProjectSongs(p ...*ProjectSong) *ProjectUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return puo.AddProjectSongIDs(ids...)
}

// Mutation returns the ProjectMutation object of the builder.
func (puo *ProjectUpdateOne) Mutation() *ProjectMutation {
	return puo.mutation
}

// ClearProjectSongs clears all "project_songs" edges to the ProjectSong entity.
func (puo *ProjectUpdateOne) ClearProjectSongs() *ProjectUpdateOne {
	puo.mutation.ClearProjectSongs()
	return puo
}

// RemoveProjectSongIDs removes the "project_songs" edge to ProjectSong entities by IDs.
func (puo *ProjectUpdateOne) RemoveProjectSongIDs(ids ...int) *ProjectUpdateOne {
	puo.mutation.RemoveProjectSongIDs(ids...)
	return puo
}

// RemoveProjectSongs removes "project_songs" edges to ProjectSong entities.
func (puo *ProjectUpdateOne) RemoveProjectSongs(p ...*ProjectSong) *ProjectUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return puo.RemoveProjectSongIDs(ids...)
}

// Where appends a list predicates to the ProjectUpdate builder.
func (puo *ProjectUpdateOne) Where(ps ...predicate.Project) *ProjectUpdateOne {
	puo.mutation.Where(ps...)
	return puo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (puo *ProjectUpdateOne) Select(field string, fields ...string) *ProjectUpdateOne {
	puo.fields = append([]string{field}, fields...)
	return puo
}

// Save executes the query and returns the updated Project entity.
func (puo *ProjectUpdateOne) Save(ctx context.Context) (*Project, error) {
	return withHooks(ctx, puo.sqlSave, puo.mutation, puo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (puo *ProjectUpdateOne) SaveX(ctx context.Context) *Project {
	node, err := puo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (puo *ProjectUpdateOne) Exec(ctx context.Context) error {
	_, err := puo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (puo *ProjectUpdateOne) ExecX(ctx context.Context) {
	if err := puo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (puo *ProjectUpdateOne) check() error {
	if v, ok := puo.mutation.Title(); ok {
		if err := project.TitleValidator(v); err != nil {
			return &ValidationError{Name: "title", err: fmt.Errorf(`ent: validator failed for field "Project.title": %w`, err)}
		}
	}
	if v, ok := puo.mutation.ShortName(); ok {
		if err := project.ShortNameValidator(v); err != nil {
			return &ValidationError{Name: "short_name", err: fmt.Errorf(`ent: validator failed for field "Project.short_name": %w`, err)}
		}
	}
	return nil
}

func (puo *ProjectUpdateOne) sqlSave(ctx context.Context) (_node *Project, err error) {
	if err := puo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(project.Table, project.Columns, sqlgraph.NewFieldSpec(project.FieldID, field.TypeInt))
	id, ok := puo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Project.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := puo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, project.FieldID)
		for _, f := range fields {
			if !project.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != project.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := puo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := puo.mutation.Title(); ok {
		_spec.SetField(project.FieldTitle, field.TypeString, value)
	}
	if value, ok := puo.mutation.ShortName(); ok {
		_spec.SetField(project.FieldShortName, field.TypeString, value)
	}
	if value, ok := puo.mutation.Config(); ok {
		_spec.SetField(project.FieldConfig, field.TypeJSON, value)
	}
	if puo.mutation.ConfigCleared() {
		_spec.ClearField(project.FieldConfig, field.TypeJSON)
	}
	if puo.mutation.ProjectSongsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ProjectSongsTable,
			Columns: []string{project.ProjectSongsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(projectsong.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.RemovedProjectSongsIDs(); len(nodes) > 0 && !puo.mutation.ProjectSongsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ProjectSongsTable,
			Columns: []string{project.ProjectSongsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(projectsong.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.ProjectSongsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   project.ProjectSongsTable,
			Columns: []string{project.ProjectSongsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(projectsong.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Project{config: puo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, puo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{project.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	puo.mutation.done = true
	return _node, nil
}
