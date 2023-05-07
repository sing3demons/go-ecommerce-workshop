package entities

import (
	"github.com/gofiber/fiber/v2"
)

type IResponse interface {
	Success(code int, data interface{}) IResponse
	Error(code int, tractId string, msg string) IResponse
	Response() error
}

type response struct {
	Status        int
	Data          interface{}
	ErrorResponse *ErrorResponse
	Context       *fiber.Ctx
	IsError       bool
}

type ErrorResponse struct {
	TraceId string `json:"trace_id"`
	Message string `json:"message"`
}

func NewResponse(ctx *fiber.Ctx) IResponse {
	return &response{
		Context: ctx,
	}
}

func (res *response) Success(code int, data interface{}) IResponse {
	res.Status = code
	res.Data = data
	return res
}
func (res *response) Error(code int, tractId string, msg string) IResponse {
	res.Status = code
	res.IsError = true
	res.ErrorResponse = &ErrorResponse{
		TraceId: tractId,
		Message: msg,
	}
	return res
}
func (res *response) Response() error {
	if res.IsError {
		return res.Context.Status(res.Status).JSON(res.ErrorResponse)
	}
	return res.Context.Status(res.Status).JSON(res.Data)
}
