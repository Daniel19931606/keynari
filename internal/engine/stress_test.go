package engine

import (
	"testing"

	"github.com/Daniel19931606/keynari/internal/layout"
	"github.com/Daniel19931606/keynari/internal/words"
)

func newFullTestEngine() *Engine {
	return New(
		layout.NewConverter(),
		words.RussianFull(),
		words.EnglishFull(),
		Options{MinWordLength: 3, UseDictionary: true, Aggressive: true},
	)
}

func runText(input string) string {
	e := newFullTestEngine()
	for _, r := range input {
		e.Type(r)
	}
	e.Flush()
	return e.Text()
}

func TestStressTexts(t *testing.T) {
	converter := layout.NewConverter()
	wrongRU := func(text string) string {
		return converter.RuToEn(text)
	}
	wrongEN := func(text string) string {
		return converter.EnToRu(text)
	}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "reported long phrase",
			input: "ghbdtn,hfccrf;b vyt rfr ndjb ltkf b xnj cj,bhftimcz ltkfnm ?",
			want:  "привет,расскажи мне как твои дела и что собираешься делать ?",
		},
		{
			name:  "plain russian sentence typed in english layout",
			input: wrongRU("сегодня я хочу проверить большое. давай потом подведем итоги!"),
			want:  "сегодня я хочу проверить большое. давай потом подведем итоги!",
		},
		{
			name:  "long russian product sentence",
			input: wrongRU("мы строим удобный переключатель раскладки, который ловит ошибки и не мешает нормальному тексту."),
			want:  "мы строим удобный переключатель раскладки, который ловит ошибки и не мешает нормальному тексту.",
		},
		{
			name:  "long russian support sentence",
			input: wrongRU("если человек быстро печатает, приложение должно спокойно исправлять каждое слово без пропусков."),
			want:  "если человек быстро печатает, приложение должно спокойно исправлять каждое слово без пропусков.",
		},
		{
			name:  "russian sentence with quotes",
			input: wrongRU("скажи мне: \"привет, как твои дела?\" и продолжай работать дальше."),
			want:  "скажи мне: \"привет, как твои дела?\" и продолжай работать дальше.",
		},
		{
			name:  "russian sentence with numbers",
			input: wrongRU("сегодня 02.05.2026 мы проверяем пятнадцать разных сценариев."),
			want:  "сегодня 02.05.2026 мы проверяем пятнадцать разных сценариев.",
		},
		{
			name:  "question with layout punctuation letters",
			input: wrongRU("скажи что надо строить дальше: простой измеримый тест и понятная документация."),
			want:  "скажи что надо строить дальше: простой измеримый тест и понятная документация.",
		},
		{
			name:  "mixed commas and semicolons inside wrong-layout words",
			input: wrongRU("я напишу команду, потому что хотел быстро поймать ошибки и исправить их."),
			want:  "я напишу команду, потому что хотел быстро поймать ошибки и исправить их.",
		},
		{
			name:  "already russian stays russian",
			input: "привет, расскажи мне как твои дела и что собираешься делать?",
			want:  "привет, расскажи мне как твои дела и что собираешься делать?",
		},
		{
			name:  "already english stays english",
			input: "hello, please keep this english text unchanged and readable.",
			want:  "hello, please keep this english text unchanged and readable.",
		},
		{
			name:  "english sentence typed in russian layout",
			input: wrongEN("hello world, please keep this english text unchanged."),
			want:  "hello world, please keep this english text unchanged.",
		},
		{
			name:  "short russian function words",
			input: wrongRU("я и ты можем пойти в дом и на работу."),
			want:  "я и ты можем пойти в дом и на работу.",
		},
		{
			name:  "last word without trailing space",
			input: wrongRU("привет как твои дела человек"),
			want:  "привет как твои дела человек",
		},
		{
			name:  "capital letters",
			input: wrongRU("Привет Как Дела"),
			want:  "Привет Как Дела",
		},
		{
			name:  "mixed valid english product name",
			input: "keynari " + wrongRU("просто работает и не мешает") + " GitHub.",
			want:  "keynari просто работает и не мешает GitHub.",
		},
		{
			name:  "english product sentence remains stable",
			input: "Keynari is a local keyboard tool for macOS and GitHub releases.",
			want:  "Keynari is a local keyboard tool for macOS and GitHub releases.",
		},
		{
			name:  "english sentence wrong layout with capitals",
			input: wrongEN("Hello GitHub, Please Keep My Text Safe."),
			want:  "Hello GitHub, Please Keep My Text Safe.",
		},
		{
			name:  "russian with brackets",
			input: wrongRU("проверь [слово] и {фразу}, потом продолжай."),
			want:  "проверь [слово] и {фразу}, потом продолжай.",
		},
		{
			name:  "russian with hyphenated phrase",
			input: wrongRU("это супер-тест для нового движка раскладки."),
			want:  "это супер-тест для нового движка раскладки.",
		},
		{
			name:  "already mixed languages stay mixed",
			input: "GitHub работает, Keynari тестируется, hello world.",
			want:  "GitHub работает, Keynari тестируется, hello world.",
		},
		{
			name:  "aggressive unknown punctuation-shaped word",
			input: "cj,bhftimcz ",
			want:  "собираешься ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runText(tt.input); got != tt.want {
				t.Fatalf("runText() = %q; want %q", got, tt.want)
			}
		})
	}
}
