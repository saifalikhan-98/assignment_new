// pkg/user/sqlite_user_manager.go

package user

import (
	"database/sql"
	"errors"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteUserManager struct {
	db  *sql.DB
	mu  sync.Mutex
	mu2 sync.Mutex // for exclusive access to SQLite operations
}

func NewSQLiteUserManager(dbPath string) (*SQLiteUserManager, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create the users table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			email TEXT
		)
	`)
	if err != nil {
		return nil, err
	}

	return &SQLiteUserManager{
		db: db,
	}, nil
}

func (m *SQLiteUserManager) Create(user User) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.mu2.Lock()
	defer m.mu2.Unlock()

	result, err := m.db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(userID), nil
}

func (m *SQLiteUserManager) Read(id int) (User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.mu2.Lock()
	defer m.mu2.Unlock()

	row := m.db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err == sql.ErrNoRows {
		return User{}, errors.New("user not found")
	} else if err != nil {
		return User{}, err
	}

	return user, nil
}

func (m *SQLiteUserManager) Update(id int, updatedUser User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.mu2.Lock()
	defer m.mu2.Unlock()

	_, err := m.db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", updatedUser.Name, updatedUser.Email, id)
	if err == sql.ErrNoRows {
		return errors.New("user not found")
	} else if err != nil {
		return err
	}

	return nil
}

func (m *SQLiteUserManager) Delete(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.mu2.Lock()
	defer m.mu2.Unlock()

	_, err := m.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return errors.New("user not found")
	} else if err != nil {
		return err
	}

	return nil
}
