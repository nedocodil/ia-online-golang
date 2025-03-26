package main

import (
	"context"
	"net/http"
	"strconv"

	"ia-online-golang/internal/config"
	"ia-online-golang/internal/lib/logger"
	"ia-online-golang/internal/storage"

	AuthService "ia-online-golang/internal/services/auth"
	EmailService "ia-online-golang/internal/services/email"
	TokenService "ia-online-golang/internal/services/token"
	UserService "ia-online-golang/internal/services/user"

	AuthController "ia-online-golang/internal/http/controllers/auth"
	UserController "ia-online-golang/internal/http/controllers/user"
	"ia-online-golang/internal/http/middleware" // Подключаем JWTMiddleware
	"ia-online-golang/internal/http/utils"
	"ia-online-golang/internal/http/validators/reg"

	"github.com/go-playground/validator/v10"
)

func main() {
	// Чтение конфига
	cfg := config.MustLoad()

	// Логирование
	log := logger.SetupLogger(cfg.Env)
	log.Info("Starting...")

	// Подключение к БД
	log.Info("Connecting database...")
	storage, err := storage.NewStorage(cfg.StorageConfig.Path)
	if err != nil {
		log.Fatal("Error connecting to storage:", err)
	}

	// Инициализация сервисов
	log.Info("Initializing services...")
	emailService := EmailService.New(cfg.EmailConfig.SMTP.Host,
		strconv.Itoa(cfg.EmailConfig.SMTP.PortSSL),
		cfg.EmailConfig.SMTP.Username,
		cfg.EmailConfig.SMTP.Password)

	tokenService := TokenService.New(
		cfg.JWTConfig.Access.SecretKey,
		cfg.JWTConfig.Refresh.SecretKey,
		int64(cfg.JWTConfig.Access.Expiration.Seconds()),
		int64(cfg.JWTConfig.Refresh.Expiration.Seconds()),
		storage,
	)
	userService := UserService.New(log, storage)

	authService := AuthService.New(log, cfg.HTTPServerConfig.Address, storage, storage, storage, storage, tokenService, emailService)

	// Инициализация валидатора
	validator := validator.New()
	validator.RegisterValidation("complexpassword", reg.PasswordValidation)

	// Инициализация контроллеров
	log.Info("Initializing controllers...")
	authController := AuthController.New(authService, log, validator)
	userController := UserController.New(log, validator, userService)

	// Создаём маршрутизатор
	mux := http.NewServeMux()

	// Открытые маршруты (не требуют авторизации)
	mux.HandleFunc("/", utils.HandleNotFound)
	mux.HandleFunc("/api/v1/auth/registration", authController.Registration)
	mux.HandleFunc("/api/v1/auth/activation/", authController.Activation)
	mux.HandleFunc("/api/v1/auth/login", authController.Login)
	mux.HandleFunc("/api/v1/auth/logout", authController.Logout)
	mux.HandleFunc("/api/v1/auth/refresh", authController.Refresh)

	// Защищённые маршруты (нужен JWT-токен)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/v1/users", userController.Users)

	// Оборачиваем защищённые маршруты в JWTMiddleware
	protectedRoutes := middleware.JWTMiddleware(context.Background(), tokenService)(protectedMux)

	// Основной серверный обработчик
	finalMux := http.NewServeMux()
	finalMux.Handle("/", mux)                         // Открытые маршруты
	finalMux.Handle("/api/v1/users", protectedRoutes) // Защищённые маршруты

	srv := &http.Server{
		Addr:         cfg.HTTPServerConfig.Address,
		Handler:      finalMux,
		ReadTimeout:  cfg.HTTPServerConfig.ReadTimeout,
		WriteTimeout: cfg.HTTPServerConfig.WriteTimeout,
		IdleTimeout:  cfg.HTTPServerConfig.IdleTimeout,
	}

	// Запускаем сервер
	log.Info("Server is running on " + cfg.HTTPServerConfig.Address)
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}
