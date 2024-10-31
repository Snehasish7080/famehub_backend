package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/snehasish7080/famehub/internal/middleware"
)

func AddUserRoutes(app *fiber.App, middleware *middleware.AuthMiddleware, controller *UserController) {
	auth := app.Group("/auth")

	auth.Post("/sign-up", controller.register)

	// verify otp token
	verifyEmail := auth.Group("/verify/otp", middleware.VerifyOtpToken)
	verifyEmail.Post("/", controller.verifyOtp)

}
