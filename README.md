# Revoluchat Go SDK 🚀

**Version**: `v1.2.4-alpha` (Group Support, Call, Attachment, & Advance Tier)

Official Go SDK for **Revoluchat**, an enterprise-grade, multi-tenant real-time chat platform. This SDK provides a seamless way to integrate your existing user database with Revoluchat using a highly secure "pointing" pattern and OpenID Connect (OIDC) compliant token generation.

[![Go Reference](https://pkg.go.dev/badge/github.com/oririfai/revoluchat-go-sdk.svg)](https://pkg.go.dev/github.com/oririfai/revoluchat-go-sdk)

## 📋 Detailed Specifications

This SDK operates via high-speed gRPC communication and provides various capabilities to support your application's needs:

### 1. Multi-Tier Architecture

- **Normal Tier**: Conversation data is stored locally in Revoluchat's internal storage (Elixir). You only need to implement the _User Provider_ and _JWT Manager_.
- **Advance Tier**: Conversation, call, group, and attachment data are stored and managed on your application side. Revoluchat acts as a forwarding agent that invokes the gRPC Providers you implement.

### 2. Security & Authentication

- **OIDC (OpenID Connect)**: Generates RS256 JWTs with dynamically computed Key IDs (KID).
- **Server-to-Server Auth**: Every inter-service request must include an `X-Server-Key` header to prevent unauthorized access to the gRPC services and the JWKS endpoint.

### 3. Services & Providers (Advance Tier)

This SDK includes comprehensive gRPC services:

- **User Service**: Retrieval of profiles, statuses, avatars, names, etc.
- **Message Service**: Storing messages, retrieving message history, marking as read/delivered, and message deletion.
- **Conversation Service**: Creation, retrieval, and listing of conversations (both 1-on-1 and groups).
- **Call Service**: Initiating calls, updating call statuses, and managing call history.
- **Group Service**: Full management for groups, adding/removing members, muting, accepting invitations, and other group administrative actions.
- **Attachment Service**: Registration and listing of image, video, and document attachment files.

---

## 📦 Installation & Setup

Run the following command at the root of your Go backend project:

```bash
go get github.com/oririfai/revoluchat-go-sdk@latest
```

### ⚡ Quick Start (Recommended)

Use our interactive CLI to bootstrap your configuration in seconds:

```bash
go run github.com/oririfai/revoluchat-go-sdk/cmd/revoluchat-init@latest
```

This process will ask for your **Tier** and **Server Key**, and then automatically generate the `revoluchat_setup.go` boilerplate file.

---

## 💻 How to Use

Integrating this SDK involves initializing the **JWT Manager** for OIDC and running the SDK's **gRPC Server** with your custom Providers.

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
	serverKey := os.Getenv("REVOLUCHAT_SERVER_KEY")

	// 1. Initialize JWT Manager (RS256)
	// Only public key is needed to expose the JWKS endpoint for Elixir
	jwtManager, err := revoluchat.NewJWTManager("keys/app.rsa.pub", serverKey)
	if err != nil {
		log.Fatalf("Failed to initialize JWT: %v", err)
	}

	// 2. Start the gRPC SDK Server
	go func() {
		err := revoluchat.Start(revoluchat.Config{
			Tier:      revoluchat.TierAdvance, // Choose TierNormal or TierAdvance
			GRPCPort:  50051,
			ServerKey: serverKey,

			// Required to be implemented in All Tiers
			UserProvider: func(ctx context.Context, id string) (*revoluchat.User, error) {
				// Query your database/service here to get the User profile
				return &revoluchat.User{
					ID:        id,
					UUID:      "your-custom-uuid",
					Name:      "User Full Name",
					Phone:     "08123456789",
					AvatarURL: "https://your-cdn.com/path/to/photo.jpg",
					Status:    "active",
					IsKYC:     true,
				}, nil
			},

			// Required to be implemented if using TierAdvance
			MessageProvider:      &MyMessageService{},      // Must implement revoluchat.MessageProvider
			ConversationProvider: &MyConversationService{}, // Must implement revoluchat.ConversationProvider
			CallProvider:         &MyCallService{},         // Must implement revoluchat.CallProvider
			GroupProvider:        &MyGroupService{},        // Must implement revoluchat.GroupProvider
			AttachmentProvider:   &MyAttachmentService{},   // Must implement revoluchat.AttachmentProvider
		})

		if err != nil {
			log.Fatalf("Failed to start Revoluchat gRPC server: %v", err)
		}
	}()

	// 3. HTTP Server Setup for the JWKS endpoint
	r := gin.Default()
	r.GET("/jwks", gin.WrapF(jwtManager.JWKSHandler))

	log.Println("Starting HTTP server on :8089")
	r.Run(":8089")
}
```

### 2. Configuration Integration to the Revoluchat Server

Set the environment variables in your Revoluchat backend to point to your gRPC IP and Port:

```bash
# Set the tier type to be uniform with your application (normal/advance)
TIER_TYPE=advance

# Point to your Backend gRPC host (only advanced tier)
CHAT_SERVICE_GRPC_ENDPOINT=your-go-backend-host:50051
USER_SERVICE_GRPC_ENDPOINT=your-go-backend-host:50051
```

---

## ℹ️ Other Information

### Data Structure `revoluchat.User`

The SDK maps the profile type with a simplified `revoluchat.User` structure so it can be transferred quickly over gRPC:

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

You can call the e.g `jwtManager.GenerateToken(userID string, appID string, privKeyPath string) (string, error)` function when a user successfully logs into your system. This token can be returned to the client-side (mobile/web app) to be used for logging into the Revoluchat WebSocket server.

## 📄 License

MIT © [Revoluchat Team](https://revolu.id)
