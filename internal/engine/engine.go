package engine

import (
	"strings"
	"unicode"
)

type Direction int

const (
	DirectionNone Direction = iota
	DirectionEnToRu
	DirectionRuToEn
)

type Converter interface {
	EnToRu(text string) string
	RuToEn(text string) string
}

type Dictionary interface {
	Contains(word string) bool
}

type Options struct {
	MinWordLength int
	UseDictionary bool
	Aggressive    bool
}

type Correction struct {
	Original    string
	Corrected   string
	ReplaceFrom int
	ReplaceTo   int
	ReplaceLen  int
	TypedLen    int
	LiveText    string
	Direction   Direction
}

type Engine struct {
	converter Converter
	ru        Dictionary
	en        Dictionary
	options   Options
	text      []rune
}

func New(converter Converter, ru, en Dictionary, options Options) *Engine {
	if options.MinWordLength <= 0 {
		options.MinWordLength = 3
	}

	return &Engine{
		converter: converter,
		ru:        ru,
		en:        en,
		options:   options,
		text:      make([]rune, 0, 256),
	}
}

func (e *Engine) Type(r rune) []Correction {
	e.text = append(e.text, r)
	if isSegmentBoundary(r) {
		return e.checkCompletedSegment()
	}
	return nil
}

func (e *Engine) Backspace() {
	if len(e.text) > 0 {
		e.text = e.text[:len(e.text)-1]
	}
}

func (e *Engine) Flush() []Correction {
	return e.checkCompletedSegment()
}

func (e *Engine) Text() string {
	return string(e.text)
}

func (e *Engine) Reset() {
	e.text = e.text[:0]
}

func (e *Engine) checkCompletedSegment() []Correction {
	start, end, typedEnd := lastSegment(e.text)
	if start == end {
		return nil
	}

	originalRunes := e.text[start:end]
	original := string(originalRunes)
	corrected, ok := e.correctSegment(original)
	if !ok || corrected == original {
		return nil
	}

	direction := DirectionEnToRu
	if scriptOf(original) == "cyrillic" {
		direction = DirectionRuToEn
	}
	trailing := string(e.text[end:typedEnd])
	e.replace(start, end, []rune(corrected))

	return []Correction{{
		Original:    original,
		Corrected:   corrected,
		ReplaceFrom: start,
		ReplaceTo:   end,
		ReplaceLen:  len([]rune(original)),
		TypedLen:    typedEnd - start,
		LiveText:    corrected + trailing,
		Direction:   direction,
	}}
}

func (e *Engine) correctSegment(segment string) (string, bool) {
	if corrected, ok := e.correctKnownWhole(segment); ok {
		return corrected, true
	}

	if corrected, ok := e.correctParts(segment); ok {
		return corrected, true
	}

	if corrected, ok := e.correctLayoutPunctuationSuffix(segment); ok {
		return corrected, true
	}

	if corrected, ok := e.correctTrimmed(segment); ok {
		return corrected, true
	}

	if hasSoftSeparator(segment) && e.hasKnownSourcePart(segment) {
		return "", false
	}

	if direction, corrected := e.detect(segment); direction != DirectionNone {
		return corrected, true
	}

	return "", false
}

func (e *Engine) correctKnownWhole(segment string) (string, bool) {
	if hasLayoutPunctuation(segment) {
		if corrected, ok := e.correctKnownCore(segment); ok {
			return corrected, true
		}
		core, suffix := trimSentencePunctuationSuffix(segment)
		if core != segment && core != "" {
			if corrected, ok := e.correctKnownCore(core); ok {
				return corrected + suffix, true
			}
		}
	}

	if hasTrailingOuterPunctuation(segment) {
		prefix, core, suffix := trimOuterPunctuation(segment)
		if core != "" {
			if corrected, ok := e.correctKnownCore(core); ok {
				return prefix + corrected + suffix, true
			}
		}
	}

	if corrected, ok := e.correctKnownCore(segment); ok {
		return corrected, true
	}

	prefix, core, suffix := trimOuterPunctuation(segment)
	if core == "" {
		return "", false
	}

	if corrected, ok := e.correctKnownCore(core); ok {
		return prefix + corrected + suffix, true
	}

	return "", false
}

func hasTrailingOuterPunctuation(segment string) bool {
	runes := []rune(segment)
	return len(runes) > 1 && isTrailingOuterPunctuation(runes[len(runes)-1])
}

