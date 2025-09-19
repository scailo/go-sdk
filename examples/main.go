package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/scailo/go-sdk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	// address of the gRPC server
	address  = ""
	username = ""
	password = ""
)

func login(ctx context.Context, conn *grpc.ClientConn) (*sdk.UserLoginResponse, error) {
	loginClient := sdk.NewLoginServiceClient(conn)
	resp, err := loginClient.LoginAsEmployeePrimary(ctx, &sdk.UserLoginRequest{
		Username:          username,
		PlainTextPassword: password,
	})
	return resp, err
}

func getPurchaseOrders(ctx context.Context, conn *grpc.ClientConn) {
	// Create a new context with the outgoing metadata
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	purchaseOrdersClient := sdk.NewPurchasesOrdersServiceClient(conn)
	resp, err := purchaseOrdersClient.Filter(ctx, &sdk.PurchasesOrdersServiceFilterReq{
		IsActive:  sdk.BOOL_FILTER_BOOL_FILTER_TRUE,
		Count:     1,
		SortOrder: sdk.SORT_ORDER_DESCENDING,
	})
	if err != nil {
		log.Fatalf("could not get purchase orders: %v", err)
	}
	fmt.Println(resp.List)
}

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()

	loginResp, err := login(ctx, conn)
	if err != nil {
		log.Fatalf("could not login: %v", err)
	}

	fmt.Println(loginResp)

	md := metadata.Pairs(
		"auth_token", loginResp.AuthToken,
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	getPurchaseOrders(ctx, conn)
}
