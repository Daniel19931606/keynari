package words

import (
	_ "embed"
	"sync"
)

//go:embed data/ru.txt
var ruWordList string

//go:embed data/ru_opencorpora.txt
var ruOpenCorporaWordList string

//go:embed data/en.txt
var enWordList string

//go:embed data/en_scowl.txt
var enScowlWordList string

var (
	russianFullOnce sync.Once
	russianFullDict Dictionary
	englishFullOnce sync.Once
	englishFullDict Dictionary
)

// RussianFull returns the embedded Russian dictionary plus project-specific seed words.
func RussianFull() Dictionary {
	russianFullOnce.Do(func() {
		russianFullDict = Merge(
			FromText(ruWordList),
			FromText(ruOpenCorporaWordList),
			NewDictionary(russianSeed...),
		)
	})
	return russianFullDict
}

// EnglishFull returns the embedded English dictionary plus project-specific seed words.
func EnglishFull() Dictionary {
	englishFullOnce.Do(func() {
		englishFullDict = Merge(
			FromText(enWordList),
			FromText(enScowlWordList),
			NewDictionary(englishSeed...),
		)
	})
	return englishFullDict
}
