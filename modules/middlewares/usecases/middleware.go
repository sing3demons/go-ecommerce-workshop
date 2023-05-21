package usecases

import "github.com/sing3demons/shop/modules/middlewares/repositories"

type IMiddlewareUsecase interface{}

type middlewareUsecase struct{
repo repositories.IMiddlewareRepository
}

func NewMiddlewareUsecase(repo repositories.IMiddlewareRepository) IMiddlewareUsecase {
	return &middlewareUsecase{
	repo: repo,
	}
}
