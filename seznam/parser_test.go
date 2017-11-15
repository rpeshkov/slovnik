package seznam_test

import (
	"os"
	"testing"

	"github.com/rpeshkov/slovnik"
	"github.com/rpeshkov/slovnik/seznam"
)

func TestParsePage(t *testing.T) {
	f, _ := os.Open("./test/sample.html")
	parser := seznam.NewParser()
	result, _ := parser.Parse(f)

	const expectedWord = "hlavní"
	w := result[0]

	if w.Word != expectedWord {
		t.Errorf("ParsePage word == %q, want %q", w.Word, expectedWord)
	}

	expectedTranslations := []string{
		"гла́вный",
		"основно́й",
		"центра́льный",
	}

	if len(w.Translations) != len(expectedTranslations) {
		t.Errorf("ParsePage len(translation) == %d, want %d", len(w.Translations), len(expectedTranslations))
	}

	for i, trans := range w.Translations {
		if trans != expectedTranslations[i] {
			t.Errorf("ParsePage translation == %q, want %q", trans, expectedTranslations[i])
		}
	}

	const expectedWordType = "přídavné jméno"
	if w.WordType != expectedWordType {
		t.Errorf("ParsePage wordType == %q, want %q", w.WordType, expectedWordType)
	}

	expectedSynonyms := []string{
		"ústřední",
		"podstatný",
		"základní",
		"zásadní",
	}

	if len(w.Synonyms) != len(expectedSynonyms) {
		t.Errorf("ParsePage len(synonyms) == %d, want %d", len(w.Synonyms), len(expectedSynonyms))
	}

	for i, synonym := range w.Synonyms {
		if synonym != expectedSynonyms[i] {
			t.Errorf("ParsePage synonym == %q, want %q", synonym, expectedSynonyms[i])
		}
	}

	expectedAntonyms := []string{
		"vedlejší",
		"podřadný",
		"podružný",
	}

	if len(w.Antonyms) != len(expectedAntonyms) {
		t.Errorf("ParsePage len(antonyms) == %d, want %d", len(w.Antonyms), len(expectedAntonyms))
	}

	for i, antonym := range w.Antonyms {
		if antonym != expectedAntonyms[i] {
			t.Errorf("ParsePage antonym == %q, want %q", antonym, expectedAntonyms[i])
		}
	}

	expectedDerivedWords := []string{
		"hlavně",
	}

	if len(w.DerivedWords) != len(expectedDerivedWords) {
		t.Errorf("ParsePage len(derivedWords) == %d, want %d", len(w.DerivedWords), len(expectedDerivedWords))
	}

	for i, derived := range w.DerivedWords {
		if derived != expectedDerivedWords[i] {
			t.Errorf("ParsePage derivedWord == %q, want %q", derived, expectedDerivedWords[i])
		}
	}

	if len(w.Samples) != 31 {
		t.Errorf("ParsePage len(Samples) == %d, want 31", len(w.Samples))
	}

	expectedSamples := []slovnik.SampleUse{
		{Keyword: "stan", Phrase: "hlavní stan", Translation: "velitelský ста́вка (главнокома́ндующего), штаб-кварти́ра"},
		{Keyword: "jídlo", Phrase: "hlavní jídlo", Translation: "основно́е/второ́е блю́до"},
		{Keyword: "město", Phrase: "hlavní město", Translation: "столи́ца"},
		{Keyword: "nádraží", Phrase: "hlavní nádraží", Translation: "гла́вный вокза́л"},
		{Keyword: "přízvuk", Phrase: "hlavní /vedlejší přízvuk", Translation: "гла́вное/второстепе́нное ударе́ние"},
		{Keyword: "rozhodčí", Phrase: "hlavní rozhodčí", Translation: "гла́вный арби́тр"},
		{Keyword: "rys", Phrase: "v hrubých/ hlavních rysech", Translation: "в о́бщих/гла́вных черта́х"},
		{Keyword: "silnice", Phrase: "hlavní /vedlejší silnice", Translation: "гла́вная/второстепе́нная доро́га"},
		{Keyword: "tah", Phrase: "hlavní tah", Translation: "гла́вная доро́га, магистра́ль"},
		{Keyword: "zeď", Phrase: "nosná/ hlavní /opěrná zeď", Translation: "несу́щая/капита́льная/опо́рная стена́"},
		{Keyword: "bod", Phrase: "hlavní bod jednání", Translation: "основно́й пункт перегово́ров"},
		{Keyword: "hlavně", Phrase: "Hlavně , že jsi přišel.", Translation: "Гла́вное, что ты здесь."},
		{Keyword: "hrát", Phrase: "hrát hlavní roli", Translation: "исполня́ть гла́вную роль"},
		{Keyword: "chod", Phrase: "hlavní chod", Translation: "горя́чее блюдо́"},
		{Keyword: "komunikace", Phrase: "hlavní /vedlejší komunikace", Translation: "гла́вная/второстепе́нная доро́га"},
		{Keyword: "mít", Phrase: "mít hlavní slovo/peníze", Translation: "име́ть гла́вное сло́во/де́ньги"},
		{Keyword: "paluba", Phrase: "hlavní paluba", Translation: "гла́вная па́луба"},
		{Keyword: "přenést", Phrase: "přenést hlavní město", Translation: "перенести́ столи́цу"},
		{Keyword: "role", Phrase: "hlavní role", Translation: "гла́вная роль"},
		{Keyword: "vchod", Phrase: "hlavní /vedlejší/zadní vchod", Translation: "пара́дный/боково́й/за́дний вход"},
		{Keyword: "výhra", Phrase: "hlavní výhra", Translation: "гла́вный вы́игрыш"},
		{Keyword: "žalobce", Phrase: "hlavní /státní/veřejný žalobce", Translation: "гла́вный/госуда́рственный/обще́ственный обвини́тель"},
		{Keyword: "второ́е", Phrase: "Co si dáte jako hlavní jídlo?", Translation: "Что зака́жете на второ́е?"},
		{Keyword: "основно́й", Phrase: "hlavní cíl", Translation: "основна́я цель"},
		{Keyword: "сезо́н", Phrase: "hlavní sezóna", Translation: "высо́кий сезо́н"},
		{Keyword: "столи́ца", Phrase: "Moskva je hlavním městem Ruska.", Translation: "Москва́ столи́ца Росси́и."},
		{Keyword: "столи́чный", Phrase: "obyvatelé hlavního města", Translation: "столи́чные жи́тели"},
		{Keyword: "центра́льный", Phrase: "hlavní ulice", Translation: "центра́льные у́лицы го́рода"},
		{Keyword: "центра́льный", Phrase: "hlavní nádraží", Translation: "центра́льный вокза́л"},
		{Keyword: "центра́льный", Phrase: "hlavní knihovna", Translation: "центра́льная библиоте́ка"},
		{Keyword: "скри́пка", Phrase: "hrát první housle, přen. mít hlavní slovo", Translation: "игра́ть пе́рвую скри́пку"},
	}

	for i, sample := range expectedSamples {
		if w.Samples[i] != sample {
			t.Errorf("ParsePage sample[%d]='%v', want '%v'", i, w.Samples[i], sample)
		}
	}
}

