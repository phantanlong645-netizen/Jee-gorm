package session

import (
	"Day2_reflect_schema/log"
	"Day2_reflect_schema/schema"
	"fmt"
	"reflect"
	"strings"
)

func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) Reftable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.Reftable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw("Create table %s (%s)", table.Name, desc).Exec()
	return err
}
func (s *Session) DropTable() error {
	table := s.Reftable()
	_, err := s.Raw("DROP TABLE IF EXIST %s", table.Name).Exec()
	return err
}

func (s *Session) HasTable() bool {
	table := s.Reftable()
	stmt, args := s.dialect.TableExist(table.Name)
	result := s.Raw(stmt, args...).QueryRow()
	log.Info("result:", result)
	var tmp string
	_ = result.Scan(&tmp)
	return tmp == table.Name
}
