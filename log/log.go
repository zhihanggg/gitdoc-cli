// Package log 命令行工具常用彩色日志输出的轻量级日志库，不适用日志服务日志的输出
package log

import (
	"errors"
	"fmt"
)

const (
	// WarnLevel TODO
	WarnLevel = "Warn"
	// DebugLevel TODO
	DebugLevel = "Debug"
	// ErrorLevel TODO
	ErrorLevel = "Error"
)

// Debug 绿色，用于打印提示信息
func Debug(format string, a ...interface{}) {
	DefaultStd.WithAddCallDepth(1).Debug(format, a...)
}

// Info 蓝色，用于打印变更信息
func Info(format string, a ...interface{}) {
	DefaultStd.WithAddCallDepth(1).Info(format, a...)
}

// Warn 黄色，打印警示信息，比如打印检查没有通过
func Warn(format string, a ...interface{}) {
	DefaultStd.WithAddCallDepth(1).Warn(format, a...)
}

// Error 红色，打印错误信息
func Error(format string, a ...interface{}) {
	DefaultStd.WithAddCallDepth(1).Error(format, a...)
}

// Trace 没有颜色，打印调试信息，只吃EnableTrace才会打印信息
func Trace(format string, a ...interface{}) {
	DefaultStd.WithAddCallDepth(1).Trace(format, a...)
}

// Normal 没有颜色，打印正常调试信息
func Normal(format string, a ...interface{}) {
	DefaultStd.WithAddCallDepth(1).Normal(format, a...)
}

// Prefix 返回一个标准的带有前缀的Printer，高频接口单独独立一个函数
func Prefix(prefix string) Printer {
	return DefaultStd.WithPrefix(prefix)
}

// Inline 返回一个标准的不换行打印的Printer，高频接口单独独立一个函数
func Inline() Printer {
	return DefaultStd.WithInline()
}

// StringsJoin 颜色字符串拼接
func StringsJoin(color TypeColor, list []string, sep string) string {
	return DefaultStd.StringsJoin(color, list, sep)
}

// Color 输出彩色字符串
func Color(color TypeColor, format string, a ...interface{}) string {
	return DefaultStd.Color(color, format, a...)
}

// Alarm 按级别告警
func Alarm(message string, level string) error {
	switch level {
	case WarnLevel:
		Warn(message)
	case DebugLevel:
		Debug(message)
	case ErrorLevel:
		Error(message)
		return errors.New(message)
	default:
		Info(message)
	}
	return nil
}

// SetColor 设置打印颜色
func SetColor(body string, color TypeColor) string {
	if color == None {
		return body
	}
	return fmt.Sprintf(FormatPre+"%v"+FormatTail, color, body)
}