func TestParseAltPage(t *testing.T) {
	f, _ := os.Open("./test/sample_issue8.html")
	parser := seznam.NewParser()
	result, _ := parser.Parse(f)
	w := result[0]

	const expectedWord = "soutěživý"

	if w.Word != expectedWord {
		t.Errorf("ParsePage word == %q, want %q", w.Word, expectedWord)
	}

	expectedTranslations := []string{
		"состяза́тельный",
	}

	if len(w.Translations) != len(expectedTranslations) {
		t.Errorf("ParsePage len(translation) == %d, want %d", len(w.Translations), len(expectedTranslations))
	}

	for i, trans := range w.Translations {
		if trans != expectedTranslations[i] {
			t.Errorf("ParsePage translation == %q, want %q", trans, expectedTranslations[i])
		}
	}

	expectedDerivedWords := []string{
		"soutěživost",
	}

	if len(w.DerivedWords) != len(expectedDerivedWords) {
		t.Errorf("ParsePage len(derivedWords) == %d, want %d", len(w.DerivedWords), len(expectedDerivedWords))
		return
	}

	for i, derived := range w.DerivedWords {
		if derived != expectedDerivedWords[i] {
			t.Errorf("ParsePage derivedWord == %q, want %q", derived, expectedDerivedWords[i])
		}
	}
}

