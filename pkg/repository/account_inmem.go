package repository

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	models "github.com/garciademarina/verse/pkg/models"
)

type inmemAccountRepo struct {
	mtx      sync.RWMutex
	accounts map[string]*models.Account
}

// NewInmemAccountRepo returns implement of account repository interface
func NewInmemAccountRepo(accounts map[string]*models.Account) AccountRepo {
	if accounts == nil {
		accounts = make(map[string]*models.Account)
	}

	return &inmemAccountRepo{
		accounts: accounts,
	}
}

func (m *inmemAccountRepo) ListAll(ctx context.Context) ([]*models.Account, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	values := make([]*models.Account, 0, len(m.accounts))
	for _, value := range m.accounts {
		values = append(values, value)
	}

	// sort map
	sortedKeys := make([]string, 0, len(m.accounts))
	for k := range m.accounts {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	accounts := make([]*models.Account, 0, len(m.accounts))
	for _, v := range sortedKeys {
		// fmt.Printf("k: %s, v: %v\n", k, myRecords[k])
		accounts = append(accounts, m.accounts[v])
	}

	return accounts, nil
}

func (m *inmemAccountRepo) FindByID(ctx context.Context, num string) (*models.Account, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	v, ok := m.accounts[num]
	if !ok {
		return v, nil
	}
	return nil, fmt.Errorf("Account number %s not found", num)
}

var (
	// ErrOriginAccountNotFound ..
	ErrOriginAccountNotFound = errors.New("Origin Account not found")
	// ErrDestinationAccountNotFound ..
	ErrDestinationAccountNotFound = errors.New("Destination account not found")
	// ErrBalanceInsufficient ..
	ErrBalanceInsufficient = errors.New("Insufficient balance")
)

func (m *inmemAccountRepo) TransferMoney(ctx context.Context, userOrigin, userDestination string, amount int64) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	var originAccount *models.Account
	var destinationAccount *models.Account
	for _, v := range m.accounts {
		if v.UserID == userOrigin {
			originAccount = v
		}
		if v.UserID == userDestination {
			destinationAccount = v
		}
	}
	if originAccount == nil {
		return ErrOriginAccountNotFound
	}
	if destinationAccount == nil {
		return ErrDestinationAccountNotFound
	}
	if originAccount.Balance < amount {
		return ErrBalanceInsufficient
	}

	originAccount.Balance = originAccount.Balance - amount
	destinationAccount.Balance = destinationAccount.Balance + amount

	return nil
}

func (m *inmemAccountRepo) FindByUserID(ctx context.Context, userID string) (*models.Account, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for _, v := range m.accounts {
		if v.UserID == userID {
			return v, nil
		}
	}

	return nil, fmt.Errorf("The User ID %s doesn't have any account", userID)
}

func (m *inmemAccountRepo) GetBalance(ctx context.Context, num string) (int64, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	v, ok := m.accounts[num]
	if !ok {
		return v.Balance, nil
	}

	return 0, fmt.Errorf("Cannot get balance, account number %s not found", num)
}

func (m *inmemAccountRepo) UpdateBalance(ctx context.Context, num string, amount int64) (*models.Account, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	var account *models.Account
	for _, v := range m.accounts {
		if v.Num == num && v.Balance >= amount {
			account = v
		}
	}
	if account != nil {
		account.Balance = account.Balance - amount
		return account, nil
	}
	return nil, fmt.Errorf("Cannot update balance")
}
