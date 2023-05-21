package userRepository

import (
	"github.com/jmoiron/sqlx"
	"github.com/sing3demons/shop/modules/users"
	userPatterns "github.com/sing3demons/shop/modules/users/repositories/patterns"
)

type IUserRepository interface {
	InsertUser(*users.UserRegisterReq, bool) (*users.UserPassport, error)
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
