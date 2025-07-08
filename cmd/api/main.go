package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/rs/cors"

	"github.com/Shakkuuu/sekai-songs-mylist/config"
	proto_auth_connect "github.com/Shakkuuu/sekai-songs-mylist/internal/gen/auth/v1/authv1connect"
	proto_master_connect "github.com/Shakkuuu/sekai-songs-mylist/internal/gen/master/masterconnect"
	proto_my_list_connect "github.com/Shakkuuu/sekai-songs-mylist/internal/gen/mylist/v1/mylistv1connect"
	proto_user_connect "github.com/Shakkuuu/sekai-songs-mylist/internal/gen/user/v1/userv1connect"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/infrastructure/db"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/infrastructure/redis"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/interface/handler"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/interface/repository"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/pkg/auth"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/pkg/googleoauth"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/usecase"
)

func main() {
	cfg, err := config.NewConfig()
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("error Config not found: \n%+v\n", err)
		os.Exit(1)
	} else if err != nil {
		log.Printf("error Failed to load config: \n%+v\n", err)
		os.Exit(1)
	}

	dbConfig := db.DBConfig{
		Host:     cfg.DBHost,
		User:     cfg.DBUserName,
		Password: cfg.DBUserPassword,
		DBName:   cfg.DBName,
		Port:     cfg.DBPort,
	}

	dbConn, queries, err := db.Init(dbConfig)
	if err != nil {
		log.Printf("Failed to initialize database: \n%+v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Printf("failed to close db connection: %v", err)
		}
	}()

	redisConfig := redis.RedisConfing{
		Host: cfg.RedisHost,
		Port: cfg.RedisPort,
	}
	rc := redis.Init(redisConfig)
	defer func() {
		if err := rc.Close(); err != nil {
			log.Printf("failed to close redis connection: %v", err)
		}
	}()

	if err := googleoauth.Init(); err != nil {
		log.Println(err)
		log.Printf("Failed to initialize mail: \n%+v\n", err)
		os.Exit(1)
	}

	redisMasterCacheRepository := repository.NewRedisMasterCacheRepository(rc)
	masterRepository := repository.NewMasterRepository(queries)
	userRepository := repository.NewUserRepository(queries)
	masterUsecase := usecase.NewMasterUsecase(masterRepository, redisMasterCacheRepository)
	userUsecase := usecase.NewUserUsecase(userRepository)
	masterHandler := handler.NewMasterHandler(masterUsecase, userUsecase)
	authUsecase := usecase.NewAuthUsecase(userRepository)
	authHandler := handler.NewAuthHandler(authUsecase, userUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	myListRepository := repository.NewMyListRepository(queries)
	myListUsecase := usecase.NewMyListUsecase(myListRepository, masterRepository, redisMasterCacheRepository)
	myListHandler := handler.NewMyListHandler(myListUsecase)
	strageHandler := handler.NewStorageHandler()

	mux := http.NewServeMux()
	mux.Handle(
		proto_master_connect.NewMasterServiceHandler(
			masterHandler,
		),
	)
	mux.Handle(
		proto_auth_connect.NewAuthServiceHandler(
			authHandler,
		),
	)
	mux.Handle(
		proto_user_connect.NewUserServiceHandler(
			userHandler,
			connect.WithInterceptors(auth.AuthInterceptor()),
		),
	)
	mux.Handle(
		proto_my_list_connect.NewMyListServiceHandler(
			myListHandler,
			connect.WithInterceptors(auth.AuthInterceptor()),
		),
	)

	mux.HandleFunc("/verify", authHandler.VerifyEmailHandler)
	mux.HandleFunc("/verify/resend", authHandler.ResendVerifyEmailHandler)
	mux.HandleFunc("/upload/thumbnail", strageHandler.UploadThumbnailHandler)
	mux.HandleFunc("/upload/attachment", strageHandler.UploadAttachmentHandler)
	mux.HandleFunc("/image", strageHandler.GetImageHandler)
	mux.HandleFunc("/delete/attachment", strageHandler.DeleteAttachmentHandler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.FrontEndURL}, // フロントエンドのURL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "connect-protocol-version", "Authorization"},
		AllowCredentials: true,
	})
	handler := corsHandler.Handler(mux)

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.ServerPort),
		Handler: handler,
	}

	log.Println("Starting server on :" + strconv.Itoa(cfg.ServerPort))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to serve server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
