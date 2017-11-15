package slovnik

type Translator interface {
	Translate(word string, language Language) ([]*Word, error)
}
