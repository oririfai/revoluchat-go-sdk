# Revoluchat Go SDK 🚀

Official Go SDK for **Revoluchat**, an enterprise-grade, multi-tenant real-time chat platform. This SDK provides a seamless way to integrate your existing user database with Revoluchat using a simple "pointing" pattern, handling all gRPC boilerplate under the hood.

[![Go Reference](https://pkg.go.dev/badge/github.com/revolu/revoluchat-go-sdk.svg)](https://pkg.go.dev/github.com/revolu/revoluchat-go-sdk)
[![License](https://img.shields.io/github/license/revolu/revoluchat-go-sdk)](LICENSE)

## ✨ Features

- **Seamless Integration**: Effortlessly connect your user service without worrying about gRPC, Protobuf, or complex interfaces.
- **Easy "Pointing" Pattern**: Simply map your internal user data to the Revoluchat user structure using a functional provider.
- **High Performance**: Built on top of `google.golang.org/grpc` for efficient, reliable communication.
- **Enterprise Ready**: Designed to work flawlessly with Revoluchat's multi-tenant architecture.

## 📦 Installation

To add the SDK to your project:

```bash
go get github.com/oririfai/revoluchat-go-sdk
```

## 🚀 Quick Start (Integration Guide)

Integrating with Revoluchat is straightforward. You only need to call `revoluchat.Start` and provide a `Provider` function that maps your user data to the Revoluchat structure.

### 1. Implementation in Your Go Backend

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/oririfai/revoluchat-go-sdk"
)

func main() {
	// Start the SDK gRPC server in a goroutine to avoid blocking your main application
	go func() {
		err := revoluchat.Start(revoluchat.Config{
			GRPCPort: 50051, // The gRPC port that Revoluchat will connect to
			Provider: func(ctx context.Context, id uint64) (*revoluchat.User, error) {
				// 1. Fetch data from your own database (GORM, SQL, etc.)
				// Example: userDB := myrepo.FindByID(id)

				// 2. "Point" values from your data to the Revoluchat struct
				return &revoluchat.User{
					ID:        id,
					Name:      "User Full Name", // Map to your internal name field
					Phone:     "08123456789",    // Map to your internal phone field
					AvatarURL: "https://your-cdn.com/path/to/photo.jpg",
					Status:    "active",
					IsKYC:     true,
				}, nil
			},
		})
		if err != nil {
			log.Fatalf("Failed to start Revoluchat SDK: %v", err)
		}
	}()

	// Continue with your main application logic (e.g., Gin/Echo HTTP server)
	fmt.Println("Main application server is running...")
	select {} // Wait indefinitely
}
```

### 2. Configuration in Revoluchat

Ensure that the environment variables in your Revoluchat instance are pointing to the gRPC port configured above:

```bash
USER_SERVICE_GRPC_ENDPOINT=your-go-backend-host:50051
```

## 🛠️ Data Structure

The SDK uses a simplified `revoluchat.User` struct for data mapping:

| Field       | Type     | Description                             |
| :---------- | :------- | :-------------------------------------- |
| `ID`        | `uint64` | Unique User ID in your system           |
| `UUID`      | `string` | Optional UUID for external references   |
| `Name`      | `string` | User's display name in the chat         |
| `Phone`     | `string` | User's phone number                     |
| `Status`    | `string` | User status (e.g., "active", "pending") |
| `IsKYC`     | `bool`   | Identity verification status            |
| `AvatarURL` | `string` | Full URL to the user's profile picture  |

## 📄 License

MIT © [Revoluchat Team](https://revolu.id)
