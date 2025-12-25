package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	//неправильный метод
	req := httptest.NewRequest("GET", "/hello/test", nil)
	rec := httptest.NewRecorder()
	ShortenHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Ожидался статус %d, получено %d", http.StatusMethodNotAllowed, res.StatusCode)
	}

	//правильный медот и корерктные данные
	jsonData := `{"url":"http://example.com/not/too/long/path"}`

	req = httptest.NewRequest("POST", "/shorten", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rec = httptest.NewRecorder()

	ShortenHandler(rec, req)
	res = rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Ожидался статус %d, получено %d", http.StatusOK, res.StatusCode)
	}
	//проверяем тип содержимого
	contentType := res.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Ожидался Content-Type application/json, получено %s", contentType)
	}

	body, _ := io.ReadAll(res.Body)
	//проверяем корректность JSON
	if !json.Valid([]byte(body)) {
		t.Errorf("Ожидалось body - JSON, получено - некорректный JSON %q", string(body))
	}

	//правильный медот и некорерктный JSON
	jsonData = `{"example":2:]}}`

	req = httptest.NewRequest("POST", "/shorten", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rec = httptest.NewRecorder()

	ShortenHandler(rec, req)
	res = rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Ожидался статус %d, получено %d", http.StatusBadRequest, res.StatusCode)
	}

	//правильный медот и некорерктный URL
	jsonData = `{"url":"/com/not/too/long/path"}`

	req = httptest.NewRequest("POST", "/shorten", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rec = httptest.NewRecorder()

	ShortenHandler(rec, req)
	res = rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Ожидался статус %d, получено %d", http.StatusBadRequest, res.StatusCode)
	}

}

func TestGetOriginalHandler(t *testing.T) {
	//неправильный метод
	req := httptest.NewRequest("POST", "/hello/test", nil)
	rec := httptest.NewRecorder()
	GetOriginalHandler(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Ожидался статус %d, получено %d", http.StatusMethodNotAllowed, res.StatusCode)
	}

	//правильный метод, неправильный путь
	req = httptest.NewRequest("GET", "/hello/test", nil)
	rec = httptest.NewRecorder()
	GetOriginalHandler(rec, req)
	res = rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Ожидался статус %d, получено %d", http.StatusBadRequest, res.StatusCode)
	}

	//правильный метод, нет записи в мапе
	req = httptest.NewRequest("GET", "/hello", nil)
	rec = httptest.NewRecorder()
	GetOriginalHandler(rec, req)
	res = rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Ожидался статус %d, получено %d", http.StatusNotFound, res.StatusCode)
	}

}
