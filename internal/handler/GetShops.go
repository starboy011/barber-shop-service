package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/starboy011/barber-shop-service/internal/db"
	"github.com/starboy011/barber-shop-service/internal/service"
)

type User struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	ShopID    string `json:"shopid"`
	ShopName  string `json:"shopname"`
}

func GetShops(w http.ResponseWriter, r *http.Request) {
	database, err := db.InitUsersDB()
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to connect to database: %v", err), http.StatusInternalServerError)
		return
	}
	defer database.Close()
	rows, err := database.Query("SELECT firstname, lastname, shopid, shopname FROM users")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying database: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users, err := service.RetrieveUsers(rows)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving users: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert users to JSON
	usersJSON, err := json.Marshal(users)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshaling users to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(usersJSON)
}
