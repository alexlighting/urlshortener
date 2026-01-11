package main

import (
	"testing"
)

func TestURLShortener_Shorten(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"валидный HTTP URL", "http://example.com", false},
		{"валидный HTTPS URL", "https://google.com/search?q=test", false},
		{"невалидная схема URL", "ftp://newHub.com/", true},
		{"невалидный URL", "not-a-url", true},
		{"невалидный URL (нет host)", "http://?q=test", true},
		{"пустая строка", "", true},
	}

	shortener := NewURLShortener()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortID, err := shortener.Shorten(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ошибка = %v, ожидали ошибку = %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(shortID) < 6 {
				t.Errorf("короткий ID слишком короткий: %s", shortID)
			}
			if !tt.wantErr && len(shortID) > 8 {
				t.Errorf("короткий ID слишком длинный: %s", shortID)
			}
		})
	}
}
