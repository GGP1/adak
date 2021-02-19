// Package rest contains all the functions related to the rest api
package rest

import (
	"net/http"
	"strings"

	"github.com/GGP1/adak/internal/response"
	"github.com/GGP1/adak/pkg/tracking"
)

// Home gives users a welcome and takes non-invasive information from them.
func Home(t tracking.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := t.Hit(ctx, r); err != nil {
			response.Error(w, http.StatusInternalServerError, err)
			return
		}

		lang := r.Header.Get("Accept-Language")

		if lang != "" {
			sentence := getLanguage(lang)
			response.HTMLText(w, http.StatusOK, sentence)
			return
		}

		response.HTMLText(w, http.StatusOK, "Welcome to the Adak home page")
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
		sentence = "Welcome to the Adak home page"
	case spanish:
		sentence = "Bienvenido a la página principal de Adak"
	case portuguese:
		sentence = "Bem-vindo ao página principal do Adak"
	case chinese:
		sentence = "歡迎來到帕洛首頁"
	case german:
		sentence = "Wilkommen auf der Adak homepage"
	case french:
		sentence = "Bienvenue sur la page d'accueil de Adak"
	case italian:
		sentence = "Benvenuti nella home page di Adak"
	case russian:
		sentence = "Добро пожаловать на домашнюю страницу Пало"
	case hindi:
		sentence = "पालो होम पेज पर आपका स्वागत है"
	case japanese:
		sentence = "パロのホームページへようこそ"
	}

	return sentence
}
