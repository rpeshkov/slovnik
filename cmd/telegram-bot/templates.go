package main

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/rpeshkov/slovnik"
)

type Template struct {
	tmpl *template.Template
}

func CreateTemplate() (*Template, error) {
	funcs := template.FuncMap{
		"join": strings.Join,
	}

	templates, err := template.New("").Funcs(funcs).ParseGlob("./templates/*.gotmpl")

	if err != nil {
		return nil, errors.Wrap(err, "unable to load templates")
	}

	return &Template{templates}, nil
}

func (t *Template) Translation(words []*slovnik.Word) string {
	var buf bytes.Buffer
	t.tmpl.ExecuteTemplate(&buf, "translation", words)
	return buf.String()
}

func (t *Template) Phrases(words []*slovnik.Word) string {
	var buf bytes.Buffer
	t.tmpl.ExecuteTemplate(&buf, "phrases", words)
	return buf.String()
}
