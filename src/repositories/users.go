package repositories

import (
	"blogger/src/models"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type users struct {
	db *sql.DB
}

// UsersRepository creates an users repository.
func UsersRepository(db *sql.DB) *users {
	return &users{db: db}
}

// PageAndPerPageValues returns page and per_age values.
func PageAndPerPageValues(params url.Values) (page uint64, per_page uint64) {
	page, err := strconv.ParseUint(params.Get("page"), 10, 64)

	if err != nil {
		page = 1
	}

	per_page, err = strconv.ParseUint(params.Get("per_page"), 10, 64)

	if err != nil {
		per_page = 15
	}

	return uint64(page), uint64(per_page)
}

// Listing returns a listing of users.
func (u users) Listing(params url.Values) ([]models.User, error) {
	page, per_page := PageAndPerPageValues(params)

	sort := params.Get("sort")

	if sort == "" {
		sort = "id ASC"
	} else if strings.Contains(sort, "-") {
		sort = sort[1:] + " DESC"
	} else {
		sort += " ASC"
	}

	userNameOrNick := params.Get("user")
	userNameOrNick = fmt.Sprintf("%%%s%%", userNameOrNick)

	records, err := u.db.Query(
		"SELECT id, name, nickname, email, created_at FROM users WHERE name LIKE ? OR nickname LIKE ? ORDER BY ? LIMIT ? OFFSET ?",
		userNameOrNick, userNameOrNick, sort, per_page, (page-1)*per_page,
	)

	if err != nil {
		return nil, err
	}

	defer records.Close()

	var users []models.User

	for records.Next() {
		var user models.User

		if err = records.Scan(
			&user.ID,
			&user.Name,
			&user.Nickname,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// Create inserts a brand new user in the database.
func (u users) Create(user models.User) (uint, error) {
	statement, err := u.db.Prepare("INSERT INTO users (name, nickname, email, password) VALUES (?, ?, ?, ?)")

	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(
		user.Name,
		user.Nickname,
		user.Email,
		user.Password,
	)

	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(lastID), nil
}

// Get finds and returns the user with the given ID, otherwise returns an error.
func (u users) Get(ID uint64) (models.User, error) {
	record, err := u.db.Query(
		"SELECT id, name, nickname, email, created_at FROM users WHERE id = ?",
		ID,
	)

	if err != nil {
		return models.User{}, nil
	}
	defer record.Close()

	var user models.User

	if record.Next() {
		if err = record.Scan(
			&user.ID,
			&user.Name,
			&user.Nickname,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return models.User{}, nil
		}
	}

	return user, nil
}

// Update finds and updates the user with the given ID, otherwise returns an error.
func (u users) Update(ID uint64, user models.User) error {
	aux, err := u.Get(ID)

	if err != nil {
		return err
	}

	if user.Name == "" {
		user.Name = aux.Name
	}
	if user.Nickname == "" {
		user.Nickname = aux.Nickname
	}
	if user.Email == "" {
		user.Email = aux.Email
	}

	statement, err := u.db.Prepare("UPDATE users SET name = ?, nickname = ?, email = ? WHERE id = ?")

	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(user.Name, user.Nickname, user.Email, ID); err != nil {
		return err
	}

	return nil
}

// Delete finds and deletes the user with the given ID, otherwise returns an error.
func (u users) Delete(ID uint64) error {
	statement, err := u.db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(ID); err != nil {
		return err
	}

	return nil
}

// Get finds and returns the user with the given ID, otherwise returns an error.
func (u users) GetByEmail(ID string) (models.User, error) {
	record, err := u.db.Query(
		"SELECT id, email, password FROM users WHERE email = ?",
		ID,
	)

	if err != nil {
		return models.User{}, nil
	}
	defer record.Close()

	var user models.User

	if record.Next() {
		if err = record.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
		); err != nil {
			return models.User{}, nil
		}
	}

	return user, nil
}

// Follow finds and links the two users.
func (u users) Follow(followID, userID uint64) error {
	statement, err := u.db.Prepare("INSERT IGNORE INTO followers(user_id, follower_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(userID, followID); err != nil {
		return err
	}

	return nil
}

// Unfollow finds and links the two users.
func (u users) Unfollow(followID, userID uint64) error {
	statement, err := u.db.Prepare("DELETE FROM followers WHERE user_id =? AND follower_id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(userID, followID); err != nil {
		return err
	}

	return nil
}

// GetFollowers retrieves all the followers of the given user.
func (u users) GetFollowers(userID uint64) ([]models.User, error) {
	records, err := u.db.Query(`
		SELECT u.id, u.name, u.nickname, u.email, u.created_at
		FROM users u 
		INNER JOIN followers f 
		ON u.id = f.follower_id 
		WHERE f.user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer records.Close()

	var followers []models.User

	for records.Next() {
		var user models.User

		if err = records.Scan(
			&user.ID,
			&user.Name,
			&user.Nickname,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		followers = append(followers, user)
	}

	return followers, nil
}

// GetFollowing retrieves all the users which the given user is following.
func (u users) GetFollowing(userID uint64) ([]models.User, error) {
	records, err := u.db.Query(`
		SELECT u.id, u.name, u.nickname, u.email, u.created_at
		FROM users u
		INNER JOIN followers f
		ON u.id = f.user_id
		WHERE f.follower_id = ?
	`, userID)
	if err != nil {
		return nil, nil
	}
	defer records.Close()

	var following []models.User

	for records.Next() {
		var user models.User

		if err = records.Scan(
			&user.ID,
			&user.Name,
			&user.Nickname,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		following = append(following, user)
	}

	return following, nil
}
