package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/batcher"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/grpc"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue/algorithms"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/registry"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/routes"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/selector"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/worker"
)

// CLI
var (
	addresses []string // nodes addrs
	queueType string
	sel       string
	loadStrat string

	// Batch
	batSize    int
	batTimeout int

	// Workers
	numWorkers int
)

// Global var
var regis = registry.NewRegistry()
var s selector.Selector
var q queue.Queue
var strat worker.LoadBalancingStrategy

var rootCmd = &cobra.Command{
	Use:     "load-manager",
	Short:   "Load Manager CLI for Distributed System",
	Long:    "A Load manager that distrubtes requests across multiple backends",
	Example: "load-manager --a host1:5000 --a host2:5000 --q FCFS --s RR --l M",
	PreRunE: preRunE,
	RunE:    runE,
}

func preRunE(cmd *cobra.Command, args []string) error {
	// Check for algos
	switch queueType {
	case "FCFS":
		q = algorithms.NewFCFSQueue()
	default:
		return fmt.Errorf("invalid queue type %s. Must be: FCFS", queueType)
	}

	// Check for load strat
	switch loadStrat {
	case "M":
		strat = worker.Mixed
	case "PR":
		strat = worker.PerResource
	case "PO":
		strat = worker.PerOperation
	case "PRO":
		strat = worker.PerResourceAndOperation
	default:
		return fmt.Errorf("invalid load strat %s. Must be: M, PR, PO, PRO", loadStrat)
	}

	// Check for selector
	switch sel {
	case "RR":
		s = selector.NewRR()
	default:
		return fmt.Errorf("invalid selector %s. Must be: RR", sel)
	}

	return nil
}

func runE(cmd *cobra.Command, args []string) error {
	// Add the addresses to registry
	err := parseAddrs(addresses)
	if err != nil {
		return err
	}

	// Health check
	go regis.HealthCheckLoop()

	// Batcher
	clients := make(map[string]*grpc.BackendClient)
	bat := batcher.NewBatcher(q, batSize, time.Duration(batTimeout)*time.Millisecond)

	// Worker
	wrk := worker.NewWorker(q, regis, s, clients, numWorkers, strat)

	// Router
	router := gin.Default()
	balancer := router.Group("balancer")

	// Users
	balancer.POST("/users", routes.CreateUser(bat))
	balancer.GET("/users", routes.GetUser(bat))
	balancer.PUT("/users", routes.UpdateUser(bat))
	balancer.DELETE("/users", routes.DeleteUser(bat))

	// Product
	balancer.POST("/products", routes.CreateProduct(bat))
	balancer.GET("/products", routes.GetProduct(bat))
	balancer.PUT("/products", routes.UpdateProduct(bat))
	balancer.DELETE("/products", routes.DeleteProduct(bat))

	// Order
	balancer.POST("/orders", routes.CreateOrder(bat))
	balancer.GET("/orders", routes.GetOrder(bat))
	balancer.PUT("/orders", routes.UpdateOrder(bat))
	balancer.DELETE("/orders", routes.DeleteOrder(bat))

	port := "8000"
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("Load manager listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	wrk.Stop()
	bat.Stop()

	return nil
}

func parseAddrs(addrs []string) error {
	for _, addr := range addrs {
		parts := strings.Split(addr, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid address format %s. Expected host:port", addr)
		}

		host := parts[0]
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}

		if port < 1 || port > 65535 {
			return fmt.Errorf("port out of range %s", addr)
		}
		regis.Add(host, port)
	}

	return nil
}

func init() {
	// []str
	rootCmd.Flags().StringSliceVarP(&addresses, "address", "a", []string{}, "Server addresses")

	// Str
	rootCmd.Flags().StringVarP(&queueType, "queue", "q", "", "Queue algorithms: FCFS\n")
	rootCmd.Flags().StringVarP(&loadStrat, "load", "l", "", "Load strategy: M\nPR\nPO\nPRO")
	rootCmd.Flags().StringVarP(&sel, "selector", "s", "", "Selector: RR\n")

	// Int
	rootCmd.Flags().IntVarP(&batSize, "batchsize", "b", 100, "Batch Size")
	rootCmd.Flags().IntVarP(&batTimeout, "batchtimeout", "t", 2, "Batch Timeout")
	rootCmd.Flags().IntVarP(&numWorkers, "workers", "w", 4, "Worker size")

	// Required
	err := rootCmd.MarkFlagRequired("address")
	if err != nil {
		log.Fatal(err)
	}
	err = rootCmd.MarkFlagRequired("queue")
	if err != nil {
		log.Fatal(err)
	}
	err = rootCmd.MarkFlagRequired("selector")
	if err != nil {
		log.Fatal(err)
	}
	err = rootCmd.MarkFlagRequired("load")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
