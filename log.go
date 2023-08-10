package log

var (
	globalLogger = newBaseLoggerWithoutLabel()
)

func SetLabel(label string) {
	globalLogger.label = label
}

func Test(args ...any) {
	globalLogger.Test(args...)
}

func Debug(args ...any) {
	globalLogger.Debug(args...)
}

func Info(args ...any) {
	globalLogger.Info(args...)
}

func Warn(args ...any) {
	globalLogger.Warn(args...)
}

func Error(args ...any) {
	globalLogger.Error(args...)
}

func Fatal(args ...any) {
	globalLogger.Fatal(args...)
}

func Testf(format string, args ...any) {
	globalLogger.Testf(format, args...)
}

func Debugf(format string, args ...any) {
	globalLogger.Debugf(format, args...)
}

func Infof(format string, args ...any) {
	globalLogger.Infof(format, args...)
}

func Warnf(format string, args ...any) {
	globalLogger.Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	globalLogger.Errorf(format, args...)
}

func Fatalf(format string, args ...any) {
	globalLogger.Fatalf(format, args...)
}
