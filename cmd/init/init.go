package init

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	impl := initImpl{}
	return &cobra.Command{
		Use:   "init",
		Short: "init 命令用来初始化环境,安装一些依赖",
		Long:  "init 命令用来初始化环境,安装一些依赖",
		RunE:  impl.run(),
	}
}

type initImpl struct {
}

func (i *initImpl) run() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}
}
