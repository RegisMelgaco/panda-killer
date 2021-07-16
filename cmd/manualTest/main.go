package main

import (
	"context"
	"fmt"
	"local/panda-killer/pkg/gateway/rpc/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(fmt.Sprintf("Failed to open connection: %v", err))
	}

	client := gen.NewPandaKillerClient(conn)

	resp, err := client.CreateAccount(context.Background(), &gen.CreateAccountRequest{
		Name:     "Valeria",
		Cpf:      "847384923301",
		Password: "secret ;)",
		Balance:  50,
	})

	reqStatus, _ := status.FromError(err)
	fmt.Println(reqStatus.Code(), reqStatus.Message())

	fmt.Println(resp)
}
