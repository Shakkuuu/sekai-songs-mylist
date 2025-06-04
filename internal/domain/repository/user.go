package repository

import (
	"context"
	"time"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/entity"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/infrastructure/db/sqlcgen"
	"github.com/google/uuid"
)

//go:generate mockgen -source=$GOFILE -destination=../../mock/$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE

type UserRepository interface {
	ListUsers(ctx context.Context) ([]*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, id uuid.UUID, emali, password string, createdAt, updatedAt, deletedAt time.Time) (*sqlcgen.User, error)
	ExistsUserByEmail(ctx context.Context, email string) (bool, error)
	ExistsUserByID(ctx context.Context, id uuid.UUID) (bool, error)
	UpdateUserEmail(ctx context.Context, id uuid.UUID, email string, updatedAt time.Time) error
	UpdateUserPassword(ctx context.Context, id uuid.UUID, password string, updatedAt time.Time) error
	SoftDeleteUser(ctx context.Context, id uuid.UUID) error
}
