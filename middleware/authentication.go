package middleware

import (
	"awesomeProject/handlers"
	"fmt"
	"net/http"
	"time"
)

var sampleSecretKey = []byte("SecretYouShouldHide")

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("You are authenticated!!!!")
		next.ServeHTTP(w, r)
	})
}

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		c, err := request.Cookie("Token")
		if err != nil {
			writer.Write([]byte("No token in cookie"))
			return
		}
		tknStr := c.Value
		if tknStr != "" {
			token, err := jwt.Parse(tknStr, func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					writer.WriteHeader(http.StatusUnauthorized)
					_, err := writer.Write([]byte("You're Unauthorized"))
					if err != nil {
						return nil, err
					}
				}
				return sampleSecretKey, nil
			})

			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				_, err2 := writer.Write([]byte("You're Unauthorized due to error parsing the JWT"))
				if err2 != nil {
					return
				}

			}
			// if there's a token
			if token.Valid {
				fmt.Println("Valid request")
				time_, _ := token.Claims.GetExpirationTime()

				if time.Until(time_.Time) < 300*time.Second {
					tokenString, err := handlers.CreateToken()
					if err != nil {
						fmt.Println(err)
						return
					}

					http.SetCookie(writer, &http.Cookie{
						Name:  "Token",
						Value: tokenString,
						Path:  "/",
					})

					fmt.Println("Token refreshed")
				}

				next.ServeHTTP(writer, request)
			} else {
				writer.WriteHeader(http.StatusBadGateway)
				_, err := writer.Write([]byte("You're Unauthorized due to invalid token"))
				if err != nil {
					return
				}
			}
		} else {
			writer.WriteHeader(http.StatusForbidden)
			_, err := writer.Write([]byte("You're Unauthorized due to No token in the header"))
			if err != nil {
				return
			}
		}
	})
}