func (e *Engine) correctKnownCore(core string) (string, bool) {
	switch scriptOf(core) {
	case "latin":
		converted := fixConvertedTypo(e.converter.EnToRu(core))
		if e.ru != nil && e.ru.Contains(converted) && (e.en == nil || !e.en.Contains(core)) {
			return converted, true
		}
	case "cyrillic":
		converted := fixConvertedTypo(e.converter.RuToEn(core))
		if e.en != nil && e.en.Contains(converted) && (e.ru == nil || !e.ru.Contains(core)) {
			return converted, true
		}
	}

	return "", false
}

func (e *Engine) correctTrimmed(segment string) (string, bool) {
	prefix, core, suffix := trimOuterPunctuation(segment)
	if core == segment || core == "" {
		return "", false
	}

	if direction, corrected := e.detect(core); direction != DirectionNone {
		return prefix + corrected + suffix, true
	}

	return "", false
}

func (e *Engine) correctParts(segment string) (string, bool) {
	parts := splitParts(segment)
	if len(parts) <= 1 {
		return "", false
	}
	if parts[0].separator && canBeLayoutLetter([]rune(parts[0].text)[0]) {
		return "", false
	}

	var out strings.Builder
	changed := false
	unresolved := false

	for _, part := range parts {
		if part.separator {
			out.WriteString(part.text)
			continue
		}

		if direction, corrected := e.detect(part.text); direction != DirectionNone {
			out.WriteString(corrected)
			changed = true
			continue
		}

		if hasLetters(part.text) && !e.isKnownSource(part.text) {
			unresolved = true
		}
		out.WriteString(part.text)
	}

	if !changed || unresolved {
		return "", false
	}

	return out.String(), true
}

func (e *Engine) hasKnownSourcePart(segment string) bool {
	for _, part := range splitParts(segment) {
		if part.separator || !hasLetters(part.text) {
			continue
		}

		switch scriptOf(part.text) {
		case "latin":
			return e.isKnownSource(part.text)
		case "cyrillic":
			return e.isKnownSource(part.text)
		}
	}
	return false
}

func (e *Engine) isKnownSource(text string) bool {
	switch scriptOf(text) {
	case "latin":
		return e.en != nil && e.en.Contains(text)
	case "cyrillic":
		return e.ru != nil && e.ru.Contains(text)
	default:
		return false
	}
}

func (e *Engine) detect(token string) (Direction, string) {
	if !hasLetters(token) {
		return DirectionNone, ""
	}

	switch scriptOf(token) {
	case "latin":
		converted := fixConvertedTypo(e.converter.EnToRu(token))
		if e.shouldCorrect(converted, token, e.ru, e.en) {
			return DirectionEnToRu, converted
		}
	case "cyrillic":
		converted := fixConvertedTypo(e.converter.RuToEn(token))
		if e.shouldCorrect(converted, token, e.en, e.ru) {
			return DirectionRuToEn, converted
		}
	}

	return DirectionNone, ""
}

func (e *Engine) shouldCorrect(converted, original string, target, source Dictionary) bool {
	if !e.options.UseDictionary {
		return true
	}
	if target == nil || source == nil {
		return false
	}

	converted = strings.ToLower(converted)
	original = strings.ToLower(original)

	if target.Contains(converted) && !source.Contains(original) {
		return true
	}

	if len([]rune(original)) < e.options.MinWordLength && target.Contains(converted) && !source.Contains(original) {
		return true
	}

	if isShortLatinRussianFunction(original, converted) {
		return true
	}

	if isShortCyrillicEnglishFunction(original, converted) {
		return true
	}

	return e.options.Aggressive && hasLayoutPunctuation(original) && !source.Contains(original)
}

func isShortLatinRussianFunction(original, converted string) bool {
	switch original {
	case "z", "b", "d", "c", "r", "j", "e", "xt", "uj", "yf", "yt", "gj", "pf", "jn", "lj", "vs", "ns", "ot", "of":
		return true
	default:
		return false
	}
}

func isShortCyrillicEnglishFunction(original, converted string) bool {
	switch strings.ToLower(converted) {
	case "a", "i", "it", "is", "in", "to", "us", "we", "my", "me", "he":
		return len([]rune(original)) <= 2
	default:
		return false
	}
}

func (e *Engine) correctLayoutPunctuationSuffix(segment string) (string, bool) {
	runes := []rune(segment)
	if len(runes) < 2 {
		return "", false
	}

	punctuation, ok := cyrillicLayoutPunctuation(runes[len(runes)-1])
	if !ok {
		return "", false
	}

	core := string(runes[:len(runes)-1])
	if scriptOf(core) != "cyrillic" {
		return "", false
	}

	converted := fixConvertedTypo(e.converter.RuToEn(core))
	if e.en != nil && e.en.Contains(converted) && (e.ru == nil || !e.ru.Contains(segment)) {
		return converted + punctuation, true
	}

	return "", false
}

