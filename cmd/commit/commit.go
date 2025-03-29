package commit

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zhihanggg/gitdoc-cli/log"
	"github.com/zhihanggg/gitdoc-cli/utils"
)

func NewCmd() *cobra.Command {
	impl := commitImpl{}
	return &cobra.Command{
		Use:   "commit",
		Short: "commit 命令用来提交变更到远端",
		Long:  "commit 命令用来提交变更到远端，会自动将doc/docx文件转换为markdown",
		RunE:  impl.run(),
	}
}

type commitImpl struct {
}

func (i *commitImpl) run() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// 检查是否支持Windows
		if err := utils.CmdNotSupportWindowsIntercept("commit"); err != nil {
			return err
		}

		// 扫描doc/docx文件
		log.Debug("开始扫描文档文件...")
		docFiles, err := utils.ScanFilesByExt(".", []string{".doc", ".docx"})
		if err != nil {
			return fmt.Errorf("扫描文档文件失败: %v", err)
		}

		if len(docFiles) == 0 {
			log.Debug("未找到需要转换的文档文件")
		} else {
			log.Debug("找到 %d 个文档文件，开始转换...", len(docFiles))

			// 转换所有文档为markdown
			for _, docFile := range docFiles {
				mdFile := strings.TrimSuffix(docFile, filepath.Ext(docFile)) + ".md"
				log.Debug("正在转换: %s -> %s", docFile, mdFile)

				if err := utils.ConvertDocToMarkdown(docFile, mdFile); err != nil {
					return fmt.Errorf("转换文件 %s 失败: %v", docFile, err)
				}
			}
		}

		// 执行git add --all
		log.Debug("执行 git add --all...")
		if _, err := utils.ExecCmd("git add --all"); err != nil {
			return fmt.Errorf("git add 失败: %v", err)
		}

		// 获取用户输入的commit信息
		log.Info("请输入本次变更信息:")
		reader := bufio.NewReader(os.Stdin)
		commitMsg, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("读取commit信息失败: %v", err)
		}
		commitMsg = strings.TrimSpace(commitMsg)

		if commitMsg == "" {
			return fmt.Errorf("commit信息不能为空")
		}

		// 执行git commit
		log.Debug("执行 git commit...")
		if _, err := utils.ExecCmd(fmt.Sprintf("git commit -m \"%s\"", commitMsg)); err != nil {
			return fmt.Errorf("git commit 失败: %v", err)
		}

		log.Debug("文档转换和提交完成")
		return nil
	}
}
