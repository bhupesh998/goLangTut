package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bhupesh998/ecommproject/controller"
	"github.com/bhupesh998/ecommproject/database"
	"github.com/bhupesh998/ecommproject/middleware"
	"github.com/bhupesh998/ecommproject/routes"
	"github.com/gin-gonic/gin"
)

func main(){
	port := os.Getenv("PORT")
	fmt.Println("port value is", port)
	if port == "" {
		port = "8000"
	}
	fmt.Println("port value is", port)
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users" ))

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server is running on localhost:8000",
		})
	})

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())


	log.Fatal(router.Run(":"+port))
}