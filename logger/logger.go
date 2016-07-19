package logger

import (
	"log"
	"os"
)

var (
	//"Info logger""
	Info *log.Logger
	//"Error logger"
	Error *log.Logger
	//Warning logger
	Warning *log.Logger
)

func init() {
	Info = log.New(os.Stdout, "Entrypoint INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "Entrypoint WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "Entrypoint Error: ", log.Ldate|log.Ltime|log.Lshortfile)
}
