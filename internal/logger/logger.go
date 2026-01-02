package logger

import "log"

func Init(appName string) {
	log.SetPrefix("[" + appName + "] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
