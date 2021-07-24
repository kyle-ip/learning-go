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

func TestWrap(t *testing.T) {
	err := service()
	//if errors.Cause(err) == sql.ErrNoRows {
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("%+v\n", err)
		return
	}
}
