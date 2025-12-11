package auth

import (
	"auth/db"
	"auth/token"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Проверка одноразового кода
func verifyCode(email, code string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dbCode string
	var used bool

	err := db.DB.QueryRow(ctx, `
		SELECT code, used FROM auth_codes
		WHERE email = $1 AND code = $2
	`, email, code).Scan(&dbCode, &used)

	if err != nil {
		return false, err
	}

	if used {
		return false, nil
	}

	// помечаем код как использованный
	_, err = db.DB.Exec(ctx, `
		UPDATE auth_codes SET used = true
		WHERE email = $1 AND code = $2
	`, email, code)

	return true, err
}

// Создание сессии на 7 дней
func createSession(userID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sessionID := uuid.New().String()
	expires := time.Now().Add(7 * 24 * time.Hour)

	_, err := db.DB.Exec(ctx, `
		INSERT INTO sessions (id, user_id, expires_at)
		VALUES ($1, $2, $3)
	`, sessionID, userID, expires)

	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// Поиск или создание нового пользователя (без роли)
func getOrCreateUser(email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id string

	// Ищем пользователя
	err := db.DB.QueryRow(ctx,
		"SELECT id FROM users WHERE email=$1 LIMIT 1",
		email,
	).Scan(&id)

	// Если найден — возвращаем
	if err == nil {
		return id, nil
	}

	// Создаём нового пользователя amigo
	newID := uuid.New().String()

	_, err = db.DB.Exec(ctx, `
		INSERT INTO users (id, email, name)
		VALUES ($1, $2, 'amigo')
	`, newID, email)

	if err != nil {
		return "", err
	}

	return newID, nil
}

// Верификация кода + создание пользователя + создание сессии
func VerifyCodeHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		email := r.FormValue("email")
		code := r.FormValue("code")

		if email == "" || code == "" {
			fmt.Fprint(w, "Заполните все поля!")
			return
		}

		ok, err := verifyCode(email, code)
		if err != nil {
			log.Println("[VerifyCodeHandler] Ошибка проверки кода:", err)
			fmt.Fprint(w, "Ошибка сервера")
			return
		}

		if !ok {
			fmt.Fprint(w, "Неверный или уже использованный код!")
			return
		}

		// 1️⃣ Создаём или ищем пользователя
		userID, err := getOrCreateUser(email)
		if err != nil {
			log.Println("[VerifyCodeHandler] Ошибка getOrCreateUser:", err)
			fmt.Fprint(w, "Ошибка сервера")
			return
		}

		// 2️⃣ Создаём сессию
		sessionID, err := createSession(userID)
		if err != nil {
			log.Println("[VerifyCodeHandler] Ошибка createSession:", err)
			fmt.Fprint(w, "Ошибка сервера")
			return
		}

		// 3️⃣ Генерируем JWT только с sessionID
		token, err := token.GenerateJWT(sessionID, token.ExpirationTime())
		if err != nil {
			log.Println("[VerifyCodeHandler] Ошибка GenerateJWT:", err)
			fmt.Fprint(w, "Ошибка сервера")
			return
		}

		// 4️⃣ Устанавливаем cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "auth",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
			MaxAge:   7 * 24 * 3600,
		})

		log.Println("[VerifyCodeHandler] JWT создан и cookie установлено")

		// 5️⃣ HTMX редирект
		w.Header().Set("HX-Redirect", "/")
		fmt.Fprint(w, "OK")
	}
}
