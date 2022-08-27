// Code generated by SQLBoiler 4.11.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Unfollowed is an object representing the database table.
type Unfollowed struct {
	UID int64 `boil:"uid" json:"uid" toml:"uid" yaml:"uid"`

	R *unfollowedR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L unfollowedL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UnfollowedColumns = struct {
	UID string
}{
	UID: "uid",
}

var UnfollowedTableColumns = struct {
	UID string
}{
	UID: "unfollowed.uid",
}

// Generated where

var UnfollowedWhere = struct {
	UID whereHelperint64
}{
	UID: whereHelperint64{field: "\"unfollowed\".\"uid\""},
}

// UnfollowedRels is where relationship names are stored.
var UnfollowedRels = struct {
}{}

// unfollowedR is where relationships are stored.
type unfollowedR struct {
}

// NewStruct creates a new relationship struct
func (*unfollowedR) NewStruct() *unfollowedR {
	return &unfollowedR{}
}

// unfollowedL is where Load methods for each relationship are stored.
type unfollowedL struct{}

var (
	unfollowedAllColumns            = []string{"uid"}
	unfollowedColumnsWithoutDefault = []string{"uid"}
	unfollowedColumnsWithDefault    = []string{}
	unfollowedPrimaryKeyColumns     = []string{"uid"}
	unfollowedGeneratedColumns      = []string{}
)

type (
	// UnfollowedSlice is an alias for a slice of pointers to Unfollowed.
	// This should almost always be used instead of []Unfollowed.
	UnfollowedSlice []*Unfollowed
	// UnfollowedHook is the signature for custom Unfollowed hook methods
	UnfollowedHook func(context.Context, boil.ContextExecutor, *Unfollowed) error

	unfollowedQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	unfollowedType                 = reflect.TypeOf(&Unfollowed{})
	unfollowedMapping              = queries.MakeStructMapping(unfollowedType)
	unfollowedPrimaryKeyMapping, _ = queries.BindMapping(unfollowedType, unfollowedMapping, unfollowedPrimaryKeyColumns)
	unfollowedInsertCacheMut       sync.RWMutex
	unfollowedInsertCache          = make(map[string]insertCache)
	unfollowedUpdateCacheMut       sync.RWMutex
	unfollowedUpdateCache          = make(map[string]updateCache)
	unfollowedUpsertCacheMut       sync.RWMutex
	unfollowedUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var unfollowedAfterSelectHooks []UnfollowedHook

var unfollowedBeforeInsertHooks []UnfollowedHook
var unfollowedAfterInsertHooks []UnfollowedHook

var unfollowedBeforeUpdateHooks []UnfollowedHook
var unfollowedAfterUpdateHooks []UnfollowedHook

var unfollowedBeforeDeleteHooks []UnfollowedHook
var unfollowedAfterDeleteHooks []UnfollowedHook

var unfollowedBeforeUpsertHooks []UnfollowedHook
var unfollowedAfterUpsertHooks []UnfollowedHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Unfollowed) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range unfollowedAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Unfollowed) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range unfollowedBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Unfollowed) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range unfollowedAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Unfollowed) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range unfollowedBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Unfollowed) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range unfollowedAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Unfollowed) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range unfollowedBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Unfollowed) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range unfollowedAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Unfollowed) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range unfollowedBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Unfollowed) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range unfollowedAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUnfollowedHook registers your hook function for all future operations.
func AddUnfollowedHook(hookPoint boil.HookPoint, unfollowedHook UnfollowedHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		unfollowedAfterSelectHooks = append(unfollowedAfterSelectHooks, unfollowedHook)
	case boil.BeforeInsertHook:
		unfollowedBeforeInsertHooks = append(unfollowedBeforeInsertHooks, unfollowedHook)
	case boil.AfterInsertHook:
		unfollowedAfterInsertHooks = append(unfollowedAfterInsertHooks, unfollowedHook)
	case boil.BeforeUpdateHook:
		unfollowedBeforeUpdateHooks = append(unfollowedBeforeUpdateHooks, unfollowedHook)
	case boil.AfterUpdateHook:
		unfollowedAfterUpdateHooks = append(unfollowedAfterUpdateHooks, unfollowedHook)
	case boil.BeforeDeleteHook:
		unfollowedBeforeDeleteHooks = append(unfollowedBeforeDeleteHooks, unfollowedHook)
	case boil.AfterDeleteHook:
		unfollowedAfterDeleteHooks = append(unfollowedAfterDeleteHooks, unfollowedHook)
	case boil.BeforeUpsertHook:
		unfollowedBeforeUpsertHooks = append(unfollowedBeforeUpsertHooks, unfollowedHook)
	case boil.AfterUpsertHook:
		unfollowedAfterUpsertHooks = append(unfollowedAfterUpsertHooks, unfollowedHook)
	}
}

