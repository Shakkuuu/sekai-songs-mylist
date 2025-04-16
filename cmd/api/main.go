package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/cockroachdb/errors"
	"google.golang.org/grpc"

	"github.com/Shakkuuu/sekai-songs-mylist/config"
	proto_master "github.com/Shakkuuu/sekai-songs-mylist/internal/gen/master"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/infrastructure/db"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/interface/handler"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/interface/repository"
	"github.com/Shakkuuu/sekai-songs-mylist/internal/usecase"
)

func main() {
	cfg, err := config.NewConfig()
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("error Config not found: %v", errors.GetReportableStackTrace(err))
		os.Exit(1)
	} else if err != nil {
		log.Printf("error Failed to load config: %v", errors.GetReportableStackTrace(err))
		os.Exit(1)
	}

	dbConfig := db.DBConfing{
		Host:     cfg.DBHost,
		User:     cfg.DBUserName,
		Password: cfg.DBUserPassword,
		DBName:   cfg.DBName,
		Port:     cfg.DBPort,
	}

    dbConn, queries, err := db.Init(dbConfig)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer dbConn.Close()

    masterRepository := repository.NewMasterRepository(queries)
    artistUsecase := usecase.NewMasterUsecase(masterRepository)
    artistHandler := handler.NewMasterHandler(artistUsecase)

    // gRPC server setup
    grpcServer := grpc.NewServer()
    proto_master.RegisterMasterServiceServer(grpcServer, artistHandler)

	// Start gRPC server
    listener, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.ServerPort))
    if err != nil {
        log.Fatalf("Failed to listen on port "+strconv.Itoa(cfg.ServerPort)+": %v", err)
    }
    log.Println("gRPC server is running on port "+strconv.Itoa(cfg.ServerPort))

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        if err := grpcServer.Serve(listener); err != nil {
            log.Fatalf("Failed to serve gRPC server: %v", err)
        }
    }()

    <-quit
    log.Println("Shutting down gRPC server...")

    grpcServer.GracefulStop()

    time.Sleep(10 * time.Second)
    log.Println("gRPC server stopped")
}
