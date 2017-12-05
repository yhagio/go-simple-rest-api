package main

import (
	"crypto/rsa"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

type Twit struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type SignUpUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

var db *sql.DB

// ============== helpers and initializer ==============

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

		token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return verifyKey, nil
			})

		if err == nil {
			if token.Valid {
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

func ValidateUserAndJWTTokenMiddleware() {}

// ============== routes and handlers ==============

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

	var user SignUpUser

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	hash, _ := HashPassword(user.Password)

	_, err = db.Exec("INSERT INTO users (USERNAME, EMAIL, PASSWORD) VALUES ($1, $2, $3)", user.Username, user.Email, hash)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var user UserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	row := db.QueryRow("SELECT * FROM users WHERE email = $1 LIMIT 1", user.Email)

	userDict := User{}
	er := row.Scan(&userDict.ID, &userDict.Username, &userDict.Email, &userDict.Password, &userDict.CreatedAt)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case er != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	if !CheckPasswordHash(user.Password, userDict.Password) {
		w.WriteHeader(http.StatusForbidden)
		http.Error(w, "Email and/or password do not match", http.StatusForbidden)
		return
	}

	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error extracting the key")
		panic(err)
	}

	tokenString, err := token.SignedString(signKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		panic(err)
	}

	response := Token{tokenString}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
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

// ============== main ==============

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
