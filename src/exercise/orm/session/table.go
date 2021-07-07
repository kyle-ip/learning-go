package session

import (
    "fmt"
    "geeorm/log"
    "reflect"
    "strings"

    "geeorm/schema"
)

// Model assigns refTable
// 为 refTable 赋值
func (s *Session) Model(value interface{}) *Session {
    // nil or different model, update refTable
    // 解析操作比较耗时，将解析的结果保存在成员变量 refTable 中，
    // 即使 Model() 被调用多次，如果传入的结构体名称不发生变化，则不会更新 refTable 的值。
    if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
        s.refTable = schema.Parse(value, s.dialect)
    }
    return s
}

// RefTable returns a Schema instance that contains all parsed fields
// 返回 refTable 的值
func (s *Session) RefTable() *schema.Schema {

    // 如果 refTable 未被赋值，则打印错误日志。
    if s.refTable == nil {
        log.Error("Model is not set")
    }

    // 返回的数据库表和字段的信息，拼接出 SQL 语句，调用原生 SQL 接口执行。
    return s.refTable
}

// CreateTable create a table in database with a model
// 创建数据库表
func (s *Session) CreateTable() error {
    table := s.RefTable()
    var columns []string
    for _, field := range table.Fields {
        columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
    }
    desc := strings.Join(columns, ",")
    _, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
    return err
}

// DropTable drops a table with the name of model
// 删除数据库表
func (s *Session) DropTable() error {
    _, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
    return err
}

// HasTable returns true of the table exists
// 判断表是否存在
func (s *Session) HasTable() bool {
    sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
    row := s.Raw(sql, values...).QueryRow()
    var tmp string
    _ = row.Scan(&tmp)
    return tmp == s.RefTable().Name
}
