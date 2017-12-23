package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/yhagio/go-twit/handler"
	"github.com/yhagio/go-twit/helper"
)

func main() {
	router := httprouter.New()
	router.GET("/", handler.AllTwits)
	router.POST("/signup", handler.Signup)
	router.POST("/login", handler.Login)
	router.GET("/logout", handler.Logout)
	router.POST("/twit", helper.ValidateJWTTokenMiddleware(handler.CreateTwit)) // Needs Authorization header 'Bearer [token]'
	router.GET("/twit/:id", helper.ValidateJWTTokenMiddleware(handler.OneTwit))
	router.PUT("/twit/:id", helper.ValidateJWTTokenMiddleware(handler.UpdateTwit))
	router.DELETE("/twit/:id", helper.ValidateJWTTokenMiddleware(handler.DeleteTwit))

	log.Fatal(http.ListenAndServe(":8080", router))
}
