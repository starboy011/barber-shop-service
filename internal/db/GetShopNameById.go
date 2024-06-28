package db

import (
	"database/sql"
)

func GetShopNameById(database *sql.DB, shopId string) (string, error) {
	var shopName string
	query := "SELECT shopname FROM users WHERE shopid = $1"

	err := database.QueryRow(query, shopId).Scan(&shopName)
	if err != nil {
		return "", err
	}
	return shopName, nil
}
