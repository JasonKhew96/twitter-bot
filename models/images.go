// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

// Image is an object representing the database table.
type Image struct {
	ID        int       `boil:"id" json:"id" toml:"id" yaml:"id"`
	HashA     string    `boil:"hash_a" json:"hash_a" toml:"hash_a" yaml:"hash_a"`
	HashB     string    `boil:"hash_b" json:"hash_b" toml:"hash_b" yaml:"hash_b"`
	HashC     string    `boil:"hash_c" json:"hash_c" toml:"hash_c" yaml:"hash_c"`
	HashD     string    `boil:"hash_d" json:"hash_d" toml:"hash_d" yaml:"hash_d"`
	ChatID    int64     `boil:"chat_id" json:"chat_id" toml:"chat_id" yaml:"chat_id"`
	MessageID int64     `boil:"message_id" json:"message_id" toml:"message_id" yaml:"message_id"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *imageR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L imageL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ImageColumns = struct {
	ID        string
	HashA     string
	HashB     string
	HashC     string
	HashD     string
	ChatID    string
	MessageID string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	HashA:     "hash_a",
	HashB:     "hash_b",
	HashC:     "hash_c",
	HashD:     "hash_d",
	ChatID:    "chat_id",
	MessageID: "message_id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

var ImageTableColumns = struct {
	ID        string
	HashA     string
	HashB     string
	HashC     string
	HashD     string
	ChatID    string
	MessageID string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "images.id",
	HashA:     "images.hash_a",
	HashB:     "images.hash_b",
	HashC:     "images.hash_c",
	HashD:     "images.hash_d",
	ChatID:    "images.chat_id",
	MessageID: "images.message_id",
	CreatedAt: "images.created_at",
	UpdatedAt: "images.updated_at",
}

// Generated where

type whereHelperint struct{ field string }

func (w whereHelperint) EQ(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint) NEQ(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint) LT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint) LTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint) GT(x int) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint) GTE(x int) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint) IN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint) NIN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod     { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod     { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod     { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperstring) LIKE(x string) qm.QueryMod   { return qm.Where(w.field+" LIKE ?", x) }
func (w whereHelperstring) NLIKE(x string) qm.QueryMod  { return qm.Where(w.field+" NOT LIKE ?", x) }
func (w whereHelperstring) ILIKE(x string) qm.QueryMod  { return qm.Where(w.field+" ILIKE ?", x) }
func (w whereHelperstring) NILIKE(x string) qm.QueryMod { return qm.Where(w.field+" NOT ILIKE ?", x) }
func (w whereHelperstring) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperstring) NIN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelperint64 struct{ field string }

func (w whereHelperint64) EQ(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint64) NEQ(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint64) LT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint64) LTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint64) GT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint64) GTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint64) IN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint64) NIN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelpertime_Time struct{ field string }

func (w whereHelpertime_Time) EQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertime_Time) NEQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertime_Time) LT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertime_Time) LTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertime_Time) GT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertime_Time) GTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var ImageWhere = struct {
	ID        whereHelperint
	HashA     whereHelperstring
	HashB     whereHelperstring
	HashC     whereHelperstring
	HashD     whereHelperstring
	ChatID    whereHelperint64
	MessageID whereHelperint64
	CreatedAt whereHelpertime_Time
	UpdatedAt whereHelpertime_Time
}{
	ID:        whereHelperint{field: "\"images\".\"id\""},
	HashA:     whereHelperstring{field: "\"images\".\"hash_a\""},
	HashB:     whereHelperstring{field: "\"images\".\"hash_b\""},
	HashC:     whereHelperstring{field: "\"images\".\"hash_c\""},
	HashD:     whereHelperstring{field: "\"images\".\"hash_d\""},
	ChatID:    whereHelperint64{field: "\"images\".\"chat_id\""},
	MessageID: whereHelperint64{field: "\"images\".\"message_id\""},
	CreatedAt: whereHelpertime_Time{field: "\"images\".\"created_at\""},
	UpdatedAt: whereHelpertime_Time{field: "\"images\".\"updated_at\""},
}

// ImageRels is where relationship names are stored.
var ImageRels = struct {
}{}

// imageR is where relationships are stored.
type imageR struct {
}

// NewStruct creates a new relationship struct
func (*imageR) NewStruct() *imageR {
	return &imageR{}
}

// imageL is where Load methods for each relationship are stored.
type imageL struct{}

var (
	imageAllColumns            = []string{"id", "hash_a", "hash_b", "hash_c", "hash_d", "chat_id", "message_id", "created_at", "updated_at"}
	imageColumnsWithoutDefault = []string{"hash_a", "hash_b", "hash_c", "hash_d", "chat_id", "message_id", "created_at", "updated_at"}
	imageColumnsWithDefault    = []string{"id"}
	imagePrimaryKeyColumns     = []string{"id"}
	imageGeneratedColumns      = []string{}
)

type (
	// ImageSlice is an alias for a slice of pointers to Image.
	// This should almost always be used instead of []Image.
	ImageSlice []*Image
	// ImageHook is the signature for custom Image hook methods
	ImageHook func(context.Context, boil.ContextExecutor, *Image) error

	imageQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	imageType                 = reflect.TypeOf(&Image{})
	imageMapping              = queries.MakeStructMapping(imageType)
	imagePrimaryKeyMapping, _ = queries.BindMapping(imageType, imageMapping, imagePrimaryKeyColumns)
	imageInsertCacheMut       sync.RWMutex
	imageInsertCache          = make(map[string]insertCache)
	imageUpdateCacheMut       sync.RWMutex
	imageUpdateCache          = make(map[string]updateCache)
	imageUpsertCacheMut       sync.RWMutex
	imageUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var imageAfterSelectMu sync.Mutex
var imageAfterSelectHooks []ImageHook

var imageBeforeInsertMu sync.Mutex
var imageBeforeInsertHooks []ImageHook
var imageAfterInsertMu sync.Mutex
var imageAfterInsertHooks []ImageHook

var imageBeforeUpdateMu sync.Mutex
var imageBeforeUpdateHooks []ImageHook
var imageAfterUpdateMu sync.Mutex
var imageAfterUpdateHooks []ImageHook

var imageBeforeDeleteMu sync.Mutex
var imageBeforeDeleteHooks []ImageHook
var imageAfterDeleteMu sync.Mutex
var imageAfterDeleteHooks []ImageHook

var imageBeforeUpsertMu sync.Mutex
var imageBeforeUpsertHooks []ImageHook
var imageAfterUpsertMu sync.Mutex
var imageAfterUpsertHooks []ImageHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Image) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Image) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Image) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Image) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Image) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Image) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Image) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Image) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Image) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range imageAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddImageHook registers your hook function for all future operations.
func AddImageHook(hookPoint boil.HookPoint, imageHook ImageHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		imageAfterSelectMu.Lock()
		imageAfterSelectHooks = append(imageAfterSelectHooks, imageHook)
		imageAfterSelectMu.Unlock()
	case boil.BeforeInsertHook:
		imageBeforeInsertMu.Lock()
		imageBeforeInsertHooks = append(imageBeforeInsertHooks, imageHook)
		imageBeforeInsertMu.Unlock()
	case boil.AfterInsertHook:
		imageAfterInsertMu.Lock()
		imageAfterInsertHooks = append(imageAfterInsertHooks, imageHook)
		imageAfterInsertMu.Unlock()
	case boil.BeforeUpdateHook:
		imageBeforeUpdateMu.Lock()
		imageBeforeUpdateHooks = append(imageBeforeUpdateHooks, imageHook)
		imageBeforeUpdateMu.Unlock()
	case boil.AfterUpdateHook:
		imageAfterUpdateMu.Lock()
		imageAfterUpdateHooks = append(imageAfterUpdateHooks, imageHook)
		imageAfterUpdateMu.Unlock()
	case boil.BeforeDeleteHook:
		imageBeforeDeleteMu.Lock()
		imageBeforeDeleteHooks = append(imageBeforeDeleteHooks, imageHook)
		imageBeforeDeleteMu.Unlock()
	case boil.AfterDeleteHook:
		imageAfterDeleteMu.Lock()
		imageAfterDeleteHooks = append(imageAfterDeleteHooks, imageHook)
		imageAfterDeleteMu.Unlock()
	case boil.BeforeUpsertHook:
		imageBeforeUpsertMu.Lock()
		imageBeforeUpsertHooks = append(imageBeforeUpsertHooks, imageHook)
		imageBeforeUpsertMu.Unlock()
	case boil.AfterUpsertHook:
		imageAfterUpsertMu.Lock()
		imageAfterUpsertHooks = append(imageAfterUpsertHooks, imageHook)
		imageAfterUpsertMu.Unlock()
	}
}

// One returns a single image record from the query.
func (q imageQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Image, error) {
	o := &Image{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for images")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Image records from the query.
func (q imageQuery) All(ctx context.Context, exec boil.ContextExecutor) (ImageSlice, error) {
	var o []*Image

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Image slice")
	}

	if len(imageAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Image records in the query.
func (q imageQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count images rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q imageQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if images exists")
	}

	return count > 0, nil
}

// Images retrieves all the records using an executor.
func Images(mods ...qm.QueryMod) imageQuery {
	mods = append(mods, qm.From("\"images\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"images\".*"})
	}

	return imageQuery{q}
}

// FindImage retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindImage(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*Image, error) {
	imageObj := &Image{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"images\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, imageObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from images")
	}

	if err = imageObj.doAfterSelectHooks(ctx, exec); err != nil {
		return imageObj, err
	}

	return imageObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Image) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no images provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(imageColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	imageInsertCacheMut.RLock()
	cache, cached := imageInsertCache[key]
	imageInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			imageAllColumns,
			imageColumnsWithDefault,
			imageColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(imageType, imageMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(imageType, imageMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"images\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"images\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into images")
	}

	if !cached {
		imageInsertCacheMut.Lock()
		imageInsertCache[key] = cache
		imageInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Image.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Image) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	imageUpdateCacheMut.RLock()
	cache, cached := imageUpdateCache[key]
	imageUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			imageAllColumns,
			imagePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update images, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"images\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, imagePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(imageType, imageMapping, append(wl, imagePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update images row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for images")
	}

	if !cached {
		imageUpdateCacheMut.Lock()
		imageUpdateCache[key] = cache
		imageUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q imageQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for images")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for images")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ImageSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), imagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"images\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, imagePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in image slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all image")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Image) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("models: no images provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(imageColumnsWithDefault, o)

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

	imageUpsertCacheMut.RLock()
	cache, cached := imageUpsertCache[key]
	imageUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			imageAllColumns,
			imageColumnsWithDefault,
			imageColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			imageAllColumns,
			imagePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert images, could not build update column list")
		}

		ret := strmangle.SetComplement(imageAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(imagePrimaryKeyColumns) == 0 {
				return errors.New("models: unable to upsert images, could not build conflict column list")
			}

			conflict = make([]string, len(imagePrimaryKeyColumns))
			copy(conflict, imagePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"images\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(imageType, imageMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(imageType, imageMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert images")
	}

	if !cached {
		imageUpsertCacheMut.Lock()
		imageUpsertCache[key] = cache
		imageUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Image record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Image) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Image provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), imagePrimaryKeyMapping)
	sql := "DELETE FROM \"images\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from images")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for images")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q imageQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no imageQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from images")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for images")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ImageSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(imageBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), imagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"images\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, imagePrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from image slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for images")
	}

	if len(imageAfterDeleteHooks) != 0 {
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
func (o *Image) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindImage(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ImageSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ImageSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), imagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"images\".* FROM \"images\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, imagePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ImageSlice")
	}

	*o = slice

	return nil
}

// ImageExists checks if the Image row exists.
func ImageExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"images\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if images exists")
	}

	return exists, nil
}

// Exists checks if the Image row exists.
func (o *Image) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return ImageExists(ctx, exec, o.ID)
}