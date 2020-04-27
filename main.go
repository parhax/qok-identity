package main

import (
	"log"
	"net/http"

	"qok.com/identity/controller"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/register", controller.RegisterHandler).
		Methods("POST")

	r.HandleFunc("/login", controller.LoginHandler).
		Methods("POST")

	r.HandleFunc("/user_info", controller.UserInfoHandler).
		Methods("GET")

	log.Fatal(http.ListenAndServe(":8585", r))
}
