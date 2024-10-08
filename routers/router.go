package routers

import (
	"UnitUser/controller"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine){
	router.POST("/signup", controller.UserSignUp)
	router.POST("/signin", controller.UserSignIn)
	router.GET("/userlist", controller.UserShow)
	router.PATCH("/user/edit/:id", controller.EditUser)
}