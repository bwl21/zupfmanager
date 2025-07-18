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
	"github.com/bwl21/zupfmanager/internal/ent/projectsong"
	"github.com/bwl21/zupfmanager/internal/ent/song"
)

// SongUpdate is the builder for updating Song entities.
type SongUpdate struct {
	config
	hooks    []Hook
	mutation *SongMutation
}

// Where appends a list predicates to the SongUpdate builder.
func (su *SongUpdate) Where(ps ...predicate.Song) *SongUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetTitle sets the "title" field.
func (su *SongUpdate) SetTitle(s string) *SongUpdate {
	su.mutation.SetTitle(s)
	return su
}

// SetNillableTitle sets the "title" field if the given value is not nil.
func (su *SongUpdate) SetNillableTitle(s *string) *SongUpdate {
	if s != nil {
		su.SetTitle(*s)
	}
	return su
}

// SetFilename sets the "filename" field.
func (su *SongUpdate) SetFilename(s string) *SongUpdate {
	su.mutation.SetFilename(s)
	return su
}

// SetNillableFilename sets the "filename" field if the given value is not nil.
func (su *SongUpdate) SetNillableFilename(s *string) *SongUpdate {
	if s != nil {
		su.SetFilename(*s)
	}
	return su
}

// SetGenre sets the "genre" field.
func (su *SongUpdate) SetGenre(s string) *SongUpdate {
	su.mutation.SetGenre(s)
	return su
}

// SetNillableGenre sets the "genre" field if the given value is not nil.
func (su *SongUpdate) SetNillableGenre(s *string) *SongUpdate {
	if s != nil {
		su.SetGenre(*s)
	}
	return su
}

// ClearGenre clears the value of the "genre" field.
func (su *SongUpdate) ClearGenre() *SongUpdate {
	su.mutation.ClearGenre()
	return su
}

// SetCopyright sets the "copyright" field.
func (su *SongUpdate) SetCopyright(s string) *SongUpdate {
	su.mutation.SetCopyright(s)
	return su
}

// SetNillableCopyright sets the "copyright" field if the given value is not nil.
func (su *SongUpdate) SetNillableCopyright(s *string) *SongUpdate {
	if s != nil {
		su.SetCopyright(*s)
	}
	return su
}

// ClearCopyright clears the value of the "copyright" field.
func (su *SongUpdate) ClearCopyright() *SongUpdate {
	su.mutation.ClearCopyright()
	return su
}

// AddProjectSongIDs adds the "project_songs" edge to the ProjectSong entity by IDs.
func (su *SongUpdate) AddProjectSongIDs(ids ...int) *SongUpdate {
	su.mutation.AddProjectSongIDs(ids...)
	return su
}

// AddProjectSongs adds the "project_songs" edges to the ProjectSong entity.
func (su *SongUpdate) AddProjectSongs(p ...*ProjectSong) *SongUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return su.AddProjectSongIDs(ids...)
}

// Mutation returns the SongMutation object of the builder.
func (su *SongUpdate) Mutation() *SongMutation {
	return su.mutation
}

// ClearProjectSongs clears all "project_songs" edges to the ProjectSong entity.
func (su *SongUpdate) ClearProjectSongs() *SongUpdate {
	su.mutation.ClearProjectSongs()
	return su
}

// RemoveProjectSongIDs removes the "project_songs" edge to ProjectSong entities by IDs.
func (su *SongUpdate) RemoveProjectSongIDs(ids ...int) *SongUpdate {
	su.mutation.RemoveProjectSongIDs(ids...)
	return su
}

