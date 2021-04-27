// Package rest contains all the functions related to the rest api
package rest

import (
	"net/http"
	"strings"

	"github.com/GGP1/adak/internal/response"
)

// Home gives users a welcome and takes non-invasive information from them.
func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get("Accept-Language")

		if lang != "" {
			message := getLanguage(lang)
			response.HTMLText(w, http.StatusOK, message)
			return
		}

		response.HTMLText(w, http.StatusOK, "Welcome to the Adak home page")
	}
}

func getLanguage(lang string) string {
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

	switch parts[0] {
	case spanish:
		return "Bienvenido a la página principal de Adak"
	case portuguese:
		return "Bem-vindo ao página principal do Adak"
	case chinese:
		return "歡迎來到帕洛首頁"
	case german:
		return "Wilkommen auf der Adak homepage"
	case french:
		return "Bienvenue sur la page d'accueil de Adak"
	case italian:
		return "Benvenuti nella home page di Adak"
	case russian:
		return "Добро пожаловать на домашнюю страницу Пало"
	case hindi:
		return "पालो होम पेज पर आपका स्वागत है"
	case japanese:
		return "パロのホームページへようこそ"
	default:
		return "Welcome to the Adak home page"
	}
}
