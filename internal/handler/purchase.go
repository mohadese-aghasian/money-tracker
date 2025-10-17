package handler

import (
	"money-tracker/internal/dto"
	"money-tracker/internal/usecase"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	PurchaseUC *usecase.PurchaseUseCase
}

func NewPurchaseHandler(uc *usecase.PurchaseUseCase) *PurchaseHandler {
	return &PurchaseHandler{PurchaseUC: uc}
}

// @Summary Create a new purchase
// @Description Creates a new purchase with the provided data.
// @Tags purchase
// @Accept json
// @Produce json
// @Param request body dto.AddPurchaseInput true "purchase creation request"
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.Response "Bad Request"
// @Security BearerAuth
// @Router /api/v0/system/purchase [post]
func (h *PurchaseHandler) CreatepurchaseHandler(c *gin.Context) {
	var req dto.AddPurchaseInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error validation input", "response": err.Error()})
		return
	}

	purchase, err := h.PurchaseUC.Add(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"response": err.Error(),
			"message":  "error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"response": purchase,
		"message":  "created",
	})
	return

}

// @Summary Get all purchases
// @Description Retrieves all purchases.
// @Tags purchase
// @Produce json
// @Param start query int false "Start offset for pagination"
// @Param limit query int false "Limit number of records"
// @Param order_by query string false "Column to order by (default: id)"
// @Param sort query string false "Sort order: ASC or DESC"
// @Param reason query string false "Filter by reason"
// @Param id query int false "Filter by ID"
// @Param category_id query int false "Filter"
// @Param status_id query int false "Filter by StatusID"
// @Param tag_ids query []int false "Filter by Tag IDs (comma-separated)"
// @Success 200 {object} dto.GetResponse
// @Failure 400 {object} dto.ErrorGetResponse "Bad Request"
// @Security BearerAuth
// @Router /api/v0/system/purchase [get]
func (h *PurchaseHandler) GetAllPurchaseHandler(c *gin.Context) {
	var req dto.PurchaseFindAll
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error validation input", "response": err.Error()})
		return
	}

	purchases, count, err := h.PurchaseUC.Get(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "purchase not found",
			"response": purchases,
			"count":    count,
			"err":      err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "purchases found",
		"response": purchases,
		"count":    count,
	})
	return

}

// @Summary Update a purchase
// @Description Updates an existing purchase with new data.
// @Tags purchase
// @Accept json
// @Produce json
// @Param request body dto.UpdatePurchaseInput true "Category update request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response "Bad Request"
// @Security BearerAuth
// @Router /api/v0/system/purchase [put]
func (h *PurchaseHandler) UpdatePurchaseHandler(c *gin.Context) {
	var req dto.UpdatePurchaseInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error validation input", "response": err.Error()})
		return
	}

	// req.Slug =

	purchase, err := h.PurchaseUC.Update(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error validation input", "response": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "updated successfully",
		"response": purchase,
	})
	return

}

// @Summary Delete a purchase
// @Description Deletes a purchase by its ID.
// @Tags purchase
// @Produce json
// @Param id path int true "purchase ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response "Bad Request"
// @Security BearerAuth
// @Router /api/v0/system/purchase/{id} [delete]
func (h *PurchaseHandler) DeleteHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": "Invalid ID format", "message": "validation error"})
		return
	}

	if err := h.PurchaseUC.Remove(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "remove purchase failed!", "response": ""})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "removed successfully",
		"response": "",
	})
	return

}
