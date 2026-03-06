package functions

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

var (
	infoColor   = color.New(color.FgGreen, color.Bold).SprintFunc()
	warnColor   = color.New(color.FgYellow, color.Bold).SprintFunc()
	errorColor  = color.New(color.FgRed, color.Bold, color.BgBlack).SprintFunc()
	debugColor  = color.New(color.FgCyan, color.Bold).SprintFunc()
	prefixColor = color.New(color.FgHiBlack, color.Italic).SprintFunc()
	msgInfo     = color.New(color.FgWhite).SprintFunc()
	msgWarn     = color.New(color.FgHiYellow).SprintFunc()
	msgError    = color.New(color.FgHiRed, color.Bold).SprintFunc()
	msgDebug    = color.New(color.FgHiCyan).SprintFunc()
)

func logWithLevel(level, msg string, args ...interface{}) {
	now := time.Now().Format("2006-01-02 15:04:05")
	var (
		levelStr string
		msgStr   string
	)
	switch level {
	case "INFO":
		levelStr = infoColor("[INFO]")
		msgStr = msgInfo(fmt.Sprintf(msg, args...))
	case "WARN":
		levelStr = warnColor("[WARN]")
		msgStr = msgWarn(fmt.Sprintf(msg, args...))
	case "ERROR":
		levelStr = errorColor("[ERROR]")
		msgStr = msgError(fmt.Sprintf(msg, args...))
	case "DEBUG":
		levelStr = debugColor("[DEBUG]")
		msgStr = msgDebug(fmt.Sprintf(msg, args...))
	default:
		levelStr = prefixColor("[LOG]")
		msgStr = fmt.Sprintf(msg, args...)
	}
	fmt.Fprintf(os.Stdout, "%s %s %s\n", prefixColor(now), levelStr, msgStr)
}

func Info(msg string, args ...interface{}) {
	logWithLevel("INFO", msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logWithLevel("WARN", msg, args...)
}

func Error(msg string, args ...interface{}) {
	logWithLevel("ERROR", msg, args...)
}

func Debug(msg string, args ...interface{}) {
	logWithLevel("DEBUG", msg, args...)
}
