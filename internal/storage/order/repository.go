package order

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/ajugalushkin/gofer-mart/internal/database"
	"github.com/ajugalushkin/gofer-mart/internal/dto"
	"github.com/ajugalushkin/gofer-mart/internal/logger"
	"github.com/ajugalushkin/gofer-mart/internal/userrors"
)

const orderTableName string = "orders"

type Repository interface {
	AddNewOrder(ctx context.Context, user dto.Order) error
	GetOrderList(ctx context.Context, userID string) (*dto.OrderList, error)
	UpdateOrder(ctx context.Context, user dto.Order) error
	CheckOrderExists(ctx context.Context, orderNumber string, userID string) error
}

func NewRepository(db *sqlx.DB) Repository {
	return &repo{
		db: db,
	}
}

type repo struct {
	db *sqlx.DB
}

func (r *repo) AddNewOrder(ctx context.Context, order dto.Order) error {
	var err error
	err = database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.StatementBuilder.
			Insert(orderTableName).
			Columns("number", "uploaded_at", "status", "user_id").
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		sb = sb.Values(
			order.Number,
			order.UploadedAt,
			order.Status,
			order.UserID,
		)

		_, err = sb.ExecContext(ctx)
		return err
	})

	if err != nil {
		if pgErr, ok := errors.Unwrap(errors.Unwrap(err)).(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			err := r.CheckOrderExists(ctx, order.Number, order.UserID)
			if err != nil {
				logger.LogFromContext(ctx).Debug("repository.AddNewOrder: order already uploaded another user")
				return errors.Wrapf(
					userrors.ErrorOrderAlreadyUploadedAnotherUser,
					"%s %s",
					userrors.ErrorOrderAlreadyUploadedAnotherUser,
					order.Number,
				)
			}
			logger.LogFromContext(ctx).Debug("repository.AddNewOrder: order already uploaded this user")
			return errors.Wrapf(
				userrors.ErrorOrderAlreadyUploadedThisUser,
				"%s %s",
				userrors.ErrorOrderAlreadyUploadedThisUser,
				order.Number,
			)
		}
		logger.LogFromContext(ctx).Debug("repository.AddNewOrder", zap.Error(err))
		return errors.Wrap(err, "repository.AddNewOrder")
	}
	return nil
}

func (r *repo) GetOrderList(ctx context.Context, user string) (*dto.OrderList, error) {
	var (
		err       error
		orderList dto.OrderList
	)

	err = database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.StatementBuilder.
			Select("number", "uploaded_at", "status", "user_id").
			From(orderTableName).
			OrderBy("uploaded_at ASC").
			Where(squirrel.Eq{"user_id": user}).
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		query, args, err := sb.ToSql()
		if err != nil {
			return err
		}

		return r.db.SelectContext(ctx, &orderList, query, args...)
	})

	if err != nil {
		return &orderList, errors.Wrap(err, "repository.GetOrderList")
	}
	return &orderList, nil
}

func (r *repo) UpdateOrder(ctx context.Context, user dto.Order) error {
	err := database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.StatementBuilder.
			Update(orderTableName).
			Set("uploaded_at", user.UploadedAt).
			Set("status", user.Status).
			Set("accrual", user.Accrual).
			Where(squirrel.Eq{"number": user.Number}).
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		_, err := sb.ExecContext(ctx)
		return err
	})
	if err != nil {
		logger.LogFromContext(ctx).Debug("repository.UpdateOrder Error", zap.Error(err))
		return errors.Wrap(err, "repository.UpdateOrder")
	}
	return nil
}

func (r *repo) CheckOrderExists(ctx context.Context, orderNumber string, login string) error {
	var (
		err       error
		orderList []dto.Order
	)

	err = database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.StatementBuilder.
			Select("number", "user_id").
			From("orders").
			Where(squirrel.Eq{"number": orderNumber, "user_id": login}).
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		query, args, err := sb.ToSql()
		if err != nil {
			return err
		}

		return r.db.SelectContext(ctx, &orderList, query, args...)
	})

	if err != nil {
		logger.LogFromContext(ctx).Debug("repository.CheckOrderExists Error", zap.Error(err))
		return errors.Wrap(err, "repository.CheckOrderExists")
	}
	if orderList == nil {
		return errors.Wrap(err, "repository.CheckOrderExists")
	}
	return nil
}
