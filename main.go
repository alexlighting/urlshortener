package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type URLStorage struct {
	Url string `json:"url"`
}

type ErrorMsg struct {
	Msg string `json:"error"`
}

type CreatedMsg struct {
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

var shortner = NewURLShortener()

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//проверяем метод, если не POST то выдаем ошибку
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorMsg{Msg: "Method not allowed"})
		return
	}
	body, error := io.ReadAll(r.Body)
	if error != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorMsg{Msg: fmt.Sprintf("HTTP body read error: %v", error)})
		return
	}
	var req URLStorage
	//проверяем что в теле пришел правильный JSON
	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorMsg{Msg: "Incorrect JSON"})
		return
	}
	shortID, error := shortner.Shorten(req.Url)
	//генерируем короткий ID, если при этом возникает ошибка передаем ее в ответе
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorMsg{Msg: error.Error()})
		return
	}
	//если все прошло нормально - упаковывеам данные в JSON и отправляем
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreatedMsg{ShortUrl: shortID, OriginalUrl: req.Url})
}

func GetOriginalHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//проверяем метод, если не POST то выдаем ошибку
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorMsg{Msg: "Method not allowed"})
		return
	}
	//выделяем путь
	str := r.URL.Path[1:]
	if strings.Contains(str, "/") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorMsg{Msg: "Incorrect path"})
		return
	}
	//получаем полный URL по короткому ID

	originalURL, error := shortner.GetOriginal(str)
	if error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorMsg{Msg: error.Error()})
		return
	}
	// w.WriteHeader(http.StatusFound)
	fmt.Printf("Found URL, redirect to : %s\n", originalURL)
	http.Redirect(w, r, originalURL, http.StatusFound)
}

func main() {
	const port = 8080
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", ShortenHandler) // POST
	mux.HandleFunc("/", GetOriginalHandler)    // GET

	log.Printf("server listening on :%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatal(err)
	}
}
