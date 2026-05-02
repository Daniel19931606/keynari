package layout

import "testing"

func TestConverter(t *testing.T) {
	c := NewConverter()

	tests := []struct {
		name string
		got  string
		want string
	}{
		{name: "en to ru", got: c.EnToRu("ghbdtn rfr ltkf"), want: "привет как дела"},
		{name: "ru to en", got: c.RuToEn("привет"), want: "ghbdtn"},
		{name: "digits preserved", got: c.EnToRu("123"), want: "123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("got %q; want %q", tt.got, tt.want)
			}
		})
	}
}
