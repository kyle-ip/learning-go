package clause

import (
	"strings"
)

// Clause contains SQL conditions
type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

// Type is the type of Clause
type Type int

// Support types for Clause
const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

// Set adds a sub clause of specific type
// 根据字句类型 Type 调用对应的 generator，生成该子句对应的 SQL 语句及参数，赋值到 Clause 上
// clause.Set(INSERT, "User", []string{"Name", "Age"}) => INSERT INTO ...
func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

// Build generate the final SQL and SQLVars
// 根据传入的 Type 的顺序，构造出最终的 SQL 语句及参数
// INSERT INTO xxx VALUES (...) 和 a, b, c, ...
func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
