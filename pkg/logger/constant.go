package logger

type LogLevel int8

const (
	LOG_LEVEL_DEBUG LogLevel = iota - 1
	LOG_LEVEL_INFO
	LOG_LEVEL_WARNING
	LOG_LEVEL_ERROR
	LOG_LEVEL_DPANIC
	LOG_LEVEL_PANIC
	LOG_LEVEL_FATAL

	minLogLevel       = LOG_LEVEL_DEBUG
	maxLogLevel       = LOG_LEVEL_FATAL
	INVALID_LOG_LEVEL = maxLogLevel + 1
)
