package Day6_transaction

import (
	"database/sql"
	"fmt"

	"Day6-transaction/dialect"
	"Day6-transaction/log"
	"Day6-transaction/session"

	_ "modernc.org/sqlite"
)

func normalizeDriver(driver string) string {
	if driver == "sqlite3" {
		return "sqlite"
	}
	return driver
}

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine, err error) {
	d, ok := dialect.GetDialect(driver)
	if !ok {
		err = fmt.Errorf("dialect %s not found", driver)
		log.Error(err)
		return
	}
	db, err := sql.Open(normalizeDriver(driver), source)
	if err != nil {
		log.Error(err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	e = &Engine{
		db:      db,
		dialect: d,
	}
	log.Info("Connect database success")
	return
}
func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error(err)
	}
	log.Info("Close database success")
}
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

type TxFunc func(session *session.Session) (interface{}, error)

func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			s.Commit()
		}
	}()
	return f(s)
}
