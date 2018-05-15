package seznam

import (
	"io"

	"github.com/PuerkitoBio/goquery"
	"github.com/rpeshkov/slovnik"

	"strings"
)

const (
	synonymsHeader     = "Synonyma"
	antonymsHeader     = "Antonyma"
	derivedWordsHeader = "Odvozená slova"
)

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
	doc, err := goquery.NewDocumentFromReader(pageBody)

	if err != nil {
		return nil, err
	}

	resultsNode := doc.Find("#results")
	results := []*slovnik.Word{}

	if class, ok := resultsNode.Attr("class"); ok && class == "transl" {
		w, err := parseOne(resultsNode)
		if err != nil {
			return nil, err
		}
		results = append(results, &w)
	} else {
		r, err := parseMultiple(resultsNode)
		if err != nil {
			return nil, err
		}
		results = r
	}

	return results, nil
}

// parseMultiple parses results node for multiple results. Multiple results are present when
// mistyped word was provided for translation
func parseMultiple(resultsNode *goquery.Selection) ([]*slovnik.Word, error) {
	result := []*slovnik.Word{}
	resultsNode.Find(".mistype li").Each(func(i int, s *goquery.Selection) {
		result = append(result, &slovnik.Word{
			Word:         strings.TrimSpace(s.Find("a").Text()),
			Translations: []string{strings.TrimSpace(s.Find("span").Text())},
		})
	})
	return result, nil
}

// parseOne parses full page of translation result and returns Word structure filled by data from page
func parseOne(resultsNode *goquery.Selection) (slovnik.Word, error) {
	w := slovnik.Word{}

	w.Word = resultsNode.Find("h1").First().Text()

	transNode := resultsNode.Find("#fastTrans #fastMeanings")

	tempTrans := ""
	transNode.Children().Each(func(i int, s *goquery.Selection) {
		if s.Is("br") || s.Is(".comma") {
			w.Translations = append(w.Translations, strings.TrimSpace(tempTrans))
			tempTrans = ""
		} else {
			tempTrans = tempTrans + " " + s.Text()
		}
	})

	w.WordType = resultsNode.Find(".morfLinks .morf").Text()

	resultsNode.Find(".hgroup").Each(func(i int, s *goquery.Selection) {
		items := strings.Split(s.NextFiltered(".other-meaning").Text(), ",")
		items = forEachString(items, strings.TrimSpace)

		head := strings.TrimSpace(s.Text())

		switch head {
		case synonymsHeader:
			w.Synonyms = items
		case antonymsHeader:
			w.Antonyms = items
		case derivedWordsHeader:
			w.DerivedWords = items
		}
	})

	resultsNode.Find("ul#fulltext li").Each(func(i int, s *goquery.Selection) {
		sample := slovnik.SampleUse{}
		sample.Keyword = strings.TrimSpace(s.Find("a").Text())
		sample.Phrase = strings.TrimSpace(s.Find(".bold").Text())

		s.Children().Remove()
		sample.Translation = strings.Trim(s.Text(), " :")
		sample.Translation = strings.TrimSpace(sample.Translation)
		w.Samples = append(w.Samples, sample)
	})
	return w, nil
}

func forEachString(items []string, fn func(s string) string) (result []string) {
	for _, i := range items {
		result = append(result, fn(i))
	}
	return
}
