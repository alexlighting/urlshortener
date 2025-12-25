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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, error := io.ReadAll(r.Body)
	if error != nil {
		http.Error(w, "HTTP body read errer:", http.StatusInternalServerError)
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
	}
	//если все прошло нормально - упаковывеам данные в JSON и отправляем
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreatedMsg{ShortUrl: shortID, OriginalUrl: req.Url})
}

func GetOriginalHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//проверяем метод, если не POST то выдаем ошибку
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//выделяем путь и проверяем его на корректность
	str := r.URL.Path
	str = strings.TrimPrefix(str, "/")
	str = strings.TrimSuffix(str, "/")
	if strings.Contains(str, "/") {
		fmt.Printf("incorrect path %s\n", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Incorrect GET path")
		return
	}
	//генерируем короткий ID
	originalURL, error := shortner.GetOriginal(str)
	if error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(error.Error())
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
