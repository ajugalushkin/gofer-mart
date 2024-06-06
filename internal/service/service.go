package service

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/ajugalushkin/gofer-mart/config"
	"github.com/ajugalushkin/gofer-mart/internal/auth"
	"github.com/ajugalushkin/gofer-mart/internal/dto"
	"github.com/ajugalushkin/gofer-mart/internal/logger"
	"github.com/ajugalushkin/gofer-mart/internal/storage/order"
	"github.com/ajugalushkin/gofer-mart/internal/storage/user"
	"github.com/ajugalushkin/gofer-mart/internal/storage/withdrawal"
	"github.com/ajugalushkin/gofer-mart/internal/userrors"
	"github.com/ajugalushkin/gofer-mart/internal/workerpool"
)

type Service struct {
	userRepo       user.Repository
	orderRepo      order.Repository
	withdrawalRepo withdrawal.Repository
}

func NewService(db *sqlx.DB) *Service {
	return &Service{
		userRepo:       user.NewRepository(db),
		orderRepo:      order.NewRepository(db),
		withdrawalRepo: withdrawal.NewRepository(db),
	}
}

func (s *Service) AddNewUser(ctx context.Context, user dto.User) error {
	var err error
	user.Password, err = auth.HashPassword(user.Password)
	if err != nil {
		logger.LogFromContext(ctx).Debug("service.AddNewUser: Failed to hash password")
		return err
	}

	return s.userRepo.AddNewUser(ctx, user)
}

func (s *Service) LoginUser(ctx context.Context, user dto.User) error {
	storageUser, err := s.userRepo.GetUser(ctx, user.Login)
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

func (s *Service) AddNewOrder(ctx context.Context, orderNumber string, login string) error {
	newOrder := dto.Order{Number: orderNumber, Status: "NEW", UploadedAt: time.Now(), UserID: login}
	err := s.orderRepo.AddNewOrder(ctx, newOrder)
	if err != nil {
		logger.LogFromContext(ctx).Debug("service.AddNewOrder: Failed to add new newOrder")
		return err
	}

	workerPool := workerpool.WorkerPoolFromContext(ctx)

	parseURL, err := url.Parse(config.FlagsFromContext(ctx).AccrualSystemAddress)
	if err != nil {
		return err
	}
	//
	accrualURL := url.URL{
		Host:   parseURL.Host,
		Path:   "/api/orders/",
		Scheme: parseURL.Scheme}

	countError := 0
	go func() {
		url := accrualURL.String() + newOrder.Number
		for {
			workerPool.AddTask(url)
			result := workerPool.GetResult()

			if result.Err != nil {
				countError++
				logger.LogFromContext(ctx).Debug(
					"service.AddNewOrder: Error getting accrual",
					zap.String("url", url),
					zap.String("order", newOrder.Number),
					zap.Error(result.Err))
				if countError > 3 {
					//workerPool.Stop()
					return
				}
			}

			logger.LogFromContext(ctx).Debug(
				"service.AddNewOrder: ОК getting accrual",
				zap.String("order", result.Data.Order),
				zap.String("Status", result.Data.Status),
				zap.String("Accrual", fmt.Sprintf("%f", result.Data.Accrual)))

			err := s.orderRepo.UpdateOrder(ctx, dto.Order{
				Number:     result.Data.Order,
				Status:     result.Data.Status,
				Accrual:    &result.Data.Accrual,
				UploadedAt: time.Now(),
				UserID:     login,
			})

			if result.Data.Status != "INVALID" && result.Data.Status != "PROCESSED" {
				logger.LogFromContext(ctx).Debug("service.AddNewOrder add order to queue again")
				continue
			}

			if err != nil {
				logger.LogFromContext(ctx).Debug(
					"service.AddNewOrder: Update accrual Error", zap.Error(err))
			}

			return
		}
	}()

	return nil
}

func (s *Service) GetOrders(ctx context.Context, login string) (*dto.OrderList, error) {
	list, err := s.orderRepo.GetOrderList(ctx, login)
	if err != nil {
		logger.LogFromContext(ctx).Debug(
			"service.GetOrders: Error getting order list",
			zap.Error(err),
		)
		return nil, err
	}
	json, err := list.MarshalJSON()
	if err != nil {
		logger.LogFromContext(ctx).Debug(
			"service.GetOrders: Error getting order list",
			zap.Error(err),
		)
		return nil, err
	}
	logger.LogFromContext(ctx).Debug(
		"app.getOrders: Ok",
		zap.String("JSON", string(json)))

	return list, nil
}

func (s *Service) UpdateOrder(ctx context.Context, order dto.Order) error {
	return s.orderRepo.UpdateOrder(ctx, order)
}

func (s *Service) GetBalance(ctx context.Context, login string) (*dto.Balance, error) {
	var balance dto.Balance

	orderList, err := s.orderRepo.GetOrderList(ctx, login)
	if err != nil {
		logger.LogFromContext(ctx).Info("service.GetBalance: Failed to get order list")
		return &balance, err
	}

	for _, orderItem := range *orderList {
		balance.Current = balance.Current + *orderItem.Accrual
	}

	withdrawalList, err := s.withdrawalRepo.GetWithdrawalList(ctx, login)
	if err != nil {
		logger.LogFromContext(ctx).Info("service.GetBalance: Failed to get withdrawal list")
		return &balance, err
	}

	for _, withdrawalItem := range *withdrawalList {
		balance.Withdrawn = balance.Withdrawn + withdrawalItem.Sum
	}

	balance.Current = balance.Current - balance.Withdrawn

	return &balance, nil
}

func (s *Service) AddNewWithdrawal(ctx context.Context, withdrawal dto.Withdraw, userID string) error {
	balance, err := s.GetBalance(ctx, userID)
	if err != nil {
		logger.LogFromContext(ctx).Debug(
			"service.AddNewWithdrawal: Failed to get account balance",
			zap.Error(err),
		)
		return err
	}
	if balance.Current < withdrawal.Sum {
		logger.LogFromContext(ctx).Debug(
			"service.AddNewWithdrawal: Account not enough withdrawal",
			zap.Error(err))
		return errors.Wrapf(userrors.ErrorInsufficientFunds,
			"%s", userrors.ErrorInsufficientFunds)
	}

	err = s.withdrawalRepo.AddNewWithdrawal(ctx,
		dto.Withdrawal{
			Number:      withdrawal.Order,
			Sum:         withdrawal.Sum,
			ProcessedAt: time.Now(),
			UserID:      userID})
	if err != nil {
		logger.LogFromContext(ctx).Debug(
			"service.AddNewWithdrawal: Failed to add withdrawal",
			zap.Error(err))
		return err
	}
	return nil
}

func (s *Service) GetWithdrawalList(ctx context.Context, login string) (*dto.WithdrawalList, error) {
	withdrawalList, err := s.withdrawalRepo.GetWithdrawalList(ctx, login)
	if err != nil {
		logger.LogFromContext(ctx).Debug(
			"service.GetWithdrawalList: Failed to get withdrawal list",
			zap.Error(err))
		return withdrawalList, err
	}

	return withdrawalList, nil
}
