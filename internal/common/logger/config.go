package logger

type Config struct {
	LogFile     string // Path to log file
	LogToStdout bool   // Whether to also log to stdout
}
