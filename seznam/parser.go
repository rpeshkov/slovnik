package seznam

import (
	"bytes"
	"fmt"
	"io"

	"github.com/rpeshkov/slovnik"

	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	synonymsHeader     = "Synonyma"
	antonymsHeader     = "Antonyma"
	derivedWordsHeader = "Odvozená slova"

	otherMeaningClass = "other-meaning"
	fastMeaningsID    = "fastMeanings"
)

// funcs map for parsing blocks
var funcs = map[string]func(*slovnik.Word, string){
	synonymsHeader:     func(w *slovnik.Word, data string) { w.Synonyms = append(w.Synonyms, data) },
	antonymsHeader:     func(w *slovnik.Word, data string) { w.Antonyms = append(w.Antonyms, data) },
	derivedWordsHeader: func(w *slovnik.Word, data string) { w.DerivedWords = append(w.DerivedWords, data) },
}

// Parser is a parser for slovnik.seznam.cz portal
type Parser struct{}

// NewParser creates new parser for seznam translation page
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses html structure from provided reader and scrap the data from it.
// There may be 2 types of pages:
// - Page with full description for specific word, like: https://slovnik.seznam.cz/cz-ru/?q=hlavn%C3%AD
// - Mistype page that shows you that you might have spelled word wrong and showing which word it can
//   actually be, like: https://slovnik.seznam.cz/cz-ru/?q=dobr
func (s *Parser) Parse(pageBody io.Reader) ([]*slovnik.Word, error) {
	doc, err := html.Parse(pageBody)

	if err != nil {
		return nil, err
	}

	results, ok := getResultsNode(doc)

	if !ok {
		return nil, fmt.Errorf("results node not found")
	}

	buf := new(bytes.Buffer)
	err = html.Render(buf, results)

	if err != nil {
		return nil, err
	}

	tokenizer := html.NewTokenizer(buf)
	attrs := attributes(results.Attr)

	if attrs.class() == "transl" {
		return processSingleWord(tokenizer), nil
	}

	return processMistype(tokenizer), nil
}

// getResultsNode parses page HTML to find node containing results of translation
func getResultsNode(document *html.Node) (results *html.Node, ok bool) {
	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.DataAtom == atom.Div && attributes(c.Attr).id() == "results" {
				results = c
				ok = true
				return
			}
			traverse(c)
		}
	}
	traverse(document)
	return
}

// processMistype is used for processing mistyped results
func processMistype(z *html.Tokenizer) (result []*slovnik.Word) {
	var w *slovnik.Word

	for z.Next() != html.ErrorToken {
		token := z.Token()

		if token.Type == html.StartTagToken && token.DataAtom != atom.Li {
			z.Next() // Advance to text
			textToken := z.Token()

			switch token.DataAtom {
			case atom.A:
				w = &slovnik.Word{Word: textToken.Data}
				result = append(result, w)
			case atom.Span:
				addTranslation(w, textToken.Data)
			}
		}
	}
	return
}

// processSingleWord is used for parsing full translation of the word with samples of synonyms, antonyms, etc.
func processSingleWord(z *html.Tokenizer) (result []*slovnik.Word) {
	blockName := ""

	var w *slovnik.Word
	var prevToken html.Token

	for z.Next() != html.ErrorToken {
		token := z.Token()
		attrs := attributes(token.Attr)

		switch {
		case token.Type == html.StartTagToken:
			if attrs.class() == otherMeaningClass {
				if f, ok := funcs[blockName]; ok {
					processBlock(z, w, f)
					blockName = ""
				}
			}

			if token.DataAtom == atom.Ul && attrs.id() == "fulltext" {
				sample := processSample(z)
				w.Samples = append(w.Samples, sample)
			}

			if attrs.id() == fastMeaningsID {
				processTranslations(z, w)
			}

		case token.Type == html.TextToken:
			if prevToken.DataAtom == atom.H1 && attributes(prevToken.Attr).lang() != "" {
				w = &slovnik.Word{Word: token.Data}
				result = append(result, w)
			}

			if prevToken.DataAtom == atom.Span && attributes(prevToken.Attr).class() == "morf" {
				w.WordType = token.Data
			}

			if prevToken.DataAtom == atom.P && attributes(prevToken.Attr).class() == "morf" {
				blockName = token.Data
			}
		}
		prevToken = token
	}
	return
}

func processTranslations(z *html.Tokenizer, w *slovnik.Word) {
	var prevClosingToken html.Token

	for z.Next() != html.ErrorToken {
		token := z.Token()

		if token.DataAtom == atom.Div {
			return
		}

		if token.Type == html.EndTagToken || token.Type == html.SelfClosingTagToken {
			prevClosingToken = token
		} else if token.Type == html.StartTagToken && attributes(token.Attr).class() != "comma" {
			z.Next() // Advance to text
			textToken := z.Token()
			trimmed := strings.TrimSpace(textToken.Data)
			if len(trimmed) > 0 {
				if prevClosingToken.DataAtom == atom.A {
					updateLastTranslation(w, trimmed)
				} else {
					addTranslation(w, trimmed)
				}
			}
		}
	}
}

// processBlock is used for parsing blocks of synonyms, antonyms etc.
//
// Here's the sample of synonyms block:
// <div class="hgroup"> <p class="morf">Synonyma</p> </div>
// <div class="other-meaning">
//   <a lang="cs" href="/cz-ru/?q=ústřední">ústřední</a>,
//   <a lang="cs" href="/cz-ru/?q=podstatný">podstatný</a>,
//   <a lang="cs" href="/cz-ru/?q=základní">základní</a>,
//   <a lang="cs" href="/cz-ru/?q=zásadní">zásadní</a>
// </div>
func processBlock(z *html.Tokenizer, w *slovnik.Word, functor func(*slovnik.Word, string)) {
	var prevToken html.Token

	for z.Next() != html.ErrorToken {
		token := z.Token()

		if token.Type == html.EndTagToken && token.DataAtom == atom.Div {
			return
		}

		if prevToken.Type == html.StartTagToken && prevToken.DataAtom == atom.A {
			functor(w, token.Data)
		}

		prevToken = token
	}
	return
}

func processSample(z *html.Tokenizer) (result slovnik.SampleUse) {
	var prevToken html.Token

	for z.Next() != html.ErrorToken {
		token := z.Token()

		if token.DataAtom == atom.Ul && token.Type == html.EndTagToken {
			return
		}

		if prevToken.DataAtom == atom.A && prevToken.Type == html.StartTagToken {
			result.Keyword = strings.TrimSpace(token.Data)
		}

		if prevToken.DataAtom == atom.Span && prevToken.Type == html.StartTagToken && len(result.Phrase) == 0 {
			result.Phrase = strings.TrimSpace(token.Data)
		}

		if token.DataAtom == atom.Li && token.Type == html.EndTagToken {
			result.Translation = strings.TrimSpace(prevToken.Data)
		}

		prevToken = token
	}

	return
}

func addTranslation(w *slovnik.Word, data string) {
	w.Translations = append(w.Translations, data)
}

func updateLastTranslation(w *slovnik.Word, data string) {
	if len(w.Translations) > 0 {
		lastTranslation := w.Translations[len(w.Translations)-1]
		lastTranslation = lastTranslation + " " + data
		w.Translations[len(w.Translations)-1] = lastTranslation
	} else {
		addTranslation(w, data)
	}
}
