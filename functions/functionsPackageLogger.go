package functions

import (
	"errors"
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	path := "./.log/"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	file, err := openLogFile("./.log/log.log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Flags are DATE , TIME, Name Element and Line Number
	InfoLogger = log.New(file, "INFO: ", log.LstdFlags|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.LstdFlags|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.LstdFlags|log.Lshortfile)
	/*
		EXAMPLE HOW TO USE
		InfoLogger.Printf("User Created: %s", userName)
		WarningLogger.Println("You did something dangerous.")
		ErrorLogger.Println(err())
	*/
}

func openLogFile(path string) (*os.File, error) {

	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}
