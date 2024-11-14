package models

import (
	"blogger/src/security"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

// User represents an user.
type User struct {
	ID        uint      `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Nickname  string    `json:"nickname,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

const (
	STEP_CREATION string = "creation"
	STEP_UPDATE   string = "update"
)

// Prepare validates and formats the request values.
func (u *User) Prepare(step string) error {
	if err := u.validate(step); err != nil {
		return err
	}

	if err := u.format(step); err != nil {
		return nil
	}

	return nil
}

func (u *User) validate(step string) error {
	if step == STEP_CREATION && u.Name == "" {
		return errors.New("the name is mandatory")
	}

	if step == STEP_CREATION && u.Nickname == "" {
		return errors.New("the nickname is mandatory")
	}

	if step == STEP_CREATION && u.Email == "" {
		return errors.New("the email is mandatory")
	}

	if err := checkmail.ValidateFormat(u.Email); err != nil {
		return err
	}

	if step == STEP_CREATION && u.Password == "" {
		return errors.New("the password is mandatory")
	}

	return nil
}

func (u *User) format(step string) error {
	u.Name = strings.TrimSpace(u.Name)
	u.Nickname = strings.TrimSpace(u.Nickname)
	u.Email = strings.TrimSpace(u.Email)

	if step == STEP_CREATION {
		hashString, err := security.Hash(u.Password)

		if err != nil {
			return nil
		}

		u.Password = string(hashString)
	}

	return nil
}
