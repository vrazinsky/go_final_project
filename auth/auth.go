package auth

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		authKey := os.Getenv("AUTH_KEY")
		if len(pass) > 0 && len(authKey) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
			tokenStr := cookie.Value
			var valid bool = true
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(authKey), nil
			})
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				valid = false
			} else {
				if claims["password"] != GetMD5Hash(pass) {
					valid = false
				}
			}
			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
