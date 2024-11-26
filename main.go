package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

const jsonFile = "clients.json"

// Carregar clientes do arquivo JSON
func loadClients() ([]Client, error) {
	var clients []Client

	file, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &clients)
	if err != nil {
		return nil, err
	}

	return clients, nil
}

// Salvar clientes no arquivo JSON
func saveClients(clients []Client) error {
	file, err := json.MarshalIndent(clients, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(jsonFile, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

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

func GetClients(w http.ResponseWriter, r *http.Request) {
	clients, err := loadClients()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(clients)
}

func GetClient(w http.ResponseWriter, r *http.Request) {
	clientId := mux.Vars(r)["clientId"]

	intClientId, err := strconv.Atoi(clientId)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	clients, err := loadClients()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, client := range clients {
		if client.ID == intClientId {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(client)
			return
		}
	}

	http.Error(w, "Client not found", http.StatusNotFound)
}

func CreateClient(w http.ResponseWriter, r *http.Request) {
	var client Client

	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	clients, err := loadClients()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client.ID = len(clients) + 1
	clients = append(clients, client)

	err = saveClients(clients)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(client)
}

func UpdateClient(w http.ResponseWriter, r *http.Request) {
	clientId := mux.Vars(r)["clientId"]

	intClientId, err := strconv.Atoi(clientId)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	clients, err := loadClients()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, client := range clients {
		if client.ID == intClientId {
			var updatedClient Client
			err := json.NewDecoder(r.Body).Decode(&updatedClient)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if updatedClient.Name != "" {
				client.Name = updatedClient.Name
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

			clients[i] = client

			err = saveClients(clients)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(client)
			return
		}
	}

	http.Error(w, "Client not found", http.StatusNotFound)
}

func DeleteClient(w http.ResponseWriter, r *http.Request) {
	clientId := mux.Vars(r)["clientId"]

	intClientId, err := strconv.Atoi(clientId)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	clients, err := loadClients()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, client := range clients {
		if client.ID == intClientId {
			clients = append(clients[:i], clients[i+1:]...)
			err = saveClients(clients)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Client not found", http.StatusNotFound)
}
