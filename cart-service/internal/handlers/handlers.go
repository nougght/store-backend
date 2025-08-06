package handlers

// import (
// 	"cart-service/internal/config"
// 	"cart-service/internal/database"
// 	"context"
// 	"fmt"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/lib/pq"
// )

// func DeleteCartItems(c *gin.Context) {
// 	config := config.GetPostgresConfig()
// 	db, _ := database.NewPostgresDB(config)

// 	var input struct {
// 		Ids []string `json:"ids"`
// 	}
// 	err := c.ShouldBindJSON(&input)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
// 		return
// 	}
// 	query := `DELETE FROM carts.cart_items WHERE id = ANY($1)`
// 	_, err = db.ExecContext(c, query, pq.Array(input.Ids))

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		fmt.Println("Error deleting cart items:", err, "ids:", input)
// 		return
// 	}
// 	fmt.Println("cart items deleted")
// 	c.JSON(http.StatusOK, gin.H{"deleted ids": input.Ids, "status": "deleted"})
// }

// func AddToCart(c *gin.Context) {
// 	config := config.GetPostgresConfig()
// 	db, _ := database.NewPostgresDB(config)
// 	var item CartItem
// 	err := c.ShouldBindJSON(&item)
// 	if err != nil {
// 		fmt.Println("Ошибка добавления в корзину: ", err)
// 	}
// 	fmt.Println(item)
// 	query :=
// 		`INSERT INTO carts.cart_items (cart_id, product_id, quantity) VALUES ($1, $2, $3) returning id`
// 	err = db.QueryRow(query, item.CartID, item.ProductID, item.Quantity).Scan(&item.ID)
// 	if err == nil {
// 		fmt.Println(item.ID)
// 		c.JSON(http.StatusCreated, gin.H{"id": item.ID})
// 	} else {
// 		fmt.Println("addToCart Error: ", err)
// 		c.JSON(http.StatusBadRequest, err)
// 	}
// }

// func UpdateCartItemQuantity(c *gin.Context) {
// 	config := config.GetPostgresConfig()
// 	db, _ := database.NewPostgresDB(config)
// 	var input struct {
// 		Id       string `db:"id" json:"id"`
// 		Quantity int    `db:"quantity" json:"quantity"`
// 	}
// 	err := c.ShouldBindJSON(&input)

// 	if err != nil {
// 		fmt.Println("Ошибка добавления в корзину: ", err)
// 	}
// 	fmt.Println(input)
// 	query :=
// 		`UPDATE carts.cart_items SET quantity = $1 WHERE id = $2`
// 	created, err := db.QueryContext(c, query, input.Quantity, input.Id)
// 	if err == nil {
// 		c.JSON(http.StatusCreated, created)
// 	} else {
// 		c.JSON(http.StatusBadRequest, err)
// 		fmt.Println("addToCart Error: ", err)
// 	}
// }

// func CartByUserID(ctx context.Context, id string) ([]string, error) {
// 	config := config.GetPostgresConfig()
// 	db, _ := database.NewPostgresDB(config)

// 	var cartId []string
// 	err := db.Select(&cartId, "SELECT cart_id FROM carts.carts WHERE user_id = $1", id)
// 	if err == nil {
// 		return cartId, nil
// 	}
// 	return []string{}, err
// }

// func DeleteFromCartById(c *gin.Context, id string) {
// 	config := config.GetPostgresConfig()
// 	db, _ := database.NewPostgresDB(config)

// 	query := `DELETE FROM carts.cart_items WHERE id = $1`
// 	_, err := db.ExecContext(c, query, id)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	fmt.Println("end")
// 	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
// }

// func GetCartItemsbyCartID(c *gin.Context, id string) {
// 	config := config.GetPostgresConfig()
// 	db, _ := database.NewPostgresDB(config)

// 	var cartItems []CartItem
// 	err := db.SelectContext(c, &cartItems, "SELECT * FROM carts.cart_items WHERE cart_id = $1 ORDER BY added_at ASC", id)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	// fmt.Println("end")
// 	c.JSON(http.StatusOK, cartItems)
// }
