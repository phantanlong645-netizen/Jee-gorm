package session

import (
	"Day3-save-query/clause"
	"database/sql"
	"strings"

	"Day3-save-query/dialect"
	"Day3-save-query/schema"
)

type Session struct {
	db       *sql.DB
	dialect  dialect.Dialect
	refTable *schema.Schema
	sql      strings.Builder
	sqlVars  []interface{}
	clause   clause.Clause
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
}

func (s *Session) DB() *sql.DB {
	return s.db
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
