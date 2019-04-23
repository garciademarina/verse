package user

import (
	"context"
	"fmt"
	"sync"

	puser "github.com/garciademarina/verse/pkg/user"
)

type inmemRepository struct {
	mtx   sync.RWMutex
	users map[string]*puser.User
}

// NewInmemoryRepository returns implement of user repository interface
func NewInmemoryRepository(users map[string]*puser.User) puser.Repository {
	if users == nil {
		users = make(map[string]*puser.User)
	}

	return &inmemRepository{
		users: users,
	}
}

func (m *inmemRepository) ListAll(ctx context.Context) ([]*puser.User, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	values := make([]*puser.User, 0, len(m.users))
	for _, value := range m.users {
		values = append(values, value)
	}
	return values, nil
}

func (m *inmemRepository) FindById(ctx context.Context, id string) (*puser.User, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for _, v := range m.users {
		if v.ID == id {
			return v, nil
		}
	}

	return nil, fmt.Errorf("The ID %s doesn't exist", id)
}
