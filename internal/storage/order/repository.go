package order

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/ajugalushkin/gofer-mart/internal/database"
	"github.com/ajugalushkin/gofer-mart/internal/dto"
	"github.com/ajugalushkin/gofer-mart/internal/userrors"
)

type Repository interface {
	AddNewOrder(ctx context.Context, user dto.Order) error
	GetOrdersByUser(ctx context.Context, userID string) (*dto.OrderList, error)
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
			Insert("orders").
			Columns("number", "uploaded_at", "status", "user_id").
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		sb = sb.Values(
			order.Number,
			time.Now(),
			"NEW",
			order.UserID,
		)

		_, err = sb.ExecContext(ctx)
		return err
	})

	if err != nil {
		if pgErr, ok := errors.Unwrap(errors.Unwrap(err)).(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			err := r.getOrderByUser(ctx, order.Number, order.UserID)
			if err != nil {
				return errors.Wrapf(
					userrors.ErrorOrderAlreadyUploadedAnotherUser,
					"%s %s",
					userrors.ErrorOrderAlreadyUploadedAnotherUser,
					order.Number,
				)
			}
			return errors.Wrapf(
				userrors.ErrorOrderAlreadyUploadedThisUser,
				"%s %s",
				userrors.ErrorOrderAlreadyUploadedThisUser,
				order.Number,
			)
		}
		return errors.Wrap(err, "repository.AddNewOrder")
	}
	return nil
}

func (r *repo) getOrderByUser(ctx context.Context, orderNumber string, userID string) error {
	var (
		err       error
		orderList []dto.Order
	)

	err = database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.StatementBuilder.
			Select("id", "number", "user_id").
			From("orders").
			Where(squirrel.Eq{"number": orderNumber, "user_id": userID}).
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		query, args, err := sb.ToSql()
		if err != nil {
			return err
		}

		return r.db.SelectContext(ctx, &orderList, query, args...)
	})

	if err != nil {
		return errors.Wrap(err, "repository.getOrderByUser")
	}
	return nil
}

func (r *repo) GetOrdersByUser(ctx context.Context, userID string) (*dto.OrderList, error) {
	var (
		err       error
		orderList dto.OrderList
	)

	err = database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.StatementBuilder.
			Select("id", "number", "uploaded_at", "status", "user_id").
			From("orders").
			OrderBy("uploaded_at DESC").
			Where(squirrel.Eq{"user_id": userID}).
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		query, args, err := sb.ToSql()
		if err != nil {
			return err
		}

		return r.db.SelectContext(ctx, &orderList, query, args...)
	})

	if err != nil {
		return &orderList, errors.Wrap(err, "repository.GetOrdersByUser")
	}
	return &orderList, nil
}
