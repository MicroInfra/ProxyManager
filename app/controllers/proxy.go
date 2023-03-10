package controllers

import (
	"encoding/json"
	"fmt"
	"log"
  "os"
  "io"
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

  filename := ""
  if err == nil {
    defer file.Close()
    filename = fmt.Sprintf("./rules/%s", handler.Filename)
    f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
      if err != nil {
        fmt.Println("Could not save file", err)
        response.WriteHeader(http.StatusBadRequest)
        return
      }
      defer f.Close()

      io.Copy(f, file)
  }

  port, err := strconv.Atoi(req.FormValue("listenPort"))
  if err != nil {
    fmt.Println("port is not a number")
		response.WriteHeader(http.StatusBadRequest)
		return
  }

  proxy :=  models.Proxy{ServiceName: req.FormValue("serviceName"), ServiceUrl: req.FormValue("serviceUrl"), ListenPort: port , ProxyType: req.FormValue("proxyType"), FilterFile: filename}

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

