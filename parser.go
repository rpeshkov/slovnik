package slovnik

import (
	"io"
)

// Parser defines an interface for any parser
type Parser interface {
	Parse(input io.ReadCloser) ([]*Word, error)
}
