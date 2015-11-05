package conf2

import (
	"os"
	"log"
)

// Don't send enything to this logger unless it's an error that deserves immediate attention
var Err *log.Logger

// Warnings, or general information.
var Info *log.Logger

func init() {
	Err = log.New(os.Stderr, "", log.Ldate | log.Ltime | log.Lshortfile)
	Info = log.New(os.Stdout, "", log.Ldate | log.Ltime)
}
