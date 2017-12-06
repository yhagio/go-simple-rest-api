package main

import (
	"crypto/rsa"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

const (
	// For simplicity these files are in the same folder as the app binary.
	// You shouldn't do this in production.
	privKeyPath = "app.rsa"     // `> openssl genrsa -out app.rsa 1024`
	pubKeyPath  = "app.rsa.pub" // `> openssl rsa -in app.rsa -pubout > app.rsa.pub`
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://yuichihagio:root@localhost/twit?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")

	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		panic(err)
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(err)
	}
	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		panic(err)
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		panic(err)
	}
}

func main() {
	router := httprouter.New()
	router.GET("/", AllTwits)
	router.POST("/signup", Signup)
	router.POST("/login", Login)
	router.GET("/logout", Logout)
	router.POST("/twit", ValidateJWTTokenMiddleware(CreateTwit)) // Needs Authorization header 'Bearer [token]'
	router.GET("/twit/:id", ValidateJWTTokenMiddleware(OneTwit))
	router.PUT("/twit/:id", ValidateJWTTokenMiddleware(UpdateTwit))
	router.DELETE("/twit/:id", ValidateJWTTokenMiddleware(DeleteTwit))

	log.Fatal(http.ListenAndServe(":8080", router))
}
