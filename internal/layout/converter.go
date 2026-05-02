package layout

import "strings"

// Converter transforms text between supported keyboard layouts.
type Converter struct{}

// NewConverter creates a Converter.
func NewConverter() Converter {
	return Converter{}
}

// EnToRu converts QWERTY keystrokes to the Russian layout.
func (Converter) EnToRu(text string) string {
	return convert(text, enToRu)
}

// RuToEn converts Russian-layout keystrokes to QWERTY.
func (Converter) RuToEn(text string) string {
	return convert(text, ruToEn)
}

func convert(text string, table map[rune]rune) string {
	var out strings.Builder
	out.Grow(len(text))

	for _, r := range text {
		if mapped, ok := table[r]; ok {
			out.WriteRune(mapped)
			continue
		}
		out.WriteRune(r)
	}

	return out.String()
}
