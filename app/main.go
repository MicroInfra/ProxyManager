package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"main/controllers"
	"main/models"
	"net/http"
	"os"
)

func main() {
	route := mux.NewRouter()
	server := controllers.Server{Proxies: models.NewAllProxies()}

	err := loadProxiesFromFile(&server.Proxies)
	if err != nil {
		log.Printf("error occured while loading proxies from file: %e", err)
	}

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

	err = http.ListenAndServe("localhost:8000", route)
	if err != nil {
		fmt.Println(err)
	}
}

func loadProxiesFromFile(p *models.Proxies) error {
	// Open the settings.example.json file
	fileContent, err := os.ReadFile("./settings.example.json")
	if err != nil {
		return err
	}

	// Unmarshal the file content to a slice of Proxy structs
	var proxies []models.ProxyRulesPlainText
	err = json.Unmarshal(fileContent, &proxies)
	if err != nil {
		return err
	}

	// Add each proxy to the Proxies map using the ServiceName field as the key
	for _, proxy := range proxies {
		path := fmt.Sprintf("./rules/%s.py", proxy.ServiceName)
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Println("Could not save file", err)
			return err
		}
		defer f.Close()
		err = os.WriteFile(path, []byte(proxy.Rules), 0666)
		if err != nil {
			return err
		}

		defaultProxy := models.Proxy{
			ServiceName: proxy.ServiceName,
			ServiceUrl:  proxy.ServiceUrl,
			ListenPort:  proxy.ListenPort,
			ProxyType:   proxy.ProxyType,
			FilterFile:  path,
		}
		err = p.Set(proxy.ServiceName, defaultProxy)
		if err != nil {
			return err
		}
	}

	return nil
}
