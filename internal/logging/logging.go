package logging

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"

	"github.com/charmbracelet/lipgloss"
	charmlog "github.com/charmbracelet/log"
)

type Level int32

const (
	TraceLevel Level = -8
	DebugLevel Level = -4
	InfoLevel  Level = 0
	WarnLevel  Level = 4
	ErrorLevel Level = 8
	OffLevel   Level = 100
)

func ParseLevel(s string) (Level, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "off":
		return OffLevel, nil
	case "trace":
		return TraceLevel, nil
	case "debug":
		return DebugLevel, nil
	case "info":
		return InfoLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "error":
		return ErrorLevel, nil
	default:
		return OffLevel, fmt.Errorf("unknown --logLevel %q (want one of: off, error, warn, info, debug, trace)", s)
	}
}

var global atomic.Pointer[charmlog.Logger]

func init() {
	global.Store(newLogger(OffLevel))
}

func Init(level Level) {
	global.Store(newLogger(level))
}

func newLogger(level Level) *charmlog.Logger {
	l := charmlog.NewWithOptions(os.Stderr, charmlog.Options{
		Level:           charmlog.Level(level),
		ReportTimestamp: true,
		TimeFormat:      "15:04:05",
	})

	styles := charmlog.DefaultStyles()
	styles.Levels[charmlog.Level(TraceLevel)] = lipgloss.NewStyle().
		SetString("TRACE").
		Bold(true).
		Foreground(lipgloss.Color("245"))
	l.SetStyles(styles)

	return l
}

func get() *charmlog.Logger { return global.Load() }

func Trace(msg any, kv ...any) { get().Log(charmlog.Level(TraceLevel), msg, kv...) }

func Debug(msg any, kv ...any) { get().Debug(msg, kv...) }

func Info(msg any, kv ...any) { get().Info(msg, kv...) }

func Warn(msg any, kv ...any) { get().Warn(msg, kv...) }

func Error(msg any, kv ...any) { get().Error(msg, kv...) }

func Fatal(msg any, kv ...any) { get().Fatal(msg, kv...) }

func With(kv ...any) *charmlog.Logger { return get().With(kv...) }
