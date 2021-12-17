package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

var log = logrus.New()

func init() {
	//TODO: next add write to file
	if !debug {
		log.SetOutput(ioutil.Discard)
	} else {

		f, err := os.OpenFile("rtsptoweb.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Printf("error opening file: %v", err)
		}
		// don't forget to close it
		// defer f.Close()

		// Output to stderr instead of stdout, could also be a file.
		log.SetOutput(f)

	}
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(Storage.ServerLogLevel())

}
