package userrors

import "github.com/pkg/errors"

var (
	ErrorDuplicateLogin                  = errors.New("login already exists")
	ErrorIncorrectLoginPassword          = errors.New("incorrect login/password pair")
	ErrorOrderAlreadyUploadedThisUser    = errors.New("order number has already been uploaded by this user")
	ErrorOrderAlreadyUploadedAnotherUser = errors.New("order number has already been uploaded by another user")
	ErrorIncorrectOrderNumber            = errors.New("incorrect order number")
	ErrorInsufficientFunds               = errors.New("insufficient funds")
)