// RemoveProjectSongs removes "project_songs" edges to ProjectSong entities.
func (su *SongUpdate) RemoveProjectSongs(p ...*ProjectSong) *SongUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return su.RemoveProjectSongIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *SongUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, su.sqlSave, su.mutation, su.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (su *SongUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *SongUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *SongUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (su *SongUpdate) check() error {
	if v, ok := su.mutation.Title(); ok {
		if err := song.TitleValidator(v); err != nil {
			return &ValidationError{Name: "title", err: fmt.Errorf(`ent: validator failed for field "Song.title": %w`, err)}
		}
	}
	if v, ok := su.mutation.Filename(); ok {
		if err := song.FilenameValidator(v); err != nil {
			return &ValidationError{Name: "filename", err: fmt.Errorf(`ent: validator failed for field "Song.filename": %w`, err)}
		}
	}
	return nil
}

func (su *SongUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := su.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(song.Table, song.Columns, sqlgraph.NewFieldSpec(song.FieldID, field.TypeInt))
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.Title(); ok {
		_spec.SetField(song.FieldTitle, field.TypeString, value)
	}
	if value, ok := su.mutation.Filename(); ok {
		_spec.SetField(song.FieldFilename, field.TypeString, value)
	}
	if value, ok := su.mutation.Genre(); ok {
		_spec.SetField(song.FieldGenre, field.TypeString, value)
	}
	if su.mutation.GenreCleared() {
		_spec.ClearField(song.FieldGenre, field.TypeString)
	}
	if value, ok := su.mutation.Copyright(); ok {
		_spec.SetField(song.FieldCopyright, field.TypeString, value)
	}
	if su.mutation.CopyrightCleared() {
		_spec.ClearField(song.FieldCopyright, field.TypeString)
	}
	if su.mutation.ProjectSongsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   song.ProjectSongsTable,
			Columns: []string{song.ProjectSongsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(projectsong.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RemovedProjectSongsIDs(); len(nodes) > 0 && !su.mutation.ProjectSongsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   song.ProjectSongsTable,
			Columns: []string{song.ProjectSongsColumn},
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
	if nodes := su.mutation.ProjectSongsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   song.ProjectSongsTable,
			Columns: []string{song.ProjectSongsColumn},
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
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{song.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	su.mutation.done = true
	return n, nil
}

// SongUpdateOne is the builder for updating a single Song entity.
type SongUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *SongMutation
}

// SetTitle sets the "title" field.
func (suo *SongUpdateOne) SetTitle(s string) *SongUpdateOne {
	suo.mutation.SetTitle(s)
	return suo
}

// SetNillableTitle sets the "title" field if the given value is not nil.
func (suo *SongUpdateOne) SetNillableTitle(s *string) *SongUpdateOne {
	if s != nil {
		suo.SetTitle(*s)
	}
	return suo
}

// SetFilename sets the "filename" field.
func (suo *SongUpdateOne) SetFilename(s string) *SongUpdateOne {
	suo.mutation.SetFilename(s)
	return suo
}

// SetNillableFilename sets the "filename" field if the given value is not nil.
func (suo *SongUpdateOne) SetNillableFilename(s *string) *SongUpdateOne {
	if s != nil {
		suo.SetFilename(*s)
	}
	return suo
}

// SetGenre sets the "genre" field.
func (suo *SongUpdateOne) SetGenre(s string) *SongUpdateOne {
	suo.mutation.SetGenre(s)
	return suo
}

// SetNillableGenre sets the "genre" field if the given value is not nil.
func (suo *SongUpdateOne) SetNillableGenre(s *string) *SongUpdateOne {
	if s != nil {
		suo.SetGenre(*s)
	}
	return suo
}

// ClearGenre clears the value of the "genre" field.
func (suo *SongUpdateOne) ClearGenre() *SongUpdateOne {
	suo.mutation.ClearGenre()
	return suo
}

// SetCopyright sets the "copyright" field.
func (suo *SongUpdateOne) SetCopyright(s string) *SongUpdateOne {
	suo.mutation.SetCopyright(s)
	return suo
}

// SetNillableCopyright sets the "copyright" field if the given value is not nil.
func (suo *SongUpdateOne) SetNillableCopyright(s *string) *SongUpdateOne {
	if s != nil {
		suo.SetCopyright(*s)
	}
	return suo
}

// ClearCopyright clears the value of the "copyright" field.
func (suo *SongUpdateOne) ClearCopyright() *SongUpdateOne {
	suo.mutation.ClearCopyright()
	return suo
}

// AddProjectSongIDs adds the "project_songs" edge to the ProjectSong entity by IDs.
func (suo *SongUpdateOne) AddProjectSongIDs(ids ...int) *SongUpdateOne {
	suo.mutation.AddProjectSongIDs(ids...)
	return suo
}

// AddProjectSongs adds the "project_songs" edges to the ProjectSong entity.
func (suo *SongUpdateOne) AddProjectSongs(p ...*ProjectSong) *SongUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return suo.AddProjectSongIDs(ids...)
}

// Mutation returns the SongMutation object of the builder.
func (suo *SongUpdateOne) Mutation() *SongMutation {
	return suo.mutation
}

// ClearProjectSongs clears all "project_songs" edges to the ProjectSong entity.
func (suo *SongUpdateOne) ClearProjectSongs() *SongUpdateOne {
	suo.mutation.ClearProjectSongs()
	return suo
}

// RemoveProjectSongIDs removes the "project_songs" edge to ProjectSong entities by IDs.
func (suo *SongUpdateOne) RemoveProjectSongIDs(ids ...int) *SongUpdateOne {
	suo.mutation.RemoveProjectSongIDs(ids...)
	return suo
}

// RemoveProjectSongs removes "project_songs" edges to ProjectSong entities.
func (suo *SongUpdateOne) RemoveProjectSongs(p ...*ProjectSong) *SongUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return suo.RemoveProjectSongIDs(ids...)
}

// Where appends a list predicates to the SongUpdate builder.
func (suo *SongUpdateOne) Where(ps ...predicate.Song) *SongUpdateOne {
	suo.mutation.Where(ps...)
	return suo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *SongUpdateOne) Select(field string, fields ...string) *SongUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Song entity.
func (suo *SongUpdateOne) Save(ctx context.Context) (*Song, error) {
	return withHooks(ctx, suo.sqlSave, suo.mutation, suo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (suo *SongUpdateOne) SaveX(ctx context.Context) *Song {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *SongUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *SongUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (suo *SongUpdateOne) check() error {
	if v, ok := suo.mutation.Title(); ok {
		if err := song.TitleValidator(v); err != nil {
			return &ValidationError{Name: "title", err: fmt.Errorf(`ent: validator failed for field "Song.title": %w`, err)}
		}
	}
	if v, ok := suo.mutation.Filename(); ok {
		if err := song.FilenameValidator(v); err != nil {
			return &ValidationError{Name: "filename", err: fmt.Errorf(`ent: validator failed for field "Song.filename": %w`, err)}
		}
	}
	return nil
}

func (suo *SongUpdateOne) sqlSave(ctx context.Context) (_node *Song, err error) {
	if err := suo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(song.Table, song.Columns, sqlgraph.NewFieldSpec(song.FieldID, field.TypeInt))
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Song.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, song.FieldID)
		for _, f := range fields {
			if !song.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != song.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := suo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := suo.mutation.Title(); ok {
		_spec.SetField(song.FieldTitle, field.TypeString, value)
	}
	if value, ok := suo.mutation.Filename(); ok {
		_spec.SetField(song.FieldFilename, field.TypeString, value)
	}
	if value, ok := suo.mutation.Genre(); ok {
		_spec.SetField(song.FieldGenre, field.TypeString, value)
	}
	if suo.mutation.GenreCleared() {
		_spec.ClearField(song.FieldGenre, field.TypeString)
	}
	if value, ok := suo.mutation.Copyright(); ok {
		_spec.SetField(song.FieldCopyright, field.TypeString, value)
	}
	if suo.mutation.CopyrightCleared() {
		_spec.ClearField(song.FieldCopyright, field.TypeString)
	}
	if suo.mutation.ProjectSongsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   song.ProjectSongsTable,
			Columns: []string{song.ProjectSongsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(projectsong.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RemovedProjectSongsIDs(); len(nodes) > 0 && !suo.mutation.ProjectSongsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   song.ProjectSongsTable,
			Columns: []string{song.ProjectSongsColumn},
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
	if nodes := suo.mutation.ProjectSongsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   song.ProjectSongsTable,
			Columns: []string{song.ProjectSongsColumn},
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
	_node = &Song{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{song.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	suo.mutation.done = true
	return _node, nil
}
