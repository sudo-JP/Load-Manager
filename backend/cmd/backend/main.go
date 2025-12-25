package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sudo-JP/Load-Manager/backend/internal/database"
	"github.com/sudo-JP/Load-Manager/backend/internal/repository"
	"github.com/sudo-JP/Load-Manager/backend/internal/routes"
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
			return "", "", fmt.Errorf("invalid flag")
		}
		if isHost {
			host = args[i]
			isHost = false
		} else if isPort {
			port = args[i]
			isPort = false
		} else if args[i] == "--port" {
			isPort = true
		} else if args[i] == "--host" {
			isHost = true
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
	defer db.Close()

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

	tcpListener, err := net.Listen("tcp", host+":"+port)
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

	// HTTP Server (for direct testing and comparison)
	httpServer := startHTTPServer(userService, productService, orderService)
	log.Printf("HTTP server listening on :9000")
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
			os.Exit(5)
		}
	}()

	// shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down...")

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	grpcServer.GracefulStop()
	os.Exit(0)
}

func startHTTPServer(
	userService service.UserServiceInterface,
	productService service.ProductServiceInterface,
	orderService service.OrderServiceInterface,
) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Create handlers
	userHandler := routes.NewUserHandler(userService)
	productHandler := routes.NewProductHandler(productService)
	orderHandler := routes.NewOrderHandler(orderService)

	// User routes
	router.POST("/users", userHandler.CreateUser)
	router.GET("/users", userHandler.GetUser)
	router.PUT("/users", userHandler.UpdateUser)
	router.DELETE("/users", userHandler.DeleteUser)

	// Product routes
	router.POST("/products", productHandler.CreateProduct)
	router.GET("/products", productHandler.GetProduct)
	router.PUT("/products", productHandler.UpdateProduct)
	router.DELETE("/products", productHandler.DeleteProduct)

	// Order routes
	router.POST("/orders", orderHandler.CreateOrder)
	router.GET("/orders", orderHandler.GetOrders)
	router.PUT("/orders", orderHandler.UpdateOrder)
	router.DELETE("/orders", orderHandler.DeleteOrder)

	return &http.Server{
		Addr:    ":9000",
		Handler: router,
	}
}
