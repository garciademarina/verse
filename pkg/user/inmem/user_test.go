package user

import (
	"context"
	"testing"

	sample "github.com/garciademarina/verse/cmd/sample-data"
	"github.com/stretchr/testify/assert"
)

func TestUserRepositoryListAllEmptyRepository(t *testing.T) {
	userRepository := NewInmemoryRepository(nil)
	got, _ := userRepository.ListAll(context.Background())
	assert.Equal(t, 0, len(got), "User repository not empty")
}

func TestUserRepositoryListAll(t *testing.T) {
	users := sample.Users

	userRepository := NewInmemoryRepository(users)
	got, _ := userRepository.ListAll(context.Background())
	assert.Equal(t, len(users), len(got), "Not the same size")

	for _, user := range got {
		i := user.ID
		assert.Equal(t, users[i], user, "Not the same")
	}
}

func TestUserRepositoryFindById(t *testing.T) {
	id := "01D3XZ3ZHCP3KG9VT4FGAD8KDR"
	expected := sample.Users[id]

	userRepository := NewInmemoryRepository(sample.Users)

	got, _ := userRepository.FindById(context.Background(), id)

	assert.Equal(t, expected, got, "Not the same")
}

func TestUserRepositoryFindByIdNotFound(t *testing.T) {
	userRepository := NewInmemoryRepository(sample.Users)
	_, err := userRepository.FindById(context.Background(), "user-id-does-not-exist")

	assert.NotNil(t, err, "Not the same")
}
