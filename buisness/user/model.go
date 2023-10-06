package user

import "errors"

var (
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrEmailAlreadyExist = errors.New("email already exist")
)

// Info holds information about a user
type Info struct {
	ID       int    `db:"id" json:"id"`
	Email    string `db:"email" json:"email" validate:"required,email"`
	Name     string `db:"name" json:"name" validate:"required"`
	Password string `db:"password" json:"-"`
}

// NewInfo hold new user info
type NewInfo struct {
	Email    string `db:"email" json:"email" validate:"required,email"`
	Name     string `db:"name" json:"name" validate:"required"`
	Password string `db:"password" json:"password"`
}

// Login hold new user info
type Login struct {
	Email    string `db:"email" json:"email" validate:"required,email"`
	Password string `db:"password" json:"password"`
}
