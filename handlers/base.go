package handlers

import (
	"html/template"
	"log"

	"net/http"
)

// --- обработчик главной страницы ---
// BaseHandler — рендерит базовую страницу с темой и именем пользователя
func BaseHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// --- определяем тему из cookie ---
		theme := "light"
		if c, err := r.Cookie("theme"); err == nil && c.Value == "dark" {
			theme = "dark"
		}

		// --- формируем данные для шаблона ---
		data := struct {
			Theme string
			Name  string
		}{
			Theme: theme,
		}

		// --- рендерим шаблон ---
		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			log.Printf("[BaseHandler] ❌ Ошибка шаблона: %v", err)
			http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
			return
		}

		log.Println("[BaseHandler] ✅ Главная страница отрендерена")
	}
}

// переключение темы через cookie
func ToggleThemeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	// читаем текущую тему
	theme := "light"
	if c, err := r.Cookie("theme"); err == nil && c.Value == "dark" {
		theme = "dark"
	}

	// переключаем тему
	if theme == "light" {
		theme = "dark"
	} else {
		theme = "light"
	}

	// ставим cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "theme",
		Value:    theme,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   30 * 24, // 1 месяц
	})

	// редирект для HTMX
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
