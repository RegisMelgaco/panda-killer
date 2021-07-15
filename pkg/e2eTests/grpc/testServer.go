package e2etest

import (
	"context"
	"local/panda-killer/pkg/gateway/rpc"
	"local/panda-killer/pkg/gateway/rpc/gen"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func getBufDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
}

func StartServerAndGetClient(ctx context.Context, api *rpc.Api) (gen.PandaKillerClient, net.Conn, error) {
	lis := bufconn.Listen(1024 * 1024)

	s := rpc.NewServer(api)
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	conn, err := lis.Dial()
	if err != nil {
		return nil, nil, err
	}

	grpcConn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(getBufDialer(lis)), grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	client := gen.NewPandaKillerClient(grpcConn)
	return client, conn, nil
}
