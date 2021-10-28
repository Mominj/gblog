package main

import (
	"fmt"
	controller "gblog/controller"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB
var DBErr error

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3030"
)

func init() {
	DB, DBErr = sqlx.Connect("postgres", "user=postgres password=momin1234 dbname=gblog sslmode=disable")
	if DBErr != nil {
		log.Fatalln("error occur when database conneting", DBErr)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/user", controller.CreateUser).Methods("POST")
	router.HandleFunc("/login", controller.Login).Methods("GET")
	router.HandleFunc("/read", controller.ReadCookie).Methods("GET")
	router.HandleFunc("/logout", controller.Logout).Methods("GET")

	router.HandleFunc("/create", controller.Createblog).Methods("POST")
	router.HandleFunc("/create/{id}/edit", controller.Edited).Methods("PATCH")

	router.HandleFunc("/category", controller.CategoryCreate).Methods("POST")
	//router.HandleFunc("/crud/{id}/delete", controller.Deleted).Methods("DELETE")

	if err := http.ListenAndServe(fmt.Sprintf(":%s", CONN_PORT), router); err != nil {
		log.Fatal("error starting server: ", err)
	}
}
