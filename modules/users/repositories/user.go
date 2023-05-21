package userRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sing3demons/shop/modules/users"
	userPatterns "github.com/sing3demons/shop/modules/users/repositories/patterns"
)

type IUserRepository interface {
	InsertUser(*users.UserRegisterReq, bool) (*users.UserPassport, error)
	FindOneUserByEmail(email string) (*users.UserCredentialCheck, error)
	InsertOauth(req *users.UserPassport) error 
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error) {
	result := userPatterns.InsertUser(r.db, req, isAdmin)

	var err error
	if isAdmin {
		result, err = result.Admin()
		if err != nil {
			return nil, err
		}
	} else {
		result, err = result.Customer()
		if err != nil {
			return nil, err
		}
	}

	// Get result from pattern
	user, err := result.Result()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
	query := `SELECT "id", "email", "password", "username", "role_id" FROM users WHERE email = $1`
	user := new(users.UserCredentialCheck)
	if err := r.db.Get(user, query, email); err != nil {
		return nil, fmt.Errorf("error get user: %w", err)
	}
	return user, nil
}

func (r *userRepository) InsertOauth(req *users.UserPassport) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	query := `INSERT INTO "oauth" ("user_id", "refresh_token", "access_token") VALUES ($1, $2, $3) RETURNING "id"`

	if err := r.db.QueryRowContext(ctx, query, req.User.Id, req.Token.RefreshToken, req.Token.AccessToken).Scan(&req.Token.Id); err != nil {
		return fmt.Errorf("error insert oauth: %w", err)
	}
	return nil
}
