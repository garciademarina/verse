package user

import "context"

// User type details
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Repository define repository interface for a user
type Repository interface {
	ListAll(ctx context.Context) ([]*User, error)
	FindById(ctx context.Context, id string) (*User, error)
}
