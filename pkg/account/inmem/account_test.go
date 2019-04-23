package account

import (
	"context"
	"sync"
	"testing"

	sample "github.com/garciademarina/verse/cmd/sample-data"
	"github.com/garciademarina/verse/pkg/account"
	"github.com/stretchr/testify/assert"
)

func TestAccountRepositoryListAllEmptyRepository(t *testing.T) {
	accountRepository := NewInmemoryRepository(nil)
	got, _ := accountRepository.ListAll(context.Background())
	assert.Equal(t, 0, len(got), "Account repository not empty")
}

func TestAccountRepositoryListAll(t *testing.T) {
	accounts := sample.Accounts

	accountRepository := NewInmemoryRepository(accounts)
	got, _ := accountRepository.ListAll(context.Background())
	assert.Equal(t, len(accounts), len(got), "Not the same size")

	for _, account := range got {
		i := account.Num
		assert.Equal(t, accounts[i], account, "Not the same")
	}
}

func TestAccountRepositoryFindByUserID(t *testing.T) {
	accountID := account.Num("D8KDR")
	expected := sample.Accounts[accountID]

	accountRepository := NewInmemoryRepository(sample.Accounts)
	got, _ := accountRepository.FindByUserID(context.Background(), expected.UserID)

	assert.Equal(t, expected.Num, got.Num, "Not the same")
}

func TestAccountRepositoryFindByUserIDNotFound(t *testing.T) {
	accountRepository := NewInmemoryRepository(sample.Accounts)
	_, err := accountRepository.FindByUserID(context.Background(), "user-id-does-not-exist")

	assert.NotNil(t, err, "Not the same")
}

func TestAccountRepositoryUpdateBalance(t *testing.T) {
	amount := int64(400)

	accountID := account.Num("D8KDR")
	expected := sample.Accounts[accountID]
	expectedBalance := sample.Accounts[accountID].Balance - (10 * amount)

	accountRepository := NewInmemoryRepository(sample.Accounts)

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(amount int64) {
			defer wg.Done()
			_, _ = accountRepository.UpdateBalance(context.Background(), account.Num("D8KDR"), amount)
		}(amount)
	}
	wg.Wait()

	account, _ := accountRepository.FindByUserID(context.Background(), expected.UserID)

	assert.Equal(t, expectedBalance, account.Balance, "Not the same")
}

func TestAccountRepositoryUpdateBalanceNotFound(t *testing.T) {
	amount := int64(400)

	accountRepository := NewInmemoryRepository(sample.Accounts)
	_, err := accountRepository.UpdateBalance(context.Background(), "user-id-does-not-exist", amount)

	assert.NotNil(t, err, "Not the same")
}
