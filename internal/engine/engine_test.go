package engine

import (
	"testing"

	"github.com/Daniel19931606/keynari/internal/layout"
	"github.com/Daniel19931606/keynari/internal/words"
)

func newTestEngine() *Engine {
	return New(
		layout.NewConverter(),
		words.Russian(),
		words.English(),
		Options{MinWordLength: 3, UseDictionary: true},
	)
}

func newAggressiveTestEngine() *Engine {
	return New(
		layout.NewConverter(),
		words.Russian(),
		words.English(),
		Options{MinWordLength: 3, UseDictionary: true, Aggressive: true},
	)
}

func TestEngineCorrectsEveryWordInPhrase(t *testing.T) {
	e := newTestEngine()
	var corrections []Correction

	for _, r := range "ghbdtn rfr ltkf " {
		corrections = append(corrections, e.Type(r)...)
	}

	if got, want := e.Text(), "привет как дела "; got != want {
		t.Fatalf("Text() = %q; want %q", got, want)
	}

	if got, want := len(corrections), 3; got != want {
		t.Fatalf("corrections = %d; want %d", got, want)
	}
}

func TestEngineCorrectsUserReportedPhrase(t *testing.T) {
	e := newTestEngine()

	for _, r := range "ghbdtn rfr ndjb ltkf xtkjdtr" {
		e.Type(r)
	}
	e.Flush()

	if got, want := e.Text(), "привет как твои дела человек"; got != want {
		t.Fatalf("Text() = %q; want %q", got, want)
	}
}

func TestEngineCorrectsLayoutLettersThatLookLikePunctuation(t *testing.T) {
	e := newTestEngine()

	for _, r := range "ghbdtn,hfccrf;b vyt rfr ndjb ltkf b xnj cj,bhftimcz ltkfnm ?" {
		e.Type(r)
	}
	e.Flush()

	if got, want := e.Text(), "привет,расскажи мне как твои дела и что собираешься делать ?"; got != want {
		t.Fatalf("Text() = %q; want %q", got, want)
	}
}

func TestEngineAggressivelyCorrectsUnknownLayoutPunctuationWord(t *testing.T) {
	e := newAggressiveTestEngine()

	for _, r := range "cj,bhftimcz " {
		e.Type(r)
	}

	if got, want := e.Text(), "собираешься "; got != want {
		t.Fatalf("Text() = %q; want %q", got, want)
	}
}

func TestEnginePreservesPunctuation(t *testing.T) {
	e := newTestEngine()

	for _, r := range "ghbdtn, rfr? " {
		e.Type(r)
	}

	if got, want := e.Text(), "привет, как? "; got != want {
		t.Fatalf("Text() = %q; want %q", got, want)
	}
}

func TestEngineBackspaceChangesCurrentToken(t *testing.T) {
	e := newTestEngine()

	for _, r := range "ghbdtnx" {
		e.Type(r)
	}
	e.Backspace()
	e.Type(' ')

	if got, want := e.Text(), "привет "; got != want {
		t.Fatalf("Text() = %q; want %q", got, want)
	}
}

func TestEngineDoesNotCorrectKnownSourceWord(t *testing.T) {
	e := newTestEngine()

	for _, r := range "hello " {
		e.Type(r)
	}

	if got, want := e.Text(), "hello "; got != want {
		t.Fatalf("Text() = %q; want %q", got, want)
	}
}

func TestEngineCanRunWithoutDictionary(t *testing.T) {
	e := New(
		layout.NewConverter(),
		nil,
		nil,
		Options{MinWordLength: 3, UseDictionary: false},
	)

	for _, r := range "ghbdtn " {
		e.Type(r)
	}

	if got, want := e.Text(), "привет "; got != want {
		t.Fatalf("Text() = %q; want %q", got, want)
	}
}
