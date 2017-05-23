package models

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/vattle/sqlboiler/strmangle"
	"gopkg.in/nullbio/null.v6"
)

// Shelf is an object representing the database table.
type Shelf struct {
	ID   int64       `boil:"id" json:"id" toml:"id" yaml:"id"`
	Area null.String `boil:"area" json:"area,omitempty" toml:"area" yaml:"area,omitempty"`

	R         *shelfR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L         shelfL  `boil:"-" json:"-" toml:"-" yaml:"-"`
	readonly  *Shelf
	whitelist []string
	operation string
}

var ShelfFieldMapping = map[string]string{
	"id":   "ID",
	"area": "Area",
}

// shelfR is where relationships are stored.
type shelfR struct {
	Books BookSlice
}

// shelfL is where Load methods for each relationship are stored.
type shelfL struct{}

var (
	shelfColumns               = []string{"id", "area"}
	shelfColumnsWithoutDefault = []string{"area"}
	shelfColumnsWithDefault    = []string{"id"}
	shelfPrimaryKeyColumns     = []string{"id"}
)

type (
	// ShelfSlice is an alias for a slice of pointers to Shelf.
	// This should generally be used opposed to []Shelf.
	ShelfSlice []*Shelf
	// ShelfHook is the signature for custom Shelf hook methods
	ShelfHook func(boil.Executor, *Shelf) error

	shelfQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	shelfType                 = reflect.TypeOf(&Shelf{})
	shelfMapping              = queries.MakeStructMapping(shelfType)
	shelfPrimaryKeyMapping, _ = queries.BindMapping(shelfType, shelfMapping, shelfPrimaryKeyColumns)
	shelfInsertCacheMut       sync.RWMutex
	shelfInsertCache          = make(map[string]insertCache)
	shelfUpdateCacheMut       sync.RWMutex
	shelfUpdateCache          = make(map[string]updateCache)
	shelfUpsertCacheMut       sync.RWMutex
	shelfUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var shelfBeforeInsertHooks []ShelfHook
var shelfBeforeUpdateHooks []ShelfHook
var shelfBeforeDeleteHooks []ShelfHook
var shelfBeforeUpsertHooks []ShelfHook

var shelfAfterInsertHooks []ShelfHook
var shelfAfterSelectHooks []ShelfHook
var shelfAfterUpdateHooks []ShelfHook
var shelfAfterDeleteHooks []ShelfHook
var shelfAfterUpsertHooks []ShelfHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Shelf) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range shelfBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Shelf) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range shelfBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Shelf) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range shelfBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Shelf) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range shelfBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Shelf) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range shelfAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Shelf) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range shelfAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Shelf) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range shelfAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Shelf) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range shelfAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Shelf) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range shelfAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddShelfHook registers your hook function for all future operations.
func AddShelfHook(hookPoint boil.HookPoint, shelfHook ShelfHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		shelfBeforeInsertHooks = append(shelfBeforeInsertHooks, shelfHook)
	case boil.BeforeUpdateHook:
		shelfBeforeUpdateHooks = append(shelfBeforeUpdateHooks, shelfHook)
	case boil.BeforeDeleteHook:
		shelfBeforeDeleteHooks = append(shelfBeforeDeleteHooks, shelfHook)
	case boil.BeforeUpsertHook:
		shelfBeforeUpsertHooks = append(shelfBeforeUpsertHooks, shelfHook)
	case boil.AfterInsertHook:
		shelfAfterInsertHooks = append(shelfAfterInsertHooks, shelfHook)
	case boil.AfterSelectHook:
		shelfAfterSelectHooks = append(shelfAfterSelectHooks, shelfHook)
	case boil.AfterUpdateHook:
		shelfAfterUpdateHooks = append(shelfAfterUpdateHooks, shelfHook)
	case boil.AfterDeleteHook:
		shelfAfterDeleteHooks = append(shelfAfterDeleteHooks, shelfHook)
	case boil.AfterUpsertHook:
		shelfAfterUpsertHooks = append(shelfAfterUpsertHooks, shelfHook)
	}
}

