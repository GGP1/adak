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
	spanish := "es-ES"
	chinese := "zh-CN"
	portuguese := "pt-BR"
	german := "de"

	langs := strings.Split(lang, ";")
	parts := strings.Split(langs[0], ",")

	var sentence string

	if parts[0] == spanish {
		sentence = "Bienvenido a la página pincipal de Palo"
	} else if parts[0] == chinese {
		sentence = "歡迎來到帕洛首頁"
	} else if parts[0] == portuguese {
		sentence = "Bem-vindo a página principal do Palo"
	} else if parts[0] == german {
		sentence = "Wilkommen auf der Palo homepage"
	}
	return sentence
}
