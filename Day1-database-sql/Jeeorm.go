package Day1_database_sql

import (
	"Day1-database-sql/log"
	"Day1-database-sql/session"
	"database/sql"

	_ "modernc.org/sqlite"
)

type Engine struct {
	db *sql.DB
}

func NewEngine(driver, source string) (e *Engine, err error) {
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
		db: db,
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
	return session.NewSession(e.db)
}
