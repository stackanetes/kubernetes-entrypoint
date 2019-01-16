package logger

import (
	"log"
	"os"
	"strings"
)

var (
	//InfoLog logger
	InfoLog *log.Logger
	//ErrorLog logger
	ErrorLog *log.Logger
	//WarningLog logger
	WarningLog *log.Logger

	//OutputJSON is used to determine whether to print human readable messages or machine reable JSON
	OutputJSON bool
)

func init() {
	InfoLog = log.New(os.Stderr, "Entrypoint INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog = log.New(os.Stderr, "Entrypoint WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(os.Stderr, "Entrypoint ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	OutputJSON = strings.EqualFold(os.Getenv("OUTPUT_JSON"), "true")
}

// Info is a thin wrapper over the builtin logger's Info()
func Info(formatString string, args ...interface{}) {
	if !OutputJSON {
		InfoLog.Printf(formatString, args...)
	}
}

// Warning is a thin wrapper over the builtin logger's Info()
func Warning(formatString string, args ...interface{}) {
	if !OutputJSON {
		WarningLog.Printf(formatString, args...)
	}
}

// Error is a thin wrapper over the builtin logger's Info()
func Error(formatString string, args ...interface{}) {
	if !OutputJSON {
		ErrorLog.Printf(formatString, args...)
	}
}
