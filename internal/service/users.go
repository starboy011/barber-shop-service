package service

import (
	"database/sql"
	"fmt"
)

type User struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	ShopID    string `json:"shopid"`
	ShopName  string `json:"shopname"`
}

func RetrieveUsers(rows *sql.Rows) ([]User, error) {
	var users []User

	// Iterate through the result set
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Firstname, &user.Lastname, &user.ShopID, &user.ShopName); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating row: %v", err)
	}

	return users, nil
}
