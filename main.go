package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

type Twit struct {
	ID         int       `json:"id"`
	User_id    int       `json:"user_id"`
	Body       string    `json:"body"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type User struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Created_at time.Time `json:"created_at"`
}

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
}

func AllTwits(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "AllTwits!\n")
}

func Signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Signup!\n")
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Login!\n")
}

func Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Logout!\n")
}

func CreateTwit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "CreateTwit!\n")
}

func OneTwit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "OneTwit!\n")
}

func UpdateTwit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "UpdateTwit!\n")
}

func DeleteTwit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "DeleteTwit!\n")
}

func main() {
	router := httprouter.New()
	router.GET("/", AllTwits)
	router.POST("/signup", Signup)
	router.POST("/login", Login)
	router.GET("/logout", Logout)
	router.POST("/twit", CreateTwit)
	router.GET("/twit/:id", OneTwit)
	router.PUT("/twit/:id", UpdateTwit)
	router.DELETE("/twit/:id", DeleteTwit)

	log.Fatal(http.ListenAndServe(":8080", router))
}
