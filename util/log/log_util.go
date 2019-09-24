// logUtil package is not going to replace the Go default log package. it is just a thin wrapper on top of it to handle different log level but underlying will still call Go log package
package logUtil

import (
	"log"
)

const (
	DEBUG = 0
	INFO  = 1
	WARN  = 2
	ERROR = 3
	FATAL = 4
)

// ValidLogLevel contain all the supported log levels.
var ValidLogLevel = map[string]int{
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"WARN":  WARN,
	"ERROR": ERROR,
	"FATAL": FATAL,
}

var logLevel = INFO

// SetLevel. level parameter valid values come from the const DEBUG,INFO,WARN,ERROR,FATAL.
//
//for performance call this method ONCE upon server startup and use it for the rest of the method calls. this is to avoid the sync.RWMutex overhead
func SetLevel(level int) {
	logLevel = level
}

// IsDebugEnabled return if the logLevel has been set to DEBUG
func IsDebugEnabled() bool {
	return logLevel == DEBUG
}

// DebugPrintf. if logLevel >= DEBUG call Go log.Printf(...)
func DebugPrintf(format string, v ...interface{}) {
	if DEBUG >= logLevel {
		log.Printf(format, v...)
	}
}

// DebugPrint. if logLevel >= DEBUG call Go log.Print(...)
func DebugPrint(v ...interface{}) {
	if DEBUG >= logLevel {
		log.Print(v...)
	}
}

// DebugPrintln. if logLevel >= DEBUG call Go log.Println(...)
func DebugPrintln(v ...interface{}) {
	if DEBUG >= logLevel {
		log.Println(v...)
	}
}

// InfoPrintf. if logLevel >= INFO call Go log.Printf(...)
func InfoPrintf(format string, v ...interface{}) {
	if INFO >= logLevel {
		log.Printf(format, v...)
	}
}

// InfoPrint. if logLevel >= INFO call Go log.Print(...)
func InfoPrint(v ...interface{}) {
	if INFO >= logLevel {
		log.Print(v...)
	}
}

// InfoPrintln. if logLevel >= INFO call Go log.Println(...)
func InfoPrintln(v ...interface{}) {
	if INFO >= logLevel {
		log.Println(v...)
	}
}

// WarnPrintf. if logLevel >= WARN call Go log.Printf(...)
func WarnPrintf(format string, v ...interface{}) {
	if WARN >= logLevel {
		log.Printf(format, v...)
	}
}

// WarnPrint. if logLevel >= WARN call Go log.Print(...)
func WarnPrint(v ...interface{}) {
	if WARN >= logLevel {
		log.Print(v...)
	}
}

// WarnPrintln. if logLevel >= WARN call Go log.Println(...)
func WarnPrintln(v ...interface{}) {
	if WARN >= logLevel {
		log.Println(v...)
	}
}

// ErrorPrintf. if logLevel >= ERROR call Go log.Printf(...)
func ErrorPrintf(format string, v ...interface{}) {
	if ERROR >= logLevel {
		log.Printf(format, v...)
	}
}

// ErrorPrint. if logLevel >= ERROR call Go log.Print(...)
func ErrorPrint(v ...interface{}) {
	if ERROR >= logLevel {
		log.Print(v...)
	}
}

// ErrorPrintln. if logLevel >= ERROR call Go log.Println(...)
func ErrorPrintln(v ...interface{}) {
	if ERROR >= logLevel {
		log.Println(v...)
	}
}

// FatalPrintf. if logLevel >= FATAL call Go log.Fatalf(...)
func FatalPrintf(format string, v ...interface{}) {
	if FATAL >= logLevel {
		log.Fatalf(format, v...)
	}
}

// FatalPrint. if logLevel >= FATAL call Go log.Fatal(...)
func FatalPrint(v ...interface{}) {
	if FATAL >= logLevel {
		log.Fatal(v...)
	}
}

// FatalPrintln. if logLevel >= FATAL call Go log.Fatalln(...)
func FatalPrintln(v ...interface{}) {
	if FATAL >= logLevel {
		log.Fatalln(v...)
	}
}

// Panicf. pass through to call Go log.Panicf(...)
func Panicf(format string, v ...interface{}) {
	log.Panicf(format, v...)
}

// Panic. pass through to call Go log.Panic(...)
func Panic(v ...interface{}) {
	log.Panic(v...)
}

// Panicln. pass through to call Go log.Panicln(...)
func Panicln(v ...interface{}) {
	log.Panicln(v...)
}
