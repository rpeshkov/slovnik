package main

import (
	"log"
)

func main() {
	config, err := InitConfig()

	if err != nil {
		log.Panic(err)
	}

	templates, err := CreateTemplate()
	if err != nil {
		log.Panic(err)
	}

	bot, err := NewBot(config, templates)

	if err != nil {
		log.Panic(err)
	}

	bot.Listen()

}
