// pkg/user/in_memory_user_manager.go

package user

import (
	"errors"
	"sync"
)

type InMemoryUserManager struct {
	mu    sync.Mutex
	users map[int]User
}

func NewInMemoryUserManager() *InMemoryUserManager {
	return &InMemoryUserManager{
		users: make(map[int]User),
	}
}

func (m *InMemoryUserManager) Create(user User) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate a unique ID (for simplicity, you can use a counter)
	user.ID = len(m.users) + 1

	// Store the user
	m.users[user.ID] = user

	return user.ID, nil
}

func (m *InMemoryUserManager) Read(id int) (User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the user with the given ID exists
	if user, ok := m.users[id]; ok {
		return user, nil
	}

	return User{}, errors.New("user not found")
}

func (m *InMemoryUserManager) Update(id int, updatedUser User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the user with the given ID exists
	if _, ok := m.users[id]; !ok {
		return errors.New("user not found")
	}

	// Update the user
	updatedUser.ID = id
	m.users[id] = updatedUser

	return nil
}

func (m *InMemoryUserManager) Delete(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the user with the given ID exists
	if _, ok := m.users[id]; !ok {
		return errors.New("user not found")
	}

	// Delete the user
	delete(m.users, id)

	return nil
}
