package log

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Logger struct {
	label  string
	config loggerConfig
}

type loggerConfig struct {
	labelMaxLen      int
	levelMaxLen      int
	callStackMaxSize int
	callStackMaxLen  int
}

func defaultLoggerConfig() loggerConfig {
	return loggerConfig{
		labelMaxLen:      12,
		levelMaxLen:      6,
		callStackMaxSize: 2,
		callStackMaxLen:  30,
	}
}

func newLogger(label string) Logger {
	l := new(Logger)
	l.label = label
	l.config = defaultLoggerConfig()
	return *l
}

func newLoggerWithoutLabel() Logger {
	l := new(Logger)
	l.config = defaultLoggerConfig()
	return *l
}

func (l *Logger) write(p []byte) (n int, err error) {
	str := string(p[:])
	reg := regexp.MustCompile(`(\n\n+)|(\s+)$`)
	str = reg.ReplaceAllString(str, "")
	str = strings.ReplaceAll(str, "\r", "")
	l.preComposeLog(NIL, PrettyMarshal, false, str)
	return len(p), nil
}

func (l *Logger) Test(args ...any) {
	l.preComposeLog(TEST, PrettyMarshal, false, args...)
}

func (l *Logger) Debug(args ...any) {
	l.preComposeLog(DEBUG, PrettyMarshal, false, args...)
}

func (l *Logger) Info(args ...any) {
	l.preComposeLog(INFO, Marshal, false, args...)
}

func (l *Logger) Warn(args ...any) {
	l.preComposeLog(WARN, Marshal, false, args...)
}

func (l *Logger) Error(args ...any) {
	l.preComposeLog(ERROR, Marshal, false, args...)
}

func (l *Logger) Fatal(args ...any) {
	l.preComposeLog(FATAL, Marshal, false, args...)
}

func (l *Logger) Testf(format string, args ...any) {
	l.preComposeLog(TEST, PrettyMarshal, true, format, args)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.preComposeLog(DEBUG, PrettyMarshal, true, format, args)
}

func (l *Logger) Infof(format string, args ...any) {
	l.preComposeLog(INFO, Marshal, true, format, args)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.preComposeLog(WARN, Marshal, true, format, args)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.preComposeLog(ERROR, Marshal, true, format, args)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.preComposeLog(FATAL, Marshal, true, format, args)
}

func (l *Logger) preComposeLog(level Level, marshalFunc MarshalFunc, isFormatMode bool, args ...any) {
	if isFormatMode {
		format := args[0].(string)
		newArgs := args[1].([]any)
		contents := fmt.Sprintf(format, newArgs...)
		l.composeLog(level, marshalFunc, []any{contents})
	} else {
		l.composeLog(level, marshalFunc, args)
	}
}

func (l *Logger) composeLog(level Level, marshalFunc MarshalFunc, contents []any) {
	// Context Tag
	label := l.label
	if label == "" {
		label = "GLOBAL"
	}

	if len(l.label) > l.config.labelMaxLen {
		label = l.label[:l.config.labelMaxLen-3] + "..."
	}

	contextSeg := sprintcf("[%"+strconv.Itoa(l.config.labelMaxLen)+"s]", CYAN, label)

	// Time Tag
	loc, _ := time.LoadLocation("Asia/Seoul")
	t := time.Now().In(loc)
	timeSeg := fmt.Sprintf("%d.%02d.%02d %02d:%02d:%02d.%06d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000)

	// Level Tag
	levelColor := getLevelColor(level)
	levelSeg := sprintcf("%6s", levelColor, string(level))

	// Call Stack Tag
	maxFileSegmentLen := l.config.callStackMaxLen
	callStackSeg := l.getCallStackSegment(l.config.callStackMaxSize)
	if len(callStackSeg) > maxFileSegmentLen {
		firstDotIndex := strings.Index(callStackSeg, ".")
		callStackSeg = callStackSeg[firstDotIndex+1:]
		callStackSeg = "..." + callStackSeg[MaxInt(len(callStackSeg)-maxFileSegmentLen, 0):]
	}
	callStackSeg = sprintcf("%-"+strconv.Itoa(maxFileSegmentLen)+"s", YELLOW, callStackSeg)

	// Build String
	builtSeg := l.buildLog([]string{contextSeg, timeSeg, levelSeg, callStackSeg, ":"}, contents, marshalFunc)
	println(builtSeg)
}

func (l *Logger) buildLog(segments []string, contents []any, marshalFunc MarshalFunc) string {
	stringBuilder := NewStringBuilder()
	stringBuilder.SetMarshaller(marshalFunc)

	for _, segment := range segments {
		stringBuilder.Append(segment).Space()
	}

	for _, content := range contents {
		stringBuilder.Append(content).Space()
	}

	return stringBuilder.Build()
}

func (l *Logger) getCallStackSegment(stack int) string {
	pcs := make([]uintptr, 10)
	n := runtime.Callers(6, pcs)
	pcs = pcs[:n]
	frames := runtime.CallersFrames(pcs)

	var callStackSegments []string
	var slicedFrames []runtime.Frame

	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		slicedFrames = append(slicedFrames, frame)
	}

	// Count of frames to slice of end of stacks
	slicedFrames = l.reverseFrameSlice(slicedFrames)
	slicedFrames = slicedFrames[MaxInt(len(slicedFrames)-stack, 0):]

	for _, frame := range slicedFrames {
		filepath := frame.File
		filepathSegmentsByDot := strings.Split(filepath, ".")
		extensionIndex := MaxInt(len(filepathSegmentsByDot)-2, 0)
		filepathExceptExtension := filepathSegmentsByDot[extensionIndex]
		filepathSegmentsBySlash := strings.Split(filepathExceptExtension, "/")

		fileSeg := filepathSegmentsBySlash[len(filepathSegmentsBySlash)-1]
		lineNumSeg := strconv.Itoa(frame.Line)

		callStackSegments = append(callStackSegments, fileSeg+"("+lineNumSeg+")")
	}

	return strings.Join(callStackSegments, ".")
}

func (l *Logger) reverseFrameSlice(slice []runtime.Frame) []runtime.Frame {
	var reversed []runtime.Frame
	sliceLen := len(slice)
	for i := 0; i < sliceLen; i++ {
		reversed = append(reversed, slice[sliceLen-1-i])
	}
	return reversed
}
