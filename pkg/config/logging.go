/*
Copyright © 2021 Rewire Group, Inc. All rights reserved.

Proprietary and confidential.

Unauthorized copying or use of this file, in any medium or form,
is strictly prohibited.
*/

package config

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"moul.io/zapgorm2"
	"strings"
)

// Zlog wraps the zap's Sugar logger.
type Zlog struct {
	*zap.SugaredLogger

	// holds global zap.Config.
	Config *zap.Config
}

// logger is local reference to initiated SugarLogger
var zlogger *Zlog

// logger contains reference to an initialized Sugar logger ready to be used. Can be used concurrently.
// you can simply start using as:
//  var zlog = config.GetLogger()
// 	zlog.Infow("Failed to fetch URL.", "url", "http://google.com", "attempt", 3, "backoff", time.Second)
//	zlog.Infof("Failed to fetch URL: %s", "http://google.com")
//
//	zlog.Info("Failed to fetch URL.",
//		// Structured context as strongly typed fields.
//		zap.String("url", "http://google.com"),
//		zap.Int("attempt", 3),
//		zap.Duration("backoff", time.Second),
//	)
func GetLogger() *Zlog {
	if zlogger == nil {
		initGlobalLogger(true)
	}
	return zlogger
}

// GORM has some specific reqs for using a logger, so using a wrapper over zap
// https://github.com/moul/zapgorm2
func GetGormLogger() zapgorm2.Logger {
	l := zapgorm2.New(zap.L())
	l.SetAsDefault()
	return l
}

// Option is an additional configuration for zap logger.
type Option func(*zap.Config) error

// WithLogPaths specifies the sink for logs it can be "stdout", "stderr" or a log file path
func WithLogPaths(paths ...string) Option {
	return func(config *zap.Config) error {
		fmt.Printf(">>> Logging to: %s \n", strings.Join(paths, ", "))
		config.OutputPaths = paths
		return nil
	}
}

// WithStacktrace specifies if the stack trace should be included in the log.
func WithStacktrace() Option {
	return func(config *zap.Config) error {
		config.DisableStacktrace = false
		return nil
	}
}

// WithCaller specifies if the caller file and line number should be included in log.
func WithCaller(enable bool) Option {
	return func(config *zap.Config) error {
		config.DisableCaller = !enable
		return nil
	}
}

// WithTime specifies the key name for time in json format.
// Also to remove time you can provide an empty key.
func WithTimeKey(key string) Option {
	return func(config *zap.Config) error {
		// zapcore.RFC3339TimeEncoder
		config.EncoderConfig.TimeKey = key
		return nil
	}
}

// WithColor specifies if logs should be colorful.
func WithColor(enable bool) Option {
	return func(config *zap.Config) error {
		if enable {
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		} else {
			config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		}
		return nil
	}
}

// BuildLogger initializes logger and returns a new zap.Logger logger instance with defaults set
// below or provided with options.
func BuildLogger(verbose bool, options ...Option) (*zap.Logger, *zap.Config) {
	var config zap.Config
	if verbose {
		config = zap.NewDevelopmentConfig()
		config.Level.SetLevel(zap.DebugLevel)
		config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	} else {
		config = zap.NewProductionConfig()
		config.Level.SetLevel(zap.InfoLevel)
	}

	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	config.DisableCaller = true     // false will add the file and line number to log
	config.DisableStacktrace = true // include full stack trace of every log line

	// collect provided options.
	for _, opt := range options {
		if err := opt(&config); err != nil {
			log.Fatalf("failed setting up logger with options, got err %v", err)
			return nil, nil
		}
	}

	l, err := config.Build()
	if err != nil {
		log.Fatalf("Error while setting up logging, got: %v", err)
	}

	return l, &config
}

// initGlobalLogger the global logger, which can be retrieved by GetLogger()
func initGlobalLogger(verbose bool, options ...Option) *Zlog {
	logger, zcfg := BuildLogger(verbose, options...)
	zap.ReplaceGlobals(logger)
	zlogger = &Zlog{logger.Sugar(), zcfg}
	return zlogger
}

// ConfigureGlobalLogger configures the global logger, which can be retrieved by GetLogger()
func ConfigureGlobalLogger(verbose bool, options ...Option) *Zlog {
	logger, zcfg := BuildLogger(verbose, options...)
	zap.ReplaceGlobals(logger)

	if zlogger == nil {
		zlogger = &Zlog{logger.Sugar(), zcfg}
	} else {
		zlogger.SugaredLogger = logger.Sugar()
		zlogger.Config = zcfg
	}
	return zlogger
}

// logWriter is an implementation of io.Writer which writes to provided logger.
type logWriter struct {
	logger   *zap.Logger
	logLevel zapcore.Level
}

func NewLogWriter(logger *zap.Logger, logLevel zapcore.Level) *logWriter {
	return &logWriter{
		logger:   logger,
		logLevel: logLevel,
	}
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	byt := bytes.TrimSpace(p)
	byt = bytes.Replace(byt, []byte("\n"), []byte("\n\t"), -1)

	if ce := w.logger.Check(w.logLevel, string(byt)); ce != nil {
		ce.Write()
	}
	return len(p), nil
}
