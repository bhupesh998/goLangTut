package routes

import (
	"go-resturant-management/controllers"

	"github.com/gin-gonic/gin"
	_ "go.mongodb.org/mongo-driver/mongo"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:user_id", controllers.GetUser())
	incomingRoutes.POST("/users/singup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
}