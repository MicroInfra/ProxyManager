package main

import (
	"fmt"
	"github.com/gorilla/mux"
  "main/controllers"
	"main/models"
	"net/http"
)

func main() {
	route := mux.NewRouter()
	server := controllers.Server{Proxies: models.NewAllProxies()}

  print("Server is start!!")
	route.HandleFunc("/proxies", server.GetAll).Methods("GET")
	route.HandleFunc("/proxies/<id>", server.Get).Methods("GET")
	route.HandleFunc("/proxies", server.Create).Methods("POST")
	route.HandleFunc("/proxies", server.Update).Methods("PUT")
	route.HandleFunc("/proxies", server.Delete).Methods("DELETE")

	err := http.ListenAndServe("localhost:8000", route)
	if err != nil {
		fmt.Println(err)
	}
}
