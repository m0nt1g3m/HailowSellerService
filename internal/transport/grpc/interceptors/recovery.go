package interceptors

import (
	"HailowSellerService/pkg/logging"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoveryInterceptor(logger *logging.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("Panic recovered in gRPC interceptor: %v", r)
				err = status.Errorf(codes.Internal, "Internal server error: %v", fmt.Sprint(r))
			}
		}()

		return handler(ctx, req)
	}
}
