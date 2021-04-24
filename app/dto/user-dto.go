package dto

type UpdateUserDTO struct {
	ID       uint64 `json:"id" form:"id" binding:"required"`
	Username string `json:"username" form:"username" binding:"required"`
	Fullname string `json:"fullname" form:"fullname" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
	Role     string `json:"role" form:"role" binding:"required"`
}
