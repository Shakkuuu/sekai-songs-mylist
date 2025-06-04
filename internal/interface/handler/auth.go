package handler

import (
	"context"
	"unicode"

	"connectrpc.com/connect"

	proto_auth "github.com/Shakkuuu/sekai-songs-mylist/internal/gen/auth/v1"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/usecase"
	"github.com/cockroachdb/errors"
)

//go:generate gotests -w -all $GOFILE

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
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

	if err := h.authUsecase.Signup(ctx, req.Msg.GetEmail(), req.Msg.GetPassword()); err != nil {
		if errors.Is(err, usecase.ErrDuplicateEmail) {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
		}
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	return connect.NewResponse(&proto_auth.SignupResponse{}), nil
}

func (h *AuthHandler) Login(ctx context.Context, req *connect.Request[proto_auth.LoginRequest]) (*connect.Response[proto_auth.LoginResponse], error) {
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
	}

	token, err := h.authUsecase.Login(ctx, req.Msg.GetEmail(), req.Msg.GetPassword())
	if err != nil {
		if errors.Is(err, usecase.ErrAlreadyDeletedUser) {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
		}
		if errors.Is(err, usecase.ErrMismatchedHashAndPassword) {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
		}
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
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
