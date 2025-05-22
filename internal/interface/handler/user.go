package handler

import (
	"context"

	"connectrpc.com/connect"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/domain/repository"
	proto_user "github.com/Shakkuuu/sekai-songs-mylist/internal/gen/user/v1"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/infrastructure/auth"
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
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("user id not found in context"))
	}

	user, err := h.userUsecase.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, errors.WithStack(err))
		}
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
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
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("user id not found in context"))
	}

	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
	}

	if err := h.userUsecase.UpdateEmail(ctx, id, req.Msg.GetEmail()); err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	user, err := h.userUsecase.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, errors.WithStack(err))
		}
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
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
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("user id not found in context"))
	}

	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
	}
	if !containsLetter(req.Msg.GetNewPassword()) || !containsDigit(req.Msg.GetNewPassword()) {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("password must include at least one letter and one number"))
	}
	if req.Msg.GetNewPassword() != req.Msg.GetNewCheckPassword() {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("password, check_password not match"))
	}

	user, err := h.userUsecase.GetUserByID(ctx, id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Msg.GetOldPassword()))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.WithStack(err))
	} else if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Msg.GetNewPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	if err := h.userUsecase.UpdatePassword(ctx, id, string(hash)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	updatedUser, err := h.userUsecase.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, errors.WithStack(err))
		}
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
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
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("user id not found in context"))
	}

	exist, err := h.userUsecase.ExistsUserByID(ctx, id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}
	if !exist {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("This id user not found"))
	}

	if err := h.userUsecase.SoftDeleteUser(ctx, id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.WithStack(err))
	}

	return connect.NewResponse(&proto_user.DeleteUserResponse{}), nil
}
