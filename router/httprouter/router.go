package httprouter

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"qok.com/identity/controller"
)

func Run(port string) {
	r := mux.NewRouter()

	r.HandleFunc("/register", controller.RegisterHandler).
		Methods("POST")

	r.HandleFunc("/login", controller.LoginHandler).
		Methods("POST")

	r.HandleFunc("/user_info", controller.UserInfoHandler).
		Methods("GET")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
