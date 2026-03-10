package session

import (
	"database/sql"
	"fmt"
	"strings"
)

type Session struct {
	db      *sql.DB
	sql     strings.Builder
	sqlVars []interface{}
}

func NewSession(db *sql.DB) *Session {
	return &Session{db: db}
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

	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		fmt.Sprintf("exec error")
	}
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
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		fmt.Sprintf("exec error")
	}
	return
}

func (s *Session) Raw(sql string, args ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, args...)
	return s
}