// OneP returns a single shelf record from the query, and panics on error.
func (q shelfQuery) OneP() *Shelf {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single shelf record from the query.
func (q shelfQuery) One() (*Shelf, error) {
	o := &Shelf{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for shelf")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all Shelf records from the query, and panics on error.
func (q shelfQuery) AllP() ShelfSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Shelf records from the query.
func (q shelfQuery) All() (ShelfSlice, error) {
	var o ShelfSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Shelf slice")
	}

	if len(shelfAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all Shelf records in the query, and panics on error.
func (q shelfQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Shelf records in the query.
func (q shelfQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count shelf rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q shelfQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q shelfQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if shelf exists")
	}

	return count > 0, nil
}

// BooksG retrieves all the book's book.
func (o *Shelf) BooksG(mods ...qm.QueryMod) bookQuery {
	return o.Books(boil.GetDB(), mods...)
}

// Books retrieves all the book's book with an executor.
func (o *Shelf) Books(exec boil.Executor, mods ...qm.QueryMod) bookQuery {
	queryMods := []qm.QueryMod{
		qm.Select("`a`.*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("`a`.`shelf_id`=?", o.ID),
	)

	query := Books(exec, queryMods...)
	queries.SetFrom(query.Query, "`book` as `a`")
	return query
}

// LoadBooks allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (shelfL) LoadBooks(e boil.Executor, singular bool, maybeShelf interface{}) error {
	var slice []*Shelf
	var object *Shelf

	count := 1
	if singular {
		object = maybeShelf.(*Shelf)
	} else {
		slice = *maybeShelf.(*ShelfSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &shelfR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &shelfR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from `book` where `shelf_id` in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load book")
	}
	defer results.Close()

	var resultSlice []*Book
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice book")
	}

	if len(bookAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Books = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.ShelfID.Int64 {
				local.R.Books = append(local.R.Books, foreign)
				break
			}
		}
	}

	return nil
}

// AddBooks adds the given related objects to the existing relationships
// of the shelf, optionally inserting them as new records.
// Appends related to o.R.Books.
// Sets related.R.Shelf appropriately.
func (o *Shelf) AddBooks(exec boil.Executor, insert bool, related ...*Book) error {
	var err error
	for _, rel := range related {
		rel.ShelfID.Int64 = o.ID
		rel.ShelfID.Valid = true
		if insert {
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			if err = rel.Update(exec, "shelf_id"); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}
		}
	}

	if o.R == nil {
		o.R = &shelfR{
			Books: related,
		}
	} else {
		o.R.Books = append(o.R.Books, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &bookR{
				Shelf: o,
			}
		} else {
			rel.R.Shelf = o
		}
	}
	return nil
}

// SetBooks removes all previously related items of the
// shelf replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Shelf's Books accordingly.
// Replaces o.R.Books with related.
// Sets related.R.Shelf's Books accordingly.
func (o *Shelf) SetBooks(exec boil.Executor, insert bool, related ...*Book) error {
	query := "update `book` set `shelf_id` = null where `shelf_id` = ?"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.Books {
			rel.ShelfID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Shelf = nil
		}

		o.R.Books = nil
	}
	return o.AddBooks(exec, insert, related...)
}

// RemoveBooks relationships from objects passed in.
// Removes related items from R.Books (uses pointer comparison, removal does not keep order)
// Sets related.R.Shelf.
func (o *Shelf) RemoveBooks(exec boil.Executor, related ...*Book) error {
	var err error
	for _, rel := range related {
		rel.ShelfID.Valid = false
		if rel.R != nil {
			rel.R.Shelf = nil
		}
		if err = rel.Update(exec, "shelf_id"); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.Books {
			if rel != ri {
				continue
			}

			ln := len(o.R.Books)
			if ln > 1 && i < ln-1 {
				o.R.Books[i] = o.R.Books[ln-1]
			}
			o.R.Books = o.R.Books[:ln-1]
			break
		}
	}

	return nil
}

// ShelvesG retrieves all records.
func ShelvesG(mods ...qm.QueryMod) shelfQuery {
	return Shelves(boil.GetDB(), mods...)
}

// Shelves retrieves all the records using an executor.
func Shelves(exec boil.Executor, mods ...qm.QueryMod) shelfQuery {
	mods = append(mods, qm.From("`shelf`"))
	return shelfQuery{NewQuery(exec, mods...)}
}

// FindShelfG retrieves a single record by ID.
func FindShelfG(id int64, selectCols ...string) (*Shelf, error) {
	return FindShelf(boil.GetDB(), id, selectCols...)
}

// FindShelfGP retrieves a single record by ID, and panics on error.
func FindShelfGP(id int64, selectCols ...string) *Shelf {
	retobj, err := FindShelf(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindShelf retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindShelf(exec boil.Executor, id int64, selectCols ...string) (*Shelf, error) {
	shelfObj := &Shelf{}
	shelfObj.readonly = &Shelf{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `shelf` where `id`=?", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(shelfObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from shelf")
	}

	// copy original obj to the readonly cache
	*shelfObj.readonly = *shelfObj
	return shelfObj, nil
}

// FindShelfP retrieves a single record by ID with an executor, and panics on error.
func FindShelfP(exec boil.Executor, id int64, selectCols ...string) *Shelf {
	retobj, err := FindShelf(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Shelf) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Shelf) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Shelf) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Shelf) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no shelf provided for insertion")
	}
	o.whitelist = whitelist
	o.operation = "INSERT"

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(shelfColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	shelfInsertCacheMut.RLock()
	cache, cached := shelfInsertCache[key]
	shelfInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			shelfColumns,
			shelfColumnsWithDefault,
			shelfColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(shelfType, shelfMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(shelfType, shelfMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO `shelf` (`%s`) VALUES (%s)", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `shelf` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, shelfPrimaryKeyColumns))
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	result, err := exec.Exec(cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into shelf")
	}

	var lastID int64
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == shelfMapping["ID"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.ID,
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, identifierCols...)
	}

	err = exec.QueryRow(cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for shelf")
	}

CacheNoHooks:
	if !cached {
		shelfInsertCacheMut.Lock()
		shelfInsertCache[key] = cache
		shelfInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Shelf record. See Update for
// whitelist behavior description.
func (o *Shelf) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Shelf record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Shelf) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Shelf, and panics on error.
// See Update for whitelist behavior description.
func (o *Shelf) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Shelf.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Shelf) Update(exec boil.Executor, whitelist ...string) error {
	o.whitelist = whitelist
	whitelist = o.Whitelist()

	o.operation = "UPDATE"
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	shelfUpdateCacheMut.RLock()
	cache, cached := shelfUpdateCache[key]
	shelfUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(shelfColumns, shelfPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update shelf, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `shelf` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, shelfPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(shelfType, shelfMapping, append(wl, shelfPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update shelf row")
	}

	if !cached {
		shelfUpdateCacheMut.Lock()
		shelfUpdateCache[key] = cache
		shelfUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q shelfQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q shelfQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for shelf")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o ShelfSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o ShelfSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o ShelfSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ShelfSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), shelfPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE `shelf` SET %s WHERE (`id`) IN (%s)",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(shelfPrimaryKeyColumns), len(colNames)+1, len(shelfPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in shelf slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Shelf) UpsertG(updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Shelf) UpsertGP(updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Shelf) UpsertP(exec boil.Executor, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Shelf) Upsert(exec boil.Executor, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no shelf provided for upsert")
	}
	o.whitelist = whitelist
	o.operation = "UPSERT"

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(shelfColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs postgres problems
	buf := strmangle.GetBuffer()
	for _, c := range updateColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range whitelist {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	shelfUpsertCacheMut.RLock()
	cache, cached := shelfUpsertCache[key]
	shelfUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			shelfColumns,
			shelfColumnsWithDefault,
			shelfColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			shelfColumns,
			shelfPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert shelf, could not build update column list")
		}

		cache.query = queries.BuildUpsertQueryMySQL(dialect, "shelf", update, whitelist)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `shelf` WHERE `id`=?",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
		)

		cache.valueMapping, err = queries.BindMapping(shelfType, shelfMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(shelfType, shelfMapping, ret)
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

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	result, err := exec.Exec(cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "models: unable to upsert for shelf")
	}

	var lastID int64
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = int64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == shelfMapping["ID"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.ID,
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, identifierCols...)
	}

	err = exec.QueryRow(cache.retQuery, identifierCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "models: unable to populate default values for shelf")
	}

CacheNoHooks:
	if !cached {
		shelfUpsertCacheMut.Lock()
		shelfUpsertCache[key] = cache
		shelfUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single Shelf record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Shelf) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Shelf record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Shelf) DeleteG() error {
	if o == nil {
		return errors.New("models: no Shelf provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Shelf record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Shelf) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Shelf record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Shelf) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Shelf provided for delete")
	}
	o.operation = "DELETE"

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), shelfPrimaryKeyMapping)
	sql := "DELETE FROM `shelf` WHERE `id`=?"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from shelf")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q shelfQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q shelfQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no shelfQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from shelf")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o ShelfSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o ShelfSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Shelf slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o ShelfSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ShelfSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Shelf slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(shelfBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), shelfPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM `shelf` WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, shelfPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(shelfPrimaryKeyColumns), 1, len(shelfPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from shelf slice")
	}

	if len(shelfAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Shelf) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Shelf) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Shelf) ReloadG() error {
	if o == nil {
		return errors.New("models: no Shelf provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Shelf) Reload(exec boil.Executor) error {
	ret, err := FindShelf(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ShelfSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ShelfSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ShelfSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty ShelfSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ShelfSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	shelves := ShelfSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), shelfPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT `shelf`.* FROM `shelf` WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, shelfPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(shelfPrimaryKeyColumns), 1, len(shelfPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&shelves)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ShelfSlice")
	}

	*o = shelves

	return nil
}

// ShelfExists checks if the Shelf row exists.
func ShelfExists(exec boil.Executor, id int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from `shelf` where `id`=? limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if shelf exists")
	}

	return exists, nil
}

