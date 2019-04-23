package user

import "context"

// User type details
type User struct {
	ID    ID     `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ID represents a user identifier
type ID string

// Repository define repository interface for a user
type Repository interface {
	ListAll(ctx context.Context) ([]*User, error)
	FindById(ctx context.Context, id ID) (*User, error)
}
