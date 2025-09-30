package main

import (
	"github.com/gin-gonic/gin"

	"open-illustrations-go/config"
	"open-illustrations-go/routes"
)

func main() {
	config.InitDatabase()
	config.InitMinio()

	r := gin.Default()
	routes.RegisterRoutes(r)

	r.Run(":8080")
}
