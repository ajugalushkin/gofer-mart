package withdrawal

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/ajugalushkin/gofer-mart/internal/database"
	"github.com/ajugalushkin/gofer-mart/internal/dto"
	"github.com/ajugalushkin/gofer-mart/internal/logger"
)

type Repository interface {
	AddNewWithdrawal(ctx context.Context, accrual dto.Withdrawal) error
	GetWithdrawalList(ctx context.Context, userID string) (*dto.WithdrawalList, error)
}

func NewRepository(db *sqlx.DB) Repository {
	return &repo{
		db: db,
	}
}

type repo struct {
	db *sqlx.DB
}

func (r *repo) AddNewWithdrawal(ctx context.Context, withdrawal dto.Withdrawal) error {
	var err error
	err = database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.StatementBuilder.
			Insert("withdrawals").
			Columns("number", "sum", "processed_at, user_id").
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		sb = sb.Values(
			withdrawal.Number,
			withdrawal.Sum,
			withdrawal.ProcessedAt,
			withdrawal.UserID,
		)

		_, err = sb.ExecContext(ctx)
		return err
	})

	if err != nil {
		logger.LogFromContext(ctx).Debug("repository.AddNewWithdrawal", zap.Error(err))
		return errors.Wrap(err, "repository.AddNewWithdrawal")
	}
	return nil
}

func (r *repo) GetWithdrawalList(ctx context.Context, userID string) (*dto.WithdrawalList, error) {
	var (
		err            error
		withdrawalList dto.WithdrawalList
	)

	err = database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.StatementBuilder.
			Select("number", "sum", "processed_at").
			From("withdrawals").
			Where(squirrel.Eq{"user_id": userID}).
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		query, args, err := sb.ToSql()
		if err != nil {
			return err
		}

		return r.db.SelectContext(ctx, &withdrawalList, query, args...)
	})

	if err != nil {
		logger.LogFromContext(ctx).Debug("repository.GetWithdrawalList Error", zap.Error(err))
		return &withdrawalList, errors.Wrap(err, "repository.GetWithdrawalList")
	}
	return &withdrawalList, nil
}
