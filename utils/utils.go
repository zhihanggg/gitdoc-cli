// Package utils TODO
package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/zhihanggg/gitdoc-cli/constant"
	"github.com/zhihanggg/gitdoc-cli/log"

	"github.com/spf13/cobra"
)

// GetLocalHostIp 通过发送udp请求来获取本机ip
func GetLocalHostIp() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.String(), ":")[0]
	return ip, nil
}

// ReadFile 读取并文件
func ReadFile(file string) (string, error) {
	log.Trace("ReadFile:", file)
	content, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("read file:%s error:%v", file, err)
	}
	return string(content), nil
}

// IsValidJson 是否json
func IsValidJson(jsonStr string) bool {
	var js interface{}
	return json.Unmarshal([]byte(jsonStr), &js) == nil
}

// GetMilliTimeHandle 毫微秒转毫秒
func GetMilliTimeHandle(nanosecondNum int64) int64 {
	milliSecondNumber := nanosecondNum / 1000000.0
	return milliSecondNumber
}

// IsContains 判断元素是否在切片中
func IsContains(slice []string, element string) bool {
	for _, s := range slice {
		if s == element {
			return true
		}
	}
	return false
}

// ContainsAnySubstring s是否有值模糊存在于v中
func ContainsAnySubstring(s []string, v string) bool {
	for _, vs := range s {
		if strings.Contains(v, vs) {
			return true
		}
	}
	return false
}

// Diff 求两个切片差集
func Diff(slice1, slice2 []string) (diff []string) {
	m := make(map[string]string)
	for _, v := range slice1 {
		m[v] = v
	}
	for _, v := range slice2 {
		if m[v] != "" {
			delete(m, v)
		}
	}
	for _, s2 := range m {
		diff = append(diff, s2)
	}
	return diff
}

// GetFrameworkTypeStr 根据整数获取对应的框架类型
func GetFrameworkTypeStr(i int64) string {
	switch i {
	case 2:
		return "trpc_java"
	case 3:
		return "trpc_go"
	case 4:
		return "fable"
	case 5:
		return "middle"
	default:
	}
	return fmt.Sprintf("unknowm framework_type: %d", i)
}

// ReadEnvFile 读取local.env
func ReadEnvFile(path string) (map[string]string, map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	// local.env 文件内分两种数据 1 命令生成文件时的参数 2 服务自身环境变量
	// 参数格式 # @env-type=DEV
	// 环境变量格式 env-type=DEV
	param := make(map[string]string)
	env := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// # @env-type=DEV
		if strings.HasPrefix(strings.TrimSpace(line), "# @") {
			key, value := parseKeyValue(strings.Replace(line, "# @", "", 1))
			param[key] = value
			continue
		}
		key, value := parseKeyValue(line)
		env[key] = value
	}
	return param, env, scanner.Err()
}

// CmdNotSupportWindowsIntercept 判断 os
func CmdNotSupportWindowsIntercept(cmd string) error {
	if runtime.GOOS == constant.OSWindows {
		return fmt.Errorf("%s 命令目前不支持windows环境", cmd)
	}
	return nil
}

// GetParamPrefix 返回命令的参数路径（忽略顶级命令的路径），如 'fit-cli run'，则返回 'run.'
func GetParamPrefix(cmd *cobra.Command) string {
	if cmd == nil {
		return ""
	}
	completePrefix := strings.Replace(cmd.CommandPath(), " ", ".", -1) + "."
	return completePrefix[strings.Index(completePrefix, ".")+1:]
}

func parseKeyValue(line string) (string, string) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// GetTimeDiffHours 获取hours小时之前的时间
func GetTimeDiffHours(hours int64) string {
	t := time.Now()
	t = t.Add(time.Duration(-hours) * time.Hour)
	return t.Format("2006-01-02 15:04:05")
}

// GetOrDefault 泛型方法：如果第一个参数是零值，返回第二个参数
func GetOrDefault[T comparable](value, defaultValue T) T {
	switch v := any(value).(type) {
	case string:
		if v == "" {
			return defaultValue
		}
	case int:
		if v == 0 {
			return defaultValue
		}
	case *struct{}:
		if v == nil {
			return defaultValue
		}
	default:
		return defaultValue
	}
	return value
}

// ParseStrToMap str 转 map
func ParseStrToMap(str string) map[string]string {
	if str == "" {
		return nil
	}
	var result map[string]string
	if err := json.Unmarshal([]byte(str), &result); err != nil {
		log.Error("%s Unmarshal to map err:%v", err)
		return nil
	}
	return result
}

// MergeMapValue 两个 map 合并
func MergeMapValue(mapA, mapB map[string]string) map[string]string {
	if mapA == nil {
		mapA = make(map[string]string)
	}
	for key, value := range mapB {
		mapA[key] = value
	}
	return mapA
}

// ExecCmd 执行命令并返回结果
func ExecCmd(cmd string) (string, error) {
	if runtime.GOOS == constant.OSWindows {
		return "", fmt.Errorf("暂不支持 Windows 系统")
	}

	// 执行命令
	output, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("执行命令失败: %v, 输出: %s", err, string(output))
	}

	return string(output), nil
}

// ScanFilesByExt 递归扫描指定目录下的特定扩展名文件
func ScanFilesByExt(root string, extensions []string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			for _, validExt := range extensions {
				if ext == validExt {
					files = append(files, path)
					break
				}
			}
		}
		return nil
	})

	return files, err
}

// ConvertDocToMarkdown 使用pandoc将doc/docx文件转换为markdown
func ConvertDocToMarkdown(docPath, mdPath string) error {
	cmd := fmt.Sprintf("pandoc -s \"%s\" --extract-media=. -t markdown -o \"%s\"", docPath, mdPath)
	_, err := ExecCmd(cmd)
	return err
}
