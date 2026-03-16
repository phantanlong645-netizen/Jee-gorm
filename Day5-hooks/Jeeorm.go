package Day5_hooks

import (
	"database/sql"
	"fmt"

	"Day5-hooks/dialect"
	"Day5-hooks/log"
	"Day5-hooks/session"

	_ "modernc.org/sqlite"
)

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
	db, err := sql.Open(driver, source)
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
