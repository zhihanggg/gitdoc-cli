package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/zhihanggg/gitdoc-cli/constant"
)

// Option 可配置项
type Option struct {
	// DisableColor 关闭颜色
	DisableColor bool
	// EnableTrace 开启打印函数位置
	EnableTrace bool
	// Prefix 打印前缀，格式为[prefix]abc
	Prefix string
	// EnableInline 不换行打印，可以用于交互输入
	EnableInline bool
}

// Printer 颜色打印器
type Printer struct {
	// Option 可配置项目
	Option

	// traceLogger 打印Trace相关日志
	traceLogger *log.Logger
	// defaultLogger 正常日志打印
	defaultLogger *log.Logger
	// traceLoggerInline 打印Trace相关日志，不换行
	traceLoggerInline *log.Logger
	// defaultLoggerInline 正常日志打印，不换行
	defaultLoggerInline *log.Logger
	// callDepth 调用深度
	callDepth int
}

// writerImp 输出
type writerImp struct {
	// EnableInline 开启行内打印
	EnableInline bool
}

// Write implement
func (p writerImp) Write(b []byte) (n int, err error) {
	if p.EnableInline && len(b) >= 1 && b[len(b)-1] == '\n' {
		return os.Stderr.Write(b[:len(b)-1])
	}
	return os.Stderr.Write(b)
}

// New 新建一个Printer
func New() *Printer {
	disableColor := false
	if runtime.GOOS == constant.OSWindows { // windows环境不支持日志输出带颜色，否则会产生乱码
		disableColor = true
	}
	p := &Printer{
		Option: Option{
			DisableColor: disableColor,
			EnableTrace:  false,
			Prefix:       "",
			EnableInline: false,
		},
		callDepth: 3,
	}
	p.defaultLogger = log.New(&writerImp{}, "", 0)
	p.traceLogger = log.New(&writerImp{}, "", log.Lshortfile|log.Ltime)
	p.defaultLoggerInline = log.New(&writerImp{EnableInline: true}, "", 0)
	p.traceLoggerInline = log.New(&writerImp{EnableInline: true}, "", log.Lshortfile|log.Ltime)
	return p
}

// NewErrorf 生成一个带有前缀的error
func (p Printer) NewErrorf(format string, a ...interface{}) error {
	p.setPrefix(&format)
	return fmt.Errorf(format, a...)
}

// setPrefix 设置前缀
func (p Printer) setPrefix(format *string) {
	if p.Prefix == "" || format == nil {
		return
	}
	*format = "[" + p.Prefix + "]" + *format
}

// DefaultStd 默认的标准Printer
var DefaultStd = New()

// WithAddCallDepth 返回一个带有修改call depth的Printer
func (p Printer) WithAddCallDepth(addDepth int) Printer {
	p.callDepth += addDepth
	return p
}

// WithPrefix 返回一个带有前缀的Printer
func (p Printer) WithPrefix(prefix string) Printer {
	p.Prefix = prefix
	return p
}

// WithEnableTrace 返回开启调试模式的Printer
func (p Printer) WithEnableTrace() Printer {
	p.EnableTrace = true
	return p
}

// WithDisableColor 返回关闭颜色打印的Printer
func (p Printer) WithDisableColor() Printer {
	p.DisableColor = true
	return p
}

// WithInline 返回用户输出不换行的信息的Printer
func (p Printer) WithInline() Printer {
	p.EnableInline = true
	return p
}

// output 公共输出
func (p Printer) output(callDepth int, color TypeColor, format string, a ...interface{}) {
	p.setPrefix(&format)
	if p.EnableTrace {
		if p.EnableInline {
			_ = p.traceLoggerInline.Output(callDepth, p.Color(color, format, a...))
			return
		}
		_ = p.traceLogger.Output(callDepth, p.Color(color, format, a...))
		return
	}
	if p.EnableInline {
		_ = p.defaultLoggerInline.Output(callDepth, p.Color(color, format, a...))
		return
	}
	_ = p.defaultLogger.Output(callDepth, p.Color(color, format, a...))
}

// Trace 用于打印开发阶段的信息
func (p Printer) Trace(format string, a ...interface{}) {
	if p.EnableTrace {
		p.output(p.callDepth, None, format, a...)
	}
}

// Debug 用于打印检查正常的信息
func (p Printer) Debug(format string, a ...interface{}) {
	p.output(p.callDepth, Green, format, a...)
}

// Info 用于打印修改操作的信息
func (p Printer) Info(format string, a ...interface{}) {
	p.output(p.callDepth, Blue, format, a...)
}

// Warn 用于打印检查不正确的信息
func (p Printer) Warn(format string, a ...interface{}) {
	p.output(p.callDepth, Yellow, format, a...)
}

// Error 用于打印错误信息
func (p Printer) Error(format string, a ...interface{}) {
	p.output(p.callDepth, Red, format, a...)
}

// Normal 用于无颜色打印
func (p Printer) Normal(format string, a ...interface{}) {
	p.output(p.callDepth, None, format, a...)
}

// Color 颜色输出
func (p Printer) Color(color TypeColor, format string, a ...interface{}) string {
	x := fmt.Sprintf(format, a...)
	if p.DisableColor || color == None {
		return x
	}
	return fmt.Sprintf(FormatPre+"%v"+FormatTail, color, x)
}

// StringsJoin 合并打印后的颜色
func (p Printer) StringsJoin(color TypeColor, list []string, sep string) string {
	res := make([]string, 0)
	for _, v := range list {
		res = append(res, p.Color(color, "%v", v))
	}
	return strings.Join(res, sep)
}
