package main

import (
	"QqVideo/config"
	"QqVideo/engine"
	"log"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalln("config loading err,err:", err.Error())
	}
	engine.GoTask() // run task
}
