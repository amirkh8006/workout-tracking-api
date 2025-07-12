package middleware

import (
	"context"
	"femProject/internal/store"
	"femProject/internal/tokens"
	"femProject/internal/utils"
	"net/http"
	"strings"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

type ContextKey string

const userContextKey = ContextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context() , userContextKey , user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(userContextKey).(*store.User)
	if !ok {
		panic("Missing User In Request")
	}

	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ") // Bearer <TOKEN>
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJson(w, http.StatusUnauthorized , utils.Envlope{"error":"Invalid authorization header"})
			return
		}

		token := headerParts[1]

		user , err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)
		if err != nil {
			utils.WriteJson(w, http.StatusUnauthorized , utils.Envlope{"error":"Invalid token"})
			return
		}

		if user == nil {
			utils.WriteJson(w, http.StatusUnauthorized , utils.Envlope{"error":"Token expired or invalid"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}


func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)

		if user.IsAnonymousUser() {
			utils.WriteJson(w, http.StatusUnauthorized , utils.Envlope{"error":"You must be logged in to access this route"})
			return
		}

		next.ServeHTTP(w, r)
	})
}