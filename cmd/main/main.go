package main

import (
	"context"
	"net/http"
	"strconv"

	"ia-online-golang/internal/config"
	"ia-online-golang/internal/lib/logger"
	"ia-online-golang/internal/storage"
	"ia-online-golang/internal/utils"

	AuthService "ia-online-golang/internal/services/auth"
	BitrixService "ia-online-golang/internal/services/bitrix"
	EmailService "ia-online-golang/internal/services/email"
	LeadService "ia-online-golang/internal/services/lead"
	PasswordCodeService "ia-online-golang/internal/services/passwordcode"
	ReferralService "ia-online-golang/internal/services/referral"
	TokenService "ia-online-golang/internal/services/token"
	UserService "ia-online-golang/internal/services/user"

	AuthController "ia-online-golang/internal/http/controllers/auth"
	BitrixController "ia-online-golang/internal/http/controllers/bitrix"
	LeadController "ia-online-golang/internal/http/controllers/lead"
	UserController "ia-online-golang/internal/http/controllers/user"
	"ia-online-golang/internal/http/middleware"
	"ia-online-golang/internal/http/validator"
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

	bitrixService := BitrixService.New(log, cfg.BitrixConfig.IncomingWebhook)

	passwordCodeService := PasswordCodeService.New(log, storage)

	userService := UserService.New(log, storage)

	leadService := LeadService.New(log, storage, userService, storage, bitrixService, storage)

	referralService := ReferralService.New(log, storage)

	tokenService := TokenService.New(
		log,
		cfg.JWTConfig.Access.SecretKey,
		cfg.JWTConfig.Refresh.SecretKey,
		int64(cfg.JWTConfig.Access.Expiration.Seconds()),
		int64(cfg.JWTConfig.Refresh.Expiration.Seconds()),
		storage,
		userService,
		leadService,
		referralService,
	)

	authService := AuthService.New(log, cfg.HTTPServerConfig.Address, storage, storage, storage, storage, tokenService, emailService, userService, passwordCodeService)

	// Инициализация валидатора
	validator := validator.New()

	// Инициализация контроллеров
	log.Info("Initializing controllers...")
	authController := AuthController.New(log, validator, authService)
	userController := UserController.New(log, validator, userService)
	leadController := LeadController.New(log, validator, leadService)
	bitrixController := BitrixController.New(log, cfg.BitrixConfig.OutgoingWebhookAuth, leadService)

	// Создаём маршрутизатор
	mux := http.NewServeMux()

	// Открытые маршруты (не требуют авторизации)
	mux.HandleFunc("/", utils.HandleNotFound)
	mux.HandleFunc("/api/v1/auth/registration", authController.Registration)
	mux.HandleFunc("/api/v1/auth/activation/", authController.Activation)
	mux.HandleFunc("/api/v1/auth/login", authController.Login)
	mux.HandleFunc("/api/v1/auth/logout", authController.Logout)
	mux.HandleFunc("/api/v1/auth/refresh", authController.Refresh)
	mux.HandleFunc("/api/v1/auth/recover", authController.SendNewPassword)

	mux.HandleFunc("/api/v1/lead/edit", bitrixController.СhangingDeal)

	// Защищённые маршруты (нужен JWT-токен)
	protectedMux := http.NewServeMux()
	protectedMux.Handle("/api/v1/users", middleware.RoleMiddleware("manager")(http.HandlerFunc(userController.Users)))
	protectedMux.Handle("/api/v1/user/", middleware.RoleMiddleware("manager")(http.HandlerFunc(userController.User)))
	protectedMux.Handle("/api/v1/user/edit", middleware.RoleMiddleware("user")(http.HandlerFunc(userController.EditUser)))

	protectedMux.Handle("/api/v1/leads", middleware.RoleMiddleware("manager", "user")(http.HandlerFunc(leadController.Leads)))
	protectedMux.Handle("/api/v1/lead/save", middleware.RoleMiddleware("user")(http.HandlerFunc(leadController.SaveLead)))

	protectedMux.Handle("/api/v1/auth/new_password", middleware.RoleMiddleware("user")(http.HandlerFunc(authController.NewPassword)))

	// Оборачиваем защищённые маршруты в JWTMiddleware
	protectedRoutes := middleware.JWTMiddleware(context.Background(), tokenService)(protectedMux)

	// Основной серверный обработчик
	finalMux := http.NewServeMux()
	finalMux.Handle("/", mux) // Открытые маршруты
	finalMux.Handle("/api/v1/users", protectedRoutes)
	finalMux.Handle("/api/v1/user", protectedRoutes)
	finalMux.Handle("/api/v1/user/edit", protectedRoutes)

	finalMux.Handle("/api/v1/leads", protectedRoutes)
	finalMux.Handle("/api/v1/lead/save", protectedRoutes)

	finalMux.Handle("/api/v1/auth/new_password", protectedRoutes)

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
