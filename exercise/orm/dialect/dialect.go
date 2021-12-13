package dialect

import "reflect"

// Dialect is an interface contains methods that a dialect has to implement
// 接口用于将多种数据库差异部分提取出来（对每一种数据库分别实现，可实现最大程度的复用和解耦）
type Dialect interface {
	// DataTypeOf 将 Go 语言的类型转换为该数据库的数据类型
	DataTypeOf(typ reflect.Value) string

	// TableExistSQL 传入表名返回某个表是否存在的 SQL 语句
	TableExistSQL(tableName string) (string, []interface{})
}

// dialectsMap 存放 Dialect 名称与其实例的映射
var dialectsMap = map[string]Dialect{}

// RegisterDialect register a dialect to the global variable
// 注册 Dialect 实例
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// GetDialect Get the dialect from global variable if it exists
// 获取 Dialect 实例
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
