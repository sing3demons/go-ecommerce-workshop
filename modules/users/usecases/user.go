package userUsecases

import (
	"github.com/sing3demons/shop/config"
	"github.com/sing3demons/shop/modules/users"
	userRepository "github.com/sing3demons/shop/modules/users/repositories"
	"github.com/sing3demons/shop/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecases interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
	RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
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

func (u *userUsecases) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	// Find user
	user, err := u.repo.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	accessToken, err := auth.NewAuth(auth.Access, u.cfg.JWT(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.NewAuth(auth.Refresh, u.cfg.JWT(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	if err != nil {
		return nil, err
	}

	passport := &users.UserPassport{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	if err := u.repo.InsertOauth(passport); err != nil {
		return nil, err
	}
	return passport, nil
}

func (u *userUsecases) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
	// Parse token
	claims, err := auth.ParseToken(u.cfg.JWT(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Check oauth
	oauth, err := u.repo.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// find profile
	profile, err := u.repo.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		Id:     profile.Id,
		RoleId: profile.RoleId,
	}

	accessToken, err := auth.NewAuth(auth.Access, u.cfg.JWT(), newClaims)
	if err != nil {
		return nil, err
	}

	refreshToken := auth.RepeatToken(u.cfg.JWT(), newClaims, claims.ExpiresAt.Unix())

	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
			Id:           oauth.Id,
		},
	}

	if err := u.repo.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}

	return passport, nil
}
