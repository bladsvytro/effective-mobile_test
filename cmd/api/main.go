package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	_ "subscription-service/docs"
	"subscription-service/internal/config"
	"subscription-service/internal/database"
	"subscription-service/internal/handlers"
	"subscription-service/internal/logging"
	"subscription-service/internal/repository"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// customValidator адаптер для validator.v10, реализующий интерфейс binding.StructValidator
type customValidator struct {
	validate *validator.Validate
}

func (cv *customValidator) ValidateStruct(obj any) error {
	return cv.validate.Struct(obj)
}

func (cv *customValidator) Engine() any {
	return cv.validate
}

// @title Subscription Service API
// @version 1.0
// @description REST-сервис для агрегации данных об онлайн подписках пользователей
// @contact.name Effective Mobile
// @contact.url https://effective‑mobile.ru
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	err = logging.Init(cfg.Logging.Level, cfg.Logging.Encoding)
	if err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}
	defer logging.Sync()

	logger := logging.Get()
	logger.Info("Конфигурация загружена",
		zap.String("server_port", cfg.Server.Port),
		zap.String("db_host", cfg.Database.Host),
	)

	db, err := database.Connect(cfg.Database)
	if err != nil {
		logger.Fatal("Ошибка подключения к базе данных", zap.Error(err))
	}
	logger.Info("База данных подключена")

	repo := repository.NewSubscriptionRepository(db)
	handler := handlers.NewSubscriptionHandler(repo, logger)

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logging.GinLogger(logger))
	router.Use(cors.Default())

	// Устанавливаем validator.v10 как дефолтный валидатор для Gin
	validate := validator.New()
	binding.Validator = &customValidator{validate}

	v1 := router.Group("/api/v1")
	{
		v1.POST("/subscriptions", handler.Create)
		v1.GET("/subscriptions", handler.List)
		v1.GET("/subscriptions/:id", handler.GetByID)
		v1.PUT("/subscriptions/:id", handler.Update)
		v1.DELETE("/subscriptions/:id", handler.Delete)

		v1.GET("/analytics/total", handler.GetTotalPrice)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Корневой маршрут
	router.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, `<!DOCTYPE html>
<html lang="ru">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Subscription Service</title>
	<style>
		* {
			margin: 0;
			padding: 0;
			box-sizing: border-box;
		}
		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
			background: #f8f9fa;
			color: #333;
			line-height: 1.6;
			padding: 2rem;
			min-height: 100vh;
			display: flex;
			flex-direction: column;
			align-items: center;
			justify-content: center;
		}
		.container {
			max-width: 800px;
			width: 100%;
			background: white;
			border-radius: 12px;
			box-shadow: 0 4px 12px rgba(0,0,0,0.1);
			padding: 2rem;
			text-align: center;
		}
		h1 {
			font-size: 2.5rem;
			margin-bottom: 1rem;
			color: #2c3e50;
		}
		.subtitle {
			font-size: 1.2rem;
			margin-bottom: 2rem;
			color: #7f8c8d;
			max-width: 600px;
			margin-left: auto;
			margin-right: auto;
		}
		.links {
			display: flex;
			flex-wrap: wrap;
			gap: 1rem;
			justify-content: center;
			margin-bottom: 2rem;
		}
		.link {
			display: inline-block;
			padding: 0.8rem 1.5rem;
			background: #3498db;
			color: white;
			text-decoration: none;
			border-radius: 6px;
			font-weight: 600;
			transition: background 0.2s;
		}
		.link:hover {
			background: #2980b9;
		}
		.link.swagger {
			background: #27ae60;
		}
		.link.swagger:hover {
			background: #219653;
		}
		.link.github {
			background: #34495e;
		}
		.link.github:hover {
			background: #2c3e50;
		}
		.footer {
			margin-top: 2rem;
			font-size: 0.9rem;
			color: #95a5a6;
			border-top: 1px solid #ecf0f1;
			padding-top: 1rem;
		}
		.footer strong {
			color: #2c3e50;
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>Subscription Service</h1>
		<p class="subtitle">REST-сервис для агрегации данных об онлайн подписках пользователей. Реализовано в рамках тестового задания для Effective Mobile.</p>
		<div class="links">
			<a class="link swagger" href="/swagger/index.html">Swagger документация</a>
			<a class="link" href="/api/v1/subscriptions">Список подписок (JSON)</a>
			<a class="link github" href="https://github.com/bladsvytro/effective-mobile_test" target="_blank">Репозиторий GitHub</a>
			<a class="link" href="/api/v1/analytics/total">Аналитика стоимости</a>
		</div>
		<div class="footer">
			<strong>Сделано Владиславом Наумовым</strong> | <em>Тестовое задание Junior Golang Developer</em>
		</div>
	</div>
</body>
</html>`)
	})

	// Настройка HTTP сервера с graceful shutdown
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// Канал для ошибок сервера
	serverErr := make(chan error, 1)

	// Запуск сервера в горутине
	go func() {
		logger.Info("Сервер запускается", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Канал для сигналов ОС
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Ожидание сигнала или ошибки сервера
	select {
	case err := <-serverErr:
		logger.Fatal("Ошибка запуска сервера", zap.Error(err))
	case sig := <-quit:
		logger.Info("Получен сигнал завершения", zap.String("signal", sig.String()))
		logger.Info("Остановка сервера...")

		// Таймаут graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatal("Принудительное завершение сервера", zap.Error(err))
		}
		logger.Info("Сервер успешно остановлен")
	}
}