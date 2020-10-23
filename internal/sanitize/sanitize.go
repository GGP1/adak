package sanitize

import (
	"regexp"
	"unicode"

	"github.com/pkg/errors"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Normalize takes a list of strings pointers and normalizes them (NFC).
func Normalize(input ...*string) error {
	for _, v := range input {
		found, err := regexp.MatchString(`[[:^graph:]]`, *v)
		if err != nil {
			return errors.Wrap(err, "failed checking for invalid input")
		}

		if found {
			// Normalize input to NFC
			t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
			result, _, _ := transform.String(t, *v)
			*v = result
		}
	}
	return nil
}

// isMn returns if the rune is in the range of unicode.Mn (nonspacing marks).
func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}