// One returns a single unfollowed record from the query.
func (q unfollowedQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Unfollowed, error) {
	o := &Unfollowed{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for unfollowed")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Unfollowed records from the query.
func (q unfollowedQuery) All(ctx context.Context, exec boil.ContextExecutor) (UnfollowedSlice, error) {
	var o []*Unfollowed

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Unfollowed slice")
	}

	if len(unfollowedAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Unfollowed records in the query.
func (q unfollowedQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count unfollowed rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q unfollowedQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if unfollowed exists")
	}

	return count > 0, nil
}

// Unfolloweds retrieves all the records using an executor.
func Unfolloweds(mods ...qm.QueryMod) unfollowedQuery {
	mods = append(mods, qm.From("\"unfollowed\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"unfollowed\".*"})
	}

	return unfollowedQuery{q}
}

// FindUnfollowed retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUnfollowed(ctx context.Context, exec boil.ContextExecutor, uID int64, selectCols ...string) (*Unfollowed, error) {
	unfollowedObj := &Unfollowed{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"unfollowed\" where \"uid\"=$1", sel,
	)

	q := queries.Raw(query, uID)

	err := q.Bind(ctx, exec, unfollowedObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from unfollowed")
	}

	if err = unfollowedObj.doAfterSelectHooks(ctx, exec); err != nil {
		return unfollowedObj, err
	}

	return unfollowedObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Unfollowed) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no unfollowed provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(unfollowedColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	unfollowedInsertCacheMut.RLock()
	cache, cached := unfollowedInsertCache[key]
	unfollowedInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			unfollowedAllColumns,
			unfollowedColumnsWithDefault,
			unfollowedColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(unfollowedType, unfollowedMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(unfollowedType, unfollowedMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"unfollowed\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"unfollowed\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into unfollowed")
	}

	if !cached {
		unfollowedInsertCacheMut.Lock()
		unfollowedInsertCache[key] = cache
		unfollowedInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Unfollowed.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Unfollowed) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	unfollowedUpdateCacheMut.RLock()
	cache, cached := unfollowedUpdateCache[key]
	unfollowedUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			unfollowedAllColumns,
			unfollowedPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update unfollowed, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"unfollowed\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, unfollowedPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(unfollowedType, unfollowedMapping, append(wl, unfollowedPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update unfollowed row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for unfollowed")
	}

	if !cached {
		unfollowedUpdateCacheMut.Lock()
		unfollowedUpdateCache[key] = cache
		unfollowedUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q unfollowedQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for unfollowed")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for unfollowed")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UnfollowedSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), unfollowedPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"unfollowed\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, unfollowedPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in unfollowed slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all unfollowed")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Unfollowed) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no unfollowed provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(unfollowedColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	unfollowedUpsertCacheMut.RLock()
	cache, cached := unfollowedUpsertCache[key]
	unfollowedUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			unfollowedAllColumns,
			unfollowedColumnsWithDefault,
			unfollowedColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			unfollowedAllColumns,
			unfollowedPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert unfollowed, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(unfollowedPrimaryKeyColumns))
			copy(conflict, unfollowedPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"unfollowed\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(unfollowedType, unfollowedMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(unfollowedType, unfollowedMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert unfollowed")
	}

	if !cached {
		unfollowedUpsertCacheMut.Lock()
		unfollowedUpsertCache[key] = cache
		unfollowedUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Unfollowed record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Unfollowed) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Unfollowed provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), unfollowedPrimaryKeyMapping)
	sql := "DELETE FROM \"unfollowed\" WHERE \"uid\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from unfollowed")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for unfollowed")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q unfollowedQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no unfollowedQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from unfollowed")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for unfollowed")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UnfollowedSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(unfollowedBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), unfollowedPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"unfollowed\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, unfollowedPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from unfollowed slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for unfollowed")
	}

	if len(unfollowedAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Unfollowed) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindUnfollowed(ctx, exec, o.UID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UnfollowedSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UnfollowedSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), unfollowedPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"unfollowed\".* FROM \"unfollowed\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, unfollowedPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in UnfollowedSlice")
	}

	*o = slice

	return nil
}

// UnfollowedExists checks if the Unfollowed row exists.
func UnfollowedExists(ctx context.Context, exec boil.ContextExecutor, uID int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"unfollowed\" where \"uid\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, uID)
	}
	row := exec.QueryRowContext(ctx, sql, uID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if unfollowed exists")
	}

	return exists, nil
}
