package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Hqqm/paygo/internal/interfaces"
)

type AuthMiddleware struct {
	authUsecases interfaces.AuthUsecases
}

func NewAuthMiddleware(authUsecases interfaces.AuthUsecases) *AuthMiddleware {
	return &AuthMiddleware{authUsecases: authUsecases}
}

func (am *AuthMiddleware) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("verify token middleware", r.URL.Path)
		header := r.Header.Get("X-Access-Token")
		header = strings.TrimSpace(header)
		if header == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
		defer cancel()

		account, err := am.authUsecases.ParseToken(ctx, header)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
		}

		ctx = context.WithValue(ctx, "account", account)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
