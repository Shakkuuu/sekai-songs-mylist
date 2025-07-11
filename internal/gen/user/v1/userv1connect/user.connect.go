// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: user/v1/user.proto

package userv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/Shakkuuu/sekai-songs-mylist/internal/gen/user/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// UserServiceName is the fully-qualified name of the UserService service.
	UserServiceName = "user.v1.UserService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// UserServiceUserInfoProcedure is the fully-qualified name of the UserService's UserInfo RPC.
	UserServiceUserInfoProcedure = "/user.v1.UserService/UserInfo"
	// UserServiceLogoutProcedure is the fully-qualified name of the UserService's Logout RPC.
	UserServiceLogoutProcedure = "/user.v1.UserService/Logout"
	// UserServiceChangeEmailProcedure is the fully-qualified name of the UserService's ChangeEmail RPC.
	UserServiceChangeEmailProcedure = "/user.v1.UserService/ChangeEmail"
	// UserServiceChangePasswordProcedure is the fully-qualified name of the UserService's
	// ChangePassword RPC.
	UserServiceChangePasswordProcedure = "/user.v1.UserService/ChangePassword"
	// UserServiceDeleteUserProcedure is the fully-qualified name of the UserService's DeleteUser RPC.
	UserServiceDeleteUserProcedure = "/user.v1.UserService/DeleteUser"
	// UserServiceIsAdminProcedure is the fully-qualified name of the UserService's IsAdmin RPC.
	UserServiceIsAdminProcedure = "/user.v1.UserService/IsAdmin"
)

// UserServiceClient is a client for the user.v1.UserService service.
type UserServiceClient interface {
	UserInfo(context.Context, *connect.Request[v1.UserInfoRequest]) (*connect.Response[v1.UserInfoResponse], error)
	Logout(context.Context, *connect.Request[v1.LogoutRequest]) (*connect.Response[v1.LogoutResponse], error)
	ChangeEmail(context.Context, *connect.Request[v1.ChangeEmailRequest]) (*connect.Response[v1.ChangeEmailResponse], error)
	ChangePassword(context.Context, *connect.Request[v1.ChangePasswordRequest]) (*connect.Response[v1.ChangePasswordResponse], error)
	DeleteUser(context.Context, *connect.Request[v1.DeleteUserRequest]) (*connect.Response[v1.DeleteUserResponse], error)
	IsAdmin(context.Context, *connect.Request[v1.IsAdminRequest]) (*connect.Response[v1.IsAdminResponse], error)
}

// NewUserServiceClient constructs a client for the user.v1.UserService service. By default, it uses
// the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewUserServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) UserServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	userServiceMethods := v1.File_user_v1_user_proto.Services().ByName("UserService").Methods()
	return &userServiceClient{
		userInfo: connect.NewClient[v1.UserInfoRequest, v1.UserInfoResponse](
			httpClient,
			baseURL+UserServiceUserInfoProcedure,
			connect.WithSchema(userServiceMethods.ByName("UserInfo")),
			connect.WithClientOptions(opts...),
		),
		logout: connect.NewClient[v1.LogoutRequest, v1.LogoutResponse](
			httpClient,
			baseURL+UserServiceLogoutProcedure,
			connect.WithSchema(userServiceMethods.ByName("Logout")),
			connect.WithClientOptions(opts...),
		),
		changeEmail: connect.NewClient[v1.ChangeEmailRequest, v1.ChangeEmailResponse](
			httpClient,
			baseURL+UserServiceChangeEmailProcedure,
			connect.WithSchema(userServiceMethods.ByName("ChangeEmail")),
			connect.WithClientOptions(opts...),
		),
		changePassword: connect.NewClient[v1.ChangePasswordRequest, v1.ChangePasswordResponse](
			httpClient,
			baseURL+UserServiceChangePasswordProcedure,
			connect.WithSchema(userServiceMethods.ByName("ChangePassword")),
			connect.WithClientOptions(opts...),
		),
		deleteUser: connect.NewClient[v1.DeleteUserRequest, v1.DeleteUserResponse](
			httpClient,
			baseURL+UserServiceDeleteUserProcedure,
			connect.WithSchema(userServiceMethods.ByName("DeleteUser")),
			connect.WithClientOptions(opts...),
		),
		isAdmin: connect.NewClient[v1.IsAdminRequest, v1.IsAdminResponse](
			httpClient,
			baseURL+UserServiceIsAdminProcedure,
			connect.WithSchema(userServiceMethods.ByName("IsAdmin")),
			connect.WithClientOptions(opts...),
		),
	}
}

