package handlers

import (
	"encoding/json"
	"fmt"
	"listservice/internal/models"
	"listservice/internal/repositories"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type ListHandler struct {
	MongoCollection *mongo.Collection
}

func (listHandler ListHandler) GetListByID(w http.ResponseWriter, r *http.Request) {
	listID := mux.Vars(r)["id"]

	listRepository := repositories.ListRepository{MongoCollection: listHandler.MongoCollection}
	list, err := listRepository.GetListByID(listID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("unable to retrieve list:", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		log.Println(err)
	}
}

func (listHandler ListHandler) CreateList(w http.ResponseWriter, r *http.Request) {
	var list models.List
	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid request body:", err)
	}

	list.ID = uuid.NewString()

	listRepository := repositories.ListRepository{MongoCollection: listHandler.MongoCollection}

	insertID, err := listRepository.CreateList(&list)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("unable to create new list:", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, insertID)

	log.Println("List created with ID", insertID)

}

func (listHandler ListHandler) DeleteList(w http.ResponseWriter, r *http.Request) {
	listID := mux.Vars(r)["id"]

	listRepository := repositories.ListRepository{MongoCollection: listHandler.MongoCollection}

	deletedCount, err := listRepository.DeleteList(listID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("unable to delete list:", err)
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("Deleted %d lists\n", deletedCount)
}

func (listHandler ListHandler) MakeListPublic(w http.ResponseWriter, r *http.Request) {
	listID := mux.Vars(r)["id"]

	listRepository := repositories.ListRepository{MongoCollection: listHandler.MongoCollection}

	updatedCount, err := listRepository.MakeListPublic(listID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("unable to make list public:", err)
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("Made %d lists public", updatedCount)
}

func (listHandler ListHandler) MakeListPrivate(w http.ResponseWriter, r *http.Request) {
	listID := mux.Vars(r)["id"]

	listRepository := repositories.ListRepository{MongoCollection: listHandler.MongoCollection}

	updatedCount, err := listRepository.MakeListPublic(listID)

	if err != nil {
		log.Println("unable to make list private:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("Made %d lists private", updatedCount)
}

func (listHandler ListHandler) RemoveDrinkFromList(w http.ResponseWriter, r *http.Request) {
	listID := mux.Vars(r)["id"]

	listRepository := repositories.ListRepository{MongoCollection: listHandler.MongoCollection}

	var requestBody map[string]string

	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		log.Println("error decoding drinkID:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	drinkID := requestBody["drinkID"]

	updatedCount, err := listRepository.RemoveDrinkFromList(listID, drinkID)

	if err != nil {
		log.Println("error removing drink from list:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("removed a drink from %d lists", updatedCount)
}

func (listHandler ListHandler) AddDrinkToList(w http.ResponseWriter, r *http.Request) {
	listID := mux.Vars(r)["id"]

	listRepository := repositories.ListRepository{MongoCollection: listHandler.MongoCollection}

	var requestBody map[string]string

	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		log.Println("error decoding drinkID:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	drinkID := requestBody["drinkID"]

	updatedCount, err := listRepository.AddDrinkToList(listID, drinkID)

	if err != nil {
		log.Println("error adding drink to list:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("added a drink to %d lists", updatedCount)
}

func (listHandler ListHandler) AddCollaboratorToList(w http.ResponseWriter, r *http.Request) {
	listID := mux.Vars(r)["id"]

	listRepository := repositories.ListRepository{MongoCollection: listHandler.MongoCollection}

	var requestBody map[string]string

	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		log.Println("error decoding userID:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := requestBody["userID"]

	updatedCount, err := listRepository.AddCollaborator(listID, userID)

	if err != nil {
		log.Println("error adding collaborator to list:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("added a collborator to %d lists", updatedCount)
}

func (listHandler ListHandler) RemoveCollaboratorFromList(w http.ResponseWriter, r *http.Request) {
	listID := mux.Vars(r)["id"]

	listRepository := repositories.ListRepository{MongoCollection: listHandler.MongoCollection}

	var requestBody map[string]string

	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		log.Println("error decoding userID:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := requestBody["userID"]

	updatedCount, err := listRepository.RemoveCollaborator(listID, userID)

	if err != nil {
		log.Println("error removing collaborator from list:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("removed a collaborator from %d lists", updatedCount)
}
