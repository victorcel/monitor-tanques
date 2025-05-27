package logger

import (
	"log"
	"os"
)

// Logger define la interfaz para el logging en la aplicación
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
}

// SimpleLogger implementa Logger con funcionalidad básica de logging
type SimpleLogger struct {
	debugLog *log.Logger
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
	fatalLog *log.Logger
}

// NewSimpleLogger crea una nueva instancia del logger simple
func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		debugLog: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		warnLog:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		fatalLog: log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// formatKeyValues formatea los pares clave-valor para el logging
func formatKeyValues(keysAndValues ...interface{}) string {
	if len(keysAndValues) == 0 {
		return ""
	}

	result := " ["
	for i := 0; i < len(keysAndValues); i += 2 {
		if i > 0 {
			result += ", "
		}

		// Añadimos la clave
		if i < len(keysAndValues) {
			result += "%v"
		}

		// Añadimos el valor si existe
		if i+1 < len(keysAndValues) {
			result += "=%v"
		}
	}
	result += "]"

	return result
}

// Debug registra un mensaje de nivel debug
func (l *SimpleLogger) Debug(msg string, keysAndValues ...interface{}) {
	format := msg + formatKeyValues(keysAndValues...)
	l.debugLog.Printf(format, keysAndValues...)
}

// Info registra un mensaje de nivel info
func (l *SimpleLogger) Info(msg string, keysAndValues ...interface{}) {
	format := msg + formatKeyValues(keysAndValues...)
	l.infoLog.Printf(format, keysAndValues...)
}

// Warn registra un mensaje de nivel warn
func (l *SimpleLogger) Warn(msg string, keysAndValues ...interface{}) {
	format := msg + formatKeyValues(keysAndValues...)
	l.warnLog.Printf(format, keysAndValues...)
}

// Error registra un mensaje de nivel error
func (l *SimpleLogger) Error(msg string, keysAndValues ...interface{}) {
	format := msg + formatKeyValues(keysAndValues...)
	l.errorLog.Printf(format, keysAndValues...)
}

// Fatal registra un mensaje de nivel fatal y termina la aplicación
func (l *SimpleLogger) Fatal(msg string, keysAndValues ...interface{}) {
	format := msg + formatKeyValues(keysAndValues...)
	l.fatalLog.Printf(format, keysAndValues...)
	os.Exit(1)
}
