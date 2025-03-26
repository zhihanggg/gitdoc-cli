package log

import (
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/suite"
)

type PrinterSuite struct {
	suite.Suite
	printer *Printer
}

func TestPrinterSuite(t *testing.T) {
	suite.Run(t, new(PrinterSuite))
}

func (p *PrinterSuite) SetupTest() {
	p.printer = DefaultStd
}

func (p *PrinterSuite) TearDownTest() {
	DefaultStd = p.printer
}

func TestDemo(t *testing.T) {
	Debug("123")
	DefaultStd.EnableTrace = true
	Debug("123")

	DefaultStd.Prefix = "123"
	Info("abc")

	n := New().WithPrefix("123")
	n.EnableTrace = true
	n.DisableColor = true
	n.Info("123")
}

func (p *PrinterSuite) TestPrinter_WithAddCallDepth() {
	p.Suite.EqualValues(6, DefaultStd.WithAddCallDepth(3).callDepth)
	p.Suite.EqualValues(3, DefaultStd.WithAddCallDepth(0).callDepth)
	p.Suite.EqualValues(2, DefaultStd.WithAddCallDepth(-1).callDepth)
}

func (p *PrinterSuite) TestEnableTrace() {
	count := 0
	logRE := regexp.MustCompile(`^[0-9][0-9]:[0-9][0-9]:[0-9][0-9] colorlog_test.go:([0-9]+): 1234\n$`)
	patches := gomonkey.ApplyMethod(reflect.TypeOf(writerImp{}), "Write",
		func(w writerImp, b []byte) (n int, err error) {
			count++
			p.T().Log(string(b))
			if !logRE.Match(b) {
				p.T().Log("GOT:", string(b), "not MATCH")
				p.T().Fail()
			}
			return os.Stderr.Write(b)
		},
	)
	defer patches.Reset()

	printer := DefaultStd.WithEnableTrace().WithDisableColor()
	printer.Trace("123%v", 4)
	printer.Warn("123%v", 4)
	printer.Info("123%v", 4)
	printer.Error("123%v", 4)
	printer.Debug("123%v", 4)
	printer.Normal("123%v", 4)

	DefaultStd.DisableColor = true
	DefaultStd.EnableTrace = true
	Trace("123%v", 4)
	Warn("123%v", 4)
	Info("123%v", 4)
	Error("123%v", 4)
	Debug("123%v", 4)
	Normal("123%v", 4)
	p.Suite.EqualValues(12, count)

	// 开启颜色，trace依然没有颜色
	DefaultStd.DisableColor = false
	Trace("123%v", 4)
}

func (p *PrinterSuite) TestDisableTrace() {
	count := 0
	logRE := regexp.MustCompile(`^1234\n$`)
	patches := gomonkey.ApplyMethod(reflect.TypeOf(writerImp{}), "Write",
		func(w writerImp, b []byte) (n int, err error) {
			count++
			p.T().Log(string(b))
			if !logRE.Match(b) {
				p.T().Log("GOT:", string(b), "not MATCH")
				p.T().Fail()
			}
			return os.Stderr.Write(b)
		},
	)
	defer patches.Reset()

	printer := DefaultStd.WithDisableColor()
	printer.Trace("123%v", 4) // 不打印
	printer.Warn("123%v", 4)
	printer.Info("123%v", 4)
	printer.Error("123%v", 4)
	printer.Debug("123%v", 4)
	printer.Normal("123%v", 4)

	DefaultStd.DisableColor = true
	Trace("123%v", 4) // 不打印
	Warn("123%v", 4)
	Info("123%v", 4)
	Error("123%v", 4)
	Debug("123%v", 4)
	Normal("123%v", 4)
	p.Suite.EqualValues(10, count)
}

func (p *PrinterSuite) TestNewErrorf() {
	err := Prefix("prefix").NewErrorf("err:%v", 123)
	p.Suite.NotNil(err)
	p.Suite.EqualValues("[prefix]err:123", err.Error())
}

func (p *PrinterSuite) TestAlarm() {
	err := Alarm("Warn", WarnLevel)
	p.Suite.Nil(err)
	err = Alarm("Error", ErrorLevel)
	p.Suite.NotNil(err)
}
