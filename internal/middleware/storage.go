package middleware

import (
	"github.com/gocql/gocql"
)

type MiddlewareStorage struct {
	session *gocql.Session
}

func NewMiddlewareStorage(session *gocql.Session) *MiddlewareStorage {
	return &MiddlewareStorage{
		session: session,
	}
}
