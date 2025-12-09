package handlers

import (
	"html/template"
	"log"

	"net/http"
)

// --- обработчик главной страницы ---
// BaseHandler — рендерит базовую страницу с темой и именем пользователя
func PersonalAccountHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// --- рендерим шаблон ---
		if err := tmpl.ExecuteTemplate(w, "personal__account", nil); err != nil {
			log.Printf("[BaseHandlerpersonal__account] ❌ Ошибка шаблона: %v", err)
			http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
			return
		}

		log.Println("[personal__account] ✅ personal__account отрендерена")
	}
}
