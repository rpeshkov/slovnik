package seznam

import (
	"github.com/rpeshkov/slovnik"
)

type Translator struct {
	client *Client
	parser *Parser
}

func NewTranslator() *Translator {
	return &Translator{
		client: NewClient(nil),
		parser: NewParser(),
	}
}

func (t *Translator) Translate(word string, language slovnik.Language) ([]*slovnik.Word, error) {
	body, err := t.client.Get(word, language)

	if err != nil {
		return nil, err
	}

	defer body.Close()

	return t.parser.Parse(body)
}
