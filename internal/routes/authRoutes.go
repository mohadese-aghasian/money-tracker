package routes

import (
	"money-tracker/internal/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(base string, router *gin.Engine) {
	h := buildHandlers()

	auth := router.Group(base + "/auth")
	{
		auth.POST("/login", h.User.LoginHandler)
		auth.POST("/signup", h.User.RegisterHandler)
		auth.GET("/signup-admin", h.User.SignAdminHandler)
	}
}

func AdminRoutes(base string, router *gin.Engine) {

	h := buildHandlers()

	admin := router.Group(base)
	admin.Use(
		middleware.AuthMiddleware([]int8{1}))
	{
		admin.POST("/category", h.Category.CreateCategoryHandler)
		admin.GET("/category", h.Category.GetAllCategoryHandler)
		admin.PUT("/category", h.Category.UpdateCategoryHandler)
		admin.DELETE("/category/:id", h.Category.DeleteHandler)

		admin.GET("/users", h.User.GetAllUsers)

	}
}
