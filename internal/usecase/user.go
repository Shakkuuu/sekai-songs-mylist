package usecase

import (
	"context"
	"time"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/entity"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/repository"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE
//go:generate gotests -w -all $GOFILE

var (
	ErrUserNotFound = errors.New("this id user not found")
)

type UserUsecase interface {
	UserInfo(ctx context.Context, id string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	ChangeEmail(ctx context.Context, id, email string) (*entity.User, error)
	ChangePassword(ctx context.Context, id, oldPassword, newPassword string) (*entity.User, error)
	DeleteUser(ctx context.Context, id string) error
	IsAdmin(ctx context.Context, id string) (bool, error)
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

func (u *userUsecase) UserInfo(ctx context.Context, id string) (*entity.User, error) {
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

func (u *userUsecase) ChangeEmail(ctx context.Context, id, email string) (*entity.User, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := u.userRepo.UpdateUserEmail(ctx, parsedID, email, time.Now()); err != nil {
		return nil, errors.WithStack(err)
	}

	user, err := u.userRepo.GetUserByID(ctx, parsedID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return user, nil
}

func (u *userUsecase) ChangePassword(ctx context.Context, id, oldPassword, newPassword string) (*entity.User, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	user, err := u.userRepo.GetUserByID(ctx, parsedID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, errors.WithStack(ErrMismatchedHashAndPassword)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := u.userRepo.UpdateUserPassword(ctx, parsedID, string(hash), time.Now()); err != nil {
		return nil, errors.WithStack(err)
	}

	updatedUser, err := u.userRepo.GetUserByID(ctx, parsedID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return updatedUser, nil
}

func (u *userUsecase) DeleteUser(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.WithStack(err)
	}

	exist, err := u.userRepo.ExistsUserByID(ctx, parsedID)
	if err != nil {
		return errors.WithStack(err)
	}
	if !exist {
		return errors.WithStack(ErrUserNotFound)
	}

	if err := u.userRepo.SoftDeleteUser(ctx, parsedID); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (u *userUsecase) IsAdmin(ctx context.Context, id string) (bool, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return false, errors.WithStack(err)
	}

	exist, err := u.userRepo.ExistsUserByID(ctx, parsedID)
	if err != nil {
		return false, errors.WithStack(err)
	}
	if !exist {
		return false, errors.WithStack(ErrUserNotFound)
	}

	isAdmin, err := u.userRepo.IsAdminByID(ctx, parsedID)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return isAdmin, nil
}
