package rpc

import (
	"context"
	"local/panda-killer/pkg/domain/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func SessionInterceptor(authUsecase usecase.AuthUsecase) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if len(md["authorization"]) > 0 {
				ctxWithClaims, _ := authUsecase.AddClaimsToContext(ctx, md["authorization"][0])
				return handler(ctxWithClaims, req)
			}
		}
		return handler(ctx, req)
	}
}
