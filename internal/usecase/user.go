package usecase

import (
	"context"
	"time"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/entity"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/repository"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE
//go:generate gotests -w -all $GOFILE

var (
	ErrDuplicateEmail = errors.New("email already exists")
)

type UserUsecase interface {
	ListUsers(ctx context.Context) ([]*entity.User, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, email, password string) error
	ExistsUserByEmail(ctx context.Context, email string) (bool, error)
	ExistsUserByID(ctx context.Context, id string) (bool, error)
	UpdateEmail(ctx context.Context, id, email string) error
	UpdatePassword(ctx context.Context, id, password string) error
	SoftDeleteUser(ctx context.Context, id string) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: repo,
	}
}

func (u *userUsecase) ListUsers(ctx context.Context) ([]*entity.User, error) {
	users, err := u.userRepo.ListUsers(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if users == nil {
		return []*entity.User{}, nil
	}

	return users, nil
}

func (u *userUsecase) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	user, err := u.userRepo.GetUserByID(ctx, parsedID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return user, nil
}

func (u *userUsecase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return user, nil
}

func (u *userUsecase) CreateUser(ctx context.Context, email, password string) error {
	userExist, err := u.userRepo.ExistsUserByEmail(ctx, email)
	if err != nil {
		return errors.WithStack(err)
	}
	if userExist {
		return errors.WithStack(ErrDuplicateEmail)
	}

	id := uuid.New()
	now := time.Now()
	createdAt := now
	updatedAt := now
	deletedAt := time.Time{}
	if _, err := u.userRepo.CreateUser(ctx, id, email, password, createdAt, updatedAt, deletedAt); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (u *userUsecase) ExistsUserByEmail(ctx context.Context, email string) (bool, error) {
	exist, err := u.userRepo.ExistsUserByEmail(ctx, email)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}

func (u *userUsecase) ExistsUserByID(ctx context.Context, id string) (bool, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return false, errors.WithStack(err)
	}

	exist, err := u.userRepo.ExistsUserByID(ctx, parsedID)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exist, nil
}

func (u *userUsecase) UpdateEmail(ctx context.Context, id, email string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := u.userRepo.UpdateEmail(ctx, parsedID, email); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (u *userUsecase) UpdatePassword(ctx context.Context, id, password string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := u.userRepo.UpdatePassword(ctx, parsedID, password); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (u *userUsecase) SoftDeleteUser(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := u.userRepo.SoftDeleteUser(ctx, parsedID); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// SoftDeleteUser、Update系、Existのusecaseを作る
// handlerでauth.go作って、loginとsignup
// signupのhandlerでpassword暗号化
// loginのhandlerでパスワード一致チェック。tokenも生成
