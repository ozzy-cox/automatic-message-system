package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

type Logger struct {
	*log.Logger
}

func NewLogger(config Config) (*Logger, error) {
	writers := []io.Writer{}

	if config.LogFile != "" {
		logDir := filepath.Dir(config.LogFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, err
		}

		file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		writers = append(writers, file)
	}

	if config.LogToStdout {
		writers = append(writers, os.Stdout)
	}

	var writer io.Writer
	if len(writers) > 1 {
		writer = io.MultiWriter(writers...)
	} else if len(writers) == 1 {
		writer = writers[0]
	} else {
		writer = os.Stdout
	}

	logger := log.New(writer, "", log.LstdFlags|log.Lshortfile)

	return &Logger{logger}, nil
}
