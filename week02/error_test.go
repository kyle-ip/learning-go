package week02

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

func dao() error {
	return errors.Wrap(sql.ErrNoRows, "dao exec failed")
}

func service() error {
	return errors.WithMessage(dao(), "service exec failed")
}

// TestWrap
// 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
func TestWrap(t *testing.T) {
	err := service()
	//if errors.Cause(err) == sql.ErrNoRows {
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("%+v\n", err)
		return
	}
}
