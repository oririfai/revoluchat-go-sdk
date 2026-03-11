# Revoluchat Go SDK 🚀

**Version**: `v1.1.0-alpha` (OIDC & Secure JWT Support)

Official Go SDK for **Revoluchat**, an enterprise-grade, multi-tenant real-time chat platform. This SDK provides a seamless way to integrate your existing user database with Revoluchat using a highly secure "pointing" pattern and OpenID Connect (OIDC) compliant token generation.

[![Go Reference](https://pkg.go.dev/badge/github.com/oririfai/revoluchat-go-sdk.svg)](https://pkg.go.dev/github.com/oririfai/revoluchat-go-sdk)

## ✨ Features

- **OIDC Complaint Authentication**: Generates RS256 JWTs with dynamically computed Key IDs (KID) based on your public keys.
- **Secure JWKS Endpoint**: Provides a built-in handler to serve your JSON Web Key Set securely, protected by a shared `Server Key` (header `X-Server-Key` or query `server_key`).
- **Seamless Integration**: Effortlessly connect your user service without worrying about gRPC boilerplate.
- **Easy "Pointing" Pattern**: Simply map your internal user data to the Revoluchat user structure using a functional provider.
- **High Performance**: Built on top of `google.golang.org/grpc` for efficient, reliable inter-service communication.

## 📦 Installation

To add the SDK to your project:

```bash
go get github.com/oririfai/revoluchat-go-sdk/revoluchat@latest
```

## 🚀 Quick Start (Integration Guide)

Integrating with Revoluchat requires two main components:
1. **JWT Manager**: For signing authentication tokens and serving the JWKS to Revoluchat servers.
2. **gRPC Server**: For Revoluchat to fetch high-speed profile data (name, avatar, KYC status) directly from your database.

### 1. Implementation in Your Go Backend

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/oririfai/revoluchat-go-sdk/revoluchat"
)

func main() {
	// 1. Initialize JWT Manager
	// Requires standard RSA private/public keys and a shared secret "Server Key"
	serverKey := os.Getenv("REVOLUCHAT_SERVER_KEY") // e.g. "my-super-secret-server-key"
	jwtManager, err := revoluchat.NewJWTManager("keys/app.rsa", "keys/app.rsa.pub", serverKey)
	if err != nil {
		log.Fatalf("Failed to initialize JWT: %v", err)
	}

	// 2. Start the SDK gRPC server in a goroutine
	go func() {
		err := revoluchat.Start(revoluchat.Config{
			GRPCPort: 50051, 
			Provider: func(ctx context.Context, id uint64) (*revoluchat.User, error) {
				// Fetch data from your own database (GORM, SQL, etc.)
				// Example: userDB := myrepo.FindByID(id)

				// "Point" values from your data to the Revoluchat struct
				return &revoluchat.User{
					ID:        id,
					Name:      "User Full Name", 
					Phone:     "08123456789",    
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

	// 3. HTTP Server Setup for Authentication & JWKS
	r := gin.Default()

	// Serve the secure JWKS endpoint natively via the Manager
	r.GET("/jwks", gin.WrapF(jwtManager.JWKSHandler))

	// Example Login Route
	r.POST("/login", func(c *gin.Context) {
		// Verify your own credentials (password, OTP, PIN) 
		// ...
		userID := "123" 
		tenantAppID := "revolu-corp"

		// Generate OIDC-compliant Token using SDK
		token, err := jwtManager.GenerateToken(userID, tenantAppID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(200, gin.H{"token": token})
	})

	fmt.Println("Backend server is running on :8089...")
	r.Run(":8089")
}
```

### 2. Configuration in Revoluchat

Ensure that the environment variables in your Revoluchat instance properly point to your new endpoints and use the exact `Server Key` defined above.

```bash
# Point to your gRPC Port
USER_SERVICE_GRPC_ENDPOINT=your-go-backend-host:50051

# Point to your JWKS Endpoint (Phoenix will automatically append ?server_key=xxx if using the Admin Dashboard dynamic keys)
JWKS_URL=http://your-go-backend-host:8089/jwks
```

*(Note: In a production environment, you should manage the dynamic `Server Key` centrally through the **Server Keys** menu within your Revoluchat Admin Dashboard, rather than hardcoding it).*

## 🛠️ Data Structure

The SDK uses a simplified `revoluchat.User` struct for data mapping over gRPC:

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