func TestParseIssue7(t *testing.T) {
	f, _ := os.Open("./test/sample_issue7.html")
	parser := seznam.NewParser()
	result, _ := parser.Parse(f)
	w := result[0]

	const expectedWord = "protože"

	if w.Word != expectedWord {
		t.Errorf("ParsePage word == %q, want %q", w.Word, expectedWord)
	}

	expectedTranslations := []string{
		"так как",
		"из-за того́",
		"потому́ что",
	}

	if len(w.Translations) != len(expectedTranslations) {
		t.Errorf("ParsePage len(translation) == %d, want %d", len(w.Translations), len(expectedTranslations))
		return
	}

	for i, trans := range w.Translations {
		if trans != expectedTranslations[i] {
			t.Errorf("ParsePage translation == %q, want %q", trans, expectedTranslations[i])
		}
	}
}

func TestMultipleResults(t *testing.T) {
	f, _ := os.Open("./test/sample_multiple_results.html")
	parser := seznam.NewParser()
	result, _ := parser.Parse(f)

	const expectedCount = 9

	if len(result) != expectedCount {
		t.Errorf("ParsePage len(w) == %d, want %d", len(result), expectedCount)

		return
	}

	expectedWords := []string{
		"dobrat se",
		"doba",
		"do",
		"dobrý",
		"dobro",
		"dobré",
		"dobrat",
		"obr",
		"bobr",
	}

	expectedTranslations := []string{
		"добра́ться",
		"вре́мя",
		"в",
		"хоро́ший",
		"добро́",
		"добро́",
		"израсхо́довать",
		"гига́нт",
		"бобр",
	}

	if len(result) != len(expectedWords) {
		t.Errorf("ParsePage len(result) == %d, want %d", len(result), len(expectedWords))
	}

	for i, w := range result {
		if w.Word != expectedWords[i] {
			t.Errorf("ParsePage w.Word == %s, want %s", w.Word, expectedWords[i])

			return
		}

		if len(w.Translations) == 0 {
			t.Errorf("ParsePage len(w.Translations) == %d, want 1", len(w.Translations))

			return
		}

		if len(w.Translations) > 1 {
			t.Errorf("ParsePage, len(w.Translations) == %d, want 1", len(w.Translations))
		}

		if w.Translations[0] != expectedTranslations[i] {
			t.Errorf("ParsePage w.Translation == %s, want %s", w.Translations[0], expectedTranslations[i])

			return
		}
	}
}

func TestParsePageIssue1(t *testing.T) {
	f, _ := os.Open("./test/sample_issue1.html")
	parser := seznam.NewParser()
	result, _ := parser.Parse(f)

	w := result[0]

	const expectedWord = "kvůli"

	if w.Word != expectedWord {
		t.Errorf("ParsePage word == %q, want %q", w.Word, expectedWord)
	}

	expectedTranslations := []string{
		"из-за",
		"ра́ди кого/чего",
	}

	if len(w.Translations) != len(expectedTranslations) {
		t.Errorf("ParsePage len(translation) == %d, want %d", len(w.Translations), len(expectedTranslations))
		return
	}

	for i, trans := range w.Translations {
		if trans != expectedTranslations[i] {
			t.Errorf("ParsePage translation == %q, want %q", trans, expectedTranslations[i])
		}
	}
}

func TestParsePageKozy(t *testing.T) {
	f, _ := os.Open("./test/sample_koza.html")
	parser := seznam.NewParser()
	result, _ := parser.Parse(f)

	w := result[0]

	const expectedWord = "koza"

	if w.Word != expectedWord {
		t.Errorf("ParsePage word == %q, want %q", w.Word, expectedWord)
	}

	expectedTranslations := []string{
		"коза́",
		"ко́злы",
		"козелки́",
		"козёл",
		"ду́ра",
		"си́ськи",
		"си́си",
		"ти́тьки",
		"буфера́",
		"сисяры́",
	}

	if len(w.Translations) != len(expectedTranslations) {
		t.Errorf("ParsePage len(translation) == %d, want %d", len(w.Translations), len(expectedTranslations))
		return
	}

	for i, trans := range w.Translations {
		if trans != expectedTranslations[i] {
			t.Errorf("ParsePage translation == %q, want %q", trans, expectedTranslations[i])
		}
	}
}
