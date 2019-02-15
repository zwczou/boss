package admin

import (
	"errors"
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
)

var (
	ContextUserId = "admin.user_id"
	ContextUser   = "admin.user"
	ContextToken  = "admin.token"

	AdminToken   = "admin.token"
	TokenExpires = 86400 * 7
)
