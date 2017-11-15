package seznam

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// attributes is a helper type for quickly accessing the attributes of HTML token or node
type attributes []html.Attribute

// get attribute by name
func (attrs attributes) get(name atom.Atom) string {
	for _, a := range attrs {
		if a.Key == name.String() {
			return a.Val
		}
	}

	return ""
}

func (attrs attributes) id() string {
	return attrs.get(atom.Id)
}

func (attrs attributes) class() string {
	return attrs.get(atom.Class)
}

func (attrs attributes) lang() string {
	return attrs.get(atom.Lang)
}
