package routes

import (
	"auth/auth"
	"auth/handlers"
	"auth/middleware"

	"auth/proxy"

	"html/template"

	"github.com/gorilla/mux"
)

// RegisterRoutes регистрирует маршруты через Gorilla mux
func RegisterRoutes(m *mux.Router, tmpl *template.Template) {

	// 🛰️ Подключаем все прокси
	proxy.AdProxy(m)

	m.Handle("/", middleware.UserContextMiddleware(handlers.BaseHandler(tmpl)))

	m.HandleFunc("/theme", handlers.ToggleThemeHandler)

	m.HandleFunc("/personal-acount", handlers.PersonalAccountHandler(tmpl))

	m.HandleFunc("/auth", auth.AuthHandler(tmpl))

	m.HandleFunc("/get-password", auth.GetPasswordHandler(tmpl))

	m.HandleFunc("/verify-code", auth.VerifyCodeHandler(tmpl))
	m.HandleFunc("/status-email", auth.CheckEmailHandler(tmpl))

}
