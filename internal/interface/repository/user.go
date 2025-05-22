package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/entity"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/repository"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/infrastructure/db/sqlcgen"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

//go:generate gotests -w -all $GOFILE

type userRepository struct {
	queries *sqlcgen.Queries
}

func NewUserRepository(queries *sqlcgen.Queries) repository.UserRepository {
	return &userRepository{queries: queries}
}

func (r *userRepository) ListUsers(ctx context.Context) ([]*entity.User, error) {
	users, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userPointers := make([]*entity.User, len(users))
	for i := range users {
		userPointers[i] = sqlToDomainUser(&users[i])
	}

	return userPointers, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.WithStack(repository.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}

	return sqlToDomainUser(&user), nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.WithStack(repository.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}

	return sqlToDomainUser(&user), nil
}

func (r *userRepository) CreateUser(ctx context.Context, id uuid.UUID, email, password string, createdAt, updatedAt, deletedAt time.Time) (*sqlcgen.User, error) {
	sqlUser := sqlcgen.InsertUserParams{
		ID:        id,
		Email:     email,
		Password:  password,
		CreatedAt: sql.NullTime{Time: createdAt, Valid: true},
		UpdatedAt: sql.NullTime{Time: updatedAt, Valid: true},
		DeletedAt: sql.NullTime{Time: deletedAt, Valid: true},
	}
	u, err := r.queries.InsertUser(ctx, sqlUser)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &u, nil
}

func (r *userRepository) ExistsUserByEmail(ctx context.Context, email string) (bool, error) {
	exist, err := r.queries.ExistsUserByEmail(ctx, email)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}

func (r *userRepository) ExistsUserByID(ctx context.Context, id uuid.UUID) (bool, error) {
	exist, err := r.queries.ExistsUserByID(ctx, id)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}

func (r *userRepository) UpdateEmail(ctx context.Context, id uuid.UUID, email string) error {
	arg := sqlcgen.UpdateEmailParams{
		Email:     email,
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		ID:        id,
	}

	if err := r.queries.UpdateEmail(ctx, arg); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	arg := sqlcgen.UpdatePasswordParams{
		Password:  password,
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		ID:        id,
	}

	if err := r.queries.UpdatePassword(ctx, arg); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *userRepository) SoftDeleteUser(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	arg := sqlcgen.SoftDeleteUserParams{
		DeletedAt: sql.NullTime{Time: now, Valid: true},
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
		ID:        id,
	}

	if err := r.queries.SoftDeleteUser(ctx, arg); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func sqlToDomainUser(sqlUser *sqlcgen.User) *entity.User {
	return &entity.User{
		ID:        sqlUser.ID.String(),
		Email:     sqlUser.Email,
		Password:  sqlUser.Password,
		CreatedAt: sqlUser.CreatedAt.Time,
		UpdatedAt: sqlUser.UpdatedAt.Time,
		DeletedAt: sqlUser.DeletedAt.Time,
	}
}
