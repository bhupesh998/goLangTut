package routes

import (
	"go-resturant-management/controllers"
	"github.com/gin-gonic/gin"
	_ "go.mongodb.org/mongo-driver/mongo"
)

func OrderRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/orders", controllers.GetOrders())
	incomingRoutes.GET("/orders/:order_id", controllers.GetOrder())
	incomingRoutes.POST("/orders", controllers.CreateOrder())
	incomingRoutes.PATCH("/orders/:order_id", controllers.UpdateOrder())
}