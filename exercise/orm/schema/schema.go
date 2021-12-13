package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

// Field represents a column of database
type Field struct {

	// Name 字段名
	Name string

	// Type 类型
	Type string

	// Tag 约束条件
	Tag string
}

// Schema represents a table of database
// 表元信息
type Schema struct {
	// Model 被映射的对象
	Model interface{}

	// Name 表名
	Name string

	// Fields 字段表
	Fields []*Field

	// Fields 字段名表
	FieldNames []string

	// fieldMap 字段映射表
	fieldMap map[string]*Field
}

// GetField returns field by name
func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// RecordValues return the values of dest's member variables
// 转换参数格式：
//      u1 := &User{Name: "Tom", Age: 18} => ("Tom", 18)
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}

// ITableName 表示表名
type ITableName interface {
	TableName() string
}

// Parse a struct to a Schema instance
// 将任意 GO 对象解析为 Schema 实例
func Parse(dest interface{}, d dialect.Dialect) *Schema {

	// 获取映射对象的类型 Type：TypeOf 类型，ValueOf 值，Indirect 指向的实例。
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()

	// 获取表名：如果 dest 是 ITableName 类型，则从 TableName 方法中获取。否则从 modelType 中反射获取。
	var tableName string
	t, ok := dest.(ITableName)
	if !ok {
		tableName = modelType.Name()
	} else {
		tableName = t.TableName()
	}

	schema := &Schema{
		Model:    dest,
		Name:     tableName,
		fieldMap: make(map[string]*Field),
	}

	// 获取实例的字段的个数。
	for i := 0; i < modelType.NumField(); i++ {

		// 通过下标获取到特定字段。
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				// 取字段名，字段类型（取 Go 类型，转换为 SQLite 类型）。
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			// 取约束条件。
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}
