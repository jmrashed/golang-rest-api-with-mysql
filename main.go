package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"./controller"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/getEmployee", controller.AllEmployee).Methods("GET")
	router.HandleFunc("/insertEmployee", controller.InsertEmployee).Methods("POST")
	http.Handle("/", router)
	fmt.Println("Connected to port 1234")
	log.Fatal(http.ListenAndServe(":1234", router))
}
