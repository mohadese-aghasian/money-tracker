package handler

import (
	"money-tracker/internal/constants"
	"money-tracker/internal/dto"
	"money-tracker/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUC *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{UserUC: uc}
}

// @Summary Signup
// @Description Creates a new user with the provided data.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "sign up request"
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.Response "Bad Request"
// @Router /api/v0/auth/signup [post]
func (h *UserHandler) RegisterHandler(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": err.Error(), "message": "validation error"})
		return
	}

	_, err := h.UserUC.Add(dto.AddUserInput{
		UserName:    req.UserName,
		Password:    req.Password,
		LevelManage: 2,
		StatusID:    1,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": err.Error(), "message": "creation error"})
		return
	}

	// token, err := h.UserUC.Login(dto.LoginRequest{UserName: user.UserName, Password: user.Password})
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"response": err.Error(), "message": "error"})
	// 	return
	// }

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		// "user":    user,
		// "response": token,
		"response": "",
	})
}

// @Summary Login
// @Description login with the provided data.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "sign up request"
// @Success 201 {object} dto.LoginReponse
// @Failure 400 {object} dto.Response "Bad Request"
// @Router /api/v0/auth/login [post]
func (h *UserHandler) LoginHandler(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": err.Error(), "message": "error validation"})
		return
	}

	token, level_manage, user_id, err := h.UserUC.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"response": err.Error(), "message": "error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged in", "response": gin.H{
		"token": token, "levelManage": level_manage, "userId": user_id,
	}})
}

// @Summary Get all users
// @Description Retrieves all users with optional filters and pagination.
// @Tags user
// @Produce json
// @Param id query int false "Filter by User ID"
// @Param name query string false "Filter by Name"
// @Param username query string false "Filter by Username"
// @Param email query string false "Filter by Email"
// @Param mobile query string false "Filter by Mobile"
// @Param level_manage query int8 false "Filter by Level Manage"
// @Param status_id query int false "Filter by Status ID"
// @Param show_inactive query bool false "Show inactive users (true/false)"
// @Param start query int false "Start offset for pagination"
// @Param limit query int false "Limit number of records"
// @Param order_by query string false "Column to order by (default: id)"
// @Param sort query string false "Sort order: ASC or DESC"
// @Success 200 {object} dto.GetResponse
// @Failure 400 {object} dto.ErrorGetResponse "Bad Request"
// @Security BearerAuth
// @Router /api/v0/admin/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	var req dto.ListUsersInput
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error(), "message": "validation input", "count": 0})
		return
	}

	items, count, err := h.UserUC.Get(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":  "user not found",
			"response": items,
			"count":    count,
			"err":      err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "users found",
		"response": items,
		"count":    count,
	})
	return

}

// @Summary Signup
// @Description Creates a new user with the provided data.
// @Tags auth
// @Accept json
// @Produce json
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.Response "Bad Request"
// @Router /api/v0/auth/signup-admin [get]
func (h *UserHandler) SignAdminHandler(c *gin.Context) {

	_, err := h.UserUC.Add(dto.AddUserInput{
		UserName:    "money.admin",
		Password:    "12345",
		LevelManage: constants.LevelManageAdmin,
		StatusID:    1,
		Email:       "",
		Mobile:      "",
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": err.Error(), "message": "validation error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		// "user":    user,
		// "response": token,
		"response": "",
	})
}

// @Summary logout
// @Description Creates a new user with the provided data.
// @Tags auth
// @Accept json
// @Produce json
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.Response "Bad Request"
// @Security BearerAuth
// @Router /api/v0/auth/logout [get]
func (h *UserHandler) Logout(c gin.Context) {
	_, err := h.UserUC.Add(dto.AddUserInput{
		UserName:    "money.admin",
		Password:    "12345",
		LevelManage: constants.LevelManageAdmin,
		StatusID:    1,
		Email:       "",
		Mobile:      "",
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": err.Error(), "message": "validation error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		// "user":    user,
		// "response": token,
		"response": "",
	})
}
