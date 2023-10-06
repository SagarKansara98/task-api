package user

import (
	"context"
	"database/sql"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

// User holds resources for User to execute business operations
type User struct {
	db           *sqlx.DB
	log          zerolog.Logger
	tokenSeceret string
}

// New initializes User
func New(db *sqlx.DB, tokenSeceret string, log zerolog.Logger) User {
	return User{
		db:           db,
		tokenSeceret: tokenSeceret,
		log:          log,
	}
}

// Create creates a new user in the database
func (u User) Create(ctx context.Context, newInfo NewInfo) (Info, error) {
	_, err := u.QueryByEmail(ctx, newInfo.Email)
	if err != ErrInvalidEmail && err != nil {
		return Info{}, errors.Wrap(err, "Create: querying user by email")
	}
	if err == nil {
		return Info{}, ErrEmailAlreadyExist
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		return Info{}, errors.Wrap(err, "Create: generating hash password")
	}
	newInfo.Password = string(hashedPassword)

	info := Info{
		Email: newInfo.Email,
		Name:  newInfo.Name,
	}
	q := `INSERT INTO users (email, name, password) VALUES (:email, :name, :password) RETURNING id`
	stmt, err := u.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return info, errors.Wrap(err, "CreateUser: preparing named statement")
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &info.ID, newInfo)
	if err != nil {
		return info, errors.Wrap(err, "CreateUser: creating user and querying last inserted id")
	}

	return info, nil
}

// QueryByID retrieves a user from the database by ID
func (u User) QueryByID(ctx context.Context, id int) (Info, error) {
	var user Info
	q := "SELECT id, email, name FROM users WHERE id = $1"
	err := u.db.GetContext(ctx, &user, q, id)
	return user, err
}

// Update updates user information in the database
func (u User) Update(ctx context.Context, userInfo Info) error {
	q := `UPDATE users
			SET email = :email,
				name = :name,
				updated_at = now()
			WHERE id = :id`
	_, err := u.db.NamedExecContext(ctx, q, userInfo)
	if err != nil {
		return errors.Wrap(err, "UpdateUser: executing update query")
	}
	return nil
}

// Delete deletes a user from the database by ID
func (u User) Delete(ctx context.Context, id int) error {
	q := `DELETE FROM users WHERE id = $1`
	_, err := u.db.ExecContext(ctx, q, id)
	if err != nil {
		return errors.Wrap(err, "DeleteUser: executing delete user query")
	}
	return nil
}

func (u User) QueryByEmail(ctx context.Context, email string) (Info, error) {
	q := "SELECT id, email, password FROM users WHERE email = $1 LIMIT 1"
	info := Info{}
	err := u.db.Get(&info, q, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return info, ErrInvalidEmail
		}
		return info, errors.Wrap(err, "Login: querying user")
	}

	return info, nil
}

func (u User) Login(ctx context.Context, newuserInfo Login) (string, error) {
	info, err := u.QueryByEmail(ctx, newuserInfo.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(info.Password), []byte(newuserInfo.Password))
	if err != nil {
		return "", ErrInvalidPassword
	}

	return u.generateToken(info.Email, info.ID)
}

func (u User) generateToken(email string, userID int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["sub"] = userID
	tokenString, err := token.SignedString([]byte(u.tokenSeceret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
