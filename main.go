package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Client struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CPF       string `json:"cpf"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Birthdate string `json:"birthdate"`
}

var clientDb = make(map[int]Client)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/client", GetClients).Methods("GET")
	router.HandleFunc("/client/{clientId}", GetClient).Methods("GET")
	router.HandleFunc("/client", CreateClient).Methods("POST")
	router.HandleFunc("/client/{clientId}", UpdateClient).Methods("PUT")
	router.HandleFunc("/client/{clientId}", DeleteClient).Methods("DELETE")

	fmt.Println("Server running on port 5000")
	http.ListenAndServe(":5000", router)
}

func CreateClient(w http.ResponseWriter, r *http.Request) {
	var client Client
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client.ID = len(clientDb) + 1
	clientDb[client.ID] = client

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(client)
}

func GetClients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(clientDb)
}

func GetClient(w http.ResponseWriter, r *http.Request) {
	var clientId = mux.Vars(r)["clientId"]

	intClientId, err := strconv.Atoi(clientId)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	client, exists := clientDb[intClientId]
	if !exists {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(client)
}

func UpdateClient(w http.ResponseWriter, r *http.Request) {
	var clientId = mux.Vars(r)["clientId"]

	intClientId, err := strconv.Atoi(clientId)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	client, exists := clientDb[intClientId]
	if !exists {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	var updatedClient Client
	err = json.NewDecoder(r.Body).Decode(&updatedClient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updatedClient.Name != "" {
		client.Name = updatedClient.Name
	}
	if updatedClient.CPF != "" {
		client.CPF = updatedClient.CPF
	}
	if updatedClient.Email != "" {
		client.Email = updatedClient.Email
	}
	if updatedClient.Phone != "" {
		client.Phone = updatedClient.Phone
	}
	if updatedClient.Birthdate != "" {
		client.Birthdate = updatedClient.Birthdate
	}

	clientDb[intClientId] = client

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(client)
}

func DeleteClient(w http.ResponseWriter, r *http.Request) {
	var clientId = mux.Vars(r)["clientId"]

	intClientId, err := strconv.Atoi(clientId)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	if _, exists := clientDb[intClientId]; exists {
		delete(clientDb, intClientId)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	http.Error(w, "Client not found", http.StatusNotFound)
}
