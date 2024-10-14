package main

import (
	"UnitUser/database"
	"UnitUser/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.CreateDB()
	router := gin.Default()
	routers.UserRoutes(router)
	router.Run(":8081")
}
