package demo

import (
	"github.com/yipwinghong/hade/framework/cobra"
	"log"
)

// InitFoo 初始化Foo命令
func InitFoo() *cobra.Command {
	FooCommand.AddCommand(Foo1Command)
	return FooCommand
}

// FooCommand 代表Foo命令
var FooCommand = &cobra.Command{
	Use:     "foo",
	Short:   "foo 的简要说明",
	Long:    "foo 的长说明",
	Aliases: []string{"fo", "f"},
	Example: "foo 命令的例子",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		log.Println(container)
		return nil
	},
}

// Foo1Command 代表Foo命令的子命令Foo1
var Foo1Command = &cobra.Command{
	Use:     "foo1",
	Short:   "foo1 的简要说明",
	Long:    "foo1 的长说明",
	Aliases: []string{"fo1", "f1"},
	Example: "foo1 命令的例子",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		log.Println(container)
		return nil
	},
}
