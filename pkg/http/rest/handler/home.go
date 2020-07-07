package handler

import (
	"net/http"
	"strings"

	"github.com/GGP1/palo/internal/response"
)

// Home page
func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get("Accept-Language")

		if lang != "" {
			sentence := getLanguage(lang)
			response.HTMLText(w, r, http.StatusOK, sentence)
			return
		}
		response.HTMLText(w, r, http.StatusOK, "Welcome to the Palo home page")
	}
}

func getLanguage(lang string) string {
	english := "en-US"
	spanish := "es-ES"
	chinese := "zh-CN"
	portuguese := "pt-BR"
	german := "de"
	french := "fr"
	italian := "it"
	russian := "ru"
	hindi := "in"
	japanese := "jp"

	langs := strings.Split(lang, ";")
	parts := strings.Split(langs[0], ",")

	var sentence string

	switch parts[0] {
	case english:
		sentence = "Welcome to the Palo home page"
	case spanish:
		sentence = "Bienvenido a la página pincipal de Palo"
	case portuguese:
		sentence = "Bem-vindo a página principal do Palo"
	case chinese:
		sentence = "歡迎來到帕洛首頁"
	case german:
		sentence = "Wilkommen auf der Palo homepage"
	case french:
		sentence = "Bienvenue sur la page d'accueil de Palo"
	case italian:
		sentence = "Benvenuti nella home page di Palo"
	case russian:
		sentence = "Добро пожаловать на домашнюю страницу Пало"
	case hindi:
		sentence = "पालो होम पेज पर आपका स्वागत है"
	case japanese:
		sentence = "パロのホームページへようこそ"
	}

	return sentence
}
