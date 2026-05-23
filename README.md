# Revoluchat Go SDK 🚀

**Version**: `v1.2.0-alpha` (Group Support & Advance Tier)

Official Go SDK for **Revoluchat**, an enterprise-grade, multi-tenant real-time chat platform. This SDK provides a seamless way to integrate your existing user database with Revoluchat using a highly secure "pointing" pattern and OpenID Connect (OIDC) compliant token generation.

[![Go Reference](https://pkg.go.dev/badge/github.com/oririfai/revoluchat-go-sdk.svg)](https://pkg.go.dev/github.com/oririfai/revoluchat-go-sdk)

## ✨ Features

- **Tiered Architecture Support**: Choose between `Normal` (local Elixir storage) and `Advance` (centralized Go storage) tiers.
- **OIDC Complaint Authentication**: Generates RS256 JWTs with dynamically computed Key IDs (KID).
- **Secure JWKS & gRPC**: Built-in `X-Server-Key` validation for all inter-service communication.
- **Group Chat Support**: Full management of groups, members, and group security in Advance Tier.
- **Easy "Pointing" Pattern**: Map your internal data to Revoluchat using simple Go interfaces or functions.
- **High Performance**: Built on top of `google.golang.org/grpc` for efficient, reliable inter-service communication.

## 📦 Installation & Setup

To add the SDK to your project:

```bash
go get github.com/oririfai/revoluchat-go-sdk@latest
```

### ⚡ Quick Start (Recommended)

Use our interactive CLI to bootstrap your Revoluchat configuration in seconds:

```bash
go run github.com/oririfai/revoluchat-go-sdk/cmd/revoluchat-init@latest
```

This will ask for your **Tier** and **Server Key**, then generate a `revoluchat_setup.go` file with all the necessary boilerplate.

## 🚀 Manual Integration Guide

Integrating with Revoluchat requires two main components:

1. **JWT Manager**: For signing authentication tokens and serving the JWKS to Revoluchat servers.
2. **gRPC Server**: For Revoluchat to fetch high-speed profile data (name, avatar, KYC status) and manage chat data in Advance Tier.

### 1. Implementation in Your Go Backend

A simple initialization example in your `main.go` file:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/oririfai/revoluchat-go-sdk/revoluchat"
)

func main() {
	// 1. Initialize JWT Manager
	serverKey := os.Getenv("REVOLUCHAT_SERVER_KEY")
	jwtManager, err := revoluchat.NewJWTManager("keys/app.rsa", "keys/app.rsa.pub", serverKey)
	if err != nil {
		log.Fatalf("Failed to initialize JWT: %v", err)
	}

	// 2. Start the SDK gRPC server
	go func() {
		err := revoluchat.Start(revoluchat.Config{
			Tier:     revoluchat.TierAdvance, // Choose Normal or Advance
			GRPCPort: 50051,
			ServerKey: serverKey,
			UserProvider: func(ctx context.Context, id uint64) (*revoluchat.User, error) {
				return &revoluchat.User{
					ID:        id,
					Name:      "User Full Name",
					Phone:     "08123456789",
					AvatarURL: "https://your-cdn.com/path/to/photo.jpg",
					Status:    "active",
					IsKYC:     true,
				}, nil
			},
			// Required for TierAdvance
			MessageProvider:      &MyMessageService{},
			ConversationProvider: &MyConversationService{},
			GroupProvider:        &MyGroupService{}, // Your implementation
		})

		if err != nil {
			log.Fatalf("Failed to start Revoluchat gRPC server: %v", err)
		}
	}()

	// 3. HTTP Server Setup for the JWKS endpoint
	r := gin.Default()
	r.GET("/jwks", gin.WrapF(jwtManager.JWKSHandler))
	r.Run(":8089")
}
```

### 2. Configuration Integration to the Revoluchat Server

Set the `TIER_TYPE` and point to your Go backend:

```bash
# Set the tier (must match SDK config)
TIER_TYPE=advance

# Point to your gRPC Port
CHAT_SERVICE_GRPC_ENDPOINT=your-go-backend-host:50051
USER_SERVICE_GRPC_ENDPOINT=your-go-backend-host:50051
```

## 🛠️ Data Structure

The SDK uses a simplified `revoluchat.User` struct for data mapping over gRPC:

| Field       | Type     | Description                                          |
| :---------- | :------- | :--------------------------------------------------- |
| `ID`        | `string` | Unique user ID in your application's database        |
| `UUID`      | `string` | Optional: UUID for additional external reference     |
| `Name`      | `string` | Display name of the user inside the chat             |
| `Phone`     | `string` | Phone number connected to the user's account         |
| `Status`    | `string` | Account status (e.g., "active", "pending", "banned") |
| `IsKYC`     | `bool`   | Indicates whether the user's identity is verified    |
| `AvatarURL` | `string` | Full URL to the user's profile picture               |

### Provider Format

When using `TierAdvance`, your structs must implement the functions wrapped by the `revoluchat` interfaces. For example, the implementation for `MessageProvider` requires you to fulfill: `InsertMessage`, `ListMessages`, `MarkRead`, `MarkDelivered`, `DeleteMessage`, `BulkDeleteMessages` according to the protobuf request and response types from `github.com/oririfai/revoluchat-go-sdk/proto/chat_v1`.

### Authentication / Generate Token Endpoint

You can call the e.g `jwtManager.GenerateToken(userID string, appID string) (string, error)` function when a user successfully logs into your system. This token can be returned to the client-side (mobile/web app) to be used for logging into the Revoluchat WebSocket server.

## 📄 License

MIT © [Revoluchat Team](https://revolu.id)
