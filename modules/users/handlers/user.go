package userHandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/shop/config"
	"github.com/sing3demons/shop/modules/entities"
	"github.com/sing3demons/shop/modules/users"
	userUsecases "github.com/sing3demons/shop/modules/users/usecases"
)

type IUserHandlers interface {
	SignUpCustomer(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
}

type userHandlers struct {
	cfg     config.IConfig
	usecase userUsecases.IUserUsecases
}

func NewUserHandlers(cfg config.IConfig, usecase userUsecases.IUserUsecases) IUserHandlers {
	return &userHandlers{cfg: cfg, usecase: usecase}
}

type userHandlersErrCode string

const (
	signUpCustomerErr     userHandlersErrCode = "users-001"
	signInErr             userHandlersErrCode = "users-002"
	refreshPassportErr    userHandlersErrCode = "users-003"
	signOutErr            userHandlersErrCode = "users-004"
	signUpAdminErr        userHandlersErrCode = "users-005"
	generateAdminTokenErr userHandlersErrCode = "users-006"
	getUserProfileErr     userHandlersErrCode = "users-007"
)

func (h *userHandlers) SignUpCustomer(c *fiber.Ctx) error {
	// Request body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, "users-001", err.Error()).Response()
	}

	if !req.IsEmail() {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, "users-001", "email pattern invalid").Response()
	}

	// Insert user
	result, err := h.usecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Response()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Response()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Response()
		}
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Response()
}

func (h *userHandlers) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(signInErr), err.Error()).Response()
	}
	user, err := h.usecase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(signInErr), err.Error()).Response()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, user).Response()
}
