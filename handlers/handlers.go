package handlers

import (
	"fmt"
	"net/http"
	"time"
)

type jwtToken struct {
	Token string `json:"token"`
}

var sampleSecretKey = []byte("SecretYouShouldHide")

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Your homepage is here!!\n"))
}
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Your Admin page is here!!\n"))
}
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "Token",
		Path:    "/",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type Claims struct {
	Username   string `json:"username"`
	Authorised int    `json:"authorised"`
	jwt.RegisteredClaims
}

func CreateToken() (string, error) {
	claims := &Claims{
		Username:   "username",
		Authorised: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}

	fmt.Println("token is", tokenString)
	return tokenString, nil
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, err := CreateToken()
	if err != nil {
		fmt.Println(err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Token",
		Value: tokenString,
		Path:  "/",
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
