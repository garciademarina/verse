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

// Num represents a account identifier
type Num string

// Account type details
type Account struct {
	Num     Num       `json:"num"`
	UserID  string    `json:"userID"`
	Name    string    `json:"name"`
	OpenAt  time.Time `json:"openAt"`
	Balance int64     `json:"balance"`
}

// Repository define repository interface for an account
type Repository interface {
	ListAll(ctx context.Context) ([]*Account, error)
	FindByID(ctx context.Context, id Num) (*Account, error)
	FindByUserID(ctx context.Context, userID string) (*Account, error)
	GetBalance(ctx context.Context, num Num) (int64, error)
	UpdateBalance(ctx context.Context, num Num, amount int64) (*Account, error)
	TransferMoney(ctx context.Context, origin, destination string, amount int64) error
}
