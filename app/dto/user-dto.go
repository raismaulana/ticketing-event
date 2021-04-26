package dto

type UpdateUserDTO struct {
	ID       uint64 `json:"id" form:"id" binding:"required"`
	Username string `json:"username,omitempty" form:"username,omitempty"`
	Fullname string `json:"fullname" form:"fullname" binding:"required"`
	Email    string `json:"email,omitempty" form:"email,omitempty" binding:"email,omitempty"`
	Password string `json:"password,omitempty" form:"password,omitempty"`
	Role     string `json:"role,omitempty" form:"role,omitempty" `
}
