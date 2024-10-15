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
	// try not changing anything
	router.Run(":8085")
}
