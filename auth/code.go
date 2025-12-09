package auth

import (
	"auth/db"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

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

	// отмечаем код как использованный
	_, err = db.DB.Exec(ctx, `
		UPDATE auth_codes SET used = true
		WHERE email = $1 AND code = $2
	`, email, code)

	return true, err
}

//---------------

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
			log.Println("[VerifyCodeHandler] ❌ Ошибка проверки кода:", err)
			fmt.Fprint(w, "Ошибка сервера")
			return
		}

		if !ok {
			fmt.Fprint(w, "Неверный или уже использованный код!")
			return
		}

		// Тут можно создать сессию или токен для входа
		fmt.Fprint(w, "✅ Вы успешно вошли!")
	}
}
