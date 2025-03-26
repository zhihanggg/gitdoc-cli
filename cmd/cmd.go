// Package cmd TODO
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/zhihanggg/gitdoc-cli/cmd/create"
	"github.com/zhihanggg/gitdoc-cli/entity/version"
	"github.com/zhihanggg/gitdoc-cli/log"
	"github.com/zhihanggg/gitdoc-cli/utils"
	"gopkg.in/op/go-logging.v1"
)

const configPath = ".gitdoc-cli.yml"

var (
	save       bool
	printTrace bool
)

var rootCmd = &cobra.Command{
	Use:          "gitdoc-cli",
	Short:        "gitdoc-cli",
	Long:         fmt.Sprintf("gitdoc-cli 是 GitDoc 的命令行工具\nVersion: %s\n", version.Version),
	SilenceUsage: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVar(&printTrace, "trace", false, "是否打印 trace 日志, 命令添加 --trace 打印 trace 日志")
	// 如果子命令定义提供了PersistentPreRun函数，那么子命令的PersistentPreRun函数需要主动调用cmd.PersistentPreRun函数
	rootCmd.PersistentPreRun = PersistentPreRun
}

// initConfig 将配置文件中的参数读取到viper中
func initConfig() {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return
	}

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Warn("读取配置文件 %s 失败 %s， 本次执行将忽略该配置文件中的参数", configPath, err.Error())
	} else {
		log.Info("读取配置文件 %s 成功；提示：命令行参数的优先级要高于配置文件中同名参数的优先级", viper.ConfigFileUsed())
	}
}

// PersistentPreRun 各个子命令需要执行的一般操作，为了能让各个子命令都能自动执行该操作，rootCmd.PersistentPreRun被赋值为该函数
func PersistentPreRun(cmd *cobra.Command, _ []string) {
	bindParams(cmd)
	setLogLevel(cmd)
}

// bindParams 将对应命令所能访问的命令行参数(Flag)绑定到viper上，后续可以通过viper访问这些参数；
// 后续若想访问这些参数，需要带上命令的特定前缀；子命令可以在preRun中加载参数
func bindParams(cmd *cobra.Command) {
	if cmd == rootCmd {
		return
	}
	prefix := utils.GetParamPrefix(cmd)
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		// 子命令 跳过 save 这个命令行参数， help Flag 会在 command.persistentPreRun 前自动引入
		if flag.Name == "save" || flag.Name == "help" {
			return
		}
		_ = viper.BindPFlag(prefix+flag.Name, flag)
	})
}

func setLogLevel(cmd *cobra.Command) {
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05} %{shortfunc} [%{level:.4s}]%{color:reset} %{message}`,
	)
	var backend = logging.AddModuleLevel(
		logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), format))
	prefix := utils.GetParamPrefix(cmd)
	if viper.GetBool(prefix + "trace") {
		backend.SetLevel(logging.DEBUG, "")
		p := log.DefaultStd.WithEnableTrace()
		log.DefaultStd = &p
	} else {
		backend.SetLevel(logging.WARNING, "")
	}
	logging.SetBackend(backend)
}

// Execute 为具体脚手架命令的执行
func Execute() {

	rootCmd.AddCommand(create.NewCmd())

	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}
