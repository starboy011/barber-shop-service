package handler

import (
	"fmt"

	"github.com/starboy011/barber-shop-service/internal/db"
)

func GetShopName(shopId string) (string, error) {
	database, err := db.InitUsersDB()
	if err != nil {
		fmt.Printf("Error in connecting Db for shopId: %s. Error: %v\n", shopId, err)
		return "", err
	}
	defer database.Close()
	shopName, err := db.GetShopNameById(database, shopId)
	if err != nil {
		fmt.Printf("Error fetching shop name for Shop ID: %s. Error: %v\n", shopId, err)
		return "", err
	}

	return shopName, nil
}
