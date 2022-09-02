package loggr

import (
	"log"
)

type Logger struct{}

func NewLogger() *Logger {
	return new(Logger)
}

func (logger *Logger) Log(a ...interface{}) {
	log.Println(a...)
}

func (logger *Logger) Debug(a ...interface{}) {
	log.Println(a...)
}

func (logger *Logger) Fatal(a ...interface{}) {
	log.Fatal(a...)
}
