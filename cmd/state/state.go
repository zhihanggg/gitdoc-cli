package state

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhihanggg/gitdoc-cli/log"
	"github.com/zhihanggg/gitdoc-cli/utils"
)

func NewCmd() *cobra.Command {
	impl := stateImpl{}
	return &cobra.Command{
		Use:   "state",
		Short: "state 命令用来查看当前项目的状态信息",
		Long:  "state 命令用来查看当前项目的状态信息",
		RunE:  impl.run(),
	}
}

type stateImpl struct {
}

func (i *stateImpl) run() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		log.Debug("开始获取项目状态信息...")

		// 获取项目在git上的链接
		remoteURL, err := utils.ExecCmd("git config --get remote.origin.url")
		if err != nil {
			return fmt.Errorf("获取远程仓库URL失败: %v", err)
		}
		log.Info("远程仓库URL: %s", remoteURL)

		// 检查是否有尚未add的修改
		statusOutput, err := utils.ExecCmd("git status --porcelain")
		if err != nil {
			return fmt.Errorf("获取git状态失败: %v", err)
		}
		if statusOutput == "" {
			log.Info("没有尚未add的修改")
		} else {
			log.Info("有尚未add的修改:\n%s", statusOutput)
		}

		// 获取本地分支
		branchOutput, err := utils.ExecCmd("git branch --show-current")
		if err != nil {
			return fmt.Errorf("获取本地分支失败: %v", err)
		}
		log.Info("当前本地分支: %s", branchOutput)

		return nil
	}
}
