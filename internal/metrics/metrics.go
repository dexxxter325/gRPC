package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
)

var isHandlerRegistered bool //default-false

func InitMetrics(port string, logger *logrus.Logger) error {

	if !isHandlerRegistered { //не можем юзать 1 хэндлер неск.раз
		http.Handle("/metrics", promhttp.Handler()) //promhttp.Handler-сборщик дкфолтных метриков для prometheus
		isHandlerRegistered = true
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: nil,
	}
	go func() {
		logger.Infof("Metrics started on port:%s", port)
		if err := srv.ListenAndServe(); err != nil {
			logger.Errorf("Run metrics failed:%s", err)
		}
	}()

	return nil
}

/*type Metrics struct {
	RequestsTotal *prometheus.CounterVec   //кол-во запросов
	Latency       *prometheus.HistogramVec //задержка в обработке запросов
	ErrorsTotal   *prometheus.CounterVec   //запросы с ошибкой
}

func CreateMetrics(port string, logger *logrus.Logger) (Metrics, error) {
	var metrics Metrics

	metrics.RequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_requests_total",           //name of metric
		Help: "Total number of gRPC requests", //descp of metric
	}, []string{"method", "status"}, //метка метрики(для отслеживания количества запросов для каждого метода gRPC и каждого статуса ответа)
	)

	metrics.Latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "grpc_latency_seconds",
		Help: "Latency of gRPC requests in seconds",
	}, []string{"method"},
	)

	metrics.ErrorsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_errors_total",
		Help: "Total number of gRPC errors",
	}, []string{"method"},
	)
	//регистрируем метрики в prometheus
	prometheus.MustRegister(metrics.RequestsTotal)
	prometheus.MustRegister(metrics.Latency)
	prometheus.MustRegister(metrics.ErrorsTotal)
	recordMetrics(metrics)
	http.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: nil,
	}
	go func() {
		logger.Infof("Metrics started on port:%s", port)
		if err := srv.ListenAndServe(); err != nil {
			logger.Errorf("Run metrics failed:%s", err)
		}
	}()
	return metrics, nil
}
func recordMetrics(metrics Metrics) { //динамическое обновление метриков
	go func() {
		for {
			time.Sleep(2 * time.Second)                                     //каждые 2с будем обновлять метрики
			metrics.RequestsTotal.WithLabelValues("method", "status").Inc() // inc-увелич.значение метрики на 1 при каждом новом запросе
		}
	}()
	go func() {
		for {
			time.Sleep(2 * time.Second)
			metrics.Latency.WithLabelValues("method").Observe(0.1)
		}

	}()
	go func() {
		for {
			time.Sleep(2 * time.Second)
			metrics.ErrorsTotal.WithLabelValues("method").Inc()
		}

	}()
}
*/
