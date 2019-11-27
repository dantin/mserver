package logutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultLogTimeFormat = "2006/01/02 15:04:05.000"
	defaultLogMaxSize    = 300 // MB
	defaultLogFormat     = "text"
	defaultLogLevel      = log.InfoLevel
)

// FileLogConfig serializes file log related config in toml.
type FileLogConfig struct {
	// Log filename, leave empty to disable file log.
	Filename string `toml:"filename"`
	// Is log compress enabled.
	LogCompress bool `toml:"log-compress"`
	// Max size for a single file, in MB.
	MaxSize int `toml:"max-size"`
	// Max log keep days.
	MaxDays int `toml:"max-days"`
	// Maximum number of old log files to retain.
	MaxBackups int `toml:"max-backups"`
}

// LogConfig serializes log related config in toml.
type LogConfig struct {
	// Log level.
	Level string `toml:"level"`
	// Log format, json, text or console.
	Format string `toml:"format"`
	// Disable automatic timestamps in output.
	DisableTimestamp bool `toml:"disable-timestamp"`
	// File log config.
	File FileLogConfig `toml:"file"`
}

// textFormatter is customized text formatter.
type textFormatter struct {
	DisableTimestamp bool
}

// Format implements logurs.Formatter.
func (f *textFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	if !f.DisableTimestamp {
		fmt.Fprintf(b, "%s ", entry.Time.Format(defaultLogTimeFormat))
	}
	if file, ok := entry.Data["file"]; ok {
		fmt.Fprintf(b, "%s:%v", file, entry.Data["line"])
	}
	fmt.Fprintf(b, "[%s] %s", entry.Level.String(), entry.Message)
	for k, v := range entry.Data {
		if k != "file" && k != "line" {
			fmt.Fprintf(b, " %v=%v", k, v)
		}
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}

// InitLogger initializes logger.
func InitLogger(cfg *LogConfig) error {
	log.SetLevel(stringToLogLevel(cfg.Level))
	formatter := stringToLogFormatter(defaultLogFormat, false)
	log.SetFormatter(formatter)
	var output io.Writer
	if len(cfg.File.Filename) != 0 {
		o, err := initFileLog(&cfg.File)
		if err != nil {
			return errors.Trace(err)
		}
		output = o
	} else {
		output = os.Stdout
	}
	log.SetOutput(output)

	return nil
}

func initFileLog(cfg *FileLogConfig) (io.Writer, error) {
	if st, err := os.Stat(cfg.Filename); err == nil {
		if st.IsDir() {
			return nil, errors.New("can't use directory as log file name")
		}
	}

	if cfg.MaxSize == 0 {
		cfg.MaxSize = defaultLogMaxSize
	}

	// use lumberjack to logrotate
	output := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxDays,
		Compress:   cfg.LogCompress,
		LocalTime:  true,
	}

	return output, nil
}

func stringToLogFormatter(format string, disableTimestamp bool) log.Formatter {
	switch strings.ToLower(format) {
	case "text":
		return &textFormatter{
			DisableTimestamp: disableTimestamp,
		}
	case "json":
		return &log.JSONFormatter{
			TimestampFormat:  defaultLogTimeFormat,
			DisableTimestamp: disableTimestamp,
		}
	case "console":
		return &log.TextFormatter{
			FullTimestamp:    true,
			TimestampFormat:  defaultLogTimeFormat,
			DisableTimestamp: disableTimestamp,
		}
	default:
		return &textFormatter{}
	}
}

func stringToLogLevel(level string) log.Level {
	switch strings.ToLower(level) {
	case "fatal":
		return log.FatalLevel
	case "error":
		return log.ErrorLevel
	case "warn", "warning":
		return log.WarnLevel
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	default:
		return defaultLogLevel
	}
}
