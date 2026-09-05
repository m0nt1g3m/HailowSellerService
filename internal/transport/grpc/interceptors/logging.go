package interceptors

import (
	"HailowSellerService/pkg/logging"
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(logger *logging.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		statusCode := codes.OK

		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code()
			} else {
				statusCode = codes.Unknown
			}
			logger.Errorf("[gRPC] %s | %s | %v | err: %v", info.FullMethod, statusCode, duration, err)
		} else {
			logger.Infof("[gRPC] %s | %s | %v", info.FullMethod, statusCode, duration)
		}

		return resp, err
	}
}
