// @title money tracker API
// @version 1.0
// @description This is a sample API with Gin and Swagger.
// @host localhost:4011
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"money-tracker/internal/config"
	"money-tracker/internal/routes"

	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	_ "money-tracker/cmd/api/docs"
)

func init() {
	config.LoadEnvVariables()
	config.ConnectToDB()

}

func main() {

	mode := os.Getenv("GIN_MODE")
	gin.SetMode(mode)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Static("/public/images", "./uploads/images")

	//---------Routes--------
	routes.AdminRoutes("/api/v0/admin", router)
	routes.AuthRoutes("/api/v0", router)
	routes.MainRoutes("/api/v0/system", router)
	//---------------------

	port := os.Getenv("PORT")
	if port != "" {
		port = "4011"
	}
	router.Run(":" + port)

}
