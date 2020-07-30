package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/GGP1/palo/internal/response"
	"github.com/GGP1/palo/pkg/tracking"
)

// Home gives users a welcome and takes non-invasive information from them.
func Home(ctx context.Context, t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := t.Hit(ctx, r)
		if err != nil {
			response.Error(w, r, http.StatusInternalServerError, err)
			return
		}

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
		sentence = "Bienvenido a la página principal de Palo"
	case portuguese:
		sentence = "Bem-vindo ao página principal do Palo"
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
