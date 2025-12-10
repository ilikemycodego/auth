package middleware

import (
	"context"
	"log"
	"net/http"

	"auth/db"
	"auth/token"
)

// структура, которую будем класть в контекст
type UserContext struct {
	UserID string
	Name   string
	Role   string
}

type ctxKey string

const UserKey ctxKey = "user"

// middleware/user_context.go
func GetUserFromContext(r *http.Request) *UserContext {
	val := r.Context().Value(UserKey)
	if val == nil {
		return nil
	}
	return val.(*UserContext)
}

// UserContextMiddleware — вытаскивает данные юзера по JWT
func UserContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("[Middleware] Проверяем cookie 'jwt'")
		c, err := r.Cookie("auth")
		if err != nil {
			log.Println("[Middleware] ❌ Cookie 'jwt' не найдена — гость")
			next.ServeHTTP(w, r)
			return
		}

		// --- разбираем токен ---
		claims, err := token.ParseJWT(c.Value)
		if err != nil {
			log.Println("[Middleware] ❌ Неверный JWT:", err)
			next.ServeHTTP(w, r)
			return
		}
		log.Println("[Middleware] JWT успешно распарсен, SessionID:", claims.SessionID, "Role:", claims.Role)

		// --- ищем сессию ---
		session, err := db.GetSession(claims.SessionID)
		if err != nil {
			log.Println("[Middleware] ❌ Сессия не найдена по SessionID:", claims.SessionID)
			next.ServeHTTP(w, r)
			return
		}
		log.Println("[Middleware] Сессия найдена, UserID:", session.UserID)

		// --- получаем пользователя по ID ---
		user, err := db.GetUserByID(session.UserID)
		if err != nil {
			log.Println("[Middleware] ❌ Пользователь не найден по ID:", session.UserID)
			next.ServeHTTP(w, r)
			return
		}
		log.Printf("[Middleware] Пользователь найден: ID=%s, Name=%s, Role=%s", user.ID, user.Name, user.Role)

		// --- кладём пользователя в контекст ---
		uc := &UserContext{
			UserID: user.ID,
			Name:   user.Name,
			Role:   user.Role,
		}
		ctx := context.WithValue(r.Context(), UserKey, uc)

		log.Println("[Middleware] Пользователь помещён в контекст")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
