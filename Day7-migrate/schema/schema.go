package schema

import (
	"Day7-migrate/dialect"
	"go/ast"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	FieldMap   map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.FieldMap[name]
}
func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldsvalues []interface{}
	for _, field := range s.Fields {
		fieldsvalues = append(fieldsvalues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldsvalues
}

type ITableName interface {
	TableName() string
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	var tableName string
	t, ok := dest.(ITableName)
	if ok {
		tableName = t.TableName()
	} else {
		tableName = modelType.Name()
	}
	schema := &Schema{
		Model:    dest,
		Name:     tableName,
		FieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		f := modelType.Field(i)
		if !f.Anonymous && ast.IsExported(f.Name) {
			F := &Field{
				Name: f.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(f.Type))),
			}
			if v, ok := f.Tag.Lookup("geeorm"); ok {
				F.Tag = v
			}
			schema.Fields = append(schema.Fields, F)
			schema.FieldNames = append(schema.FieldNames, f.Name)
			schema.FieldMap[f.Name] = F
		}
	}
	return schema
}
