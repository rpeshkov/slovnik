package seznam

import (
	"github.com/rpeshkov/slovnik"
)

// Translator represents seznam translator type
type Translator struct {
	client *Client
	parser *Parser
}

// NewTranslator creates new seznam translator instance
func NewTranslator() *Translator {
	return &Translator{
		client: NewClient(nil),
		parser: NewParser(),
	}
}

// Translate translates provided word and returns results
func (t *Translator) Translate(word string, language slovnik.Language) ([]*slovnik.Word, error) {
	body, err := t.client.Get(word, language)

	if err != nil {
		return nil, err
	}

	defer body.Close()

	return t.parser.Parse(body)
}
