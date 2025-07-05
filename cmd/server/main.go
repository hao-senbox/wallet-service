package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wallet-service/config"
	"wallet-service/internal/exchange"
	"wallet-service/internal/user"
	"wallet-service/internal/wallet"
	"wallet-service/pkg/consul"
	"wallet-service/pkg/zap"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.LoadConfig()

	// Initialize logger
	logger, err := zap.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	consulConn := consul.NewConsulConn(logger, cfg)
	consulClient := consulConn.Connect()
	defer consulConn.Deregister()

	// Connect to MongoDB
	mongoClient, err := connectToMongoDB(cfg.MongoURI)
	if err != nil {
		logger.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			logger.Fatal(err)
		}
	}()

	// Set up router with Gin
	router := gin.Default()
	userService := user.NewUserService(consulClient)

	exchangeRateCollection := mongoClient.Database(cfg.MongoDB).Collection("exchange_rates")
	exchangeRepository := exchange.NewExchangeRepository(exchangeRateCollection)
	exchangeService := exchange.NewExchangeService(exchangeRepository)
	exchangeHandler := exchange.NewExchangeHandler(exchangeService)
	
	walletCollection := mongoClient.Database(cfg.MongoDB).Collection("wallets")
	transactionCollection := mongoClient.Database(cfg.MongoDB).Collection("transactions")
	walletRepository := wallet.NewWalletRepository(walletCollection, transactionCollection)
	walletService := wallet.NewWalletService(walletRepository, exchangeService, userService)
	walletHandler := wallet.NewWalletHandler(walletService)
	wallet.RegisterRoutes(router, walletHandler)
	exchange.RegisterRoutes(router, exchangeHandler)

	// Initialize HTTP server
	server := &http.Server{	
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Run server in a separate goroutine
	go func() {
		logger.Infof("Server running on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Error starting server: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Error shutting down server: %v", err)
	}
	logger.Info("Server stopped")
}

func connectToMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("Failed to connect to MongoDB")
		return nil, err
	}

	// Check connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("Failed to ping to MongoDB")
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}
