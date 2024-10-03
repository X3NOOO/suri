package parser

import (
	"errors"
	"slices"

	"github.com/X3NOOO/llamaparse-go"
)

type Parser interface {
	Parse(file []byte, mime string) (string, error)
}

var ErrUnsupportedMIME = errors.New("unsupported MIME type")

type DefaultParser struct{}

func (DefaultParser) Parse(file []byte, mime string) (string, error) {
	switch mime {
	case "text/plain":
		return string(file), nil
	default:
		if slices.Contains(llamaparse.SUPPORTED_MIME_TYPES, mime) {
			return llamaparse.Parse(file, llamaparse.MARKDOWN, nil, nil, nil, nil)
		}
	}

	return "", ErrUnsupportedMIME
}
