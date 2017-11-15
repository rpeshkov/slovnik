package slovnik_test

import (
	"testing"

	"github.com/rpeshkov/slovnik"
)

func TestDetectLanguage(t *testing.T) {
	cases := []struct {
		in   string
		lang slovnik.Language
	}{
		{"hlavní", slovnik.Cz},
		{"привет", slovnik.Ru},
		{"sиniy", slovnik.Ru},
	}

	for _, c := range cases {
		got := slovnik.DetectLanguage(c.in)
		if got != c.lang {
			t.Errorf("DetectLanguage(%q) == %q, want %q", c.in, got, c.lang)
		}
	}
}
