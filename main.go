package main

import (
	"kaero/application"
	"log"
)

const (
	version = "0.0.0-dev"
)

func main() {
	app, err := application.New(version)
	if err != nil {
		log.Fatalln(err.Error())
	}

	defer func() {
		if x := recover(); x != nil {
			app.Stop()
			panic(x)
		}
	}()
	err = app.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
