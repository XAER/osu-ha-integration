package osu

import "log"

type Logger interface {
	Info(msg string)
}

type StdLogger struct{}

func (l *StdLogger) Info(msg string) {
	log.Printf("[osu-cache] %s", msg)
}
