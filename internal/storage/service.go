package storage

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ajugalushkin/gofer-mart/config"
	"github.com/ajugalushkin/gofer-mart/internal/auth"
	"github.com/ajugalushkin/gofer-mart/internal/database"
	"github.com/ajugalushkin/gofer-mart/internal/dto"
	"github.com/ajugalushkin/gofer-mart/internal/storage/user"
	"github.com/ajugalushkin/gofer-mart/internal/userrors"
)

// var defaultStorage Storage
var defaultUserStorage user.Repository

func Init(ctx context.Context) {
	cfg := config.FlagsFromContext(ctx)
	if cfg.DataBaseURI != "" {
		db, _ := database.NewConnection("pgx", cfg.DataBaseURI)
		//if err != nil {
		//log.Error("storage.GetStorage Error:", zap.Error(err))
		//}
		defaultUserStorage = user.NewRepository(db)
	}
}

func AddNewUser(ctx context.Context, user dto.User) error {
	var err error
	user.Password, err = auth.HashPassword(user.Password)
	if err != nil {
		return err
	}
	return defaultUserStorage.AddNewUser(ctx, user)
}

func LoginUser(ctx context.Context, user dto.User) error {
	storageUser, err := defaultUserStorage.GetUser(ctx, user)
	if err != nil {
		return err
	}

	if !auth.CheckPasswordHash(user.Password, storageUser.Password) {
		return errors.Wrapf(userrors.ErrorIncorrectLoginPassword, "%s", userrors.ErrorIncorrectLoginPassword)
	}
	return nil
}
