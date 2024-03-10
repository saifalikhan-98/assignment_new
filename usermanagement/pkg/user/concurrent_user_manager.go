// pkg/user/concurrent_user_manager.go

package user

import (
	"errors"
	"sync"
)

type ConcurrentUserManager struct {
	mu InMemoryUserManager
}

func NewConcurrentUserManager() *ConcurrentUserManager {
	return &ConcurrentUserManager{
		mu: InMemoryUserManager{
			users: make(map[int]User),
		},
	}
}

func (m *ConcurrentUserManager) Create(user User) (int, error) {
	m.mu.mu.Lock()
	defer m.mu.mu.Unlock()

	// Use goroutine for concurrent create logic
	var wg sync.WaitGroup
	wg.Add(1)

	var userID int
	var err error

	go func() {
		defer wg.Done()
		userID, err = m.mu.Create(user)
	}()

	wg.Wait()

	return userID, err
}

func (m *ConcurrentUserManager) Read(id int) (User, error) {
	m.mu.mu.Lock()
	defer m.mu.mu.Unlock()

	// Use goroutine for concurrent read logic
	var wg sync.WaitGroup
	wg.Add(1)

	var result User
	var err error

	go func() {
		defer wg.Done()
		result, err = m.mu.Read(id)
	}()

	wg.Wait()

	return result, err
}

func (m *ConcurrentUserManager) Update(id int, updatedUser User) error {
	m.mu.mu.Lock()
	defer m.mu.mu.Unlock()

	// Use goroutine for concurrent update logic
	var wg sync.WaitGroup
	wg.Add(1)

	var err error

	go func() {
		defer wg.Done()
		err = m.mu.Update(id, updatedUser)
	}()

	wg.Wait()

	return err
}

func (m *ConcurrentUserManager) Delete(id int) error {
	m.mu.mu.Lock()
	defer m.mu.mu.Unlock()

	// Use goroutine for concurrent delete logic
	var wg sync.WaitGroup
	wg.Add(1)

	var err error

	go func() {
		defer wg.Done()
		err = m.mu.Delete(id)
	}()

	wg.Wait()

	return err
}
