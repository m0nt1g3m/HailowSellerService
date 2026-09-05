package interceptors

import (
	"HailowSellerService/internal/domain"
	"HailowSellerService/internal/transport/grpc/response/errorcode"
	"HailowSellerService/pkg/logging"
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ContextKey string

const (
	ContextSellerIDKey     ContextKey = "seller_id"
	ContextRefreshTokenKey ContextKey = "refresh_token"
)

func AuthInterceptor(logger *logging.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if strings.HasSuffix(info.FullMethod, "SignIn") || strings.HasSuffix(info.FullMethod, "SignUp") || strings.HasSuffix(info.FullMethod, "RefreshTokens") {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errorcode.ToStatus(domain.ErrMDNotFound)
		}

		var sellerID string
		if sellerIDs := md.Get("x-seller-id"); len(sellerIDs) > 0 && sellerIDs[0] != "" {
			sellerID = sellerIDs[0]
		} else if sellerIDs := md.Get("seller-id"); len(sellerIDs) > 0 && sellerIDs[0] != "" {
			sellerID = sellerIDs[0]
		}

		if sellerID == "" {
			return nil, errorcode.ToStatus(domain.ErrSellerIDNotFoundMD)
		}

		ctx = context.WithValue(ctx, ContextSellerIDKey, sellerID)

		return handler(ctx, req)
	}
}

func RefreshTokenInterceptor(logger *logging.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if strings.HasSuffix(info.FullMethod, "RefreshTokens") || strings.HasSuffix(info.FullMethod, "Logout") {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, errorcode.ToStatus(domain.ErrRefreshTokenMD)
			}

			tokens := md.Get("refresh-token")
			if len(tokens) == 0 || tokens[0] == "" {
				return nil, errorcode.ToStatus(domain.ErrRefreshTokenMD)
			}

			ctx = context.WithValue(ctx, ContextRefreshTokenKey, tokens[0])
		}
		return handler(ctx, req)
	}
}
