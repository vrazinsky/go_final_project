package auth

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	pass    string
	authKey string
}

func NewAuthService(pass string, authKey string) *AuthService {
	return &AuthService{pass: pass, authKey: authKey}
}

func (s *AuthService) Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(s.pass) > 0 && len(s.authKey) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
			tokenStr := cookie.Value
			var valid bool = true
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(s.authKey), nil
			})
			if err != nil {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				valid = false
			} else {
				if claims["password"] != GetMD5Hash(s.pass) {
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
