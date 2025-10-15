package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"MongoService/controllers"
	"MongoService/middleware"
	"MongoService/models"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type LogEntry struct {
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp string `json:"@timestamp"`
}

var esClient *elasticsearch.Client

func init() {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://elasticsearch:9200"},
	}
	var err error
	esClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Ошибка создания клиента Elasticsearch: %v", err)
	}
}

func logToElasticsearch(level, message string) {
	logEntry := LogEntry{
		Level:     level,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	body, _ := json.Marshal(logEntry)

	req := esapi.IndexRequest{
		Index:   "go-app-logs",
		Body:    bytes.NewReader(body),
		Refresh: "true",
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		log.Printf("Ошибка отправки лога в Elasticsearch: %v", err)
		return
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Printf("Ошибка ответа Elasticsearch: %s", res.String())
	}
}

func main() {
	r := gin.Default()

	// Логирование при старте
	logToElasticsearch("info", "Приложение запущено")

	r.Use(middleware.PrometheusMiddleware())
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Отправляем логи Gin в Elasticsearch
		logMessage := fmt.Sprintf("%s - %s %s %s %d %s",
			param.TimeStamp.Format(time.RFC3339),
			param.ClientIP,
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
		logToElasticsearch("info", logMessage)
		return ""
	}))

	r.GET("/ping", func(c *gin.Context) {
		logToElasticsearch("info", "Получен запрос на /ping")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	models.ConnectMongo()

	userGroup := r.Group("/users")
	{
		userGroup.POST("/create", controllers.InsertUser)
		userGroup.POST("/create_all", controllers.InsertUsers)
		userGroup.PUT("/update/:id", controllers.UpdateUser)
		userGroup.DELETE("/delete/:id", controllers.DeleteUser)
		userGroup.GET("/:id", controllers.FindUserById)
		userGroup.GET("/all", controllers.ListAllUsers)
		userGroup.DELETE("/delete/all", controllers.DeleteAll)
	}

	metricsGroup := r.Group("/metrics")
	{
		metricsGroup.GET("/", gin.WrapH(promhttp.Handler()))
	}

	r.Run(":8081")
}
