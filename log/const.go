package log

// TypeColor 颜色种类
type TypeColor int

const (
	// Green 绿色
	Green TypeColor = 32
	// Blue 蓝色
	Blue TypeColor = 34
	// Red 红色
	Red TypeColor = 31
	// Yellow 黄色
	Yellow TypeColor = 33
	// None 无颜色
	None TypeColor = -1

	// FormatPre 颜色开始
	FormatPre = "\x1b[0;%dm"
	// FormatTail 颜色重置
	FormatTail = "\x1b[0m"
)
