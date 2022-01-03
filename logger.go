package logger

import (
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// config represents the setting for zap logger.
type config struct {
	ZapConfig zap.Config        `json:"zap_config" yaml:"zap_config"`
	LogRotate lumberjack.Logger `json:"log_rotate" yaml:"log_rotate"`
}

const (
	defaultFileName  = "config.logger"
	overrideFileName = "config.logger.override"
)

type Logger interface {
	Desugar() *zap.Logger
	Named(name string) *zap.SugaredLogger
	With(args ...interface{}) *zap.SugaredLogger
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	DPanic(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	DPanicf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	DPanicw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Sync() error
}

// NewLogger constructs new zap logger with given configuration.
// Configuration must be of type yaml or json.
func NewLogger(configPath string) Logger {
	loggerConfig := readConfig(configPath)
	zapLogger, err := build(loggerConfig)

	if err != nil {
		fmt.Printf("Failed to compose zap logger : %s", err)
		os.Exit(2)
	}

	logger := zapLogger.Sugar()
	_ = zapLogger.Sync()

	return logger
}

// readConfig of the default configuration from the given
// config path. This configuration is required.
// If a .override config exists the given default will be overwritten.
func readConfig(configPath string) *config {
	configType := ".yml"
	configDir := filepath.Dir(configPath)

	configFilePath := fmt.Sprintf("%s/%s%s", configDir, defaultFileName, configType)

	configYamlFile, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		fmt.Printf("failed to read default %s configuration. error: %s", defaultFileName, err)
		os.Exit(2)
	}

	var loggerConfig *config

	if err = yaml.Unmarshal(configYamlFile, &loggerConfig); err != nil {
		fmt.Printf("unable to decode config. error: %s", err)
		os.Exit(2)
	}

	configOverrideFilePath := fmt.Sprintf("%s/%s%s", configDir, overrideFileName, configType)

	// check if override file exists and use it if exists
	if _, err := os.Stat(configOverrideFilePath); err == nil {
		configFilePath = configOverrideFilePath
	}

	configYamlFile, err = ioutil.ReadFile(configFilePath)

	if err != nil {
		fmt.Printf("failed to read default %s configuration. error: %s", defaultFileName, err)
		os.Exit(2)
	}

	if err = yaml.Unmarshal(configYamlFile, &loggerConfig); err != nil {
		fmt.Printf("unable to decode config. error: %s", err)
		os.Exit(2)
	}

	return loggerConfig
}
