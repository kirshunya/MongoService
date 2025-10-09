package middleware

import (
	"MongoService/metrics"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Запоминаем время начала запроса
		start := time.Now()

		// Продолжаем обработку запроса
		c.Next()

		// После обработки запроса собираем метрики
		duration := time.Since(start).Seconds()
		method := c.Request.Method
		endpoint := c.Request.URL.Path
		status := strconv.Itoa(c.Writer.Status())

		// Инкрементируем счетчик запросов
		metrics.HttpRequestsTotal.With(prometheus.Labels{
			"method":   method,
			"endpoint": endpoint,
			"status":   status,
		}).Inc()

		// Записываем время обработки
		metrics.HttpRequestDuration.With(prometheus.Labels{
			"method":   method,
			"endpoint": endpoint,
		}).Observe(duration)
	}
}
