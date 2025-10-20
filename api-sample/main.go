package main

import (
	"api-sample/config"
	"api-sample/database"
	"api-sample/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func main() {
	cfg := config.Load()

	shutdown, err := config.InitOTel(cfg)
	if err != nil {
		log.Fatal("Failed to initialize OpenTelemetry:", err)
	}
	defer shutdown()

	log.Println("OpenTelemetry initialized")

	db, err := database.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer db.Disconnect()

	log.Println("Connected to MongoDB")

	userHandler := handlers.NewUserHandler(db)

	router := mux.NewRouter()
	router.Use(otelmux.Middleware(cfg.ServiceName))
	
	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
	router.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	log.Printf("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal("Server failed:", err)
	}
}
