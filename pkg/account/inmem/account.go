package account

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/garciademarina/verse/pkg/account"
)

// InmemRepository inmemory implementation of account.Repository
type InmemRepository struct {
	mtx      sync.RWMutex
	accounts map[account.Num]*account.Account
}

// NewInmemoryRepository returns implement of user repository interface
func NewInmemoryRepository(accounts map[account.Num]*account.Account) *InmemRepository {
	if accounts == nil {
		accounts = make(map[account.Num]*account.Account)
	}

	return &InmemRepository{
		accounts: accounts,
	}
}

// ListAll returns all accounts
func (m *InmemRepository) ListAll(ctx context.Context) ([]*account.Account, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	// sort map
	sortedKeys := make([]account.Num, 0, len(m.accounts))
	for k := range m.accounts {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Slice(sortedKeys, func(i, j int) bool {
		return sortedKeys[i] < sortedKeys[j]
	})

	accounts := make([]*account.Account, 0, len(m.accounts))
	for _, v := range sortedKeys {
		accounts = append(accounts, m.accounts[v])
	}

	return accounts, nil
}

func (m *InmemRepository) FindByID(ctx context.Context, num account.Num) (*account.Account, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	v, ok := m.accounts[num]
	if !ok {
		return v, nil
	}
	return nil, fmt.Errorf("Account number %s not found", num)
}

func (m *InmemRepository) TransferMoney(ctx context.Context, userOrigin, userDestination string, amount int64) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	var originAccount *account.Account
	var destinationAccount *account.Account
	for _, v := range m.accounts {
		if v.UserID == userOrigin {
			originAccount = v
		}
		if v.UserID == userDestination {
			destinationAccount = v
		}
	}
	if originAccount == nil {
		return account.ErrOriginAccountNotFound
	}
	if destinationAccount == nil {
		return account.ErrDestinationAccountNotFound
	}
	if originAccount.Balance < amount {
		return account.ErrBalanceInsufficient
	}

	originAccount.Balance = originAccount.Balance - amount
	destinationAccount.Balance = destinationAccount.Balance + amount

	return nil
}

func (m *InmemRepository) FindByUserID(ctx context.Context, userID string) (*account.Account, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for _, v := range m.accounts {
		if v.UserID == userID {
			return v, nil
		}
	}

	return nil, fmt.Errorf("The User ID %s doesn't have any account", userID)
}

func (m *InmemRepository) GetBalance(ctx context.Context, num account.Num) (int64, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	v, ok := m.accounts[num]
	if !ok {
		return v.Balance, nil
	}

	return 0, fmt.Errorf("Cannot get balance, account number %s not found", num)
}

func (m *InmemRepository) UpdateBalance(ctx context.Context, num account.Num, amount int64) (*account.Account, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	var account *account.Account
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
