package log

type Level string

const (
	NIL   = "NIL"
	TEST  = "TEST"
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
	FATAL = "FATAL"
)

func getLevelColor(level Level) Color {
	switch level {
	case NIL:
		return WHITE
	case TEST:
		return BLUE
	case DEBUG:
		return GREEN
	case INFO:
		return CYAN
	case WARN:
		return YELLOW
	case ERROR:
		return RED
	case FATAL:
		return RED_BOLD
	}
	return ""
}
