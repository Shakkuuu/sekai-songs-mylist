package handler

import (
	"context"
	"log"
	"net/http"
	"time"
	"unicode"

	"connectrpc.com/connect"

	proto_auth "github.com/Shakkuuu/sekai-songs-mylist/internal/gen/auth/v1"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/usecase"
	"github.com/cockroachdb/errors"
)

//go:generate gotests -w -all $GOFILE

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	userUsecase usecase.UserUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase, userUsecase usecase.UserUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase, userUsecase: userUsecase}
}

func (h *AuthHandler) Signup(ctx context.Context, req *connect.Request[proto_auth.SignupRequest]) (*connect.Response[proto_auth.SignupResponse], error) {
	if err := req.Msg.Validate(); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
	}
	if !containsLetter(req.Msg.GetPassword()) || !containsDigit(req.Msg.GetPassword()) {
		err := errors.New("password must include at least one letter and one number")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
	}
	if req.Msg.GetPassword() != req.Msg.GetCheckPassword() {
		err := errors.New("password, check_password not match")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
	}

	if err := h.authUsecase.Signup(ctx, req.Msg.GetEmail(), req.Msg.GetPassword()); err != nil {
		if errors.Is(err, usecase.ErrDuplicateEmail) {
			cerr := errors.WithStack(err)
			log.Printf("%+v\n", cerr)
			return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
		}
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInternal, cerr)
	}

	return connect.NewResponse(&proto_auth.SignupResponse{}), nil
}

func (h *AuthHandler) Login(ctx context.Context, req *connect.Request[proto_auth.LoginRequest]) (*connect.Response[proto_auth.LoginResponse], error) {
	if err := req.Msg.Validate(); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
	}

	token, err := h.authUsecase.Login(ctx, req.Msg.GetEmail(), req.Msg.GetPassword())
	if err != nil {
		if errors.Is(err, usecase.ErrAlreadyDeletedUser) {
			cerr := errors.WithStack(err)
			log.Printf("%+v\n", cerr)
			return nil, connect.NewError(connect.CodeInvalidArgument, cerr)
		}
		if errors.Is(err, usecase.ErrNotVerify) {
			cerr := errors.WithStack(err)
			log.Printf("%+v\n", cerr)
			return nil, connect.NewError(connect.CodeUnauthenticated, cerr)
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

func (h *AuthHandler) VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	toEmail := r.URL.Query().Get("email")

	ctx := r.Context()

	user, err := h.userUsecase.GetUserByEmail(ctx, toEmail)
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, "無効なメールアドレス", http.StatusBadRequest)
		return
	}

	if user.VerifyToken != token {
		http.Error(w, "無効なトークン", http.StatusBadRequest)
		return
	}

	if time.Now().After(user.TokenExpiresAt) {
		http.Error(w, "トークンの有効期限が切れています", http.StatusBadRequest)
		return
	}

	if err := h.authUsecase.VerifyUser(ctx, user.ID, true, ""); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, "認証処理に失敗しました", http.StatusBadRequest)
		return
	}

	if _, err := w.Write([]byte("認証が完了しました。ログインをしてください。")); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) ResendVerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	toEmail := r.URL.Query().Get("email")

	ctx := r.Context()

	user, err := h.userUsecase.GetUserByEmail(ctx, toEmail)
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, "無効なメールアドレス", http.StatusBadRequest)
		return
	}

	if user.VerifyToken != token {
		http.Error(w, "無効な旧トークン", http.StatusBadRequest)
		return
	}

	if err := h.authUsecase.ResendVerifyUser(ctx, user.ID, user.Email); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, "トークン更新処理に失敗しました", http.StatusBadRequest)
		return
	}

	if _, err := w.Write([]byte("再度認証メールを送信しました。メールを確認してください。")); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}
