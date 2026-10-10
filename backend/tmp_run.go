package main

import (
	"fmt"
	"pos-backend/internal/database"
	"pos-backend/internal/models"
)

func main() {
	_, err := database.InitDB("sqlite", "", "pos.db")
	if err != nil {
		fmt.Println("DB err:", err)
		return
	}
	var product models.Product
	id := "a1b6b2f1-f3f8-419e-a912-18663a950975"
	err = database.DB.First(&product, "id = ?", id).Error
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Success:", product.Name)
	}
}
