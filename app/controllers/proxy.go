package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"main/models"
	"net/http"
)

type Server struct {
	Proxies models.AllProxies
}

// GetAll Function to get all proxies
func (s *Server) GetAll(response http.ResponseWriter, req *http.Request) {
	log.Printf("handling get task at %s\n", req.URL.Path)

	//proxy := models.Proxy{ServiceUrl: "http://httpd:8000", ListenPort: 8000, ProxyType: "http", FilterFile: nil}

	jsonResponse, jsonError := json.Marshal(s.Proxies)

	if jsonError != nil {
		fmt.Println("Unable to encode JSON")
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonResponse)
}


// Get Function to get a proxy by id
func (s *Server) Get(response http.ResponseWriter, req *http.Request, int id) {
	log.Printf("handling get task at %s\n", req.URL.Path)

	proxy := models.Proxy{ServiceUrl: "http://httpd:8000", ListenPort: 8000, ProxyType: "http", FilterFile: nil}
 // Get proxy by id TODO

	jsonResponse, jsonError := json.Marshal(proxy)

	if jsonError != nil {
		fmt.Println("Unable to encode JSON")
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonResponse)
}


// Create Function to create a new proxy from a json request
func (s *Server) Create(response http.ResponseWriter, req *http.Request) {
	log.Printf("Creating a new proxy. URL %s\n", req.URL.Path)

	// Decode the request body into a proxy struct
	decoder := json.NewDecoder(req.Body)
	var proxy models.Proxy
	err := decoder.Decode(&proxy)
	if err != nil {
		fmt.Println("Unable to decode JSON")
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	// Add the proxy to the list of proxies
	s.Proxies.Set(proxy.ServiceName, proxy)

	// Return the proxy to the user
	jsonResponse, jsonError := json.Marshal(proxy)
	if jsonError != nil {
		fmt.Println("Unable to encode JSON")
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonResponse)
}
