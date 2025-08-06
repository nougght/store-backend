package handlers

// import (
// 	"database/sql"
// 	"fmt"
// 	"net/http"
// 	"product-service/internal/config"
// 	"product-service/internal/database"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/lib/pq"
// )

// type product struct {
// 	ID           string          `json:"id"`
// 	Name         string          `json:"name"`
// 	Description  string          `json:"description"`
// 	Price        float64         `json:"price"`
// 	CategoryId   string          `db:"category_id" json:"category_id"`
// 	Images       *[]string       `json:"images"`
// 	Quantity     float64         `json:"quantity"`
// 	Unit         string          `json:"unit"`
// 	Stock        int             `json:"stock"`
// 	CreationDate time.Time       `db:"created_at" json:"created_at"`
// 	UpdatedDate  time.Time       `db:"updated_at" json:"updated_at"`
// 	IsActive     bool            `db:"is_active" json:"is_active"`
// 	Weight       sql.NullFloat64 `json:"weight"`
// }

// func GetProducts(c *gin.Context) {
// 	config := config.GetPostgresConfig()
// 	db, err := database.NewPostgresDB(config)

// 	if err == nil {
// 		idsParam := c.Query("ids")
// 		var products []product
// 		// fmt.Println("params : ", idsParam)
// 		if idsParam == "" {

// 			err = db.SelectContext(c, &products, "SELECT * FROM products.products ")
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 				// return
// 			}

// 			c.JSON(http.StatusOK, products)
// 		} else {

// 			ids := strings.Split(idsParam, ",")
// 			products, err := ProductsByIDs(ids)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"errorrrr": err.Error()})
// 				fmt.Println("errorrrr", err)
// 				return
// 			}

// 			c.JSON(http.StatusOK, products)
// 		}
// 	} else {
// 		fmt.Println("ERROR", err)
// 		c.JSON(404, 1)
// 	}
// }

// func GetProductsByIDs(c *gin.Context) {
// 	idsParam := c.Query("ids")
// 	if idsParam == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ids parameter"})
// 		return
// 	}

// 	// Разделяем строку на массив ID
// 	ids := strings.Split(idsParam, ",")

// 	// Валидация UUID
// 	// for _, id := range ids {
// 	//     if _, err := uuid.Parse(id); err != nil {
// 	//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
// 	//         return
// 	//     }
// 	// }

// 	// Получаем товары из БД
// 	products, err := ProductsByIDs(ids)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, products)
// }

// func ProductsByIDs(ids []string) ([]product, error) {

// 	config := config.GetPostgresConfig()
// 	db, _ := database.NewPostgresDB(config)

// 	if len(ids) == 0 {
// 		return []product{}, nil
// 	}

// 	// Создаем плейсхолдеры ($1, $2, ...)
// 	placeholders := make([]string, len(ids))
// 	args := make([]interface{}, len(ids))
// 	for i, id := range ids {
// 		placeholders[i] = fmt.Sprintf("$%d", i+1)
// 		args[i] = id
// 	}

// 	query := fmt.Sprintf(`
//     SELECT *
//     FROM products.products
//     WHERE id = ANY($1::uuid[])
//     ORDER BY array_position($1::uuid[], id)`,
// 	) // была проблема с порядком, поэтому добавил ORDER BY array_position

// 	var products []product
// 	fmt.Println(len(ids))
// 	fmt.Println("query: ", query, "args: ", args)
// 	err := db.Select(&products, query, pq.Array(ids))

// 	if err == nil {
// 		fmt.Println("Products found: ", products)
// 		return products, nil
// 	} else {
// 		fmt.Println("Ошибка получения продуктов по id:", query, args, ids)
// 		return []product{}, err
// 	}
// }
