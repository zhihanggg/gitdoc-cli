package constant

// 操作系统
const (
	// OSLinux linux 操作系统
	OSLinux = "linux"
	// OSMac mac 操作系统
	OSMac = "darwin"
	// OSWindows windows操作系统
	OSWindows = "windows"
)

const (
	// 日志格式
	LogFormat = `%{color}%{time:15:04:05} %{shortfunc} [%{level:.4s}]%{color:reset} %{message}`
)
