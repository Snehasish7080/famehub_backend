package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/snehasish7080/famehub/internal/middleware"
)

func AddUserRoutes(app *fiber.App, middleware *middleware.AuthMiddleware, controller *UserController) {
	auth := app.Group("/auth")

	// user sign up
	auth.Post("/sign-up", controller.register)

	// user login
	auth.Post("/login", controller.loginUser)

	// verify otp token
	verifyOtp := auth.Group("/verify/otp", middleware.VerifyOtpToken)
	verifyOtp.Post("/", controller.verifyOtp)

}
