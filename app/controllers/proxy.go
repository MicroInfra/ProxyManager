package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"main/models"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
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
	_, err := response.Write(jsonResponse)
	if err != nil {
		fmt.Println("unable to write data to response ", err)
	}
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
	_, err := response.Write(jsonResponse)
	if err != nil {
		fmt.Println("unable to write data to response ", err)
	}
}

// Write function creates or updates a proxy from a json request
func (s *Server) Write(response http.ResponseWriter, req *http.Request) {
	log.Printf("Write proxy URL: %s\n", req.URL.Path)

	file, handler, err := req.FormFile("filterFile") // Get the file from the form data

	filterRules := ""
	if err == nil {
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				fmt.Println("could not close file ", err)
			}
		}(file)
		filterRules = fmt.Sprintf("./rules/%s", handler.Filename)
		f, err := os.OpenFile(filterRules, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Println("Could not save file", err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				fmt.Println("could not close file ", err)
			}
		}(f)

		_, err = io.Copy(f, file)
		if err != nil {
			fmt.Println("could not copy to file ", err)
		}
	}

	port, err := strconv.Atoi(req.FormValue("listenPort"))
	if err != nil {
		fmt.Println("port is not a number")
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	proxy := models.Proxy{
		ServiceName: req.FormValue("serviceName"),
		ServiceUrl:  req.FormValue("serviceUrl"),
		ListenPort:  port,
		ProxyType:   req.FormValue("proxyType"),
		FilterFile:  filterRules}

	if req.Method == "PUT" {
		s.Proxies.Delete(proxy.ServiceName)
	}
	// Add the proxy to the list of proxies
	err = s.Proxies.Set(proxy.ServiceName, proxy)
	if err != nil {
		log.Printf("Could not create proxy: %e", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the proxy to the user
	jsonResponse, jsonError := json.Marshal(proxy)
	if jsonError != nil {
		fmt.Println("Unable to encode JSON")
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	_, err = response.Write(jsonResponse)
	if err != nil {
		fmt.Println("unable to write data to response ", err)
	}
}

// Delete proxy by name
func (s *Server) Delete(response http.ResponseWriter, req *http.Request) {
	log.Printf("DELETE %s\n", req.URL.Path)
	vars := mux.Vars(req)
	name := vars["name"]

	s.Proxies.Delete(name)

	response.WriteHeader(http.StatusNoContent)
}

// Options Function to get a proxy by name
func (s *Server) Options(response http.ResponseWriter, req *http.Request) {
	log.Printf("OPTIONS %s\n", req.URL.Path)
	response.Header().Set("Allow", "GET, POST, PUT, DELETE")

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusNoContent)
}
