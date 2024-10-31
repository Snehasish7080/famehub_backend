package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/snehasish7080/famehub/pkg/jwtclaim"
)

type AuthMiddleware struct {
	storage *MiddlewareStorage
}

func NewAuthMiddleware(storage *MiddlewareStorage) *AuthMiddleware {
	return &AuthMiddleware{
		storage: storage,
	}
}

func (a *AuthMiddleware) VerifyOtpToken(c *fiber.Ctx) error {
	reqToken := c.Request().Header.Peek("Authorization")

	userId, valid := jwtclaim.ExtractId(string(reqToken))

	if !valid {
		return c.Status(fiber.StatusUnauthorized).SendString("unauthorized access")
	}
	c.Locals("userId", userId)
	return c.Next()
}
