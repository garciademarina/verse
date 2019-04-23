package account

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrOriginAccountNotFound ..
	ErrOriginAccountNotFound = errors.New("Origin Account not found")
	// ErrDestinationAccountNotFound ..
	ErrDestinationAccountNotFound = errors.New("Destination account not found")
	// ErrBalanceInsufficient ..
	ErrBalanceInsufficient = errors.New("Insufficient balance")
)

// Account type details
type Account struct {
	Num     string    `json:"num"`
	UserID  string    `json:"userID"`
	Name    string    `json:"name"`
	OpenAt  time.Time `json:"openAt"`
	Balance int64     `json:"balance"`
}

// Repository define repository interface for an account
type Repository interface {
	ListAll(ctx context.Context) ([]*Account, error)
	FindByID(ctx context.Context, id string) (*Account, error)
	FindByUserID(ctx context.Context, userID string) (*Account, error)
	GetBalance(ctx context.Context, num string) (int64, error)
	UpdateBalance(ctx context.Context, num string, amount int64) (*Account, error)
	TransferMoney(ctx context.Context, origin, destination string, amount int64) error
}
