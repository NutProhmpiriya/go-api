package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"socialnetwork/docs"
	"socialnetwork/internal/delivery/http/handler"
	"socialnetwork/internal/repository/postgres"
	"socialnetwork/internal/usecase"
	"socialnetwork/pkg/database"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Social Network API
// @version         1.0
// @description     A Social Network API with authentication and post management.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Swagger initialization
	docs.SwaggerInfo.BasePath = "/api"

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

	// JWT secret key
	secretKey := "your-secret-key"

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, secretKey)
	userUseCase := usecase.NewUserUseCase(userRepo)
	postUseCase := usecase.NewPostUseCase(postRepo, userRepo)

	// Initialize HTTP handlers
	authHandler := handler.NewAuthHandler(authUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	postHandler := handler.NewPostHandler(postUseCase, secretKey)

	// Initialize Gin router
	router := gin.Default()

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve static files
	router.Static("/static", "./web/static")
	router.StaticFile("/", "./web/static/index.html")

	// API routes
	api := router.Group("/api")
	{
		authHandler.Register(api)
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

		// Shutdown signal with grace period of 30 seconds
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
