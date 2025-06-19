package usecase

import (
	"context"
	"time"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/repository"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/pkg/auth"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/pkg/mail"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/$GOFILE -package=mock_$GOPACKAGE
//go:generate gotests -w -all $GOFILE

var (
	ErrDuplicateEmail            = errors.New("email already exists")
	ErrAlreadyDeletedUser        = errors.New("This user already deleted")
	ErrNotVerify                 = errors.New("This user not mail verify")
	ErrMismatchedHashAndPassword = errors.New("mismatched hash and password")
)

type AuthUsecase interface {
	Signup(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
	VerifyUser(ctx context.Context, id string, isVerified bool, verifyToken string) error
	ResendVerifyUser(ctx context.Context, id, email string) error
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
	isVerified := false
	verifyToken := uuid.New().String()
	tokenExpiresAt := time.Now().Add(3 * time.Hour)
	isAdmin := false
	now := time.Now()
	createdAt := now
	updatedAt := now
	deletedAt := time.Time{}
	if _, err := u.userRepo.CreateUser(ctx, id, email, string(hash), isVerified, verifyToken, tokenExpiresAt, isAdmin, createdAt, updatedAt, deletedAt); err != nil {
		return errors.WithStack(err)
	}

	if err := mail.SendVerificationEmail(ctx, email, verifyToken); err != nil {
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

	if !user.IsVerified {
		return "", errors.WithStack(ErrNotVerify)
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

func (u *authUsecase) VerifyUser(ctx context.Context, id string, isVerified bool, verifyToken string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := u.userRepo.UpdateUserIsVerified(ctx, parsedID, isVerified, time.Now()); err != nil {
		return errors.WithStack(err)
	}

	if err := u.userRepo.UpdateUserVerifyToken(ctx, parsedID, verifyToken, time.Now()); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (u *authUsecase) ResendVerifyUser(ctx context.Context, id, email string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.WithStack(err)
	}

	verifyToken := uuid.New().String()
	tokenExpiresAt := time.Now().Add(3 * time.Hour)

	if err := u.userRepo.UpdateUserVerifyToken(ctx, parsedID, verifyToken, time.Now()); err != nil {
		return errors.WithStack(err)
	}

	if err := u.userRepo.UpdateUserTokenExpiresAt(ctx, parsedID, tokenExpiresAt, time.Now()); err != nil {
		return errors.WithStack(err)
	}

	if err := mail.SendVerificationEmail(ctx, email, verifyToken); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
