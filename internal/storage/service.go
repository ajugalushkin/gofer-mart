package storage

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/ajugalushkin/gofer-mart/config"
	"github.com/ajugalushkin/gofer-mart/internal/auth"
	"github.com/ajugalushkin/gofer-mart/internal/database"
	"github.com/ajugalushkin/gofer-mart/internal/dto"
	"github.com/ajugalushkin/gofer-mart/internal/logger"
	"github.com/ajugalushkin/gofer-mart/internal/queue"
	"github.com/ajugalushkin/gofer-mart/internal/storage/order"
	"github.com/ajugalushkin/gofer-mart/internal/storage/user"
	"github.com/ajugalushkin/gofer-mart/internal/storage/withdrawal"
	"github.com/ajugalushkin/gofer-mart/internal/userrors"
)

var defaultUserStorage user.Repository
var defaultOrderStorage order.Repository
var defaultWithdrawalStorage withdrawal.Repository

func Init(ctx context.Context) {
	cfg := config.FlagsFromContext(ctx)
	if cfg.DataBaseURI != "" {
		db, err := database.NewConnection("pgx", cfg.DataBaseURI)
		if err != nil {
			logger.LogFromContext(ctx).Debug("storage.Init", zap.Error(err))
		}
		defaultUserStorage = user.NewRepository(db)
		defaultOrderStorage = order.NewRepository(db)
		defaultWithdrawalStorage = withdrawal.NewRepository(db)
	}
}

func AddNewUser(ctx context.Context, user dto.User) error {
	var err error
	user.Password, err = auth.HashPassword(user.Password)
	if err != nil {
		logger.LogFromContext(ctx).Debug("service.AddNewUser: Failed to hash password")
		return err
	}

	return defaultUserStorage.AddNewUser(ctx, user)
}

func LoginUser(ctx context.Context, user dto.User) error {
	storageUser, err := defaultUserStorage.GetUser(ctx, user.Login)
	if err != nil {
		logger.LogFromContext(ctx).Debug("service.LoginUser: Failed to get user")
		return err
	}

	if !auth.CheckPasswordHash(user.Password, storageUser.Password) {
		logger.LogFromContext(ctx).Debug("service.LoginUser: Invalid password")
		return errors.Wrapf(userrors.ErrorIncorrectLoginPassword, "%s", userrors.ErrorIncorrectLoginPassword)
	}
	return nil
}

func AddNewOrder(ctx context.Context, orderNumber string, login string) error {
	order := dto.Order{Number: orderNumber, UploadedAt: time.Now(), UserID: login}
	err := defaultOrderStorage.AddNewOrder(ctx, order)
	if err != nil {
		logger.LogFromContext(ctx).Debug("service.AddNewOrder: Failed to add new order")
		return err
	}

	queue.AddOrder(&order)
	return nil
}

func GetOrders(ctx context.Context, login string) (*dto.OrderList, error) {
	user, err := defaultUserStorage.GetUser(ctx, login)
	if err != nil {
		return &dto.OrderList{}, err
	}
	return defaultOrderStorage.GetOrderList(ctx, user.Login)
}

func UpdateOrder(ctx context.Context, order dto.Order) error {
	return defaultOrderStorage.UpdateOrder(ctx, order)
}

//func AddAccrual(ctx context.Context, accrual float64, order string) error {
//	return defaultAccrualStorage.AddNewAccrual(ctx,
//		dto.Accrual{Accrual: accrual, ProcessedAt: time.Now(), OrderID: order})
//}

func GetBalance(ctx context.Context, login string) (*dto.Balance, error) {
	var balance dto.Balance

	orderList, err := defaultOrderStorage.GetOrderList(ctx, login)
	if err != nil {
		return &balance, err
	}

	var orders = make([]string, 0)
	for _, order := range *orderList {
		balance.Current = balance.Current + order.Accrual
		orders = append(orders, order.Number)
	}

	withdrawalList, err := defaultWithdrawalStorage.GetWithdrawalList(ctx, orders)
	if err != nil {
		return &balance, err
	}

	for _, withdrawal := range *withdrawalList {
		balance.Withdrawn = balance.Withdrawn + withdrawal.Sum
	}

	return &balance, nil
}

func AddNewWithdrawal(ctx context.Context, withdrawal dto.Withdraw, login string) error {
	balance, err := GetBalance(ctx, login)
	if err != nil {
		return err
	}
	if balance.Current < withdrawal.Sum {
		return errors.Wrapf(userrors.ErrorInsufficientFunds,
			"%s", userrors.ErrorInsufficientFunds)
	}

	err = defaultOrderStorage.CheckOrderExists(ctx, withdrawal.Order, login)
	if err != nil {
		return errors.Wrapf(userrors.ErrorIncorrectOrderNumber,
			"%s", userrors.ErrorIncorrectOrderNumber)
	}

	return defaultWithdrawalStorage.AddNewWithdrawal(ctx,
		dto.Withdrawal{
			Number:      withdrawal.Order,
			Sum:         withdrawal.Sum,
			ProcessedAt: time.Now()})
}

func GetWithdrawalList(ctx context.Context, login string) (*dto.WithdrawalList, error) {
	orderList, err := defaultOrderStorage.GetOrderList(ctx, login)
	if err != nil {
		return &dto.WithdrawalList{}, err
	}

	var orders = make([]string, 0)
	for _, order := range *orderList {
		orders = append(orders, order.Number)
	}

	withdrawalList, err := defaultWithdrawalStorage.GetWithdrawalList(ctx, orders)
	if err != nil {
		return withdrawalList, err
	}

	return withdrawalList, nil
}
