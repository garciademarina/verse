package repository

import (
	"context"

	"github.com/garciademarina/verse/pkg/models"
)

// UserRepo explain...
type UserRepo interface {
	ListAll(ctx context.Context) ([]*models.User, error)
	FindById(ctx context.Context, id string) (*models.User, error)
}

// AccountRepo explain...
type AccountRepo interface {
	ListAll(ctx context.Context) ([]*models.Account, error)
	FindByID(ctx context.Context, id string) (*models.Account, error)
	FindByUserID(ctx context.Context, userID string) (*models.Account, error)
	GetBalance(ctx context.Context, num string) (int64, error)
	UpdateBalance(ctx context.Context, num string, amount int64) (*models.Account, error)
	TransferMoney(ctx context.Context, origin, destination string, amount int64) error
}
