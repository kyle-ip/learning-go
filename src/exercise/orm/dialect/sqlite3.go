package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type sqlite3 struct{}

// 匿名变量确保 sqlite3 实现 Dialect 接口所有方法
var _ Dialect = (*sqlite3)(nil)

// init 第一次加载时，将 sqlite3 实例注册到全局
func init() {
	RegisterDialect("sqlite3", &sqlite3{})
}

// DataTypeOf Get Data Type for sqlite3 Dialect
// 传入 Go 类型，返回 SQLite 类型的字符串
func (s *sqlite3) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Uint,
		reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

// TableExistSQL returns SQL that judge whether the table exists in database
func (s *sqlite3) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", args
}
