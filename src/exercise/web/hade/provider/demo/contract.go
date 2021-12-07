package demo

const Key = "hade:demo"

type Service interface {
	GetFoo() Foo
}

type Foo struct {
	Name string
}
