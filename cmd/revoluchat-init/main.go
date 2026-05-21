package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/template"
)

const banner = `
=========================================
    REVOLUCHAT SDK INITIALIZER 🚀
=========================================
`

const setupTemplate = `package main

import (
	"context"
	"fmt"
	"log"

	"github.com/oririfai/revoluchat-go-sdk/revoluchat"
	{{ if eq .Tier "advance" }}pb_chat "github.com/oririfai/revoluchat-go-sdk/proto/chat_v1"{{ end }}
)

// InitRevoluchat sets up the Revoluchat SDK based on your configuration.
func InitRevoluchat() {
	config := revoluchat.Config{
		Tier:      revoluchat.Tier{{ .TierTitle }},
		GRPCPort:  {{ .Port }},
		ServerKey: "{{ .ServerKey }}",
		UserProvider: func(ctx context.Context, id uint64) (*revoluchat.User, error) {
			// TODO: Fetch user from your database
			return &revoluchat.User{
				ID:   id,
				Name: "Example User",
			}, nil
		},
		{{ if eq .Tier "advance" }}
		// Advance Tier Providers
		MessageProvider:      &ChatHandler{},
		ConversationProvider: &ChatHandler{},
		CallProvider:         &CallHandler{},
		GroupProvider:        &GroupHandler{},
		{{ end }}
	}

	go func() {
		if err := revoluchat.Start(config); err != nil {
			log.Fatalf("Failed to start Revoluchat SDK: %v", err)
		}
	}()
}

{{ if eq .Tier "advance" }}
// --- ADVANCE TIER HANDLERS ---

type ChatHandler struct{}

func (h *ChatHandler) InsertMessage(ctx context.Context, req *pb_chat.InsertMessageRequest) (*pb_chat.Message, error) {
	return nil, fmt.Errorf("unimplemented")
}
func (h *ChatHandler) ListMessages(ctx context.Context, req *pb_chat.ListMessagesRequest) ([]*pb_chat.Message, error) {
	return nil, nil
}
func (h *ChatHandler) MarkRead(ctx context.Context, req *pb_chat.MarkReadRequest) (*pb_chat.Message, error) {
	return nil, nil
}
func (h *ChatHandler) DeleteMessage(ctx context.Context, req *pb_chat.DeleteMessageRequest) (*pb_chat.Message, error) {
	return nil, nil
}
func (h *ChatHandler) BulkDeleteMessages(ctx context.Context, req *pb_chat.BulkDeleteMessagesRequest) (uint32, error) {
	return 0, nil
}
func (h *ChatHandler) CreateConversation(ctx context.Context, req *pb_chat.CreateConversationRequest) (*pb_chat.CreateConversationResponse, error) {
	return nil, nil
}
func (h *ChatHandler) ListConversations(ctx context.Context, req *pb_chat.ListConversationsRequest) ([]*pb_chat.Conversation, error) {
	return nil, nil
}

type CallHandler struct{}

func (h *CallHandler) InitiateCall(ctx context.Context, req *pb_chat.InitiateCallRequest) (*pb_chat.Call, error) {
	return nil, nil
}
func (h *CallHandler) UpdateCallStatus(ctx context.Context, req *pb_chat.UpdateCallStatusRequest) (*pb_chat.Call, error) {
	return nil, nil
}
func (h *CallHandler) ListCallHistory(ctx context.Context, req *pb_chat.ListCallHistoryRequest) ([]*pb_chat.CallHistoryRecord, error) {
	return nil, nil
}
func (h *CallHandler) DeleteCallHistory(ctx context.Context, req *pb_chat.DeleteCallHistoryRequest) (uint32, error) {
	return 0, nil
}

type GroupHandler struct{}

func (h *GroupHandler) CreateGroup(ctx context.Context, req *pb_chat.CreateGroupRequest) (*pb_chat.Group, error) {
	return nil, nil
}
func (h *GroupHandler) AddMembers(ctx context.Context, req *pb_chat.AddMembersRequest) error {
	return nil
}
func (h *GroupHandler) RemoveMember(ctx context.Context, req *pb_chat.RemoveMemberRequest) error {
	return nil
}
func (h *GroupHandler) UpdateGroup(ctx context.Context, req *pb_chat.UpdateGroupRequest) (*pb_chat.Group, error) {
	return nil, nil
}
func (h *GroupHandler) LeaveGroup(ctx context.Context, req *pb_chat.LeaveGroupRequest) error {
	return nil
}
func (h *GroupHandler) DeleteGroup(ctx context.Context, req *pb_chat.DeleteGroupRequest) error {
	return nil
}
func (h *GroupHandler) MuteGroup(ctx context.Context, req *pb_chat.MuteGroupRequest) error {
	return nil
}
{{ end }}
`

type Config struct {
	Tier      string
	TierTitle string
	Port      string
	ServerKey string
}

func main() {
	fmt.Print(banner)
	reader := bufio.NewReader(os.Stdin)

	// 1. Choose Tier
	fmt.Println("\nChoose your Revoluchat Tier:")
	fmt.Println("1) Normal  - Local Elixir persistence (default)")
	fmt.Println("2) Advance - Centralized Go persistence")
	fmt.Print("Selection (1/2): ")
	tierChoice, _ := reader.ReadString('\n')
	tierChoice = strings.TrimSpace(tierChoice)

	tier := "normal"
	tierTitle := "Normal"
	if tierChoice == "2" {
		tier = "advance"
		tierTitle = "Advance"
	}

	// 2. Port
	fmt.Print("gRPC Port (default 50051): ")
	port, _ := reader.ReadString('\n')
	port = strings.TrimSpace(port)
	if port == "" {
		port = "50051"
	}

	// 3. Server Key
	fmt.Print("Server Key (Shared Secret): ")
	serverKey, _ := reader.ReadString('\n')
	serverKey = strings.TrimSpace(serverKey)

	config := Config{
		Tier:      tier,
		TierTitle: tierTitle,
		Port:      port,
		ServerKey: serverKey,
	}

	// 4. Generate File
	fileName := "revoluchat_setup.go"
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("❌ Failed to create file: %v\n", err)
		return
	}
	defer f.Close()

	tmpl, _ := template.New("setup").Parse(setupTemplate)
	err = tmpl.Execute(f, config)
	if err != nil {
		fmt.Printf("❌ Failed to generate template: %v\n", err)
		return
	}

	fmt.Printf("\n✅ Successfully generated %s!\n", fileName)
	fmt.Println("Follow these steps to complete setup:")
	fmt.Println("1. Implement the providers in the generated file.")
	fmt.Printf("2. Call InitRevoluchat() in your main.go.\n\n")
}
