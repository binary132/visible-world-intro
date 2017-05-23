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

// Book is an object representing the database table.
type Book struct {
	ID      int64       `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name    null.String `boil:"name" json:"name,omitempty" toml:"name" yaml:"name,omitempty"`
	Author  null.String `boil:"author" json:"author,omitempty" toml:"author" yaml:"author,omitempty"`
	ShelfID null.Int64  `boil:"shelf_id" json:"shelf_id,omitempty" toml:"shelf_id" yaml:"shelf_id,omitempty"`

	R         *bookR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L         bookL  `boil:"-" json:"-" toml:"-" yaml:"-"`
	readonly  *Book
	whitelist []string
	operation string
}

var BookFieldMapping = map[string]string{
	"id":       "ID",
	"name":     "Name",
	"author":   "Author",
	"shelf_id": "ShelfID",
}

// bookR is where relationships are stored.
type bookR struct {
	Shelf *Shelf
}

// bookL is where Load methods for each relationship are stored.
type bookL struct{}

var (
	bookColumns               = []string{"id", "name", "author", "shelf_id"}
	bookColumnsWithoutDefault = []string{"name", "author", "shelf_id"}
	bookColumnsWithDefault    = []string{"id"}
	bookPrimaryKeyColumns     = []string{"id"}
)

type (
	// BookSlice is an alias for a slice of pointers to Book.
	// This should generally be used opposed to []Book.
	BookSlice []*Book
	// BookHook is the signature for custom Book hook methods
	BookHook func(boil.Executor, *Book) error

	bookQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	bookType                 = reflect.TypeOf(&Book{})
	bookMapping              = queries.MakeStructMapping(bookType)
	bookPrimaryKeyMapping, _ = queries.BindMapping(bookType, bookMapping, bookPrimaryKeyColumns)
	bookInsertCacheMut       sync.RWMutex
	bookInsertCache          = make(map[string]insertCache)
	bookUpdateCacheMut       sync.RWMutex
	bookUpdateCache          = make(map[string]updateCache)
	bookUpsertCacheMut       sync.RWMutex
	bookUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var bookBeforeInsertHooks []BookHook
var bookBeforeUpdateHooks []BookHook
var bookBeforeDeleteHooks []BookHook
var bookBeforeUpsertHooks []BookHook

var bookAfterInsertHooks []BookHook
var bookAfterSelectHooks []BookHook
var bookAfterUpdateHooks []BookHook
var bookAfterDeleteHooks []BookHook
var bookAfterUpsertHooks []BookHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Book) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range bookBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Book) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range bookBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Book) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range bookBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Book) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range bookBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Book) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range bookAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Book) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range bookAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Book) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range bookAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Book) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range bookAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Book) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range bookAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddBookHook registers your hook function for all future operations.
func AddBookHook(hookPoint boil.HookPoint, bookHook BookHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		bookBeforeInsertHooks = append(bookBeforeInsertHooks, bookHook)
	case boil.BeforeUpdateHook:
		bookBeforeUpdateHooks = append(bookBeforeUpdateHooks, bookHook)
	case boil.BeforeDeleteHook:
		bookBeforeDeleteHooks = append(bookBeforeDeleteHooks, bookHook)
	case boil.BeforeUpsertHook:
		bookBeforeUpsertHooks = append(bookBeforeUpsertHooks, bookHook)
	case boil.AfterInsertHook:
		bookAfterInsertHooks = append(bookAfterInsertHooks, bookHook)
	case boil.AfterSelectHook:
		bookAfterSelectHooks = append(bookAfterSelectHooks, bookHook)
	case boil.AfterUpdateHook:
		bookAfterUpdateHooks = append(bookAfterUpdateHooks, bookHook)
	case boil.AfterDeleteHook:
		bookAfterDeleteHooks = append(bookAfterDeleteHooks, bookHook)
	case boil.AfterUpsertHook:
		bookAfterUpsertHooks = append(bookAfterUpsertHooks, bookHook)
	}
}

// OneP returns a single book record from the query, and panics on error.
func (q bookQuery) OneP() *Book {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single book record from the query.
func (q bookQuery) One() (*Book, error) {
	o := &Book{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for book")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all Book records from the query, and panics on error.
func (q bookQuery) AllP() BookSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Book records from the query.
func (q bookQuery) All() (BookSlice, error) {
	var o BookSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Book slice")
	}

	if len(bookAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all Book records in the query, and panics on error.
func (q bookQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Book records in the query.
func (q bookQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count book rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q bookQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q bookQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if book exists")
	}

	return count > 0, nil
}

// ShelfG pointed to by the foreign key.
func (o *Book) ShelfG(mods ...qm.QueryMod) shelfQuery {
	return o.ShelfF(boil.GetDB(), mods...)
}

// Shelf pointed to by the foreign key.
func (o *Book) ShelfF(exec boil.Executor, mods ...qm.QueryMod) shelfQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.ShelfID),
	}

	queryMods = append(queryMods, mods...)

	query := Shelves(exec, queryMods...)
	queries.SetFrom(query.Query, "`shelf`")

	return query
}

// LoadShelf allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (bookL) LoadShelf(e boil.Executor, singular bool, maybeBook interface{}) error {
	var slice []*Book
	var object *Book

	count := 1
	if singular {
		object = maybeBook.(*Book)
	} else {
		slice = *maybeBook.(*BookSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &bookR{}
		}
		args[0] = object.ShelfID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &bookR{}
			}
			args[i] = obj.ShelfID
		}
	}

	query := fmt.Sprintf(
		"select * from `shelf` where `id` in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Shelf")
	}
	defer results.Close()

	var resultSlice []*Shelf
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Shelf")
	}

	if len(bookAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if singular && len(resultSlice) != 0 {
		object.R.Shelf = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ShelfID.Int64 == foreign.ID {
				local.R.Shelf = foreign
				break
			}
		}
	}

	return nil
}

// SetShelf of the book to the related item.
// Sets o.R.Shelf to related.
// Adds o to related.R.Books.
func (o *Book) SetShelf(exec boil.Executor, insert bool, related *Shelf) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `book` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"shelf_id"}),
		strmangle.WhereClause("`", "`", 0, bookPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ShelfID.Int64 = related.ID
	o.ShelfID.Valid = true

	if o.R == nil {
		o.R = &bookR{
			Shelf: related,
		}
	} else {
		o.R.Shelf = related
	}

	if related.R == nil {
		related.R = &shelfR{
			Books: BookSlice{o},
		}
	} else {
		related.R.Books = append(related.R.Books, o)
	}

	return nil
}

// RemoveShelf relationship.
// Sets o.R.Shelf to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *Book) RemoveShelf(exec boil.Executor, related *Shelf) error {
	var err error

	o.ShelfID.Valid = false
	if err = o.Update(exec, "shelf_id"); err != nil {
		o.ShelfID.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.Shelf = nil
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.Books {
		if o.ShelfID.Int64 != ri.ShelfID.Int64 {
			continue
		}

		ln := len(related.R.Books)
		if ln > 1 && i < ln-1 {
			related.R.Books[i] = related.R.Books[ln-1]
		}
		related.R.Books = related.R.Books[:ln-1]
		break
	}
	return nil
}

// BooksG retrieves all records.
func BooksG(mods ...qm.QueryMod) bookQuery {
	return Books(boil.GetDB(), mods...)
}

// Books retrieves all the records using an executor.
func Books(exec boil.Executor, mods ...qm.QueryMod) bookQuery {
	mods = append(mods, qm.From("`book`"))
	return bookQuery{NewQuery(exec, mods...)}
}

// FindBookG retrieves a single record by ID.
func FindBookG(id int64, selectCols ...string) (*Book, error) {
	return FindBook(boil.GetDB(), id, selectCols...)
}

// FindBookGP retrieves a single record by ID, and panics on error.
func FindBookGP(id int64, selectCols ...string) *Book {
	retobj, err := FindBook(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindBook retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindBook(exec boil.Executor, id int64, selectCols ...string) (*Book, error) {
	bookObj := &Book{}
	bookObj.readonly = &Book{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `book` where `id`=?", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(bookObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from book")
	}

	// copy original obj to the readonly cache
	*bookObj.readonly = *bookObj
	return bookObj, nil
}

// FindBookP retrieves a single record by ID with an executor, and panics on error.
func FindBookP(exec boil.Executor, id int64, selectCols ...string) *Book {
	retobj, err := FindBook(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Book) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Book) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Book) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Book) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no book provided for insertion")
	}
	o.whitelist = whitelist
	o.operation = "INSERT"

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(bookColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	bookInsertCacheMut.RLock()
	cache, cached := bookInsertCache[key]
	bookInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			bookColumns,
			bookColumnsWithDefault,
			bookColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(bookType, bookMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(bookType, bookMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO `book` (`%s`) VALUES (%s)", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `book` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, bookPrimaryKeyColumns))
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
		return errors.Wrap(err, "models: unable to insert into book")
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
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == bookMapping["ID"] {
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
		return errors.Wrap(err, "models: unable to populate default values for book")
	}

CacheNoHooks:
	if !cached {
		bookInsertCacheMut.Lock()
		bookInsertCache[key] = cache
		bookInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Book record. See Update for
// whitelist behavior description.
func (o *Book) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Book record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Book) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Book, and panics on error.
// See Update for whitelist behavior description.
func (o *Book) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Book.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Book) Update(exec boil.Executor, whitelist ...string) error {
	o.whitelist = whitelist
	whitelist = o.Whitelist()

	o.operation = "UPDATE"
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	bookUpdateCacheMut.RLock()
	cache, cached := bookUpdateCache[key]
	bookUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(bookColumns, bookPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update book, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `book` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, bookPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(bookType, bookMapping, append(wl, bookPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update book row")
	}

	if !cached {
		bookUpdateCacheMut.Lock()
		bookUpdateCache[key] = cache
		bookUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q bookQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q bookQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for book")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o BookSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o BookSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o BookSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o BookSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), bookPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE `book` SET %s WHERE (`id`) IN (%s)",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(bookPrimaryKeyColumns), len(colNames)+1, len(bookPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in book slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Book) UpsertG(updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Book) UpsertGP(updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Book) UpsertP(exec boil.Executor, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Book) Upsert(exec boil.Executor, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no book provided for upsert")
	}
	o.whitelist = whitelist
	o.operation = "UPSERT"

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(bookColumnsWithDefault, o)

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

	bookUpsertCacheMut.RLock()
	cache, cached := bookUpsertCache[key]
	bookUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			bookColumns,
			bookColumnsWithDefault,
			bookColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			bookColumns,
			bookPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert book, could not build update column list")
		}

		cache.query = queries.BuildUpsertQueryMySQL(dialect, "book", update, whitelist)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `book` WHERE `id`=?",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
		)

		cache.valueMapping, err = queries.BindMapping(bookType, bookMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(bookType, bookMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert for book")
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
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == bookMapping["ID"] {
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
		return errors.Wrap(err, "models: unable to populate default values for book")
	}

CacheNoHooks:
	if !cached {
		bookUpsertCacheMut.Lock()
		bookUpsertCache[key] = cache
		bookUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single Book record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Book) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Book record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Book) DeleteG() error {
	if o == nil {
		return errors.New("models: no Book provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Book record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Book) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Book record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Book) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Book provided for delete")
	}
	o.operation = "DELETE"

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), bookPrimaryKeyMapping)
	sql := "DELETE FROM `book` WHERE `id`=?"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from book")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q bookQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q bookQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no bookQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from book")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o BookSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o BookSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Book slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o BookSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o BookSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Book slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(bookBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), bookPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM `book` WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, bookPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(bookPrimaryKeyColumns), 1, len(bookPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from book slice")
	}

	if len(bookAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Book) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Book) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Book) ReloadG() error {
	if o == nil {
		return errors.New("models: no Book provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Book) Reload(exec boil.Executor) error {
	ret, err := FindBook(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *BookSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *BookSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *BookSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty BookSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *BookSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	books := BookSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), bookPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT `book`.* FROM `book` WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, bookPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(bookPrimaryKeyColumns), 1, len(bookPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&books)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in BookSlice")
	}

	*o = books

	return nil
}

// BookExists checks if the Book row exists.
func BookExists(exec boil.Executor, id int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from `book` where `id`=? limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if book exists")
	}

	return exists, nil
}

// BookExistsG checks if the Book row exists.
func BookExistsG(id int64) (bool, error) {
	return BookExists(boil.GetDB(), id)
}

// BookExistsGP checks if the Book row exists. Panics on error.
func BookExistsGP(id int64) bool {
	e, err := BookExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// BookExistsP checks if the Book row exists. Panics on error.
func BookExistsP(exec boil.Executor, id int64) bool {
	e, err := BookExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Book) Changes() (ch *Changeset, err error) {
	ch = &Changeset{Table: "book",
		Changes: []*ChangeItem{}, Operation: o.Operation()}
	v := reflect.Indirect(reflect.ValueOf(o.readonly))
	vnew := reflect.Indirect(reflect.ValueOf(o))

	for _, c := range o.Whitelist() {
		if f, ok := BookFieldMapping[c]; ok {
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
func (o *Book) Whitelist() (wl []string) {
	if len(o.whitelist) > 0 {
		return o.whitelist
	}

	// Calculates changed columns as whitelist
	v := reflect.Indirect(reflect.ValueOf(o.readonly))
	vnew := reflect.Indirect(reflect.ValueOf(o))

	for _, c := range bookColumns {
		if f, ok := BookFieldMapping[c]; ok {
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

func (o *Book) Operation() string {
	return o.operation
}

// Generated change history hook for models
func init() {
	chFunc := func(exec boil.Executor, s *Book) error {
		if s == nil || exec == nil {
			return nil
		}

		ch, _ := s.Changes()
		if changeable, ok := exec.(Changeable); ok {
			changeable.AddChange(ch)
		}

		return nil
	}

	afterSel := func(exec boil.Executor, s *Book) error {
		if s == nil || exec == nil {
			return nil
		}

		s.readonly = &Book{}
		*s.readonly = *s
		return nil
	}

	beforeInsert := func(exec boil.Executor, s *Book) error {
		if s == nil || exec == nil {
			return nil
		}

		s.operation = "INSERT"
		return nil
	}

	beforeUpdate := func(exec boil.Executor, s *Book) error {
		if s == nil || exec == nil {
			return nil
		}

		s.operation = "UPDATE"
		return nil
	}

	beforeUpsert := func(exec boil.Executor, s *Book) error {
		if s == nil || exec == nil {
			return nil
		}

		s.operation = "UPSERT"
		return nil
	}

	afterDelete := func(exec boil.Executor, s *Book) error {
		if s == nil || exec == nil {
			return nil
		}

		s.operation = "DELETE"
		return nil
	}

	AddBookHook(boil.AfterSelectHook, afterSel)
	AddBookHook(boil.BeforeInsertHook, beforeInsert)
	AddBookHook(boil.AfterInsertHook, chFunc)
	AddBookHook(boil.BeforeUpdateHook, beforeUpdate)
	AddBookHook(boil.AfterUpdateHook, chFunc)
	AddBookHook(boil.BeforeUpsertHook, beforeUpsert)
	AddBookHook(boil.AfterUpsertHook, chFunc)
	AddBookHook(boil.AfterDeleteHook, afterDelete)
	AddBookHook(boil.AfterDeleteHook, chFunc)
}
