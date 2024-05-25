package storage

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ajugalushkin/gofer-mart/config"
	"github.com/ajugalushkin/gofer-mart/internal/auth"
	"github.com/ajugalushkin/gofer-mart/internal/database"
	"github.com/ajugalushkin/gofer-mart/internal/dto"
	"github.com/ajugalushkin/gofer-mart/internal/storage/order"
	"github.com/ajugalushkin/gofer-mart/internal/storage/user"
	"github.com/ajugalushkin/gofer-mart/internal/userrors"
)

var defaultUserStorage user.Repository
var defaultOrderStorage order.Repository

func Init(ctx context.Context) {
	cfg := config.FlagsFromContext(ctx)
	if cfg.DataBaseURI != "" {
		db, _ := database.NewConnection("pgx", cfg.DataBaseURI)
		//if err != nil {
		//log.Error("storage.GetStorage Error:", zap.Error(err))
		//}
		defaultUserStorage = user.NewRepository(db)
		defaultOrderStorage = order.NewRepository(db)
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

func AddNewOrder(ctx context.Context, order string, login string) error {
	user, err := defaultUserStorage.GetUser(ctx, dto.User{Login: login})
	if err != nil {
		return err
	}

	return defaultOrderStorage.AddNewOrder(ctx, dto.Order{Number: order,
		UserID: user.ID})
}

func GetOrders(ctx context.Context, login string) (*dto.OrderList, error) {
	user, err := defaultUserStorage.GetUser(ctx, dto.User{Login: login})
	if err != nil {
		return &dto.OrderList{}, err
	}
	return defaultOrderStorage.GetOrdersByUser(ctx, user.ID)
}
