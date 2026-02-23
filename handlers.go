// http handlers

package main

import (
	"encoding/json"
	"net/http"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = db.NewInsert().Model(&user).Exec(r.Context())

	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []User
	err := db.NewSelect().Model(&users).Scan(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	result, err := db.NewDelete().Model((*User)(nil)).Where("id=?", id).Exec(r.Context())
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Write([]byte("User deleted successfully"))
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	result, err := db.NewUpdate().Model(&user).Where("id=?", id).Exec(r.Context())
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Write([]byte("User updated successfully"))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createUserHandler(w, r)
	case http.MethodGet:
		getUsersHandler(w, r)
	case http.MethodDelete:
		deleteUserHandler(w, r)
	case http.MethodPut:
		updateUserHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}
