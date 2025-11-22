package middleware

import (
	"context"
	"log"
	"time"

	"authService/internal/monitoring"
	"authService/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MetricsInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		method := utils.GetMethodName(info.FullMethod)

		log.Printf("Processing request: %s.%s", serviceName, method)

		monitoring.RequestsInFlight.WithLabelValues(serviceName, method).Inc()
		defer monitoring.RequestsInFlight.WithLabelValues(serviceName, method).Dec()

		start := time.Now()

		resp, err := handler(ctx, req)

		duration := float64(time.Since(start).Microseconds())

		monitoring.RequestDuration.WithLabelValues(serviceName, method).Observe(duration)

		log.Printf("Method %s took %.3f ms", method, duration)

		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			} else {
				code = codes.Unknown
			}
			log.Printf("Method %s failed with code: %s", method, code.String())
		} else {
			log.Printf("Method %s completed successfully", method)
		}

		monitoring.RequestsTotal.WithLabelValues(serviceName, method, code.String()).Inc()

		return resp, err
	}
}
