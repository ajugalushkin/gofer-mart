package storage

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/ajugalushkin/gofer-mart/config"
	"github.com/ajugalushkin/gofer-mart/internal/database"
	"github.com/ajugalushkin/gofer-mart/internal/dto"
	"github.com/ajugalushkin/gofer-mart/internal/storage/user"
	"github.com/ajugalushkin/gofer-mart/internal/user_errors"
)

//type Storage interface {
//	Save(full string) (string, error)
//	//Get(short string) (string, error)
//	//Ping() error
//}

// var defaultStorage Storage
var defaultUserStorage user.Repository

func Init(ctx context.Context) {
	cfg := config.FlagsFromContext(ctx)
	if cfg.DataBaseURI != "" {
		db, err := database.NewConnection("pgx", cfg.DataBaseURI)
		if err != nil {
			//log.Error("storage.GetStorage Error:", zap.Error(err))
			//return nil
		}
		defaultUserStorage = user.NewRepository(db)
	}
}

func AddNewUser(ctx context.Context, user dto.User) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return defaultUserStorage.AddNewUser(ctx, user)
}

func LoginUser(ctx context.Context, user dto.User) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)

	storageUser, err := defaultUserStorage.GetUser(ctx, user)
	if err != nil {
		return err
	}

	if storageUser.Password != user.Password {
		return errors.Wrapf(user_errors.ErrorLoginAlreadyTaken, "%s %s", user_errors.ErrorLoginAlreadyTaken, user.Login)
	}
	return nil
}
