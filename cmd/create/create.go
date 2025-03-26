package create

import (
	"github.com/spf13/cobra"
	"github.com/zhihanggg/gitdoc-cli/log"
)

var (
	projectName string
)

// NewCmd 返回 create 相关子命令
func NewCmd() *cobra.Command {
	impl := createImpl{}
	createCmd := &cobra.Command{
		Use:   "create --project-name=<your_project_name>",
		Short: "create 命令用来创建新的项目, 使用方式可以执行 gitdoc-cli help create",
		Long:  "create 命令用来创建新的项目, 更多说明可以参考文档",
		RunE:  impl.run(),
	}
	createCmd.PersistentFlags().StringVar(&projectName, "project-name", "", "项目英文名")
	return createCmd
}

type createImpl struct {
}

func (c *createImpl) run() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		log.Info("create project: %s", projectName)
		return nil
	}
}
