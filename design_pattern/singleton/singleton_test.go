package main

import (
	"sync"
)

type singleton struct {
}

var ins *singleton

var once sync.Once

func GetInsOr() *singleton {
	once.Do(func() {
		ins = &singleton{}
	})
	return ins
}
