package models

type Password struct {
	New     string `json:"new_password"`
	Current string `json:"current_password"`
}
