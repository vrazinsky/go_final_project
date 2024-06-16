package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vrazinsky/go-final-project/auth"
	"github.com/vrazinsky/go-final-project/models"
)

func (h *Handlers) HandleSignIn(res http.ResponseWriter, req *http.Request) {
	var passwordInput models.SignInInput
	var buf bytes.Buffer
	authKey := os.Getenv("AUTH_KEY")
	if len(authKey) == 0 {
		res.Write(ErrorResponse(nil, "internal error"))
		return
	}
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		res.Write(ErrorResponse(err, ""))
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &passwordInput); err != nil {
		res.Write(ErrorResponse(err, ""))
		return
	}
	envPassword := os.Getenv("TODO_PASSWORD")
	if passwordInput.Password != envPassword {
		res.Write(ErrorResponse(nil, "incorrect password"))

		return
	}
	claims := jwt.MapClaims{}
	claims["password"] = auth.GetMD5Hash(envPassword)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(authKey))
	if err != nil {
		res.Write(ErrorResponse(err, ""))
		return
	}
	response := models.SignInResponse{Token: signedToken}
	data, _ := json.Marshal(response)
	res.Write(data)
}
