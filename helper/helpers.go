package helper

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/julienschmidt/httprouter"
	"github.com/yhagio/go-twit/config"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidateJWTTokenMiddleware(next httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// If we use cookie:
		// _, er := r.Cookie("Auth")
		// if er != nil {
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	fmt.Fprint(w, "Token is not valid")
		// }

		token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return config.VerifyKey, nil
			})

		if err == nil {
			if token.Valid {
				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok || !token.Valid {
					w.WriteHeader(http.StatusUnauthorized)
					fmt.Fprint(w, "Token is not valid")
				}

				// We want to pass user_id extracted from JWT token to next handler:
				// Take the context out from the request
				ctx := r.Context()
				// Get new context with key-value "user_id" -> "2" for example
				ctx = context.WithValue(ctx, "user_id", claims["user_id"])
				// Get new http.Request with the new context
				r = r.WithContext(ctx)

				next(w, r, ps)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Token is not valid")
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized access to this resource")
		}

	}
}
