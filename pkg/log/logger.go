package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"guduo/pkg/constant"
	"path"
)

// 初始化log
func InitLogger() {
	log.AddHook(newLfsHook())
}

// 使用local file system 的 hook， 将日志文件存于本地
func newLfsHook() *lfshook.LfsHook {

	writer, _ := rotatelogs.New(
		path.Join(constant.AppPath, "log/go-%Y-%m-%d.log"),
	)

	writerWarn, _ := rotatelogs.New(
		path.Join(constant.AppPath, "log/go-warn-%Y-%m-%d.log"),
	)

	writerError, _ := rotatelogs.New(
		path.Join(constant.AppPath, "log/go-error-%Y-%m-%d.log"),
	)

	writerFatal, _ := rotatelogs.New(
		path.Join(constant.AppPath, "log/go-fatal-%Y-%m-%d.log"),
	)

	writerPanic, _ := rotatelogs.New(
		path.Join(constant.AppPath, "log/go-panic-%Y-%m-%d.log"),
	)

	h := lfshook.NewHook(
		lfshook.WriterMap{
			log.PanicLevel: writerPanic,
			log.FatalLevel: writerFatal,
			log.ErrorLevel: writerError,
			log.WarnLevel:  writerWarn,
			log.InfoLevel:  writer,
			log.DebugLevel: writer,
			log.TraceLevel: writer,
		},
		&log.TextFormatter{},
	)

	return h
}
