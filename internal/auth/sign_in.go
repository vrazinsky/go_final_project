package auth

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vrazinsky/go-final-project/internal/handlers"
	"github.com/vrazinsky/go-final-project/internal/models"
)

func (a *AuthService) HandleSignIn(res http.ResponseWriter, req *http.Request) {
	var passwordInput models.SignInInput
	var buf bytes.Buffer
	if len(a.authKey) == 0 {
		logWriteErr(res.Write(handlers.ErrorResponse(nil, "internal error")))
		return
	}
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		logWriteErr(res.Write(handlers.ErrorResponse(err, "")))
	}
	if err = json.Unmarshal(buf.Bytes(), &passwordInput); err != nil {
		logWriteErr(res.Write(handlers.ErrorResponse(err, "")))
		return
	}
	if passwordInput.Password != a.pass {
		logWriteErr(res.Write(handlers.ErrorResponse(nil, "incorrect password")))

		return
	}
	claims := jwt.MapClaims{}
	claims["password"] = GetMD5Hash(a.pass)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(a.authKey))
	if err != nil {
		logWriteErr(res.Write(handlers.ErrorResponse(err, "")))
		return
	}
	response := models.SignInResponse{Token: signedToken}
	data, _ := json.Marshal(response)
	logWriteErr(res.Write(data))
}

func logWriteErr(_ int, err error) {
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}
