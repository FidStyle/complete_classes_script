package logger

import (
	"log"
	"os"
)

var (
	accessLogger *log.Logger
	errorLogger  *log.Logger
)

func createOutputFile(file string) (*os.File, error) {
	output, err := os.OpenFile(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// the return val should be closed later
func getOutput(accessLog, errorLog string) (*os.File, *os.File, error) {
	accessOutput, err := createOutputFile(accessLog)
	if err != nil {
		return nil, nil, err
	}
	errorOutput, err := createOutputFile(errorLog)
	if err != nil {
		return nil, nil, err
	}

	return accessOutput, errorOutput, nil
}

func Init(accessLog string, errorLog string) {
	accessOutput, errorOutput, err := getOutput(accessLog, errorLog)
	if err != nil {
		panic(err)
	}

	accessLogger = log.New(accessOutput, "", log.LstdFlags)
	errorLogger = log.New(errorOutput, "[Error] ", log.LstdFlags)
}

func Infof(format string, a ...any) {
	accessLogger.Printf(format, a...)
}

func Info(a ...any) {
	accessLogger.Println(a...)
}

func Errorf(format string, a ...any) {
	errorLogger.Printf(format, a...)
}

func Error(a ...any) {
	errorLogger.Println(a...)
}