func (e *Engine) replace(start, end int, replacement []rune) {
	next := make([]rune, 0, len(e.text)-end+start+len(replacement))
	next = append(next, e.text[:start]...)
	next = append(next, replacement...)
	next = append(next, e.text[end:]...)
	e.text = next
}

func lastSegment(text []rune) (int, int, int) {
	end := len(text)
	for end > 0 && isSegmentBoundary(text[end-1]) {
		end--
	}
	typedEnd := len(text)

	start := end
	for start > 0 && !isSegmentBoundary(text[start-1]) {
		start--
	}

	return start, end, typedEnd
}

func isSegmentBoundary(r rune) bool {
	return unicode.IsSpace(r)
}

func isSoftSeparator(r rune) bool {
	return r == '.' || r == ',' || r == '!' ||
		r == '?' || r == ':' ||
		r == '(' || r == ')' || r == '"' || r == '\'' ||
		r == '-'
}

func canBeLayoutLetter(r rune) bool {
	switch r {
	case ',', ';', '[', ']', '.', '`', '<', '>', ':', '"', '{', '}', '~':
		return true
	default:
		return false
	}
}

func cyrillicLayoutPunctuation(r rune) (string, bool) {
	switch r {
	case 'б':
		return ",", true
	case 'Б':
		return "<", true
	case 'ю':
		return ".", true
	case 'Ю':
		return ">", true
	case 'ж':
		return ";", true
	case 'Ж':
		return ":", true
	case 'э':
		return "'", true
	case 'Э':
		return "\"", true
	default:
		return "", false
	}
}

func hasSoftSeparator(text string) bool {
	for _, r := range text {
		if isSoftSeparator(r) {
			return true
		}
	}
	return false
}

type segmentPart struct {
	text      string
	separator bool
}

func splitParts(segment string) []segmentPart {
	runes := []rune(segment)
	parts := make([]segmentPart, 0, 4)

	for i := 0; i < len(runes); {
		start := i
		separator := isSoftSeparator(runes[i])
		for i < len(runes) && isSoftSeparator(runes[i]) == separator {
			i++
		}
		parts = append(parts, segmentPart{
			text:      string(runes[start:i]),
			separator: separator,
		})
	}

	return parts
}

func trimOuterPunctuation(segment string) (string, string, string) {
	runes := []rune(segment)
	start := 0
	end := len(runes)

	for start < end && isLeadingOuterPunctuation(runes[start]) {
		start++
	}
	for end > start && isTrailingOuterPunctuation(runes[end-1]) {
		end--
	}

	return string(runes[:start]), string(runes[start:end]), string(runes[end:])
}

func trimSentencePunctuationSuffix(segment string) (string, string) {
	runes := []rune(segment)
	end := len(runes)
	for end > 0 {
		switch runes[end-1] {
		case '.', '!', '?':
			end--
		default:
			return string(runes[:end]), string(runes[end:])
		}
	}
	return "", segment
}

func isLeadingOuterPunctuation(r rune) bool {
	return r == '(' || r == '"' || r == '[' || r == '{'
}

func isTrailingOuterPunctuation(r rune) bool {
	return r == '.' || r == '!' || r == '?' || r == ':' ||
		r == ',' || r == ')' || r == '"' ||
		r == ']' || r == '}'
}

func hasLayoutPunctuation(text string) bool {
	for _, r := range text {
		switch r {
		case ',', ';', '[', ']', '.', '`', '<', '>', ':', '"', '{', '}', '~':
			return true
		}
	}
	return false
}

func fixConvertedTypo(text string) string {
	switch text {
	case "теексту":
		return "тексту"
	case "динные":
		return "длинные"
	case "прадщаложения":
		return "предложения"
	case "сво,и":
		return "свои"
	case "своби":
		return "свои"
	case "ёто":
		return "это"
	case "eghlish":
		return "english"
	case "nrmal":
		return "normal"
	case "euglish":
		return "english"
	case "prnctuation":
		return "punctuation"
	case "braking":
		return "breaking"
	default:
		return text
	}
}

func scriptOf(text string) string {
	latin := 0
	cyrillic := 0

	for _, r := range text {
		if unicode.Is(unicode.Latin, r) {
			latin++
		} else if unicode.Is(unicode.Cyrillic, r) {
			cyrillic++
		}
	}

	if latin > 0 && cyrillic == 0 {
		return "latin"
	}
	if cyrillic > 0 && latin == 0 {
		return "cyrillic"
	}
	return "mixed"
}

func hasLetters(text string) bool {
	for _, r := range text {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}
