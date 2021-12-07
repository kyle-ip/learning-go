package console

import (
	"github.com/yipwinghong/hade/app/console/command/demo"
	"github.com/yipwinghong/hade/framework"
	"github.com/yipwinghong/hade/framework/cobra"
	"github.com/yipwinghong/hade/framework/command"
)

// RunCommand  初始化根Command并运行
func RunCommand(container framework.Container) error {
	// 根Command
	var rootCmd = &cobra.Command{
		// 定义根命令的关键字。
		Use: "hade",
		// 简短介绍。
		Short: "hade 命令",
		// 详细介绍。
		Long: "hade 框架提供的命令行工具，使用这个命令行工具能很方便执行框架自带命令，也能很方便编写业务命令。",
		// 根命令的执行函数
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpFlag()
			return cmd.Help()
		},
		// 不需要出现cobra默认的completion子命令
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	// 为根 Command 设置服务容器。
	rootCmd.SetContainer(container)

	// 绑定框架的命令。
	command.AddKernelCommands(rootCmd)

	// 绑定业务的命令。
	AddAppCommand(rootCmd)

	// 执行 RootCommand。
	return rootCmd.Execute()
}

// AddAppCommand 绑定业务的命令。
func AddAppCommand(rootCmd *cobra.Command) {
	//  demo 例子
	rootCmd.AddCommand(demo.InitFoo())
}