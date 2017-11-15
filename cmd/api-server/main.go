package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rpeshkov/slovnik/seznam"

	"github.com/gorilla/handlers"
	"github.com/rpeshkov/slovnik"

	"github.com/gorilla/mux"
)

func main() {
	translator := seznam.NewTranslator()

	router := mux.NewRouter().StrictSlash(true)

	router.
		Methods(http.MethodGet).
		Path("/api/translate/{word}").
		HandlerFunc(translate(translator))

	cors := handlers.CORS()
	err := http.ListenAndServe(":8080", cors(router))

	if err != nil {
		log.Fatal(err)
	}
}

func translate(translator slovnik.Translator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		word := vars["word"]

		lang := slovnik.DetectLanguage(word)
		translations, err := translator.Translate(word, lang)

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		err = json.NewEncoder(w).Encode(translations)

		if err != nil {
			fmt.Fprintln(w, err)
		}
	}
}
