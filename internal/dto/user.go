package dto

type ListUsersInput struct {
	ID          uint   `form:"id"`
	UserName    string `form:"username"`
	LevelManage int8   `form:"level_manage"`
	StatusID    uint   `form:"status_id"`
	Start       int    `form:"start"`
	Limit       int    `form:"limit"`
	OrderBy     string `form:"order_by"`
	Sort        string `form:"sort"`
}

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	ID          uint   `json:"id" binding:"required"`
	UserName    string `json:"username"`
	LevelManage int8   `json:"level_manage"`
	StatusID    uint   `json:"status_id"`
	Password    string `json:"password"`
}

type AddUserInput struct {
	UserName    string
	Password    string
	LevelManage int8
	StatusID    uint
	Email       string
	Mobile      string
}

type RegisterRequest struct {
	Name     string `json:"name"`
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	StatusID uint   `json:"status_id"`
}

type LogoutInput struct {
	Token  string `json:"token"`
	UserId *uint  `json:"user_id"`
}
