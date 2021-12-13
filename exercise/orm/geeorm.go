package geeorm

import (
	"database/sql"
	"fmt"
	"geeorm/dialect"
	"geeorm/log"
	"geeorm/session"
	"strings"
)

// Engine is the main struct of geeorm, manages all db sessions and transactions.
type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// NewEngine create a instance of Engine
// connect database and ping it to test whether it's alive
func NewEngine(driver, source string) (e *Engine, err error) {

	// 连接数据库：输入驱动名称、数据库名称，返回 sql.DB 实例指针。
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}

	// Send a ping to make sure the database connection is alive.
	// 测试数据库连接是否正常可用。
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	// make sure the specific dialect exists
	// 确保 driver 特定的 Dialect 的实现存在。
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}
	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

// Close database connection
// 关闭数据库连接。
func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
		return
	}
	log.Info("Close database success")
}

// NewSession creates a new session for next operations
func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}

// TxFunc will be called between tx.Begin() and tx.Commit()
// https://stackoverflow.com/questions/16184238/database-sql-tx-detecting-commit-or-rollback
type TxFunc func(*session.Session) (interface{}, error)

// Transaction executes sql wrapped in a transaction, then automatically commit if no error occurs
// 事务管理：将所有的操作放到回调函数中，作为入参传递给 engine.Transaction()，发生任何错误，自动回滚，如果没有错误发生，则提交。
func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = s.Rollback() // err is non-nil; don't change it
		} else {
			err = s.Commit() // err is nil; if Commit returns error update err
		}
	}()
	return f(s)
}

// difference returns a - b
// 计算前后两个字段切片的差集。新表 - 旧表 = 新增字段，旧表 - 新表 = 删除字段。
func difference(a []string, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}
	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

// Migrate table
func (engine *Engine) Migrate(value interface{}) error {
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		if !s.Model(value).HasTable() {
			log.Infof("table %s doesn't exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}
		table := s.RefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("added cols %v, deleted cols %v", addCols, delCols)

		// 使用 ALTER 语句新增字段。
		for _, col := range addCols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, f.Name, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		if len(delCols) == 0 {
			return
		}

		// 使用创建新表并以重命名的方式删除字段。
		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
		_, err = s.Exec()
		return
	})
	return err
}
