package main

import (
	"errors"
	"log"
	"net"
	"os"
	"syscall"
	"os/signal"

	"github.com/sudo-JP/Load-Manager/backend/internal/database"
	"github.com/sudo-JP/Load-Manager/backend/internal/repository"
	"github.com/sudo-JP/Load-Manager/backend/internal/server"
	"github.com/sudo-JP/Load-Manager/backend/internal/service"

	// grpc
	pbOrder "github.com/sudo-JP/Load-Manager/backend/api/proto/order"
	pbProduct "github.com/sudo-JP/Load-Manager/backend/api/proto/product"
	pbUser "github.com/sudo-JP/Load-Manager/backend/api/proto/user"
	"google.golang.org/grpc"
)

// host,port,error
func parseCLI(args []string) (string, string, error) {
	host := "" 
	port := ""

	isHost := false 
	isPort := false

	for i := 1; i < len(args); i++ {
		if isHost && host != "" || isPort && port != "" {
			return "", "", errors.New("Invalid flag")
		}
		if isHost {
			host = args[i]
			isHost = false
		} else if isPort {
			port = args[i]
			isPort = false
		} else if args[i] == "--port" {
			isHost = true 
		} else if args[i] == "--host" {
			isPort = true
		} 
	}
	return host, port, nil
}

func main() {
	// Port and host 
	// --port 
	// --host 
	host, port, err := parseCLI(os.Args)
	if err != nil {
		log.Fatalf("Failed to Parse Args: %v", err)
		os.Exit(1)
	}

	// Database 
	db, err := database.DatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to Database: %v", err)
		os.Exit(2)
	}

	// Repository
	orderRepo := repository.NewOrderRepository(db)
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)

	// Services 
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo, userService, productService)

	grpcServer := grpc.NewServer()
	pbUser.RegisterUserServiceServer(grpcServer, server.NewUserServer(userService)) 
	pbOrder.RegisterOrderServiceServer(grpcServer, server.NewOrderServer(orderService))
	pbProduct.RegisterProductServiceServer(grpcServer, server.NewProductServer(productService))

	tcpListener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		os.Exit(3)
	}

	log.Printf("gRPC server listening on :%s", port)
	go func() {
		if err := grpcServer.Serve(tcpListener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
			os.Exit(4)
		}
	}()

	// shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Shutting down...")
	grpcServer.GracefulStop()
	os.Exit(0)
}