// ShelfExistsG checks if the Shelf row exists.
func ShelfExistsG(id int64) (bool, error) {
	return ShelfExists(boil.GetDB(), id)
}

// ShelfExistsGP checks if the Shelf row exists. Panics on error.
func ShelfExistsGP(id int64) bool {
	e, err := ShelfExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// ShelfExistsP checks if the Shelf row exists. Panics on error.
func ShelfExistsP(exec boil.Executor, id int64) bool {
	e, err := ShelfExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Shelf) Changes() (ch *Changeset, err error) {
	ch = &Changeset{Table: "shelf",
		Changes: []*ChangeItem{}, Operation: o.Operation()}
	v := reflect.Indirect(reflect.ValueOf(o.readonly))
	vnew := reflect.Indirect(reflect.ValueOf(o))

	for _, c := range o.Whitelist() {
		if f, ok := ShelfFieldMapping[c]; ok {
			var before, after interface{}
			if v.IsValid() {
				before = v.FieldByName(f).Interface()
			}

			if vnew.IsValid() {
				after = vnew.FieldByName(f).Interface()
			}

			chitem := &ChangeItem{Name: c}
			if o.operation == "DELETE" {
				chitem.Before = before
			} else {
				chitem.Before = before
				chitem.After = after
			}

			if !reflect.DeepEqual(chitem.Before, chitem.After) {
				ch.Changes = append(ch.Changes, chitem)
			}
		}
	}

	return
}

// Calculates changed columns on the object
func (o *Shelf) Whitelist() (wl []string) {
	if len(o.whitelist) > 0 {
		return o.whitelist
	}

	// Calculates changed columns as whitelist
	v := reflect.Indirect(reflect.ValueOf(o.readonly))
	vnew := reflect.Indirect(reflect.ValueOf(o))

	for _, c := range shelfColumns {
		if f, ok := ShelfFieldMapping[c]; ok {
			var before, after interface{}
			if v.IsValid() {
				before = v.FieldByName(f).Interface()
			}

			if vnew.IsValid() {
				after = vnew.FieldByName(f).Interface()
			}
			if !reflect.DeepEqual(before, after) || o.operation == "DELETE" {
				wl = append(wl, c)
			}
		}
	}

	return
}

func (o *Shelf) Operation() string {
	return o.operation
}

// Generated change history hook for models
func init() {
	chFunc := func(exec boil.Executor, s *Shelf) error {
		if s == nil || exec == nil {
			return nil
		}

		ch, _ := s.Changes()
		if changeable, ok := exec.(Changeable); ok {
			changeable.AddChange(ch)
		}

		return nil
	}

	afterSel := func(exec boil.Executor, s *Shelf) error {
		if s == nil || exec == nil {
			return nil
		}

		s.readonly = &Shelf{}
		*s.readonly = *s
		return nil
	}

	beforeInsert := func(exec boil.Executor, s *Shelf) error {
		if s == nil || exec == nil {
			return nil
		}

		s.operation = "INSERT"
		return nil
	}

	beforeUpdate := func(exec boil.Executor, s *Shelf) error {
		if s == nil || exec == nil {
			return nil
		}

		s.operation = "UPDATE"
		return nil
	}

	beforeUpsert := func(exec boil.Executor, s *Shelf) error {
		if s == nil || exec == nil {
			return nil
		}

		s.operation = "UPSERT"
		return nil
	}

	afterDelete := func(exec boil.Executor, s *Shelf) error {
		if s == nil || exec == nil {
			return nil
		}

		s.operation = "DELETE"
		return nil
	}

	AddShelfHook(boil.AfterSelectHook, afterSel)
	AddShelfHook(boil.BeforeInsertHook, beforeInsert)
	AddShelfHook(boil.AfterInsertHook, chFunc)
	AddShelfHook(boil.BeforeUpdateHook, beforeUpdate)
	AddShelfHook(boil.AfterUpdateHook, chFunc)
	AddShelfHook(boil.BeforeUpsertHook, beforeUpsert)
	AddShelfHook(boil.AfterUpsertHook, chFunc)
	AddShelfHook(boil.AfterDeleteHook, afterDelete)
	AddShelfHook(boil.AfterDeleteHook, chFunc)
}
