package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Defining a type for the API functions
type apiFunc func(http.ResponseWriter, *http.Request) error

type APIServer struct {
	listenAddress string
	apiVersion string
	apiBaseUrl string
}

type Response struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data interface{} `json:"data"`
}

type APIError struct {
	Error string `json:"error"`
}


func WriteJsonResponse(writer http.ResponseWriter, status int, response interface{} ) error {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	return json.NewEncoder(writer).Encode(response)
}


func createHttpHandler(f apiFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := f(writer, request)
		if err != nil {
			log.Fatal(err)
			WriteJsonResponse(writer, http.StatusInternalServerError, APIError{Error: err.Error()})
		}
	}
}


func NewAPIServer(listenAddress string, apiVersion string, apiBaseUrl string ) *APIServer {
	// Returning a Pointer to the APIServer
	return &APIServer{
		listenAddress: listenAddress,
		apiVersion: apiVersion,
		apiBaseUrl: apiBaseUrl,
	}
}

func (server *APIServer) Start() {
	router := mux.NewRouter()

	// handle get request for /api/v1/coupon and api/v1/coupon/{id}
	router.HandleFunc(fmt.Sprintf("%s/%s/coupon", server.apiBaseUrl, server.apiVersion), createHttpHandler(server.handleCoupon)).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("%s/%s/coupon/{id}", server.apiBaseUrl, server.apiVersion), createHttpHandler(server.handleCoupon)).Methods(http.MethodGet)

	// handle post request for /api/v1/coupon
	router.HandleFunc(fmt.Sprintf("%s/%s/coupon", server.apiBaseUrl, server.apiVersion), createHttpHandler(server.handleCoupon)).Methods(http.MethodPost)

	// handle put request for /api/v1/coupon/{id}
	router.HandleFunc(fmt.Sprintf("%s/%s/coupon/{id}", server.apiBaseUrl, server.apiVersion), createHttpHandler(server.handleCoupon)).Methods(http.MethodPut)

	// handle delete request for /api/v1/coupon/{id}
	router.HandleFunc(fmt.Sprintf("%s/%s/coupon/{id}", server.apiBaseUrl, server.apiVersion), createHttpHandler(server.handleCoupon)).Methods(http.MethodDelete)

	log.Println("Starting API Server on", server.listenAddress)

	http.ListenAndServe(server.listenAddress, router)
}


func (server *APIServer) handleCoupon(writer http.ResponseWriter, request *http.Request) error {
	log.Printf("Coupon Handler, handling %s request", request.Method)

	switch request.Method {
		case http.MethodGet:
			return server.handleGetCoupon(writer, request)
		case http.MethodPost:
			return server.handleCreateCoupon(writer, request)
		case http.MethodPut:
			return server.handleUpdateCoupon(writer, request)
		case http.MethodDelete:
			return server.handleDeleteCoupon(writer, request)
	}

	return fmt.Errorf("method %s not allowed", request.Method)
}


func (server * APIServer) handleGetCoupon(writer http.ResponseWriter, request *http.Request) error {
	log.Println("Get Coupon Handler")
	pathVariables := mux.Vars(request)
	log.Println("Path Variables", pathVariables)
	coupon := NewCoupon("10PERCENT", "percentage", 10.0)
	return WriteJsonResponse(writer, http.StatusOK, 
		Response{
			Success: true, 
			Data: coupon,
			Message: "Coupon Retrieved Successfully",
		},
	)
}


func (server *APIServer) handleCreateCoupon(writer http.ResponseWriter, request *http.Request) error {
	log.Println("Create Coupon Handler")
	return WriteJsonResponse(writer, http.StatusCreated,
		Response{
			Success: true,
			Message: "Coupon Created Successfully",
			Data: NewCoupon("20PERCENT", "percentage", 20.0),
		},
	)
}


func (server *APIServer) handleUpdateCoupon(writer http.ResponseWriter, request *http.Request) error {
	log.Println("Update Coupon Handler")

	log.Println("Path Variables", mux.Vars(request))

	return WriteJsonResponse(writer, http.StatusOK, Response{
		Success: true,
		Message: "Coupon Updated Successfully",
		Data: NewCoupon("30PERCENT", "percentage", 30.0),
	})
}

func (server *APIServer) handleDeleteCoupon(writer http.ResponseWriter, request *http.Request) error {
	log.Println("Delete Coupon Handler")
	log.Println("Path Variables", mux.Vars(request))
	return WriteJsonResponse(writer, http.StatusOK, Response{
		Success: true,
		Message: "Coupon Deleted Successfully",
		Data: nil,
	})
}