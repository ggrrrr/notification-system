package config

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type PlainFormatter struct {
	TimestampFormat string
	LevelDesc       []string
}

func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	dataStr := ""
	for k := range entry.Data {
		v := entry.Data[k]
		dataStr = fmt.Sprintf("%s%s=%+v ", dataStr, k, v)
	}

	timestamp := fmt.Sprint(entry.Time.Format(f.TimestampFormat))
	fileA := strings.Split(entry.Caller.File, "/")
	fileAL := len(fileA)

	functA := strings.Split(entry.Caller.Function, "/")
	funct := functA[len(functA)-1]
	filename := fmt.Sprintf("%s/%s", fileA[fileAL-2], fileA[fileAL-1])
	return []byte(fmt.Sprintf("%5s %s %17s:%d %10s %s %s\n",
		f.LevelDesc[entry.Level],
		timestamp,
		filename,
		entry.Caller.Line,
		funct,
		entry.Message,
		dataStr,
	)), nil
}

func text() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			fileA := strings.Split(f.File, "/")
			funct := strings.Split(f.Function, "/")
			fileAL := len(fileA)
			filename := fmt.Sprintf("%v/%v", fileA[fileAL-2], fileA[fileAL-1])
			return funct[len(funct)-1], fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})

}

func plain() {
	plainFormatter := new(PlainFormatter)
	plainFormatter.TimestampFormat = "2006-01-02T15:04:05.000000"
	plainFormatter.LevelDesc = []string{"PANIC", "FATAL", "ERROR", "WARN", "INFO", "DEBUG"}
	logrus.SetFormatter(plainFormatter)
}

func init() {
	logrus.SetReportCaller(true)
	formater := strings.TrimSpace(os.Getenv("LOG_FORMATTER"))
	level := strings.TrimSpace(os.Getenv("LOG_LEVEL"))

	switch strings.ToLower(formater) {
	case "text":
		text()
	default:
		plain()
	}

	switch strings.ToUpper(level) {
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.DebugLevel)
	}

}
