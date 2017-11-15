package slovnik

import "strings"

// Language of the input string
type Language int

const (
	// Ru represents Russian language
	Ru Language = iota
	// Cz represents Czech language
	Cz
)

// Russian alphabet
const rusSymbols = "абвгдеёжзийклмнопрстуфхцчшщьыъэюя"

// DetectLanguage used to find out which language is used for the input string
func DetectLanguage(input string) Language {
	for _, ch := range input {
		if strings.Contains(rusSymbols, strings.ToLower(string(ch))) {
			return Ru
		}
	}
	return Cz
}
