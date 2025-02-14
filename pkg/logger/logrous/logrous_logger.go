package logrous

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/nolleh/caption_json_formatter"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type logrusLogger struct {
	level    string
	encoding string
	logger   *logrus.Logger
}

// For mapping config logger
var loggerLevelMap = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
	"panic": logrus.PanicLevel,
	"fatal": logrus.FatalLevel,
}

func (l *logrusLogger) GetLoggerLevel() logrus.Level {
	level, exist := loggerLevelMap[l.level]
	if !exist {
		return logrus.DebugLevel
	}

	return level
}

// NewLogrusLogger creates a new logrus logger
func NewLogrusLogger(cfg *logger.LogConfig) logger.Logger {
	logrusLogger := &logrusLogger{level: cfg.LogLevel}
	logrusLogger.initLogger()

	return logrusLogger
}

// InitLogger Init logger
func (l *logrusLogger) initLogger() {
	logLevel := l.GetLoggerLevel()

	// Create a new instance of the logger. You can have any number of instances.
	logrusLogger := logrus.New()

	logrusLogger.SetLevel(logLevel)

	// Output to stdout instead of the defaultLogger stderr
	// Can be any io.Writer, see below for File example
	logrusLogger.SetOutput(os.Stdout)

	if core.IsDevelopment() {
		logrusLogger.SetReportCaller(false)
		logrusLogger.SetFormatter(&logrus.TextFormatter{
			DisableColors: false,
			ForceColors:   true,
			FullTimestamp: true,
		})
	} else {
		logrusLogger.SetReportCaller(false)
		//https://github.com/nolleh/caption_json_formatter
		logrusLogger.SetFormatter(&caption_json_formatter.Formatter{PrettyPrint: true})
	}

	l.logger = logrusLogger
}

func (l *logrusLogger) LogType() logger.LogType {
	return logger.Logrus
}

func (l *logrusLogger) Configure(cfg func(internalLog interface{})) {
	cfg(l.logger)
}

func (l *logrusLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *logrusLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *logrusLogger) Debugw(msg string, fields logger.Fields) {
	entry := l.mapToFields(fields)
	entry.Debug(msg)
}

func (l *logrusLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *logrusLogger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *logrusLogger) Infow(msg string, fields logger.Fields) {
	entry := l.mapToFields(fields)
	entry.Info(msg)
}

func (l *logrusLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *logrusLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *logrusLogger) WarnMsg(msg string, err error) {
	l.logger.Warn(msg, logrus.WithField("error", err.Error()))
}

func (l *logrusLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *logrusLogger) Errorw(msg string, fields logger.Fields) {
	entry := l.mapToFields(fields)
	entry.Error(msg)
}

func (l *logrusLogger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l *logrusLogger) Err(msg string, err error) {
	l.logger.Error(msg, logrus.WithField("error", err.Error()))
}

func (l *logrusLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *logrusLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

func (l *logrusLogger) Printf(template string, args ...interface{}) {
	l.logger.Printf(template, args...)
}

func (l *logrusLogger) WithName(name string) {
	l.logger.WithField(constants.NAME, name)
}

func (l *logrusLogger) GrpcMiddlewareAccessLogger(method string, time time.Duration, metaData map[string][]string, err error) {
	l.Info(
		constants.GRPC,
		logrus.WithField(constants.METHOD, method),
		logrus.WithField(constants.TIME, time),
		logrus.WithField(constants.METADATA, metaData),
		logrus.WithError(err),
	)
}

func (l *logrusLogger) GrpcClientInterceptorLogger(method string, req interface{}, reply interface{}, time time.Duration, metaData map[string][]string, err error) {
	l.Info(
		constants.GRPC,
		logrus.WithField(constants.METHOD, method),
		logrus.WithField(constants.REQUEST, req),
		logrus.WithField(constants.REPLY, reply),
		logrus.WithField(constants.TIME, time),
		logrus.WithField(constants.METADATA, metaData),
		logrus.WithError(err),
	)
}

func (l *logrusLogger) KafkaProcessMessage(topic string, partition int, message string, workerID int, offset int64, time time.Time) {
	l.Debug(
		"Processing Kafka message",
		logrus.WithField(constants.Topic, topic),
		logrus.WithField(constants.Partition, partition),
		logrus.WithField(constants.Message, message),
		logrus.WithField(constants.WorkerID, workerID),
		logrus.WithField(constants.Offset, offset),
		logrus.WithField(constants.Time, time),
	)
}

func (l *logrusLogger) KafkaLogCommittedMessage(topic string, partition int, offset int64) {
	l.Info(
		"Committed Kafka message",
		logrus.WithField(constants.Topic, topic),
		logrus.WithField(constants.Partition, partition),
		logrus.WithField(constants.Offset, offset),
	)
}

func (l *logrusLogger) mapToFields(fields map[string]interface{}) *logrus.Entry {
	return l.logger.WithFields(logrus.Fields{"ss": 1})
}
