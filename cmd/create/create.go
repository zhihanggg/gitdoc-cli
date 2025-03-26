package create

import (
	"github.com/spf13/cobra"
	"github.com/zhihanggg/gitdoc-cli/log"
)

var (
	moduleName string
)

// NewCmd 返回 create 相关子命令
func NewCmd() *cobra.Command {
	impl := createImpl{}
	createCmd := &cobra.Command{
		Use:   "create --module-name=<your_module_name>",
		Short: "create 命令用来生成代码相关，使用方式可以执行 fit-cli help create",
		Long:  "create 命令用来生成代码相关，更多说明可以参考文档 https://git.woa.com/Fit-scaffold/fit-cli/tree/master/docs/create",
		RunE:  impl.run(),
	}
	createCmd.PersistentFlags().StringVar(&moduleName, "module-name", "", "模块英文名")
	return createCmd
}

type createImpl struct {
}

func (c *createImpl) run() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		log.Info("create module: %s", moduleName)
		return nil
	}
}
