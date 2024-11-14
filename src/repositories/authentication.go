package repositories

import (
	"blogger/src/models"
	"database/sql"
)

type authentication struct {
	db *sql.DB
}

// UsersRepository creates a users repository.
func AuthenticationRepository(db *sql.DB) *authentication {
	return &authentication{db: db}
}

// GetPassword retrieves the hashed password of the given user.
func (a authentication) GetPassword(userID uint64) (string, error) {
	record, err := a.db.Query("SELECT password FROM users WHERE id = ?", userID)
	if err != nil {
		return "", err
	}
	defer record.Close()

	var user models.User

	if record.Next() {
		if err = record.Scan(&user.Password); err != nil {
			return "", err

		}
	}

	return user.Password, nil
}

// ChangePassword changes the password of the given user.
func (a authentication) ChangePassword(userID uint64, newPassword string) error {
	record, err := a.db.Prepare("UPDATE users SET password = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer record.Close()

	if _, err = record.Exec(newPassword, userID); err != nil {
		return err
	}

	return nil
}
