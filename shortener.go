package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"time"
)

const MAX_STEPS = 200

type URLShortener struct {
	urls map[string]string
	mu   sync.RWMutex
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewURLShortener() *URLShortener {
	return &URLShortener{
		urls: make(map[string]string),
	}
}

// Shorten создает короткий идентификатор для URL
func (us *URLShortener) Shorten(originalURL string) (string, error) {
	// TODO: валидация URL
	// TODO: генерация короткого ID
	// TODO: сохранение в map
	if !isValidURL(originalURL) {
		return string(""), fmt.Errorf("%s is not a valid URL", originalURL)
	}
	steps := 0
	exist := true
	shortID := ""
	for exist {
		shortID = generateShortID()
		exist = us.shortExist(shortID)
		if steps > MAX_STEPS {
			return string(""), fmt.Errorf("Can't get unic shortID in %d steps, exit", MAX_STEPS)
		}
		steps++
	}
	us.mu.Lock()
	us.urls[shortID] = originalURL
	us.mu.Unlock()
	return shortID, nil
}

// GetOriginal возвращает оригинальный URL по короткому ID
func (us *URLShortener) GetOriginal(shortID string) (string, error) {
	// TODO: поиск в map
	// TODO: обработка отсутствующих ключей
	us.mu.RLock()
	original, exists := us.urls[shortID]
	us.mu.RUnlock()
	if exists {
		return original, nil
	} else {
		return string(""), fmt.Errorf("No URL with shortID = %s", shortID)
	}
}

// generateShortID генерирует случайный короткий идентификатор
func generateShortID() string {
	// TODO: генерация случайной строки 6-8 символов
	shortIDSize := rand.Intn(8-6) + 6
	b := make([]byte, shortIDSize)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// isValidURL проверяет корректность URL
func isValidURL(str string) bool {
	// TODO: валидация URL
	parsedUrl, err := url.Parse(str)
	return err == nil && parsedUrl.Host != "" && (parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https")
}

func (us *URLShortener) shortExist(str string) bool {
	us.mu.RLock()
	_, exists := us.urls[str]
	us.mu.RUnlock()
	return exists
}
