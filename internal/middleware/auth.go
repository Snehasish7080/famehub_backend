package middleware

type AuthMiddleware struct {
	storage *MiddlewareStorage
}

func NewAuthMiddleware(storage *MiddlewareStorage) *AuthMiddleware {
	return &AuthMiddleware{
		storage: storage,
	}
}
