package handler

import (
	"context"
	"log"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/repository"
	proto_user "github.com/Shakkuuu/sekai-songs-mylist/internal/gen/user/v1"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/pkg/auth"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/usecase"
	"github.com/cockroachdb/errors"
)

//go:generate gotests -w -all $GOFILE

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) UserInfo(ctx context.Context, req *connect.Request[proto_user.UserInfoRequest]) (*connect.Response[proto_user.UserInfoResponse], error) {
	id, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		err := errors.New("user id not found in context")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeUnauthenticated, cerr)
	}

	user, err := h.userUsecase.UserInfo(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			cerr := errors.WithStack(err)
			log.Printf("%+v\n", cerr)
			return nil, connect.NewError(connect.CodeNotFound, cerr)
		}
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	return connect.NewResponse(&proto_user.UserInfoResponse{
		Id:        user.ID,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}), nil
}

func (h *UserHandler) Logout(ctx context.Context, req *connect.Request[proto_user.LogoutRequest]) (*connect.Response[proto_user.LogoutResponse], error) {
	// id, ok := ctx.Value(auth.UserIDKey).(string)
	// if !ok {
	// 	return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("user id not found in context"))
	// }

	// 	exist, err := h.UserUsecase.ExistsUserByID(ctx, id)
	// 	if err != nil {
	// 		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	// 	}
	// 	if !exist {
	// 		return nil, connect.NewError(connect.CodeNotFound, errors.New("This id user not found"))
	// 	}
	return connect.NewResponse(&proto_user.LogoutResponse{}), nil
}

func (h *UserHandler) ChangeEmail(ctx context.Context, req *connect.Request[proto_user.ChangeEmailRequest]) (*connect.Response[proto_user.ChangeEmailResponse], error) {
	id, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		err := errors.New("user id not found in context")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeUnauthenticated, cerr)
	}

	if err := req.Msg.Validate(); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
	}

	user, err := h.userUsecase.ChangeEmail(ctx, id, req.Msg.GetEmail())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			cerr := errors.WithStack(err)
			log.Printf("%+v\n", cerr)
			return nil, connect.NewError(connect.CodeNotFound, cerr)
		}
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	return connect.NewResponse(&proto_user.ChangeEmailResponse{
		Id:        user.ID,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}), nil
}

func (h *UserHandler) ChangePassword(ctx context.Context, req *connect.Request[proto_user.ChangePasswordRequest]) (*connect.Response[proto_user.ChangePasswordResponse], error) {
	id, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		err := errors.New("user id not found in context")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeUnauthenticated, cerr)
	}

	if err := req.Msg.Validate(); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
	}
	if !containsLetter(req.Msg.GetNewPassword()) || !containsDigit(req.Msg.GetNewPassword()) {
		err := errors.New("password must include at least one letter and one number")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
	}
	if req.Msg.GetNewPassword() != req.Msg.GetNewCheckPassword() {
		err := errors.New("password, check_password not match")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
	}

	updatedUser, err := h.userUsecase.ChangePassword(ctx, id, req.Msg.GetOldPassword(), req.Msg.GetNewPassword())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			cerr := errors.WithStack(err)
			log.Printf("%+v\n", cerr)
			return nil, connect.NewError(connect.CodeNotFound, cerr)
		}
		if errors.Is(err, usecase.ErrMismatchedHashAndPassword) {
			cerr := errors.WithStack(err)
			log.Printf("%+v\n", cerr)
			return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
		}
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	return connect.NewResponse(&proto_user.ChangePasswordResponse{
		Id:        updatedUser.ID,
		Email:     updatedUser.Email,
		CreatedAt: timestamppb.New(updatedUser.CreatedAt),
	}), nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *connect.Request[proto_user.DeleteUserRequest]) (*connect.Response[proto_user.DeleteUserResponse], error) {
	id, ok := ctx.Value(auth.UserIDKey).(string)
	if !ok {
		err := errors.New("user id not found in context")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeUnauthenticated, cerr)
	}

	if err := h.userUsecase.DeleteUser(ctx, id); err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			cerr := errors.WithStack(err)
			log.Printf("%+v\n", cerr)
			return nil, connect.NewError(connect.CodeNotFound, cerr)
		}
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	return connect.NewResponse(&proto_user.DeleteUserResponse{}), nil
}
