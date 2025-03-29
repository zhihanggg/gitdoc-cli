package init

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zhihanggg/gitdoc-cli/log"
	"github.com/zhihanggg/gitdoc-cli/utils"
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
		// 检查git安装情况
		if err := CheckAndInstallGit(); err != nil {
			return err
		}

		// 检查pandoc安装情况
		if err := CheckAndInstallPandoc(); err != nil {
			return err
		}

		// 检查并设置git用户信息
		if err := CheckAndSetupGitConfig(); err != nil {
			return err
		}

		return nil
	}
}

// CheckAndSetupGitConfig 检测是否设置git用户信息，如果没有则提示用户设置
func CheckAndSetupGitConfig() error {
	log.Info("开始检查git用户配置...")

	// 检查git user.name是否已设置
	userName, err := utils.ExecCmd("git config --global user.name")
	if err != nil || strings.TrimSpace(userName) == "" {
		log.Warn("未检测到git user.name配置")
		log.Info("请输入您的git用户名: ")
		var inputName string
		fmt.Scanln(&inputName)

		if inputName == "" {
			return fmt.Errorf("git用户名不能为空")
		}

		_, err := utils.ExecCmd(fmt.Sprintf("git config --global user.name \"%s\"", inputName))
		if err != nil {
			return fmt.Errorf("设置git user.name失败: %v", err)
		}
		log.Info("git user.name设置成功")
	} else {
		log.Info("git user.name已配置: %s", strings.TrimSpace(userName))
	}

	// 检查git user.email是否已设置
	userEmail, err := utils.ExecCmd("git config --global user.email")
	if err != nil || strings.TrimSpace(userEmail) == "" {
		log.Warn("未检测到git user.email配置")
		log.Info("请输入您的git邮箱: ")
		var inputEmail string
		fmt.Scanln(&inputEmail)

		if inputEmail == "" {
			return fmt.Errorf("git邮箱不能为空")
		}

		_, err := utils.ExecCmd(fmt.Sprintf("git config --global user.email \"%s\"", inputEmail))
		if err != nil {
			return fmt.Errorf("设置git user.email失败: %v", err)
		}
		log.Info("git user.email设置成功")
	} else {
		log.Info("git user.email已配置: %s", strings.TrimSpace(userEmail))
	}

	return nil
}

// CheckAndInstallGit 检测是否安装git，如果没有则安装
func CheckAndInstallGit() error {
	log.Info("开始检查git安装情况...")
	// 检查是否已安装git
	_, err := utils.ExecCmd("which git")
	if err == nil {
		// git已安装
		log.Info("git已安装")
		return nil
	}

	// git未安装，尝试安装
	log.Info("检测到系统未安装git, 正在尝试安装...")

	// 在macOS上使用brew安装git
	_, err = utils.ExecCmd("brew install git")
	if err != nil {
		return fmt.Errorf("安装git失败: %v", err)
	}

	log.Info("git安装成功")
	return nil
}

// CheckAndInstallPandoc 检测是否安装pandoc，如果没有则安装
func CheckAndInstallPandoc() error {
	log.Info("开始检查pandoc安装情况...")
	// 检查是否已安装pandoc
	_, err := utils.ExecCmd("which pandoc")
	if err == nil {
		// pandoc已安装
		log.Info("pandoc已安装")
		return nil
	}

	// pandoc未安装，尝试安装
	log.Info("检测到系统未安装pandoc, 正在尝试安装...")

	// 在macOS上使用brew安装pandoc
	_, err = utils.ExecCmd("brew install pandoc")
	if err != nil {
		return fmt.Errorf("安装pandoc失败: %v", err)
	}

	log.Info("pandoc安装成功")
	return nil
}
