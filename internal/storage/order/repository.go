package order

import (
	"context"

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
	//GetUser(ctx context.Context, user dto.User) (*dto.User, error)
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
			Columns("order_number", "user_id").
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		sb = sb.Values(
			order.OrderNumber,
			order.UserID,
		)

		_, err = sb.ExecContext(ctx)
		return err
	})

	if err != nil {
		if pgErr, ok := errors.Unwrap(errors.Unwrap(err)).(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			err := r.getOrderByUser(ctx, order.OrderNumber, order.UserID)
			if err != nil {
				return errors.Wrapf(
					userrors.ErrorOrderAlreadyUploadedAnotherUser,
					"%s %s",
					userrors.ErrorOrderAlreadyUploadedAnotherUser,
					order.OrderNumber,
				)
			}
			return errors.Wrapf(
				userrors.ErrorOrderAlreadyUploadedThisUser,
				"%s %s",
				userrors.ErrorOrderAlreadyUploadedThisUser,
				order.OrderNumber,
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
			Select("order_id", "order_number", "user_id").
			From("orders").
			Where(squirrel.Eq{"order_number": orderNumber, "user_id": userID}).
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
