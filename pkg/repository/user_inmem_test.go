package repository

import (
	"context"
	"fmt"
	"testing"

	sample "github.com/garciademarina/verse/cmd/sample-data"
	"github.com/stretchr/testify/assert"
)

func TestUserRepositoryListAllEmptyRepository(t *testing.T) {
	userRepository := NewInmemUserRepo(nil)
	got, _ := userRepository.ListAll(context.Background())
	assert.Equal(t, 0, len(got), "User repository not empty")
}

func TestUserRepositoryListAll(t *testing.T) {
	users := sample.Users

	userRepository := NewInmemUserRepo(users)
	got, _ := userRepository.ListAll(context.Background())
	assert.Equal(t, len(users), len(got), "Not the same size")

	for _, value := range got {
		i := value.ID
		assert.Equal(t, fmt.Sprintf("%+v", users[i]), fmt.Sprintf("%+v", value), "Not the same")
	}
}

func TestUserRepositoryFindById(t *testing.T) {
	id := "01D3XZ3ZHCP3KG9VT4FGAD8KDR"
	expected := sample.Users[id]

	userRepository := NewInmemUserRepo(sample.Users)

	got, _ := userRepository.FindById(context.Background(), id)

	assert.Equal(t, expected, got, "Not the same")
}

func TestUserRepositoryFindByIdNotFound(t *testing.T) {
	userRepository := NewInmemUserRepo(sample.Users)
	_, err := userRepository.FindById(context.Background(), "user-id-does-not-exist")

	assert.NotNil(t, err, "Not the same")
}
