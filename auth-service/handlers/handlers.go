package handlers

import (
	"auth-service/internal/config"
	"auth-service/internal/database"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type FavouriteItem struct {
	UserID    string    `db:"user_id" json:"user_id"`
	ProductID string    `db:"product_id" json:"product_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func AddToFavourites(c *gin.Context, userID string) {
	config := config.GetPostgresConfig()
	db, _ := database.NewPostgresDB(config)
	var input struct {
		ProductID string `bd:"product_id" json:"product_id"`
	}
	err := c.ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	fmt.Println("Adding to favourites:", userID, input)
	query := `INSERT INTO users.user_favourites (user_id, product_id) VALUES ($1, $2) RETURNING user_id, product_id, created_at`
	var item FavouriteItem
	err = db.QueryRowContext(c, query, userID, input.ProductID).Scan(&item.UserID, &item.ProductID, &item.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println("Error adding to favourites:", err)
		return
	}
	c.JSON(http.StatusOK, item)

}

func GetFavouritesByUserID(c *gin.Context, userID string) {
	config := config.GetPostgresConfig()
	db, _ := database.NewPostgresDB(config)

	var favouriteItems []FavouriteItem
	query := `SELECT user_id, product_id, created_at FROM users.user_favourites WHERE user_id = $1 ORDER BY created_at ASC`
	err := db.SelectContext(c, &favouriteItems, query, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println("Error fetching favourites:", err)
		return
	}
	c.JSON(http.StatusOK, favouriteItems)
}

func DeleteFromFavourites(c *gin.Context, userID, productID string) {
	config := config.GetPostgresConfig()
	db, _ := database.NewPostgresDB(config)

	query := `DELETE FROM users.user_favourites WHERE user_id = $1 AND product_id = $2`
	_, err := db.ExecContext(c, query, userID, productID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println("Error deleting from favourites:", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
