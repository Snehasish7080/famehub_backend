package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	storage *UserStorage
}

func NewUserController(storage *UserStorage) *UserController {
	return &UserController{
		storage: storage,
	}
}

var validate = validator.New()

type signUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type signUpResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (u *UserController) register(c *fiber.Ctx) error {
	var req signUpRequest

	c.BodyParser(&req)
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(signUpResponse{
			Message: "Invalid request body",
			Success: false,
		})
	}

	err := validate.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(signUpResponse{
			Message: "Invalid request body",
			Success: false,
		})
	}

	token, err := u.storage.signUp(req.Email, req.Password, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(signUpResponse{
			Message: err.Error(),
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(signUpResponse{
		Token:   token,
		Success: true,
		Message: "Otp sent successfully",
	})
}
