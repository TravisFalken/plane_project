package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//All of the routers

func main() {
	log.Println("Server starting...") //for testing
	//creating the router
	router := mux.NewRouter().StrictSlash(true)

	//the routes
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/home", homePage).Methods("GET")
	router.HandleFunc("/signup", showSignUp).Methods("GET")
	router.HandleFunc("/signup", signUp).Methods("POST")
	router.HandleFunc("/login", loginPage).Methods("GET")
	router.HandleFunc("/login", loginUser).Methods("POST")
	router.HandleFunc("/logout", logout).Methods("GET")
	router.HandleFunc("/createPlan", createPlanPage).Methods("GET")
	router.HandleFunc("/createPlan", createPlan).Methods("POST")
	router.HandleFunc("/showPlan/{id}", showPlan).Methods("GET")
	router.HandleFunc("/showPlan/{id}", editPlan).Methods("POST")
	router.HandleFunc("/deletePlan/{id}", deleteUserPlan).Methods("POST")

	//Start Server
	log.Fatal(http.ListenAndServe(":8080", router))
}
