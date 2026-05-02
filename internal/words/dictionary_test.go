package words

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDictionaryIsCaseInsensitive(t *testing.T) {
	d := NewDictionary("Привет")

	if !d.Contains("привет") {
		t.Fatal("dictionary should match lowercase word")
	}

	if !d.Contains("ПРИВЕТ") {
		t.Fatal("dictionary should match uppercase word")
	}
}

func TestDictionaryNormalizesYo(t *testing.T) {
	d := NewDictionary("ёлка")

	if !d.Contains("елка") {
		t.Fatal("dictionary should match е/ё variants")
	}
}

func TestDictionaryLoadsExtraFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "extra.txt")
	if err := os.WriteFile(path, []byte("кастомслово\n"), 0644); err != nil {
		t.Fatalf("write extra dictionary: %v", err)
	}

	d, err := FromFile(path)
	if err != nil {
		t.Fatalf("FromFile() error = %v", err)
	}

	if !d.Contains("кастомслово") {
		t.Fatal("dictionary should contain word from file")
	}
}

func TestDictionaryMerge(t *testing.T) {
	d := Merge(NewDictionary("первое"), NewDictionary("второе"))

	if !d.Contains("первое") || !d.Contains("второе") {
		t.Fatal("merged dictionary should contain words from all inputs")
	}
}

func TestSeedDictionaries(t *testing.T) {
	if !Russian().Contains("дела") {
		t.Fatal("Russian seed dictionary should contain дела")
	}

	if !English().Contains("hello") {
		t.Fatal("English seed dictionary should contain hello")
	}
}

func TestFullDictionaries(t *testing.T) {
	if !RussianFull().Contains("человек") {
		t.Fatal("full Russian dictionary should contain человек")
	}

	if !RussianFull().Contains("переключателя") {
		t.Fatal("full Russian dictionary should contain morphology forms from OpenCorpora")
	}

	if !EnglishFull().Contains("keyboard") {
		t.Fatal("full English dictionary should contain keyboard")
	}

	if !EnglishFull().Contains("punctuation") {
		t.Fatal("full English dictionary should contain SCOWL words")
	}
}
