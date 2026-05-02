// domain/auth.go  ← add this new file
package domain

import "go-ticket/pkg/jwt"

type AuthService interface {
	Register(name, email, password string) (*User, error)
	Login(email, password string) (string, *User, error)
	GetProfile(userID int64) (*User, error)
	ChangePassword(userID int64, oldPassword, newPassword string) error
	ValidateToken(token string) (*jwt.Claims, error)
}
