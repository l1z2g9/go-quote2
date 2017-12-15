package util

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func init() {
	Info = log.New(ioutil.Discard,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(ioutil.Discard,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func ShowLog() {
	writer, err := os.OpenFile("web.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	//defer f.Close()
	Info = log.New(writer,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(writer,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
