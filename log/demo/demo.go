package main

import (
	"fmt"

	"github.com/zhihanggg/gitdoc-cli/log"
)

func main() {
	// 使用默认的log
	log.Info("123")
	log.Warn("告警")
	log.Prefix("step1").Info("ok") // 带有前缀打印
	log.Inline().Debug("123")      // 不换行打印

	// 自定义创建一个，不与默认的共享
	printer := log.New().WithPrefix("xingbengwen")
	printer.EnableTrace = true  // 修改默认配置，开启Trace
	printer.DisableColor = true // 关闭颜色
	printer.Prefix = "999"      // 指定前缀
	printer.Debug("1233")

	printer.EnableInline = true
	printer.Debug("123")
	printer.Info("123")
	printer.EnableInline = false

	// 自定义颜色打印
	printer.Normal(log.StringsJoin(log.Green, []string{"1", "2"}, ","))
	printer.Normal(log.Color(log.Green, "1"))

	// 演示不换行用于交互输入
	inline := log.New().WithInline().WithEnableTrace()
	inline.Normal("please input a number:") // 不换行
	var input int
	_, _ = fmt.Scanln(&input)
	printer.Normal("input:%v", input)
}