// userServiceClient implements UserServiceClient.
type userServiceClient struct {
	userInfo       *connect.Client[v1.UserInfoRequest, v1.UserInfoResponse]
	logout         *connect.Client[v1.LogoutRequest, v1.LogoutResponse]
	changeEmail    *connect.Client[v1.ChangeEmailRequest, v1.ChangeEmailResponse]
	changePassword *connect.Client[v1.ChangePasswordRequest, v1.ChangePasswordResponse]
	deleteUser     *connect.Client[v1.DeleteUserRequest, v1.DeleteUserResponse]
	isAdmin        *connect.Client[v1.IsAdminRequest, v1.IsAdminResponse]
}

// UserInfo calls user.v1.UserService.UserInfo.
func (c *userServiceClient) UserInfo(ctx context.Context, req *connect.Request[v1.UserInfoRequest]) (*connect.Response[v1.UserInfoResponse], error) {
	return c.userInfo.CallUnary(ctx, req)
}

// Logout calls user.v1.UserService.Logout.
func (c *userServiceClient) Logout(ctx context.Context, req *connect.Request[v1.LogoutRequest]) (*connect.Response[v1.LogoutResponse], error) {
	return c.logout.CallUnary(ctx, req)
}

// ChangeEmail calls user.v1.UserService.ChangeEmail.
func (c *userServiceClient) ChangeEmail(ctx context.Context, req *connect.Request[v1.ChangeEmailRequest]) (*connect.Response[v1.ChangeEmailResponse], error) {
	return c.changeEmail.CallUnary(ctx, req)
}

// ChangePassword calls user.v1.UserService.ChangePassword.
func (c *userServiceClient) ChangePassword(ctx context.Context, req *connect.Request[v1.ChangePasswordRequest]) (*connect.Response[v1.ChangePasswordResponse], error) {
	return c.changePassword.CallUnary(ctx, req)
}

// DeleteUser calls user.v1.UserService.DeleteUser.
func (c *userServiceClient) DeleteUser(ctx context.Context, req *connect.Request[v1.DeleteUserRequest]) (*connect.Response[v1.DeleteUserResponse], error) {
	return c.deleteUser.CallUnary(ctx, req)
}

// IsAdmin calls user.v1.UserService.IsAdmin.
func (c *userServiceClient) IsAdmin(ctx context.Context, req *connect.Request[v1.IsAdminRequest]) (*connect.Response[v1.IsAdminResponse], error) {
	return c.isAdmin.CallUnary(ctx, req)
}

// UserServiceHandler is an implementation of the user.v1.UserService service.
type UserServiceHandler interface {
	UserInfo(context.Context, *connect.Request[v1.UserInfoRequest]) (*connect.Response[v1.UserInfoResponse], error)
	Logout(context.Context, *connect.Request[v1.LogoutRequest]) (*connect.Response[v1.LogoutResponse], error)
	ChangeEmail(context.Context, *connect.Request[v1.ChangeEmailRequest]) (*connect.Response[v1.ChangeEmailResponse], error)
	ChangePassword(context.Context, *connect.Request[v1.ChangePasswordRequest]) (*connect.Response[v1.ChangePasswordResponse], error)
	DeleteUser(context.Context, *connect.Request[v1.DeleteUserRequest]) (*connect.Response[v1.DeleteUserResponse], error)
	IsAdmin(context.Context, *connect.Request[v1.IsAdminRequest]) (*connect.Response[v1.IsAdminResponse], error)
}

// NewUserServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewUserServiceHandler(svc UserServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	userServiceMethods := v1.File_user_v1_user_proto.Services().ByName("UserService").Methods()
	userServiceUserInfoHandler := connect.NewUnaryHandler(
		UserServiceUserInfoProcedure,
		svc.UserInfo,
		connect.WithSchema(userServiceMethods.ByName("UserInfo")),
		connect.WithHandlerOptions(opts...),
	)
	userServiceLogoutHandler := connect.NewUnaryHandler(
		UserServiceLogoutProcedure,
		svc.Logout,
		connect.WithSchema(userServiceMethods.ByName("Logout")),
		connect.WithHandlerOptions(opts...),
	)
	userServiceChangeEmailHandler := connect.NewUnaryHandler(
		UserServiceChangeEmailProcedure,
		svc.ChangeEmail,
		connect.WithSchema(userServiceMethods.ByName("ChangeEmail")),
		connect.WithHandlerOptions(opts...),
	)
	userServiceChangePasswordHandler := connect.NewUnaryHandler(
		UserServiceChangePasswordProcedure,
		svc.ChangePassword,
		connect.WithSchema(userServiceMethods.ByName("ChangePassword")),
		connect.WithHandlerOptions(opts...),
	)
	userServiceDeleteUserHandler := connect.NewUnaryHandler(
		UserServiceDeleteUserProcedure,
		svc.DeleteUser,
		connect.WithSchema(userServiceMethods.ByName("DeleteUser")),
		connect.WithHandlerOptions(opts...),
	)
	userServiceIsAdminHandler := connect.NewUnaryHandler(
		UserServiceIsAdminProcedure,
		svc.IsAdmin,
		connect.WithSchema(userServiceMethods.ByName("IsAdmin")),
		connect.WithHandlerOptions(opts...),
	)
	return "/user.v1.UserService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case UserServiceUserInfoProcedure:
			userServiceUserInfoHandler.ServeHTTP(w, r)
		case UserServiceLogoutProcedure:
			userServiceLogoutHandler.ServeHTTP(w, r)
		case UserServiceChangeEmailProcedure:
			userServiceChangeEmailHandler.ServeHTTP(w, r)
		case UserServiceChangePasswordProcedure:
			userServiceChangePasswordHandler.ServeHTTP(w, r)
		case UserServiceDeleteUserProcedure:
			userServiceDeleteUserHandler.ServeHTTP(w, r)
		case UserServiceIsAdminProcedure:
			userServiceIsAdminHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedUserServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedUserServiceHandler struct{}

func (UnimplementedUserServiceHandler) UserInfo(context.Context, *connect.Request[v1.UserInfoRequest]) (*connect.Response[v1.UserInfoResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.v1.UserService.UserInfo is not implemented"))
}

func (UnimplementedUserServiceHandler) Logout(context.Context, *connect.Request[v1.LogoutRequest]) (*connect.Response[v1.LogoutResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.v1.UserService.Logout is not implemented"))
}

func (UnimplementedUserServiceHandler) ChangeEmail(context.Context, *connect.Request[v1.ChangeEmailRequest]) (*connect.Response[v1.ChangeEmailResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.v1.UserService.ChangeEmail is not implemented"))
}

func (UnimplementedUserServiceHandler) ChangePassword(context.Context, *connect.Request[v1.ChangePasswordRequest]) (*connect.Response[v1.ChangePasswordResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.v1.UserService.ChangePassword is not implemented"))
}

func (UnimplementedUserServiceHandler) DeleteUser(context.Context, *connect.Request[v1.DeleteUserRequest]) (*connect.Response[v1.DeleteUserResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.v1.UserService.DeleteUser is not implemented"))
}

func (UnimplementedUserServiceHandler) IsAdmin(context.Context, *connect.Request[v1.IsAdminRequest]) (*connect.Response[v1.IsAdminResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.v1.UserService.IsAdmin is not implemented"))
}
