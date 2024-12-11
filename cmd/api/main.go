package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"socialnetwork/internal/delivery/http/handler"
	"socialnetwork/internal/repository/postgres"
	"socialnetwork/internal/usecase"
	"socialnetwork/pkg/database"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connection
	db, err := database.NewPostgresConnection(&database.PostgresConfig{
		Host:     "localhost",
		Port:     5433,
		User:     "postgres",
		Password: "postgres",
		DBName:   "socialnetwork",
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	postRepo := postgres.NewPostRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)
	postUseCase := usecase.NewPostUseCase(postRepo, userRepo)

	// Initialize HTTP handlers
	userHandler := handler.NewUserHandler(userUseCase)
	postHandler := handler.NewPostHandler(postUseCase)

	// Initialize Gin router
	router := gin.Default()

	// Serve static files
	router.Static("/static", "./web/static")
	router.StaticFile("/", "./web/static/index.html")

	// API routes
	api := router.Group("/api")
	{
		userHandler.Register(api)
		postHandler.Register(api)
	}

	// Create server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	fmt.Println("Server is running on http://localhost:8080")
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
