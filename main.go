package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

type Twit struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
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
	// We only accept 'GET' method here
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	// Get all twits from DB
	rows, err := db.Query("SELECT * FROM twit")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	// Close the db connection at the end
	defer rows.Close()

	// Create twit object list
	twits := make([]Twit, 0)
	for rows.Next() {
		twit := Twit{}
		err := rows.Scan(&twit.ID, &twit.UserId, &twit.Body, &twit.CreatedAt, &twit.UpdatedAt) // order matters
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		twits = append(twits, twit)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Returns as JSON (List of Twit objects)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(twits); err != nil {
		panic(err)
	}
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

func OneTwit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// We only accept 'GET' method here
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	twitID := ps.ByName("id")

	// Get the specific twit from DB
	row := db.QueryRow("SELECT * FROM twit WHERE id = $1", twitID)

	// Create twit object
	twit := Twit{}
	err := row.Scan(&twit.ID, &twit.UserId, &twit.Body, &twit.CreatedAt, &twit.UpdatedAt)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// Returns as JSON (single Twit object)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(twit); err != nil {
		panic(err)
	}
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
