package handler

import (
	"money-tracker/internal/dto"
	"money-tracker/internal/usecase"
	"money-tracker/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	CategoryUC *usecase.CategoryUseCase
}

func NewCategoryHandler(uc *usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{CategoryUC: uc}
}


// @Summary Create a new category
// @Description Creates a new category with the provided data.
// @Tags category
// @Accept json
// @Produce json
// @Param request body dto.CreateCategoryRequest true "Category creation request"
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.Response "Bad Request"
// @Security BearerAuth
// @Router /api/v0/system/category [post]
func (h *CategoryHandler) CreateCategoryHandler(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error validation input", "response": err.Error()})
		return
	}

	category, err := h.CategoryUC.Add(dto.AddCategoryInput{
		Title:    req.Title,
		StatusID: req.StatusID,
		// TagIDs:   req.TagIDs,
		Slug:  utils.GenerateSlugUnicode(req.Title),
		Color: req.Color,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"response": err.Error(),
			"message":  "error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"response": category,
		"message":  "created",
	})
	return

}

// @Summary Get all categories
// @Description Retrieves all categories.
// @Tags category
// @Produce json
// @Param start query int false "Start offset for pagination"
// @Param limit query int false "Limit number of records"
// @Param order_by query string false "Column to order by (default: id)"
// @Param sort query string false "Sort order: ASC or DESC"
// @Param title query string false "Filter by title"
// @Param id query int false "Filter by ID"
// @Param parent_id query int false "Filter by ParentID"
// @Param status_id query int false "Filter by StatusID"
// @Param slug query string false "Filter by slug"
// @Param tag_ids query []int false "Filter by Tag IDs (comma-separated)"
// @Success 200 {object} dto.GetResponse
// @Failure 400 {object} dto.ErrorGetResponse "Bad Request"
// @Security BearerAuth
// @Router /api/v0/system/category [get]
func (h *CategoryHandler) GetAllCategoryHandler(c *gin.Context) {
	var req dto.ListCategoriesInput
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error validation input", "response": err.Error()})
		return
	}

	cateories, count, err := h.CategoryUC.Get(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "categories not found",
			"response": cateories,
			"count":    count,
			"err":      err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "categories found",
		"response": cateories,
		"count":    count,
	})
	return

}

// @Summary Get all categories
// @Description Retrieves all categories.
// @Tags category
// @Produce json
// @Param start query int false "Start offset for pagination"
// @Param limit query int false "Limit number of records"
// @Param order_by query string false "Column to order by (default: id)"
// @Param sort query string false "Sort order: ASC or DESC"
// @Param title query string false "Filter by title"
// @Param id query int false "Filter by ID"
// @Param parent_id query int false "Filter by ParentID"
// @Param slug query string false "Filter by slug"
// @Param tag_ids query []int false "Filter by Tag IDs (comma-separated)"
// @Success 200 {object} dto.GetResponse
// @Failure 400 {object} dto.ErrorGetResponse "Bad Request"
// @Security BearerAuth
// @Router /api/v0/system/category [get]
func (h *CategoryHandler) GetAllPublicCategoryHandler(c *gin.Context) {
	var req dto.ListCategoriesInput
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error validation input", "response": err.Error()})
		return
	}

	req.StatusID = 1

	cateories, count, err := h.CategoryUC.Get(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "categories not found",
			"response": cateories,
			"count":    count,
			"err":      err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "categories found",
		"response": cateories,
		"count":    count,
	})
	return

}

// @Summary Update a category
// @Description Updates an existing category with new data.
// @Tags category
// @Accept json
// @Produce json
// @Param request body dto.UpdateCategoryRequest true "Category update request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response "Bad Request"
// @Security BearerAuth
// @Router /api/v0/system/category [put]
func (h *CategoryHandler) UpdateCategoryHandler(c *gin.Context) {
	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error validation input", "response": err.Error()})
		return
	}

	// req.Slug =

	category, err := h.CategoryUC.Update(dto.UpdateCategoryInput{
		ID:       req.ID,
		Title:    req.Title,
		StatusID: req.StatusID,
		TagIDs:   req.TagIDs,
		Color:    req.Color,
		Slug:     utils.GenerateSlugUnicode(req.Title),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error validation input", "response": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "updated successfully",
		"response": category,
	})
	return

}

// @Summary Delete a category
// @Description Deletes a category by its ID.
// @Tags category
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response "Bad Request"
// @Security BearerAuth
// @Router /api/v0/system/category/{id} [delete]
func (h *CategoryHandler) DeleteHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": "Invalid ID format", "message": "validation error"})
		return
	}

	if err := h.CategoryUC.Remove(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "remove category failed!", "response": ""})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "removed successfully",
		"response": "",
	})
	return

}
