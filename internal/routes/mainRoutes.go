package routes

import (
	"money-tracker/internal/config"
	"money-tracker/internal/handler"
	// "money-tracker/internal/middleware"
	"money-tracker/internal/repository"
	"money-tracker/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Category *handler.CategoryHandler
	User     *handler.UserHandler
	Tag      *handler.TagHandler
	Purchase *handler.PurchaseHandler
}

func buildHandlers() *Handlers {
	// repositories

	repoCat := repository.NewRepositoryGorm(config.DB)
	repoUser := repository.NewUserRepositoryGorm(config.DB)
	repoToken := repository.NewUserTokenRepositoryGorm(config.DB)
	repoTag := repository.NewTagRepoGorm(config.DB)
	repoPurchase := repository.NewPurchaseRepo(config.DB)
	// use cases

	ucCategory := usecase.NewCategoryUseCase(repoCat)
	ucUser := usecase.NewUserUseCase(repoUser, repoToken)
	ucTag := usecase.NewTagUseCase(repoTag)
	ucPurchase := usecase.NewPurchaseUseCase(repoPurchase, repoTag, repoCat)

	// handlers
	return &Handlers{
		Category: handler.NewCategoryHandler(ucCategory),
		User:     handler.NewUserHandler(ucUser),
		Tag:      handler.NewTagHandler(ucTag),
		Purchase: handler.NewPurchaseHandler(ucPurchase),
	}

}

func MainRoutes(base string, router *gin.Engine) {
	h := buildHandlers()

	api := router.Group(base)
	// api.Use(middleware.AuthMiddleware([]int8{2}))
	{
		api.POST("/tag", h.Tag.CreateTagHandler)
		api.GET("/tag", h.Tag.GetAllTagsHandler)
		api.PUT("/tag", h.Tag.UpdateTagHandler)
		api.DELETE("/tag/:id", h.Tag.DeleteHandler)

		api.GET("/category", h.Category.GetAllPublicCategoryHandler)
		api.POST("/category", h.Category.CreateCategoryHandler)
		api.PUT("/category", h.Category.UpdateCategoryHandler)
		api.DELETE("/category/:id", h.Category.DeleteHandler)

		api.GET("/purchase", h.Purchase.GetAllPurchaseHandler)
		api.POST("/purchase", h.Purchase.CreatepurchaseHandler)
		api.PUT("/purchase", h.Purchase.UpdatePurchaseHandler)
		api.DELETE("/purchase/:id", h.Purchase.DeleteHandler)

	}

}
