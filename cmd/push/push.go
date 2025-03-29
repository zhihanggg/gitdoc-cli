package push

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhihanggg/gitdoc-cli/log"
	"github.com/zhihanggg/gitdoc-cli/utils"
)

func NewCmd() *cobra.Command {
	impl := pushImpl{}
	return &cobra.Command{
		Use:   "push",
		Short: "push 命令用来推送变更到远端",
		Long:  "push 命令用来推送变更到远端",
		RunE:  impl.run(),
	}
}

type pushImpl struct {
}

func (i *pushImpl) run() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		log.Debug("开始执行 git push...")
		output, err := utils.ExecCmd("git push")
		if err != nil {
			return fmt.Errorf("git push 失败: %v", err)
		}

		log.Info("git push 成功执行")
		log.Debug("git push 输出: %s", output)
		return nil
	}
}
