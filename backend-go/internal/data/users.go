package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/saiharsha/money-manager/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"username"`
	Password  password  `json:"_"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	Activated bool      `json:"activated"`
	Version   int       `json:"version"`
}

type UserModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m UserModel) CreateUser(user *User) error {
	query := `
		INSERT INTO users (name, email, password_hash, activated)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{user.Name, user.Email, user.Password.hashed, user.Activated}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) GetUserByMail(email string) (*User, error) {
	query := `
		SELECT id, created_at, name, email, password_hash, activated, version, role
		FROM users
		WHERE email = $1
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hashed,
		&user.Activated,
		&user.Version,
		&user.Role,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m *UserModel) UpdateUser(user *User) (*User, error) {
	query := `
	UPDATE users 
	SET name = $1, email = $2, role = $3, activated = $4, password_hash = $5, version = version + 1
	WHERE id = $6 AND version = $7
	RETURNING version
	`

	args := []interface{}{
		user.Name,
		user.Email,
		user.Role,
		user.Activated,
		user.Password.hashed,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrEditConflict
		default:
			return nil, err
		}
	}

	return user, nil
}

type password struct {
	plaintext *string
	hashed    []byte
}

func (p *password) SetPasswordHash(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.hashed = hash
	p.plaintext = &plaintextPassword
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hashed, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func CheckPassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "should be not empty")
	v.Check(len(password) > 8, "password", "should be atleast 8 chars length")
	v.Check(len(password) < 63, "password", "should be atmost 63 chars length")
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be valid email address")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "username", "should be not empty")
	v.Check(len(user.Name) < 100, "username", "length of user name should be less than 100")
	ValidateEmail(v, user.Email)
	CheckPassword(v, *user.Password.plaintext)
}
