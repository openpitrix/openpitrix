package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/gocraft/dbr/dialect"
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
	JoinCount int // for join filter
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

type UpsertQuery struct {
	table string
	*dbr.Session
	whereConds   map[string]string
	upsertValues map[string]interface{}
}

// SelectQuery
// Example: Select().From().Where().Limit().Offset().OrderDir().Load()
//          Select().From().Where().Limit().Offset().OrderDir().LoadOne()
//          Select().From().Where().Count()
//          SelectAll().From().Where().Limit().Offset().OrderDir().Load()
//          SelectAll().From().Where().Limit().Offset().OrderDir().LoadOne()
//          SelectAll().From().Where().Count()

func (db *Database) Select(columns ...string) *SelectQuery {
	return &SelectQuery{db.Session.Select(columns...), 0}
}

func (db *Database) SelectBySql(query string, value ...interface{}) *SelectQuery {
	return &SelectQuery{db.Session.SelectBySql(query, value...), 0}
}

func (db *Database) SelectAll(columns ...string) *SelectQuery {
	return &SelectQuery{db.Session.Select("*"), 0}
}

func (b *SelectQuery) Join(table, on interface{}) *SelectQuery {
	b.SelectBuilder.Join(table, on)
	return b
}

func (b *SelectQuery) JoinAs(table string, alias string, on interface{}) *SelectQuery {
	b.SelectBuilder.Join(dbr.I(table).As(alias), on)
	return b
}

func (b *SelectQuery) From(table string) *SelectQuery {
	b.SelectBuilder.From(table)
	return b
}

func (b *SelectQuery) Where(query interface{}, value ...interface{}) *SelectQuery {
	b.SelectBuilder.Where(query, value...)
	return b
}

func (b *SelectQuery) GroupBy(col ...string) *SelectQuery {
	b.SelectBuilder.GroupBy(col...)
	return b
}

func (b *SelectQuery) Distinct() *SelectQuery {
	b.SelectBuilder.Distinct()
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

func getColumns(dbrColumns []interface{}) string {
	var columns []string
	for _, column := range dbrColumns {
		if c, ok := column.(string); ok {
			columns = append(columns, c)
		}
	}
	return strings.Join(columns, ", ")
}

func (b *SelectQuery) Count() (count uint32, err error) {
	// cache SelectStmt
	selectStmt := b.SelectStmt

	limit := selectStmt.LimitCount
	offset := selectStmt.OffsetCount
	column := selectStmt.Column
	isDistinct := selectStmt.IsDistinct
	order := selectStmt.Order

	b.SelectStmt.LimitCount = -1
	b.SelectStmt.OffsetCount = -1
	b.SelectStmt.Column = []interface{}{"COUNT(*)"}
	b.SelectStmt.Order = []dbr.Builder{}

	if isDistinct {
		b.SelectStmt.Column = []interface{}{fmt.Sprintf("COUNT(DISTINCT %s)", getColumns(column))}
		b.SelectStmt.IsDistinct = false
	}

	err = b.LoadOne(&count)
	// fallback SelectStmt
	selectStmt.LimitCount = limit
	selectStmt.OffsetCount = offset
	selectStmt.Column = column
	selectStmt.IsDistinct = isDistinct
	selectStmt.Order = order
	b.SelectStmt = selectStmt
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

// UpsertQuery
// Example: Upsert().Where().Set().Exec()
//          Upsert().WhereMap().SetMap().Exec()

func (db *Database) Upsert(table string) *UpsertQuery {
	return &UpsertQuery{
		table:        table,
		Session:      db.Session,
		whereConds:   make(map[string]string),
		upsertValues: make(map[string]interface{}),
	}
}

func (b *UpsertQuery) Where(column, value string) *UpsertQuery {
	b.whereConds[column] = value
	return b
}

func (b *UpsertQuery) WhereMap(m map[string]string) *UpsertQuery {
	for column, value := range m {
		b.Where(column, value)
	}
	return b
}

func (b *UpsertQuery) Set(column string, value interface{}) *UpsertQuery {
	b.upsertValues[column] = value
	return b
}

func (b *UpsertQuery) SetMap(m map[string]interface{}) *UpsertQuery {
	for column, value := range m {
		b.Set(column, value)
	}
	return b
}

func (b *UpsertQuery) Exec() (sql.Result, error) {
	var columns []string
	var values []interface{}
	for column, value := range b.whereConds {
		columns = append(columns, column)
		values = append(values, value)
	}
	for column, value := range b.upsertValues {
		columns = append(columns, column)
		values = append(values, value)
	}
	d := dialect.MySQL
	buf := dbr.NewBuffer()
	stmt := dbr.InsertInto(b.table).Columns(columns...).Values(values...)
	stmt.Build(d, buf)

	if len(b.upsertValues) > 0 {
		buf.WriteString(" ON DUPLICATE KEY UPDATE ")

		i := 0
		for col, v := range b.upsertValues {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(d.QuoteIdent(col))
			buf.WriteString(" = ")
			buf.WriteString("?")

			buf.WriteValue(v)
			i++
		}
	}
	query, err := dbr.InterpolateForDialect(buf.String(), buf.Value(), d)
	if err != nil {
		return nil, err
	}

	return b.InsertBySql(query).Exec()
}
