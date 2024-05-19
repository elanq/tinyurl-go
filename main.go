package main

import (
	"fmt"
	"os"
	"time"

	"github.com/elanq/tinyurl-go/connection"
	"github.com/elanq/tinyurl-go/handler"
	"github.com/elanq/tinyurl-go/repository"
	"github.com/elanq/tinyurl-go/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	servicePort := os.Getenv("SERVICE_PORT")
	redisConnection := connection.NewRedis()
	postgreConnection := connection.NewPostgre()

	urlRepository := repository.NewURL(postgreConnection.DB())
	cacheService := service.NewRedisCache(redisConnection)
	urlService := service.NewURL(urlRepository, cacheService)

	urlHandler := handler.NewURL(urlService)
	e := echo.New()
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 5 * time.Second,
	}))

	e.GET("/:url", urlHandler.GetByShortURL)
	e.POST("/", urlHandler.Create)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", servicePort)))
}
