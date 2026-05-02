// dto/auth_dto.go
package dto

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,min=10,max=13"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
