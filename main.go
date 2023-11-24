package main

import (
	"log"
	"net/http"

	"github.com/autumnleaf-ra/gorilla-mux/controllers/authcontroller"
	"github.com/autumnleaf-ra/gorilla-mux/controllers/productcontroller"
	"github.com/autumnleaf-ra/gorilla-mux/middlewares"
	"github.com/autumnleaf-ra/gorilla-mux/models"
	"github.com/gorilla/mux"
)

func main() {
	models.ConnectDB()
	r := mux.NewRouter()

	r.HandleFunc("/login", authcontroller.Login).Methods("POST")
	r.HandleFunc("/register", authcontroller.Register).Methods("POST")
	r.HandleFunc("/logout", authcontroller.Logout).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/products", productcontroller.Index).Methods("GET")
	api.HandleFunc("/product/{id}", productcontroller.Show).Methods("GET")
	api.HandleFunc("/product", productcontroller.Create).Methods("POST")
	api.HandleFunc("/product/{id}", productcontroller.Update).Methods("PUT")
	api.HandleFunc("/product", productcontroller.Delete).Methods("DELETE")
	api.Use(middlewares.JWTMiddleWare)

	log.Fatal(http.ListenAndServe(":8000", r))

}
