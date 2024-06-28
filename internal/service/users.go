package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type UserData struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	ShopID    string `json:"shopid"`
	ShopName  string `json:"shopname"`
}

type ShopData struct {
	ShopID string   `json:"shopid"`
	Data   UserData `json:"data"`
}

func RetrieveUsers(rows *sql.Rows) (string, error) {
	shopMap := make(map[string]ShopData)
	// Iterate through the result set
	for rows.Next() {
		var user UserData
		if err := rows.Scan(&user.Firstname, &user.Lastname, &user.ShopID, &user.ShopName); err != nil {
			return "", fmt.Errorf("error scanning row: %v", err)
		}
		shopMap[user.ShopID] = ShopData{
			ShopID: user.ShopID,
			Data:   user,
		}
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating row: %v", err)
	}

	// Convert the map to the desired JSON structure
	result, err := json.MarshalIndent(shopMap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling result: %v", err)
	}

	return string(result), nil
}
