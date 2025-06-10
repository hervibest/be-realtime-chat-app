package repository

import (
	"be-realtime-chat-app/services/user-svc/internal/entity"
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"go.uber.org/zap"
)

type UserRepository interface {
	// Insert(ctx context.Context, db Querier, user *entity.User) (*entity.User, error)
	// FindByUUID(ctx context.Context, db Querier, uuid string) (*entity.User, error)
	// FindByID(ctx context.Context, db Querier, id string) (*entity.User, error)
	// FindByEmail(ctx context.Context, db Querier, email string) (*entity.User, error)
	// DeleteByID(ctx context.Context, db Querier, id string) error
	// DeleteByUUID(ctx context.Context, db Querier, uuid string) error
	// DeleteByEmail(ctx context.Context, db Querier, email string) error
	// ExistsByNameOrEmail(ctx context.Context, db Querier, name, email string) (bool, error)
	// UpdateName(ctx context.Context, db Querier, user *entity.User) (*entity.User, error)
}
type userRepositoryImpl struct {
	log *zap.Logger
}

func NewUserRepository(log *zap.Logger) UserRepository {
	return &userRepositoryImpl{log: log}
}

// func (r *userRepositoryImpl) Insert(ctx context.Context, db Querier, user *entity.User) (*entity.User, error) {
// 	query := `
// 	INSERT INTO users
// 		(id, name, email, password, uuid, created_at, updated_at
// 	VALUES
// 		($1, $2, $3, $4, $5, $6, $7)`

// 	_, err := db.Exec(ctx, query, user.ID, user.Name, user.Email, user.Password, user.UUID, user.CreatedAt, user.UpdatedAt)
// 	if err != nil {
// 		r.log.Error("failed to exec insert query", zap.String("query", query), zap.Error(err))
// 		return nil, err
// 	}
// 	return user, nil
// }

func (r *userRepositoryImpl) FindByUUID(ctx context.Context, db Querier, uuid string) (*entity.User, error) {
	user := new(entity.User)
	query := `SELECT * FROM users WHERE uuid = $1 AND deleted_at IS NULL`
	if err := pgxscan.Get(ctx, db, user, query, uuid); err != nil {
		r.log.Error("failed to get query", zap.String("query", query), zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, db Querier, id string) (*entity.User, error) {
	user := new(entity.User)
	query := `SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL`
	if err := pgxscan.Get(ctx, db, user, query, id); err != nil {
		r.log.Error("failed to get query", zap.String("query", query), zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, db Querier, email string) (*entity.User, error) {
	user := new(entity.User)
	query := `SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL`
	if err := pgxscan.Get(ctx, db, user, query, email); err != nil {
		r.log.Error("failed to get query", zap.String("query", query), zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (r *userRepositoryImpl) DeleteByID(ctx context.Context, db Querier, id string) error {
	query := `UPDATE users SET deleted_at = now() WHERE id = $1 AND deleted_at IS NOT NULL`
	_, err := db.Exec(ctx, query, id)
	if err != nil {
		r.log.Error("failed to exec delete query", zap.String("query", query), zap.Error(err))
		return err
	}
	return nil
}

func (r *userRepositoryImpl) DeleteByUUID(ctx context.Context, db Querier, uuid string) error {
	query := `UPDATE users SET deleted_at = now() WHERE uuid = $1 AND deleted_at IS NOT NULL`
	row, err := db.Exec(ctx, query, uuid)
	if err != nil {
		r.log.Error("failed to exec delete query", zap.String("query", query), zap.Error(err))
		return err
	}

	if row.RowsAffected() == 0 {
		return errors.New("row has been deleted")
	}

	return nil
}

func (r *userRepositoryImpl) DeleteByEmail(ctx context.Context, db Querier, email string) error {
	query := `UPDATE users SET deleted_at = now() WHERE email = $1 AND deleted_at IS NOT NULL`
	_, err := db.Exec(ctx, query, email)
	if err != nil {
		r.log.Error("failed to exec delete query", zap.String("query", query), zap.Error(err))
		return err

	}
	return nil
}

func (r *userRepositoryImpl) ExistsByNameOrEmail(ctx context.Context, db Querier, name, email string) (bool, error) {
	var total int
	query := `SELECT COUNT(*) FROM users WHERE name = $1 OR email = $2 AND deleted_at IS NULL`
	if err := pgxscan.Get(ctx, db, &total, query, name, email); err != nil {
		r.log.Error("failed to get query", zap.String("query", query), zap.Error(err))
		return false, err
	}

	if total > 0 {
		return true, nil
	}

	return false, nil
}

// func (r *userRepositoryImpl) UpdateName(ctx context.Context, db Querier, user *entity.User, oldUpdateAt *time.Time) (*entity.User, error) {
// 	query := "UPDATE users SET name = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL AND updated_at = $4"
// 	row, err := db.Exec(ctx, query, user.Name, user.UpdatedAt, user.ID, oldUpdateAt)
// 	if err != nil {
// 		r.log.Error("failed to exec delete query", zap.String("query", query), zap.Error(err))
// 		return nil, err
// 	}

// 	if row.RowsAffected() == 0 {
// 		return nil, errors.New("row has been deleted")
// 	}
// 	return user, nil
// }
