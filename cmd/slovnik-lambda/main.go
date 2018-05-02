package main

import (
	"github.com/rpeshkov/slovnik/seznam"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rpeshkov/slovnik"
)

// Request defines lambda function request
type Request struct {
	// Word contains the word that need to be translated
	Word string `json:"word"`
}

func translate(request Request) ([]*slovnik.Word, error) {
	translator := seznam.NewTranslator()
	lang := slovnik.DetectLanguage(request.Word)
	return translator.Translate(request.Word, lang)
}

func main() {
	lambda.Start(translate)
}
