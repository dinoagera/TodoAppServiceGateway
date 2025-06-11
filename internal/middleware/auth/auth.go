package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const (
	errorKey contextKey = "errorKey"
	uidKey   contextKey = "uidKey"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}
	return splitToken[1]
}

func New(log *slog.Logger, secretKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractBearerToken(r)
			if tokenStr == "" {
				http.Error(w, "Authorization token required", http.StatusUnauthorized)
				return
			}
			keyFunc := func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			}
			token, err := jwt.Parse(tokenStr, keyFunc)
			if err != nil {
				log.Warn("failed to parse token", "error", err)
				ctx := context.WithValue(r.Context(), errorKey, ErrInvalidToken)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				log.Warn("invalid token claims")
				ctx := context.WithValue(r.Context(), errorKey, ErrInvalidToken)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			uid, ok := claims["uid"].(float64)
			if !ok {
				log.Warn("uid not found in token claims or wrong type")
				ctx := context.WithValue(r.Context(), errorKey, ErrInvalidToken)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			uidInt := int64(uid)
			ctx := context.WithValue(r.Context(), uidKey, uidInt)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UIDFromContext(ctx context.Context) (int64, bool) {
	uid, ok := ctx.Value(uidKey).(int64)
	return uid, ok
}

func ErrorFromContext(ctx context.Context) (error, bool) {
	err, ok := ctx.Value(errorKey).(error)
	return err, ok
}
