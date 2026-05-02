package dto

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Profile  string `json:"profile" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Profile  string `json:"profile" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type UpdateProfileRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Profile string `json:"profile"`
}

type UserResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Profile string `json:"profile"`
	Role    string `json:"role"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type DeleteUserRequest struct {
	ID int64 `json:"id" binding:"required"`
}
