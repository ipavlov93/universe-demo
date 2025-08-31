package std

import (
	"io"
	"log"
)

type Logger struct {
	w      io.Writer
	logger *log.Logger
}

func NewWriterLogger(w io.Writer) *Logger {
	return &Logger{
		w:      w,
		logger: log.New(w, "", log.LstdFlags|log.LUTC),
	}
}

func (l Logger) Log(msg string) {
	l.logger.Println(msg)
}
