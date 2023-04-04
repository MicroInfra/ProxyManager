package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"main/controllers"
	"main/models"
	"net/http"
)

func main() {
	route := mux.NewRouter()
	server := controllers.Server{Proxies: models.NewAllProxies()}

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	route.Use(handlers.CORS(headers, methods, origins))

	route.HandleFunc("/proxies", server.GetAll).Methods("GET")
	route.HandleFunc("/proxies/{name}", server.Get).Methods("GET")
	route.HandleFunc("/proxies", server.Write).Methods("POST")
	route.HandleFunc("/proxies", server.Write).Methods("PUT")
	route.HandleFunc("/proxies/{name}", server.Delete).Methods("DELETE")
	route.HandleFunc("/proxies/{name}", server.Options).Methods("OPTIONS")

	err := http.ListenAndServe("localhost:8000", route)
	if err != nil {
		fmt.Println(err)
	}
}
