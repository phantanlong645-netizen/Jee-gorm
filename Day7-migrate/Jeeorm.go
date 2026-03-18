package Day7_migrate

import (
	"database/sql"
	"fmt"
	"strings"

	"Day7-migrate/dialect"
	"Day7-migrate/log"
	"Day7-migrate/session"

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
func difference(a, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}
	for _, v := range a {
		if !mapB[v] {
			diff = append(diff, v)
		}
	}
	return
}
func (engine *Engine) Migrate(value interface{}) error {
	_, err := engine.Transaction(func(session *session.Session) (result interface{}, err error) {
		if !session.Model(value).HasTable() {
			log.Infof("table %s doesn't exist", session.Reftable().Name)
			return nil, session.CreateTable()
		}
		table := session.Reftable()
		rows, err := session.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		if err != nil {
			return nil, err
		}
		columns, err := rows.Columns()
		_ = rows.Close()
		if err != nil {
			return nil, err
		}
		addcolumns := difference(table.FieldNames, columns)
		delcolumns := difference(columns, table.FieldNames)
		log.Infof("added cols %v, deleted cols %v", addcolumns, delcolumns)
		for _, col := range addcolumns {
			f := table.GetField(col)
			sqlstr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table.Name, f.Name, f.Type)
			_, err = session.Raw(sqlstr).Exec()
			if err != nil {
				return nil, err
			}
		}

		if len(delcolumns) == 0 {
			return nil, nil
		}

		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ",")
		var columnDefs []string
		for _, field := range table.Fields {
			columnDefs = append(columnDefs, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
		}
		if _, err = session.Raw(fmt.Sprintf("CREATE TABLE %s (%s)", tmp, strings.Join(columnDefs, ","))).Exec(); err != nil {
			return nil, err
		}
		if _, err = session.Raw(fmt.Sprintf("INSERT INTO %s (%s) SELECT %s FROM %s", tmp, fieldStr, fieldStr, table.Name)).Exec(); err != nil {
			return nil, err
		}
		if _, err = session.Raw(fmt.Sprintf("DROP TABLE %s", table.Name)).Exec(); err != nil {
			return nil, err
		}
		if _, err = session.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", tmp, table.Name)).Exec(); err != nil {
			return nil, err
		}
		return nil, nil

	})
	return err
}
