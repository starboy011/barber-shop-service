package handler

import (
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/starboy011/barber-shop-service/internal/db"
	"github.com/starboy011/barber-shop-service/internal/service"
)

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

	usersJSON, err := service.RetrieveUsers(rows)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving users: %v", err), http.StatusInternalServerError)
		return
	}
	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(usersJSON))
}
