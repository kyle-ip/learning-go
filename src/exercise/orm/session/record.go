package session

import (
	"errors"
	"geeorm/clause"
	"reflect"
)

// ORM 主要 API，如：
// 		s.Where("Name = ?", "Tom").Delete()
//		s.Where("Name = ?", "Tom").Update("Age", 30)

// Insert one or more records in database
// 插入一或多条记录到数据库中。
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)

	// 多次调用 clause.Set() 构造好每一个子句。
	for _, value := range values {
		s.CallMethod(BeforeInsert, value)
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}
	s.clause.Set(clause.VALUES, recordValues...)

	// 调用 clause.Build() 按照传入的顺序构造出最终的 SQL 语句。
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)

	// 执行构造好的 SQL 语句。
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	// 执行钩子函数。
	s.CallMethod(AfterInsert, nil)

	// 返回结果集。
	return result.RowsAffected()
}

// Find gets all eligible records
// Find 传入一个切片指针（“泛型”），把查询的结果保存在切片中。
func (s *Session) Find(values interface{}) error {

	// 调用钩子函数。
	s.CallMethod(BeforeQuery, nil)

	// ========== 构建 Schema ==========

	// 获取指针指向的切片，以及切片元素的类型。
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()

	// 创建 destType 实例，传入 Model，调用 RefTable() 映射出表结构 Schema。
	modelValue := reflect.New(destType).Elem().Interface()
	schema := s.Model(modelValue).RefTable()

	// ========== 构造 SELECT 语句并执行 ==========

	// 根据表结构，使用 clause 构造出 SELECT 语句，查询到所有符合条件的记录 rows。
	s.clause.Set(clause.SELECT, schema.Name, schema.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	// ========== 采集查询返回的数据 ==========

	// 遍历每一行记录，利用反射创建 destType 的实例 dest，将 dest 的所有字段平铺开，构造切片 values。
	for rows.Next() {

		dest := reflect.New(destType).Elem()

		// 存放一行记录（的所有字段的值）。
		var values []interface{}
		for _, name := range schema.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}

		// 将该行记录每一列的值依次赋值给 values 中的每一个字段。
		if err := rows.Scan(values...); err != nil {
			return err
		}

		// 调用钩子函数。
		s.CallMethod(AfterQuery, dest.Addr().Interface())

		// 将 dest 添加到切片 destSlice 中。循环直到所有的记录都添加到切片 destSlice 中。
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

// First gets the 1st row
// 获取第一条记录：
//      u := &User{}
//      _ = s.OrderBy("Age DESC").First(u)
func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}

// Limit adds limit condition to clause
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

// Where adds limit condition to clause
func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

// OrderBy adds order by condition to clause
func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

// Update records with where clause
// support map[string]interface{}
// also support kv list: "Name", "Tom", "Age", 18, ....
// Update 接受 2 种入参，平铺开来的键值对和 map 类型的键值对。
func (s *Session) Update(kv ...interface{}) (int64, error) {

	// 执行钩子函数。
	s.CallMethod(BeforeUpdate, nil)

	// 如果 kv 不是 map 类型的参数，则遍历、转换为 map。
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterUpdate, nil)
	return result.RowsAffected()
}

// Delete records with where clause
func (s *Session) Delete() (int64, error) {
	s.CallMethod(BeforeDelete, nil)
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterDelete, nil)
	return result.RowsAffected()
}

// Count records with where clause
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}
