package user

import (
	"context"
	"fmt"
	"sync"

	"github.com/garciademarina/verse/pkg/user"
)

// InmemRepository inmemory implementation of user.Repository
type InmemRepository struct {
	mtx   sync.RWMutex
	users map[user.ID]*user.User
}

// NewInmemoryRepository returns implement of user repository interface
func NewInmemoryRepository(users map[user.ID]*user.User) *InmemRepository {
	if users == nil {
		users = make(map[user.ID]*user.User)
	}

	return &InmemRepository{
		users: users,
	}
}

func (m *InmemRepository) ListAll(ctx context.Context) ([]*user.User, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	values := make([]*user.User, 0, len(m.users))
	for _, value := range m.users {
		values = append(values, value)
	}
	return values, nil
}

func (m *InmemRepository) FindById(ctx context.Context, id user.ID) (*user.User, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for _, v := range m.users {
		if v.ID == id {
			return v, nil
		}
	}

	return nil, fmt.Errorf("The ID %s doesn't exist", id)
}
