package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
)

const (
	DefaultSelectLimit = 200
)

func GetLimit(n uint64) uint64 {
	if n < 0 {
		n = 0
	}
	if n > DefaultSelectLimit {
		n = DefaultSelectLimit
	}
	return n
}

func GetOffset(n uint64) uint64 {
	if n < 0 {
		n = 0
	}
	return n
}

type Database struct {
	*dbr.Session
}

type SelectQuery struct {
	*dbr.SelectBuilder
}

type InsertQuery struct {
	*dbr.InsertBuilder
}

type DeleteQuery struct {
	*dbr.DeleteBuilder
}

type UpdateQuery struct {
	*dbr.UpdateBuilder
}

// SelectQuery
// Example: Select().From().Where().Limit().Offset().OrderDir().Load()
//          Select().From().Where().Limit().Offset().OrderDir().LoadOne()
//          Select().From().Where().Count()

func (db *Database) Select(columns ...string) *SelectQuery {
	return &SelectQuery{db.Session.Select(columns...)}
}

func (b *SelectQuery) From(table string) *SelectQuery {
	b.SelectBuilder.From(table)
	return b
}

func (b *SelectQuery) Where(query interface{}, value ...interface{}) *SelectQuery {
	b.SelectBuilder.Where(query, value...)
	return b
}

func (b *SelectQuery) Limit(n uint64) *SelectQuery {
	n = GetLimit(n)
	b.SelectBuilder.Limit(n)
	return b
}

func (b *SelectQuery) Offset(n uint64) *SelectQuery {
	n = GetLimit(n)
	b.SelectBuilder.Offset(n)
	return b
}

func (b *SelectQuery) OrderDir(col string, isAsc bool) *SelectQuery {
	b.SelectBuilder.OrderDir(col, isAsc)
	return b
}

func (b *SelectQuery) Load(value interface{}) (int, error) {
	return b.SelectBuilder.Load(value)
}

func (b *SelectQuery) LoadOne(value interface{}) error {
	return b.SelectBuilder.LoadOne(value)
}

func (b *SelectQuery) Count() (count uint32, err error) {
	// cache SelectStmt
	selectStmt := b.SelectBuilder.SelectStmt
	b.SelectBuilder.SelectStmt = dbr.Select("count(*)").From(selectStmt.Table)
	err = b.SelectBuilder.LoadOne(&count)
	// fallback SelectStmt
	b.SelectBuilder.SelectStmt = selectStmt
	return
}

// InsertQuery
// Example: InsertInto().Columns().Record().Exec()

func (db *Database) InsertInto(table string) *InsertQuery {
	return &InsertQuery{db.Session.InsertInto(table)}
}

func (b *InsertQuery) Exec() (sql.Result, error) {
	return b.InsertBuilder.Exec()
}

func (b *InsertQuery) Columns(columns ...string) *InsertQuery {
	b.InsertBuilder.Columns(columns...)
	return b
}

func (b *InsertQuery) Record(structValue interface{}) *InsertQuery {
	b.InsertBuilder.Record(structValue)
	return b
}

// DeleteQuery
// Example: DeleteFrom().Where().Limit().Exec()

func (db *Database) DeleteFrom(table string) *DeleteQuery {
	return &DeleteQuery{db.Session.DeleteFrom(table)}
}

func (b *DeleteQuery) Where(query interface{}, value ...interface{}) *DeleteQuery {
	b.DeleteBuilder.Where(query, value...)
	return b
}

func (b *DeleteQuery) Limit(n uint64) *DeleteQuery {
	b.DeleteBuilder.Limit(n)
	return b
}

func (b *DeleteQuery) Exec() (sql.Result, error) {
	return b.DeleteBuilder.Exec()
}

// UpdateQuery
// Example: Update().Set().Where().Exec()

func (db *Database) Update(table string) *UpdateQuery {
	return &UpdateQuery{db.Session.Update(table)}
}

func (b *UpdateQuery) Exec() (sql.Result, error) {
	return b.UpdateBuilder.Exec()
}

func (b *UpdateQuery) Set(column string, value interface{}) *UpdateQuery {
	b.UpdateBuilder.Set(column, value)
	return b
}

func (b *UpdateQuery) SetMap(m map[string]interface{}) *UpdateQuery {
	b.UpdateBuilder.SetMap(m)
	return b
}

func (b *UpdateQuery) Where(query interface{}, value ...interface{}) *UpdateQuery {
	b.UpdateBuilder.Where(query, value...)
	return b
}

func (b *UpdateQuery) Limit(n uint64) *UpdateQuery {
	b.UpdateBuilder.Limit(n)
	return b
}
