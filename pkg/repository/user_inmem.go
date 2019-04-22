package repository

import (
	"context"
	"fmt"
	"sync"

	models "github.com/garciademarina/verse/pkg/models"
)

// NewInmemUserRepo returns implement of user repository interface
func NewInmemUserRepo(users map[string]*models.User) UserRepo {
	if users == nil {
		users = make(map[string]*models.User)
	}

	return &inmemUserRepo{
		users: users,
	}
}

type inmemUserRepo struct {
	mtx   sync.RWMutex
	users map[string]*models.User
}

func (m *inmemUserRepo) ListAll(ctx context.Context) ([]*models.User, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	values := make([]*models.User, 0, len(m.users))
	for _, value := range m.users {
		values = append(values, value)
	}
	return values, nil
}

func (m *inmemUserRepo) FindById(ctx context.Context, id string) (*models.User, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for _, v := range m.users {
		if v.ID == id {
			return v, nil
		}
	}

	return nil, fmt.Errorf("The ID %s doesn't exist", id)
}
