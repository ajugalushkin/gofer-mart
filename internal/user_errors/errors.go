package user_errors

import "github.com/pkg/errors"

var (
	ErrorDuplicateLogin    = errors.New("login already exists")
	ErrorLoginAlreadyTaken = errors.New("login is already taken")
)
