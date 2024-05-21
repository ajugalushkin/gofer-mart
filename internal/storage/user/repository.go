package user

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/ajugalushkin/gofer-mart/internal/database"
	"github.com/ajugalushkin/gofer-mart/internal/dto"
	"github.com/ajugalushkin/gofer-mart/internal/userrors"
)

type Repository interface {
	AddNewUser(ctx context.Context, user dto.User) error
	GetUser(ctx context.Context, user dto.User) (*dto.User, error)
}

func NewRepository(db *sqlx.DB) Repository {
	return &repo{
		db: db,
	}
}

type repo struct {
	db *sqlx.DB
}

func (r *repo) AddNewUser(ctx context.Context, user dto.User) error {
	var err error
	err = database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.StatementBuilder.
			Insert("users").
			Columns("login", "password_hash").
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		sb = sb.Values(
			user.Login,
			user.Password,
		)

		_, err = sb.ExecContext(ctx)
		return err
	})

	if err != nil {
		if pgErr, ok := errors.Unwrap(errors.Unwrap(err)).(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return errors.Wrapf(userrors.ErrorDuplicateLogin, "%s %s", userrors.ErrorDuplicateLogin, user.Login)
		}
		return errors.Wrap(err, "repository.AddNewUser")
	}
	return nil
}

func (r *repo) GetUser(ctx context.Context, user dto.User) (*dto.User, error) {
	var (
		err             error
		storageUserList []dto.User
		storageUser     dto.User
	)

	err = database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.StatementBuilder.
			Select("login", "password_hash").
			From("users").
			Where("login = ?", user.Login).
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		query, args, err := sb.ToSql()
		if err != nil {
			return err
		}

		return r.db.SelectContext(ctx, &storageUserList, query, args...)
	})

	if err != nil {
		return &storageUser, errors.Wrap(err, "repository.GetUser")
	}
	storageUser = storageUserList[0]
	return &storageUser, nil
}
