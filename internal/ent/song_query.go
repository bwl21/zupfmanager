// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/bwl21/zupfmanager/internal/ent/predicate"
	"github.com/bwl21/zupfmanager/internal/ent/projectsong"
	"github.com/bwl21/zupfmanager/internal/ent/song"
)

// SongQuery is the builder for querying Song entities.
type SongQuery struct {
	config
	ctx              *QueryContext
	order            []song.OrderOption
	inters           []Interceptor
	predicates       []predicate.Song
	withProjectSongs *ProjectSongQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the SongQuery builder.
func (sq *SongQuery) Where(ps ...predicate.Song) *SongQuery {
	sq.predicates = append(sq.predicates, ps...)
	return sq
}

// Limit the number of records to be returned by this query.
func (sq *SongQuery) Limit(limit int) *SongQuery {
	sq.ctx.Limit = &limit
	return sq
}

// Offset to start from.
func (sq *SongQuery) Offset(offset int) *SongQuery {
	sq.ctx.Offset = &offset
	return sq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (sq *SongQuery) Unique(unique bool) *SongQuery {
	sq.ctx.Unique = &unique
	return sq
}

// Order specifies how the records should be ordered.
func (sq *SongQuery) Order(o ...song.OrderOption) *SongQuery {
	sq.order = append(sq.order, o...)
	return sq
}

// QueryProjectSongs chains the current query on the "project_songs" edge.
func (sq *SongQuery) QueryProjectSongs() *ProjectSongQuery {
	query := (&ProjectSongClient{config: sq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := sq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := sq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(song.Table, song.FieldID, selector),
			sqlgraph.To(projectsong.Table, projectsong.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, song.ProjectSongsTable, song.ProjectSongsColumn),
		)
		fromU = sqlgraph.SetNeighbors(sq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Song entity from the query.
// Returns a *NotFoundError when no Song was found.
func (sq *SongQuery) First(ctx context.Context) (*Song, error) {
	nodes, err := sq.Limit(1).All(setContextOp(ctx, sq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{song.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (sq *SongQuery) FirstX(ctx context.Context) *Song {
	node, err := sq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Song ID from the query.
// Returns a *NotFoundError when no Song ID was found.
func (sq *SongQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = sq.Limit(1).IDs(setContextOp(ctx, sq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{song.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (sq *SongQuery) FirstIDX(ctx context.Context) int {
	id, err := sq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Song entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Song entity is found.
// Returns a *NotFoundError when no Song entities are found.
func (sq *SongQuery) Only(ctx context.Context) (*Song, error) {
	nodes, err := sq.Limit(2).All(setContextOp(ctx, sq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{song.Label}
	default:
		return nil, &NotSingularError{song.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (sq *SongQuery) OnlyX(ctx context.Context) *Song {
	node, err := sq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Song ID in the query.
// Returns a *NotSingularError when more than one Song ID is found.
// Returns a *NotFoundError when no entities are found.
func (sq *SongQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = sq.Limit(2).IDs(setContextOp(ctx, sq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{song.Label}
	default:
		err = &NotSingularError{song.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (sq *SongQuery) OnlyIDX(ctx context.Context) int {
	id, err := sq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Songs.
func (sq *SongQuery) All(ctx context.Context) ([]*Song, error) {
	ctx = setContextOp(ctx, sq.ctx, ent.OpQueryAll)
	if err := sq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Song, *SongQuery]()
	return withInterceptors[[]*Song](ctx, sq, qr, sq.inters)
}

// AllX is like All, but panics if an error occurs.
func (sq *SongQuery) AllX(ctx context.Context) []*Song {
	nodes, err := sq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Song IDs.
func (sq *SongQuery) IDs(ctx context.Context) (ids []int, err error) {
	if sq.ctx.Unique == nil && sq.path != nil {
		sq.Unique(true)
	}
	ctx = setContextOp(ctx, sq.ctx, ent.OpQueryIDs)
	if err = sq.Select(song.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (sq *SongQuery) IDsX(ctx context.Context) []int {
	ids, err := sq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (sq *SongQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, sq.ctx, ent.OpQueryCount)
	if err := sq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, sq, querierCount[*SongQuery](), sq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (sq *SongQuery) CountX(ctx context.Context) int {
	count, err := sq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (sq *SongQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, sq.ctx, ent.OpQueryExist)
	switch _, err := sq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (sq *SongQuery) ExistX(ctx context.Context) bool {
	exist, err := sq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the SongQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (sq *SongQuery) Clone() *SongQuery {
	if sq == nil {
		return nil
	}
	return &SongQuery{
		config:           sq.config,
		ctx:              sq.ctx.Clone(),
		order:            append([]song.OrderOption{}, sq.order...),
		inters:           append([]Interceptor{}, sq.inters...),
		predicates:       append([]predicate.Song{}, sq.predicates...),
		withProjectSongs: sq.withProjectSongs.Clone(),
		// clone intermediate query.
		sql:  sq.sql.Clone(),
		path: sq.path,
	}
}

// WithProjectSongs tells the query-builder to eager-load the nodes that are connected to
// the "project_songs" edge. The optional arguments are used to configure the query builder of the edge.
func (sq *SongQuery) WithProjectSongs(opts ...func(*ProjectSongQuery)) *SongQuery {
	query := (&ProjectSongClient{config: sq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	sq.withProjectSongs = query
	return sq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Title string `json:"title,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Song.Query().
//		GroupBy(song.FieldTitle).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (sq *SongQuery) GroupBy(field string, fields ...string) *SongGroupBy {
	sq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &SongGroupBy{build: sq}
	grbuild.flds = &sq.ctx.Fields
	grbuild.label = song.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Title string `json:"title,omitempty"`
//	}
//
//	client.Song.Query().
//		Select(song.FieldTitle).
//		Scan(ctx, &v)
func (sq *SongQuery) Select(fields ...string) *SongSelect {
	sq.ctx.Fields = append(sq.ctx.Fields, fields...)
	sbuild := &SongSelect{SongQuery: sq}
	sbuild.label = song.Label
	sbuild.flds, sbuild.scan = &sq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a SongSelect configured with the given aggregations.
func (sq *SongQuery) Aggregate(fns ...AggregateFunc) *SongSelect {
	return sq.Select().Aggregate(fns...)
}

func (sq *SongQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range sq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, sq); err != nil {
				return err
			}
		}
	}
	for _, f := range sq.ctx.Fields {
		if !song.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if sq.path != nil {
		prev, err := sq.path(ctx)
		if err != nil {
			return err
		}
		sq.sql = prev
	}
	return nil
}

func (sq *SongQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Song, error) {
	var (
		nodes       = []*Song{}
		_spec       = sq.querySpec()
		loadedTypes = [1]bool{
			sq.withProjectSongs != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Song).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Song{config: sq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, sq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := sq.withProjectSongs; query != nil {
		if err := sq.loadProjectSongs(ctx, query, nodes,
			func(n *Song) { n.Edges.ProjectSongs = []*ProjectSong{} },
			func(n *Song, e *ProjectSong) { n.Edges.ProjectSongs = append(n.Edges.ProjectSongs, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (sq *SongQuery) loadProjectSongs(ctx context.Context, query *ProjectSongQuery, nodes []*Song, init func(*Song), assign func(*Song, *ProjectSong)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Song)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(projectsong.FieldSongID)
	}
	query.Where(predicate.ProjectSong(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(song.ProjectSongsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.SongID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "song_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (sq *SongQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := sq.querySpec()
	_spec.Node.Columns = sq.ctx.Fields
	if len(sq.ctx.Fields) > 0 {
		_spec.Unique = sq.ctx.Unique != nil && *sq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, sq.driver, _spec)
}

func (sq *SongQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(song.Table, song.Columns, sqlgraph.NewFieldSpec(song.FieldID, field.TypeInt))
	_spec.From = sq.sql
	if unique := sq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if sq.path != nil {
		_spec.Unique = true
	}
	if fields := sq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, song.FieldID)
		for i := range fields {
			if fields[i] != song.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := sq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := sq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := sq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := sq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (sq *SongQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(sq.driver.Dialect())
	t1 := builder.Table(song.Table)
	columns := sq.ctx.Fields
	if len(columns) == 0 {
		columns = song.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if sq.sql != nil {
		selector = sq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if sq.ctx.Unique != nil && *sq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range sq.predicates {
		p(selector)
	}
	for _, p := range sq.order {
		p(selector)
	}
	if offset := sq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := sq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// SongGroupBy is the group-by builder for Song entities.
type SongGroupBy struct {
	selector
	build *SongQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (sgb *SongGroupBy) Aggregate(fns ...AggregateFunc) *SongGroupBy {
	sgb.fns = append(sgb.fns, fns...)
	return sgb
}

// Scan applies the selector query and scans the result into the given value.
func (sgb *SongGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, sgb.build.ctx, ent.OpQueryGroupBy)
	if err := sgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*SongQuery, *SongGroupBy](ctx, sgb.build, sgb, sgb.build.inters, v)
}

func (sgb *SongGroupBy) sqlScan(ctx context.Context, root *SongQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(sgb.fns))
	for _, fn := range sgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*sgb.flds)+len(sgb.fns))
		for _, f := range *sgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*sgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := sgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// SongSelect is the builder for selecting fields of Song entities.
type SongSelect struct {
	*SongQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ss *SongSelect) Aggregate(fns ...AggregateFunc) *SongSelect {
	ss.fns = append(ss.fns, fns...)
	return ss
}

// Scan applies the selector query and scans the result into the given value.
func (ss *SongSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ss.ctx, ent.OpQuerySelect)
	if err := ss.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*SongQuery, *SongSelect](ctx, ss.SongQuery, ss, ss.inters, v)
}

func (ss *SongSelect) sqlScan(ctx context.Context, root *SongQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ss.fns))
	for _, fn := range ss.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ss.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ss.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
