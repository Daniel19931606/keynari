package words

import (
	"bufio"
	"os"
	"strings"
)

// Dictionary checks whether a word is known.
type Dictionary struct {
	words map[string]struct{}
}

var russianSeed = []string{
	"привет", "как", "дела", "пока", "спасибо", "да", "нет",
	"тест", "слово", "текст", "сегодня", "завтра", "можно",
	"и", "а", "в", "с", "к", "о", "у", "на", "не", "по", "за", "от", "до",
	"я", "ты", "он", "она", "оно", "мы", "вы", "они",
	"мне", "тебе", "ему", "ей", "нам", "вам", "им", "меня", "тебя",
	"мой", "моя", "мое", "мои", "твой", "твоя", "твое", "твои",
	"его", "ее", "наш", "наша", "наше", "наши", "ваш", "ваша", "ваше", "ваши",
	"это", "этот", "эта", "эти", "тот", "та", "те", "там", "тут", "здесь",
	"что", "кто", "где", "когда", "почему", "зачем", "какой", "какая", "какие",
	"человек", "люди", "друг", "дом", "работа", "время", "день", "ночь",
	"хорошо", "плохо", "нормально", "супер", "очень", "просто", "тоже",
	"будет", "было", "есть", "нету", "хочу", "можешь", "можем", "надо",
	"расскажи", "рассказать", "собираешься", "собираюсь", "собирается", "делать",
	"строим", "удобный", "переключатель", "раскладки", "который", "ловит",
	"ошибки", "мешает", "нормальному", "тексту", "если", "быстро", "печатает",
	"приложение", "должно", "спокойно", "исправлять", "каждое", "без", "пропусков",
	"работать", "проверяем", "пятнадцать", "разных", "сценариев", "строить",
	"дальше", "простой", "измеримый", "понятная", "документация", "напишу",
	"команду", "потому", "хотел", "поймать", "исправить", "их", "работу",
	"работает", "продолжай", "проверь", "слово", "фразу", "потом", "нового",
	"движка", "тестируется", "подведем", "итоги", "проверить", "большое",
	"подвижем",
	"даник", "тестов", "проверим", "длинные", "предложения", "порчи",
	"комы", "точки", "свои", "тексту",
	"че", "чето", "чтоб", "щас", "ща", "норм", "нормас", "нормально",
	"кринж", "кринжово", "кринжовый", "имба", "имбовый", "рофл", "рофлю",
	"рофлишь", "рофлить", "лол", "ору", "угар", "угарно", "жиза", "жизу",
	"изи", "чилл", "чилю", "чилить", "токсик", "токсично", "вайб", "вайбовый",
	"хайп", "хайповый", "зашквар", "зашкварно", "краш", "крашиха",
	"душнила", "душно", "флекс", "флексить", "флекшу", "мем", "мемный",
	"треш", "трешак", "капец", "офигеть", "офигел", "офигела", "офигенно",
	"прикол", "прикольно", "прикинь", "реально", "жестко", "жесть",
	"го", "погнали", "забей", "забыл", "забыла", "заценить", "зацени",
	"лайк", "лайкнул", "лайкнула", "репост", "подрубить", "подруби",
	"юзать", "юзаю", "юзаешь", "запушить", "пушить", "пушу", "пушим",
	"задеплоить", "деплой", "деплою", "фиксить", "фикшу", "фиксишь",
	"баг", "баги", "бага", "багов", "фича", "фичи", "фичу", "релиз",
	"релиза", "релизы", "аппа", "аппу", "аппка", "иконка", "иконку",
	"нажимаю", "нажимаешь", "нажать", "кликаю", "кликнуть", "закрываю",
	"закрываешь", "закрыть", "включено", "запущено",
}

var englishSeed = []string{
	"hello", "world", "test", "word", "text", "today", "tomorrow",
	"yes", "no", "thanks", "please",
	"a", "i", "is", "for", "and", "my", "keep", "this", "unchanged",
	"readable", "safe", "local", "keyboard", "tool", "macos", "github",
	"keynari", "releases",
	"danik", "english", "us", "to", "normal", "program", "should", "allow",
	"write", "long", "sentences", "punctuation", "preserve", "words",
	"want", "we", "it", "still", "without", "breaking",
}

// NewDictionary creates a lowercase word set.
func NewDictionary(items ...string) Dictionary {
	words := make(map[string]struct{}, len(items))
	for _, item := range items {
		addWord(words, item)
	}
	return Dictionary{words: words}
}

// FromText creates a dictionary from a newline-delimited word list.
func FromText(text string, extra ...string) Dictionary {
	words := make(map[string]struct{})
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Buffer(make([]byte, 1024), 1024*1024)

	for scanner.Scan() {
		addWord(words, scanner.Text())
	}

	for _, item := range extra {
		addWord(words, item)
	}

	return Dictionary{words: words}
}

// FromFile creates a dictionary from a newline-delimited word-list file.
func FromFile(path string, extra ...string) (Dictionary, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return Dictionary{}, err
	}

	return FromText(string(bytes), extra...), nil
}

// Merge returns a dictionary containing all words from the given dictionaries.
func Merge(dicts ...Dictionary) Dictionary {
	total := 0
	for _, dict := range dicts {
		total += len(dict.words)
	}

	words := make(map[string]struct{}, total)
	for _, dict := range dicts {
		for word := range dict.words {
			words[word] = struct{}{}
		}
	}

	return Dictionary{words: words}
}

// Contains reports whether word exists in the dictionary.
func (d Dictionary) Contains(word string) bool {
	_, ok := d.words[normalize(word)]
	return ok
}

// Len returns the number of unique words.
func (d Dictionary) Len() int {
	return len(d.words)
}

func addWord(words map[string]struct{}, word string) {
	word = normalize(word)
	if word == "" {
		return
	}
	words[word] = struct{}{}
}

func normalize(word string) string {
	word = strings.TrimSpace(strings.ToLower(word))
	word = strings.ReplaceAll(word, "ё", "е")
	return word
}

// Russian returns a small built-in Russian seed dictionary for engine tests and demos.
func Russian() Dictionary {
	return NewDictionary(russianSeed...)
}

// English returns a small built-in English seed dictionary for engine tests and demos.
func English() Dictionary {
	return NewDictionary(englishSeed...)
}
