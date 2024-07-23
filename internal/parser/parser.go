package parser

import (
	"errors"
	"slices"

	"github.com/X3NOOO/llamaparse-go"
)

var ErrUnsupportedMIME = errors.New("unsupported MIME type")

func Parse(file []byte, mime string) (string, error) {
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
