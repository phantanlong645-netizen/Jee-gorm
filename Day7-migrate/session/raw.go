package session

import (
	"Day7-migrate/clause"
	"database/sql"
	"strings"

	"Day7-migrate/dialect"
	"Day7-migrate/schema"
)

type Session struct {
	db       *sql.DB
	dialect  dialect.Dialect
	refTable *schema.Schema
	sql      strings.Builder
	sqlVars  []interface{}
	clause   clause.Clause
	tx       *sql.Tx
}

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	result, err = s.DB().Exec(s.sql.String(), s.sqlVars...)
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	var rows *sql.Row
	rows = s.DB().QueryRow(s.sql.String(), s.sqlVars...)
	return rows
}
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	rows, err = s.DB().Query(s.sql.String(), s.sqlVars...)
	return
}

func (s *Session) Raw(sql string, args ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, args...)
	return s
}
