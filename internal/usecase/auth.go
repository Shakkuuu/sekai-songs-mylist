package usecase

import (
	"context"
	"time"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/repository"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/infrastructure/auth"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE
//go:generate gotests -w -all $GOFILE

var (
	ErrDuplicateEmail            = errors.New("email already exists")
	ErrAlreadyDeletedUser        = errors.New("This user already deleted")
	ErrMismatchedHashAndPassword = errors.New("mismatched hash and password")
)

type AuthUsecase interface {
	Signup(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
}

type authUsecase struct {
	userRepo repository.UserRepository
}

func NewAuthUsecase(repo repository.UserRepository) AuthUsecase {
	return &authUsecase{
		userRepo: repo,
	}
}

func (u *authUsecase) Signup(ctx context.Context, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.WithStack(err)
	}

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
	if _, err := u.userRepo.CreateUser(ctx, id, email, string(hash), createdAt, updatedAt, deletedAt); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if !user.DeletedAt.IsZero() {
		return "", errors.WithStack(ErrAlreadyDeletedUser)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", errors.WithStack(ErrMismatchedHashAndPassword)
	} else if err != nil {
		return "", errors.WithStack(err)
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return token, nil
}
