package main

import (
	"log"
)

func main() {
	// c, err := NewSlovnikClient("http://localhost:8080", nil)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// words, err := c.Translate("hlavni", slovnik.Cz)

	// if err != nil {
	// 	log.Panic(err)
	// }

	// for _, w := range words {
	// 	fmt.Println(w)
	// }
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
