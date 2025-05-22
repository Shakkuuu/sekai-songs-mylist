package handler

import (
	"context"
	"unicode"

	"connectrpc.com/connect"
	"golang.org/x/crypto/bcrypt"

	proto_auth "github.com/Shakkuuu/sekai-songs-mylist/internal/gen/auth/v1"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/infrastructure/auth"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/usecase"
	"github.com/cockroachdb/errors"
)

//go:generate gotests -w -all $GOFILE

type AuthHandler struct {
	userUsecase usecase.UserUsecase
}

func NewAuthHandler(userUsecase usecase.UserUsecase) *AuthHandler {
	return &AuthHandler{userUsecase: userUsecase}
}

func (h *AuthHandler) Signup(ctx context.Context, req *connect.Request[proto_auth.SignupRequest]) (*connect.Response[proto_auth.SignupResponse], error) {
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
	}
	if !containsLetter(req.Msg.GetPassword()) || !containsDigit(req.Msg.GetPassword()) {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("password must include at least one letter and one number"))
	}
	if req.Msg.GetPassword() != req.Msg.GetCheckPassword() {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("password, check_password not match"))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Msg.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	if err := h.userUsecase.CreateUser(ctx, req.Msg.GetEmail(), string(hash)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	return connect.NewResponse(&proto_auth.SignupResponse{}), nil
}

func (h *AuthHandler) Login(ctx context.Context, req *connect.Request[proto_auth.LoginRequest]) (*connect.Response[proto_auth.LoginResponse], error) {
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
	}

	user, err := h.userUsecase.GetUserByEmail(ctx, req.Msg.GetEmail())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	if !user.DeletedAt.IsZero() {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("this user already delete"))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Msg.GetPassword()))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
	} else if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return connect.NewResponse(&proto_auth.LoginResponse{
		Token: token,
	}), nil
}

func containsLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func containsDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}
