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

	tests := []struct {
		name        string
		method      string
		path        string
		jsonData    string
		status      int
		correctJSON bool
	}{
		{name: "неправильный метод", method: "GET", path: "hello/test", jsonData: "", status: http.StatusMethodNotAllowed, correctJSON: false},
		{name: "правильный медот и корерктные данные", method: "POST", path: "shorten", jsonData: `{"url":"http://example.com/not/too/long/path"}`, status: http.StatusCreated, correctJSON: true},
		{name: "правильный медот и некорерктный JSON", method: "POST", path: "shorten", jsonData: `{"example":2:]}}`, status: http.StatusBadRequest, correctJSON: false},
		{name: "правильный медот и некорерктный URL", method: "POST", path: "shorten", jsonData: `{"url":"/com/not/too/long/path"}`, status: http.StatusBadRequest, correctJSON: false},
	}
	// var shortner = NewURLShortener()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/"+tt.path, strings.NewReader(tt.jsonData))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			ShortenHandler(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			if res.StatusCode != tt.status {
				t.Errorf("Ожидался статус %d, получено %d", tt.status, res.StatusCode)
			}
			contentType := res.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Ожидался Content-Type application/json, получено %s", contentType)
			}
			body, _ := io.ReadAll(res.Body)
			//проверяем корректность JSON
			if !json.Valid([]byte(body)) {
				t.Errorf("Ожидалось body - JSON, получено - некорректный JSON %q", string(body))
			} else {
				if tt.method == "POST" {
					var jsonContent CreatedMsg
					//проверяем что в теле пришел правильный JSON
					if err := json.Unmarshal(body, &jsonContent); err != nil && !tt.correctJSON {
						t.Errorf("Ожидалось correct_json = %t, получено - correct_JSON  = %t", tt.correctJSON, err != nil)
					}
				}
			}
		})

	}
}

func TestGetOriginalHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		path         string
		status       int
		needShortURL bool
		url          string
	}{
		{name: "неправильный метод", method: "POST", path: "hello/test", status: http.StatusMethodNotAllowed, needShortURL: false, url: ""},
		{name: "правильный метод, неправильный путь", method: "GET", path: "hello/test", status: http.StatusBadRequest, needShortURL: false, url: ""},
		{name: "правильный метод, нет записи в мапе", method: "GET", path: "hello", status: http.StatusNotFound, needShortURL: false, url: ""},
		{name: "редирект", method: "GET", path: "hello", status: http.StatusFound, needShortURL: true, url: "http://kopilka.info"},
	}
	// var shortner = NewURLShortener()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.needShortURL {
				tt.path, _ = shortner.Shorten(tt.url)
			}
			req := httptest.NewRequest(tt.method, "/"+tt.path, strings.NewReader(""))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			GetOriginalHandler(rec, req)
			res := rec.Result()
			defer res.Body.Close()
			if res.StatusCode != tt.status {
				t.Errorf("Ожидался статус %d, получено %d", tt.status, res.StatusCode)
			}
			if res.StatusCode == http.StatusFound {
				if res.Header.Get("Location") != tt.url {
					t.Errorf("Ожидался редирект на %s, получено %s", tt.url, res.Header.Get("Location"))
				}
			}
		})

	}
}
