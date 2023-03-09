package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/gorilla/mux"
  "strconv"
	"main/models"
	"net/http"
)

type Server struct {
	Proxies models.Proxies
}

// GetAll Function to get all proxies
func (s *Server) GetAll(response http.ResponseWriter, req *http.Request) {
	log.Printf("GET %s\n", req.URL.Path)

	jsonResponse, jsonError := json.Marshal(s.Proxies)

	if jsonError != nil {
		fmt.Println("Unable to encode JSON")
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonResponse)
}


// Get Function to get a proxy by name
func (s *Server) Get(response http.ResponseWriter, req *http.Request) {
	log.Printf("GET %s\n", req.URL.Path)
  vars := mux.Vars(req)
  name := vars["name"]

  proxy := s.Proxies.Get(name)

	jsonResponse, jsonError := json.Marshal(proxy)
	if jsonError != nil {
		fmt.Println("Unable to encode JSON")
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonResponse)
}


// Write function creates or updates a proxy from a json request
func (s *Server) Write(response http.ResponseWriter, req *http.Request) {
  log.Printf("Write proxy URL: %s\n", req.URL.Path)

  file, handler, err := req.FormFile("filterFile") // Get the file from the form data
  if err != nil {
      fmt.Println(err)
      return
  }
  defer file.Close()

  fmt.Println("File name: %v\n", handler.Filename) // Write the file name to the response
  port, err := strconv.Atoi(req.FormValue("listenPort"))
  if err != nil {
    fmt.Println("port is not a number")
		response.WriteHeader(http.StatusBadRequest)
		return
  }
  proxy :=  models.Proxy{ServiceName: req.FormValue("serviceName"), ServiceUrl: req.FormValue("serviceUrl"), ListenPort: port , ProxyType: req.FormValue("proxyType"), FilterFile: nil}

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


// Get Function to get a proxy by name
func (s *Server) Delete(response http.ResponseWriter, req *http.Request) {
	log.Printf("DELETE %s\n", req.URL.Path)
  vars := mux.Vars(req)
  name := vars["name"]

  s.Proxies.Delete(name)

	response.WriteHeader(http.StatusNoContent)
}

