package handler

import (
	"money-tracker/internal/dto"
	"money-tracker/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	TagUC *usecase.TagUseCase
}

func NewTagHandler(uc *usecase.TagUseCase) *TagHandler {
	return &TagHandler{TagUC: uc}
}

// @Summary Create a new tag
// @Description Creates a new tag with the provided data.
// @Tags tag
// @Accept json
// @Produce json
// @Param request body dto.CreateTagRequest true "Article creation request"
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.Response "Bad Request"
// @Security BearerAuth
// @Router /api/v0/system/tag [post]
func (h *TagHandler) CreateTagHandler(c *gin.Context) {
	var req dto.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error validation input", "response": err.Error()})
		return
	}

	tag, err := h.TagUC.Add(req.Title, req.StatusID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "created",
		"response": tag,
	})
	return

}

// @Summary Get all tags
// @Description Retrieves all tags.
// @Tags tag
// @Produce json
// @Param start query int false "Start"
// @Param limit query int false "limit"
// @Param order_by query string false "orderBy"
// @Param sort query string false "Sort"
// @Param id query int false "id"
// @Param status_id query int false "status_id"
// @Param title query string false "title"
// @Success 200 {object} dto.GetResponse
// @Failure 400 {object} dto.ErrorGetResponse "Bad Request"
// @Security BearerAuth
// @Router /api/v0/system/tag [get]
func (h *TagHandler) GetAllTagsHandler(c *gin.Context) {
	var req dto.ListTagsInput
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error validation input", "response": err.Error()})
		return
	}

	tags, count, err := h.TagUC.Get(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "tags not found",
			"response": tags,
			"count":    count,
			"err":      err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "tags found",
		"response": tags,
		"count":    count,
	})
	return

}

// @Summary Update a tag
// @Description Updates an existing tag with new data.
// @Tags tag
// @Accept json
// @Produce json
// @Param request body dto.UpdateTagRequest true "Tag update request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response "Bad Request"
// @Security BearerAuth
// @Router /api/v0/system/tag [put]
func (h *TagHandler) UpdateTagHandler(c *gin.Context) {
	var req dto.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error validation input", "response": err.Error()})
		return
	}

	tag, err := h.TagUC.Update(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": err.Error(),
			"message":  "error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "updated successfully",
		"response": tag,
	})
	return

}

// @Summary Delete a tag
// @Description Deletes a tag by its ID.
// @Tags tag
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response "Bad Request"
// @Security BearerAuth
// @Router /api/v0/system/tag/{id} [delete]
func (h *TagHandler) DeleteHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format", "response": ""})
		return
	}

	if err := h.TagUC.Remove(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "remove tag failed!",
			"response": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "removed successfully",
		"response": "",
	})
	return

}