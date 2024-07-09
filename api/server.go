package api

import (
	"context"
	"listservice/internal/handlers"
	"listservice/internal/repositories"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func StartServer() {

	repositories.ConnectDB()
	defer func() {
		if err := repositories.MongoClient.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	listCollection := repositories.MongoClient.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLLECTION_NAME"))

	listHandler := handlers.ListHandler{MongoCollection: listCollection}

	router := mux.NewRouter()

	router.HandleFunc("/list/{id}", listHandler.GetListByID).Methods("GET")
	router.HandleFunc("/list", listHandler.CreateList).Methods("POST")
	router.HandleFunc("list/{id}", listHandler.DeleteList).Methods("DELETE")
	router.HandleFunc("list/{id}/drinks", listHandler.AddDrinkToList).Methods("POST")
	router.HandleFunc("list/{id}/drinks", listHandler.RemoveDrinkFromList).Methods("POST")
	router.HandleFunc("list/{id}/drinks", listHandler.AddCollaboratorToList).Methods("POST")
	router.HandleFunc("list/{id}/drinks", listHandler.RemoveCollaboratorFromList).Methods("POST")
	router.HandleFunc("list/{id}/private", listHandler.MakeListPrivate).Methods("PATCH")
	router.HandleFunc("list/{id}/public", listHandler.MakeListPublic).Methods("PATCH")

	log.Fatal(http.ListenAndServe(":8081", router))
}
