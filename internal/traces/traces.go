package traces

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func InitTraces(url, serviceName string) (func(), error) {
	// Create Jaeger exporter-отправляет собранные opentelemtry трейсы в Jaeger UI
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url))) //send to Jaeger Collector-место хранения трейсов./api/trace указывает ,что прийдут трейсы
	if err != nil {
		return nil, fmt.Errorf("create new jaeger failed:%s", err)
	}
	// Create trace provider-отправляет трейсы в jaeger
	traceProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter), //отправляет данные экспортеру пачками(не по отдельности), повышая производительность
		tracesdk.WithResource(resource.NewWithAttributes( //атрибуты ,которые описывают мой сервис
			semconv.SchemaURL, //схема для дефолтных данных,которые будут включены в трейсы(HTTP метод,url,Статус HTTP,имя сервиса и тд.).Спаны их потом заполнят
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	// Set global trace provider-активация моего провайдера
	otel.SetTracerProvider(traceProvider)
	shutdownFunc := func() {
		if err = traceProvider.Shutdown(context.Background()); err != nil {
			logrus.Errorf("Error shutting down tracer provider: %s", err)
		}
	}

	return shutdownFunc, nil
}
