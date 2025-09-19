<p align="center">
  <a href="https://scailo.com" target="_blank">
    <img src="https://pub-fbb2435be97c492d8ece0578844483ea.r2.dev/scailo-logo.png" alt="Scailo Logo" width="200">
  </a>
</p>

<h1 align="center">Scailo Official Go SDK</h1>

[![Go Reference](https://pkg.go.dev/badge/github.com/scailo/go-sdk.svg)](https://pkg.go.dev/github.com/scailo/go-sdk)

Welcome to the official Go SDK for the Scailo API. This repository contains a generated gRPC package that allows you to seamlessly integrate your Go applications with the full suite of Scailo services.

## About Scailo

Scailo is a powerful, modern ERP solution designed to be the foundation for your business needs. It provides a wide range of customizable business applications that cover everything from e-commerce, accounting, and CRM to order management, manufacturing, and human resources. With Scailo, you can streamline operations and unify your business processes on a single, scalable platform.

To learn more about what Scailo can do for your business, visit [scailo.com](https://scailo.com).

## Installation

To get started, use `go get` to add the package to your project:

```sh
go get github.com/scailo/go-sdk
```

## Getting Started & Usage
Interacting with the Scailo API is done through gRPC. The following examples will guide you through connecting to the server, authenticating, and making API calls.

#### 1. Authentication
First, you need to authenticate to get an auth_token. This token must be included in the metadata of all subsequent requests.

The example below shows how to create a gRPC client, call the LoginAsEmployeePrimary method, and retrieve the authentication token.

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/scailo/go-sdk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
    // The address of the Scailo gRPC server
    address  = "your-scailo-instance.com:443"
    username = "your-username"
    password = "your-password"
)

// login authenticates the user and returns the login response containing the auth token.
func login(ctx context.Context, conn *grpc.ClientConn) (*sdk.UserLoginResponse, error) {
    loginClient := sdk.NewLoginServiceClient(conn)
    
    resp, err := loginClient.LoginAsEmployeePrimary(ctx, &sdk.UserLoginRequest{
        Username:          username,
        PlainTextPassword: password,
    })

    if err != nil {
        return nil, fmt.Errorf("login failed: %w", err)
    }

    fmt.Println("Successfully logged in!")
    return resp, nil
}

func main() {
    // Establish a connection to the gRPC server.
    // For production, you should use secure credentials.
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

    // You can now use loginResp.AuthToken for subsequent API calls.
    fmt.Printf("Received Auth Token: %s\n", loginResp.AuthToken)
}
```

#### 2. Making Authenticated Requests
Once you have the auth_token, you must embed it into the context of your next requests using gRPC metadata.

Here is a complete example that logs in, creates a new context with the authentication token, and uses it to fetch a list of purchase orders.

```go
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
	// The address of the Scailo gRPC server
	address  = "your-scailo-instance.com:443"
	username = "your-username"
	password = "your-password"
)

// login authenticates the user.
func login(ctx context.Context, conn *grpc.ClientConn) (*sdk.UserLoginResponse, error) {
	loginClient := sdk.NewLoginServiceClient(conn)
	return loginClient.LoginAsEmployeePrimary(ctx, &sdk.UserLoginRequest{
		Username:          username,
		PlainTextPassword: password,
	})
}

// getPurchaseOrders fetches the latest active purchase order.
func getPurchaseOrders(ctx context.Context, conn *grpc.ClientConn) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	purchaseOrdersClient := sdk.NewPurchasesOrdersServiceClient(conn)
	resp, err := purchaseOrdersClient.Filter(ctx, &sdk.PurchasesOrdersServiceFilterReq{
		IsActive:  sdk.BOOL_FILTER_BOOL_FILTER_TRUE,
		Count:     1, // Get only one record
		SortOrder: sdk.SORT_ORDER_DESCENDING, // Get the latest one
	})
	if err != nil {
		log.Fatalf("could not get purchase orders: %v", err)
	}

    fmt.Println("Successfully fetched purchase orders:")
	fmt.Println(resp.List)
}

func main() {
	// 1. Establish a connection. Use secure credentials in production.
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()

	// 2. Login to get the authentication token.
	loginResp, err := login(ctx, conn)
	if err != nil {
		log.Fatalf("could not login: %v", err)
	}
	fmt.Printf("Login successful. AuthToken received.\n")

	// 3. Create metadata with the auth token.
	md := metadata.Pairs(
		"auth_token", loginResp.AuthToken,
	)
	
    // 4. Create a new context with the metadata attached.
	authedCtx := metadata.NewOutgoingContext(ctx, md)

	// 5. Use the new authenticated context for subsequent API calls.
	getPurchaseOrders(authedCtx, conn)
}
```

## API Use Cases

The Scailo API is extensive and allows you to build powerful integrations. Some common use cases include:

- E-commerce Integration: Sync orders, customer data, and inventory levels between Scailo and platforms like Shopify or WooCommerce.

- Automate Business Processes: Automatically transfer data from a CRM or Warehouse Management System (WMS) directly into the ERP.

- Financial Management: Connect Scailo with accounting systems to automate invoice generation and financial reporting.

- Custom Workflows: Build custom applications and workflows tailored to your specific business logic.

For more detailed information on what you can build, please see our [API](https://scailo.com/api) documentation.