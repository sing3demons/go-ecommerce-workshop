package userUsecases

import (
	"github.com/sing3demons/shop/config"
	"github.com/sing3demons/shop/modules/users"
	userRepository "github.com/sing3demons/shop/modules/users/repositories"
)

type IUserUsecases interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
}

type userUsecases struct {
	cfg  config.IConfig
	repo userRepository.IUserRepository
}

func NewUserUsecases(cfg config.IConfig, repo userRepository.IUserRepository) IUserUsecases {
	return &userUsecases{cfg: cfg, repo: repo}
}

func (u *userUsecases) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// Insert user
	result, err := u.repo.InsertUser(req, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}
