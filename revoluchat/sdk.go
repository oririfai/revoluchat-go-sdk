package revoluchat

import (
	"context"
	"fmt"
	"net"

	pb_user "github.com/oririfai/revoluchat-go-sdk/proto/user_v1"
	pb_chat "github.com/oririfai/revoluchat-go-sdk/proto/chat_v1"
	pb_admin "github.com/oririfai/revoluchat-go-sdk/proto/admin_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Tier represents the Revoluchat architecture tier.
type Tier string

const (
	TierNormal  Tier = "normal"
	TierAdvance Tier = "advance"
)

// User represents the user data structure required by Revoluchat.
type User struct {
	ID        string // Numeric ID (stringified)
	UUID      string // UUID
	Name      string
	Phone     string
	Status    string
	IsKYC     bool
	AvatarURL string
}

// UserProvider is a function that returns a User given a UUID.
type UserProvider func(ctx context.Context, id string) (*User, error)

// --- ADVANCE TIER PROVIDERS ---

type MessageProvider interface {
	InsertMessage(ctx context.Context, req *pb_chat.InsertMessageRequest) (*pb_chat.InsertMessageResponse, error)
	ListMessages(ctx context.Context, req *pb_chat.ListMessagesRequest) (*pb_chat.ListMessagesResponse, error)
	MarkRead(ctx context.Context, req *pb_chat.MarkReadRequest) (*pb_chat.MarkReadResponse, error)
	MarkDelivered(ctx context.Context, req *pb_chat.MarkDeliveredRequest) (*pb_chat.MarkDeliveredResponse, error)
	DeleteMessage(ctx context.Context, req *pb_chat.DeleteMessageRequest) (*pb_chat.DeleteMessageResponse, error)
	BulkDeleteMessages(ctx context.Context, req *pb_chat.BulkDeleteMessagesRequest) (*pb_chat.BulkDeleteMessagesResponse, error)
}

type ConversationProvider interface {
	CreateConversation(ctx context.Context, req *pb_chat.CreateConversationRequest) (*pb_chat.CreateConversationResponse, error)
	ListConversations(ctx context.Context, req *pb_chat.ListConversationsRequest) (*pb_chat.ListConversationsResponse, error)
	GetConversation(ctx context.Context, req *pb_chat.GetConversationRequest) (*pb_chat.GetConversationResponse, error)
	DeleteConversation(ctx context.Context, req *pb_chat.DeleteConversationRequest) (*pb_chat.ActionResponse, error)
	ArchiveConversation(ctx context.Context, req *pb_chat.ArchiveConversationRequest) (*pb_chat.ActionResponse, error)
	UnarchiveConversation(ctx context.Context, req *pb_chat.UnarchiveConversationRequest) (*pb_chat.ActionResponse, error)
	ListArchivedConversations(ctx context.Context, req *pb_chat.ListConversationsRequest) (*pb_chat.ListConversationsResponse, error)
}

type CallProvider interface {
	InitiateCall(ctx context.Context, req *pb_chat.InitiateCallRequest) (*pb_chat.InitiateCallResponse, error)
	UpdateCallStatus(ctx context.Context, req *pb_chat.UpdateCallStatusRequest) (*pb_chat.UpdateCallStatusResponse, error)
	GetCall(ctx context.Context, req *pb_chat.GetCallRequest) (*pb_chat.GetCallResponse, error)
	ListCallHistory(ctx context.Context, req *pb_chat.ListCallHistoryRequest) (*pb_chat.ListCallHistoryResponse, error)
	DeleteCallHistory(ctx context.Context, req *pb_chat.DeleteCallHistoryRequest) (*pb_chat.DeleteCallHistoryResponse, error)
}

type GroupProvider interface {
	CreateGroup(ctx context.Context, req *pb_chat.CreateGroupRequest) (*pb_chat.CreateGroupResponse, error)
	GetGroup(ctx context.Context, req *pb_chat.GetGroupRequest) (*pb_chat.GetGroupResponse, error)
	AddMembers(ctx context.Context, req *pb_chat.AddMembersRequest) (*pb_chat.ActionResponse, error)
	RemoveMember(ctx context.Context, req *pb_chat.RemoveMemberRequest) (*pb_chat.ActionResponse, error)
	UpdateGroup(ctx context.Context, req *pb_chat.UpdateGroupRequest) (*pb_chat.UpdateGroupResponse, error)
	LeaveGroup(ctx context.Context, req *pb_chat.LeaveGroupRequest) (*pb_chat.ActionResponse, error)
	DeleteGroup(ctx context.Context, req *pb_chat.DeleteGroupRequest) (*pb_chat.ActionResponse, error)
	MuteGroup(ctx context.Context, req *pb_chat.MuteGroupRequest) (*pb_chat.ActionResponse, error)
	AcceptGroupInvitation(ctx context.Context, req *pb_chat.AcceptGroupInvitationRequest) (*pb_chat.ActionResponse, error)
}

type AttachmentProvider interface {
	RegisterAttachment(ctx context.Context, req *pb_chat.RegisterAttachmentRequest) (*pb_chat.RegisterAttachmentResponse, error)
	ListAttachmentsByIds(ctx context.Context, req *pb_chat.ListAttachmentsByIdsRequest) (*pb_chat.ListAttachmentsByIdsResponse, error)
}

// Config holds the SDK configuration.
type Config struct {
	Tier         Tier
	GRPCPort     int
	ServerKey    string // Shared secret for server-to-server auth
	UserProvider UserProvider
	
	// Required for TierAdvance
	MessageProvider      MessageProvider
	ConversationProvider ConversationProvider
	CallProvider         CallProvider
	GroupProvider        GroupProvider
	AttachmentProvider   AttachmentProvider
	ContactProvider      ContactProvider
	AdminProvider        AdminProvider
}

type ContactProvider interface {
	SearchUserByPhone(ctx context.Context, req *pb_user.SearchUserByPhoneRequest) (*pb_user.UserResponse, error)
	AddContact(ctx context.Context, req *pb_user.AddContactRequest) (*pb_user.ActionResponse, error)
	ListContacts(ctx context.Context, req *pb_user.ListContactsRequest) (*pb_user.ListContactsResponse, error)
	RemoveContact(ctx context.Context, req *pb_user.RemoveContactRequest) (*pb_user.ActionResponse, error)
}

type AdminProvider interface {
	ListUsers(ctx context.Context, req *pb_admin.ListUsersRequest) (*pb_admin.ListUsersResponse, error)
	SuspendUser(ctx context.Context, req *pb_admin.SuspendUserRequest) (*pb_admin.ActionResponse, error)
	UnsuspendUser(ctx context.Context, req *pb_admin.UnsuspendUserRequest) (*pb_admin.ActionResponse, error)
}

// --- ADMIN SERVICE ---

func (s *server) ListUsers(ctx context.Context, req *pb_admin.ListUsersRequest) (*pb_admin.ListUsersResponse, error) {
	if s.config.AdminProvider == nil {
		return nil, status.Error(codes.Unimplemented, "admin provider not configured")
	}
	return s.config.AdminProvider.ListUsers(ctx, req)
}

func (s *server) SuspendUser(ctx context.Context, req *pb_admin.SuspendUserRequest) (*pb_admin.ActionResponse, error) {
	if s.config.AdminProvider == nil {
		return nil, status.Error(codes.Unimplemented, "admin provider not configured")
	}
	return s.config.AdminProvider.SuspendUser(ctx, req)
}

func (s *server) UnsuspendUser(ctx context.Context, req *pb_admin.UnsuspendUserRequest) (*pb_admin.ActionResponse, error) {
	if s.config.AdminProvider == nil {
		return nil, status.Error(codes.Unimplemented, "admin provider not configured")
	}
	return s.config.AdminProvider.UnsuspendUser(ctx, req)
}

type server struct {
	pb_user.UnimplementedUserServiceServer
	pb_chat.UnimplementedMessageServiceServer
	pb_chat.UnimplementedConversationServiceServer
	pb_chat.UnimplementedCallServiceServer
	pb_chat.UnimplementedGroupServiceServer
	pb_chat.UnimplementedAttachmentServiceServer
	pb_admin.UnimplementedAdminServiceServer
	config Config
}

// --- USER SERVICE ---

func (s *server) GetUser(ctx context.Context, req *pb_user.GetUserRequest) (*pb_user.GetUserResponse, error) {
	user, err := s.config.UserProvider(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb_user.GetUserResponse{
		Id:        user.ID,
		Uuid:      user.UUID,
		Name:      user.Name,
		Phone:     user.Phone,
		Status:    user.Status,
		IsKyc:     user.IsKYC,
		AvatarUrl: user.AvatarURL,
	}, nil
}

func (s *server) SearchUserByPhone(ctx context.Context, req *pb_user.SearchUserByPhoneRequest) (*pb_user.UserResponse, error) {
	if s.config.ContactProvider == nil {
		return nil, status.Error(codes.Unimplemented, "contact provider not configured")
	}
	return s.config.ContactProvider.SearchUserByPhone(ctx, req)
}

func (s *server) AddContact(ctx context.Context, req *pb_user.AddContactRequest) (*pb_user.ActionResponse, error) {
	if s.config.ContactProvider == nil {
		return nil, status.Error(codes.Unimplemented, "contact provider not configured")
	}
	return s.config.ContactProvider.AddContact(ctx, req)
}

func (s *server) ListContacts(ctx context.Context, req *pb_user.ListContactsRequest) (*pb_user.ListContactsResponse, error) {
	if s.config.ContactProvider == nil {
		return nil, status.Error(codes.Unimplemented, "contact provider not configured")
	}
	return s.config.ContactProvider.ListContacts(ctx, req)
}

func (s *server) RemoveContact(ctx context.Context, req *pb_user.RemoveContactRequest) (*pb_user.ActionResponse, error) {
	if s.config.ContactProvider == nil {
		return nil, status.Error(codes.Unimplemented, "contact provider not configured")
	}
	return s.config.ContactProvider.RemoveContact(ctx, req)
}

// --- CHAT SERVICE (ADVANCE TIER) ---

func (s *server) InsertMessage(ctx context.Context, req *pb_chat.InsertMessageRequest) (*pb_chat.InsertMessageResponse, error) {
	if s.config.MessageProvider == nil {
		return nil, status.Error(codes.Unimplemented, "message provider not configured")
	}
	return s.config.MessageProvider.InsertMessage(ctx, req)
}

func (s *server) ListMessages(ctx context.Context, req *pb_chat.ListMessagesRequest) (*pb_chat.ListMessagesResponse, error) {
	if s.config.MessageProvider == nil {
		return nil, status.Error(codes.Unimplemented, "message provider not configured")
	}
	return s.config.MessageProvider.ListMessages(ctx, req)
}

func (s *server) MarkRead(ctx context.Context, req *pb_chat.MarkReadRequest) (*pb_chat.MarkReadResponse, error) {
	if s.config.MessageProvider == nil {
		return nil, status.Error(codes.Unimplemented, "message provider not configured")
	}
	return s.config.MessageProvider.MarkRead(ctx, req)
}

func (s *server) MarkDelivered(ctx context.Context, req *pb_chat.MarkDeliveredRequest) (*pb_chat.MarkDeliveredResponse, error) {
	if s.config.MessageProvider == nil {
		return nil, status.Error(codes.Unimplemented, "message provider not configured")
	}
	return s.config.MessageProvider.MarkDelivered(ctx, req)
}

func (s *server) DeleteMessage(ctx context.Context, req *pb_chat.DeleteMessageRequest) (*pb_chat.DeleteMessageResponse, error) {
	if s.config.MessageProvider == nil {
		return nil, status.Error(codes.Unimplemented, "message provider not configured")
	}
	return s.config.MessageProvider.DeleteMessage(ctx, req)
}

func (s *server) BulkDeleteMessages(ctx context.Context, req *pb_chat.BulkDeleteMessagesRequest) (*pb_chat.BulkDeleteMessagesResponse, error) {
	if s.config.MessageProvider == nil {
		return nil, status.Error(codes.Unimplemented, "message provider not configured")
	}
	return s.config.MessageProvider.BulkDeleteMessages(ctx, req)
}

// --- CONVERSATION SERVICE (ADVANCE TIER) ---

func (s *server) CreateConversation(ctx context.Context, req *pb_chat.CreateConversationRequest) (*pb_chat.CreateConversationResponse, error) {
	if s.config.ConversationProvider == nil {
		return nil, status.Error(codes.Unimplemented, "conversation provider not configured")
	}
	return s.config.ConversationProvider.CreateConversation(ctx, req)
}

func (s *server) ListConversations(ctx context.Context, req *pb_chat.ListConversationsRequest) (*pb_chat.ListConversationsResponse, error) {
	if s.config.ConversationProvider == nil {
		return nil, status.Error(codes.Unimplemented, "conversation provider not configured")
	}
	return s.config.ConversationProvider.ListConversations(ctx, req)
}

func (s *server) GetConversation(ctx context.Context, req *pb_chat.GetConversationRequest) (*pb_chat.GetConversationResponse, error) {
	if s.config.ConversationProvider == nil {
		return nil, status.Error(codes.Unimplemented, "conversation provider not configured")
	}
	return s.config.ConversationProvider.GetConversation(ctx, req)
}

func (s *server) DeleteConversation(ctx context.Context, req *pb_chat.DeleteConversationRequest) (*pb_chat.ActionResponse, error) {
	if s.config.ConversationProvider == nil {
		return nil, status.Error(codes.Unimplemented, "conversation provider not configured")
	}
	return s.config.ConversationProvider.DeleteConversation(ctx, req)
}

func (s *server) ArchiveConversation(ctx context.Context, req *pb_chat.ArchiveConversationRequest) (*pb_chat.ActionResponse, error) {
	if s.config.ConversationProvider == nil {
		return nil, status.Error(codes.Unimplemented, "conversation provider not configured")
	}
	return s.config.ConversationProvider.ArchiveConversation(ctx, req)
}

func (s *server) UnarchiveConversation(ctx context.Context, req *pb_chat.UnarchiveConversationRequest) (*pb_chat.ActionResponse, error) {
	if s.config.ConversationProvider == nil {
		return nil, status.Error(codes.Unimplemented, "conversation provider not configured")
	}
	return s.config.ConversationProvider.UnarchiveConversation(ctx, req)
}

func (s *server) ListArchivedConversations(ctx context.Context, req *pb_chat.ListConversationsRequest) (*pb_chat.ListConversationsResponse, error) {
	if s.config.ConversationProvider == nil {
		return nil, status.Error(codes.Unimplemented, "conversation provider not configured")
	}
	return s.config.ConversationProvider.ListArchivedConversations(ctx, req)
}

// --- CALL SERVICE (ADVANCE TIER) ---

func (s *server) InitiateCall(ctx context.Context, req *pb_chat.InitiateCallRequest) (*pb_chat.InitiateCallResponse, error) {
	if s.config.CallProvider == nil {
		return nil, status.Error(codes.Unimplemented, "call provider not configured")
	}
	return s.config.CallProvider.InitiateCall(ctx, req)
}

func (s *server) UpdateCallStatus(ctx context.Context, req *pb_chat.UpdateCallStatusRequest) (*pb_chat.UpdateCallStatusResponse, error) {
	if s.config.CallProvider == nil {
		return nil, status.Error(codes.Unimplemented, "call provider not configured")
	}
	return s.config.CallProvider.UpdateCallStatus(ctx, req)
}

func (s *server) GetCall(ctx context.Context, req *pb_chat.GetCallRequest) (*pb_chat.GetCallResponse, error) {
	if s.config.CallProvider == nil {
		return nil, status.Error(codes.Unimplemented, "call provider not configured")
	}
	return s.config.CallProvider.GetCall(ctx, req)
}

func (s *server) ListCallHistory(ctx context.Context, req *pb_chat.ListCallHistoryRequest) (*pb_chat.ListCallHistoryResponse, error) {
	if s.config.CallProvider == nil {
		return nil, status.Error(codes.Unimplemented, "call provider not configured")
	}
	return s.config.CallProvider.ListCallHistory(ctx, req)
}

func (s *server) DeleteCallHistory(ctx context.Context, req *pb_chat.DeleteCallHistoryRequest) (*pb_chat.DeleteCallHistoryResponse, error) {
	if s.config.CallProvider == nil {
		return nil, status.Error(codes.Unimplemented, "call provider not configured")
	}
	return s.config.CallProvider.DeleteCallHistory(ctx, req)
}

// --- GROUP SERVICE (ADVANCE TIER) ---

func (s *server) CreateGroup(ctx context.Context, req *pb_chat.CreateGroupRequest) (*pb_chat.CreateGroupResponse, error) {
	if s.config.GroupProvider == nil {
		return nil, status.Error(codes.Unimplemented, "group provider not configured")
	}
	return s.config.GroupProvider.CreateGroup(ctx, req)
}

func (s *server) GetGroup(ctx context.Context, req *pb_chat.GetGroupRequest) (*pb_chat.GetGroupResponse, error) {
	if s.config.GroupProvider == nil {
		return nil, status.Error(codes.Unimplemented, "group provider not configured")
	}
	return s.config.GroupProvider.GetGroup(ctx, req)
}

func (s *server) AddMembers(ctx context.Context, req *pb_chat.AddMembersRequest) (*pb_chat.ActionResponse, error) {
	if s.config.GroupProvider == nil {
		return nil, status.Error(codes.Unimplemented, "group provider not configured")
	}
	return s.config.GroupProvider.AddMembers(ctx, req)
}

func (s *server) RemoveMember(ctx context.Context, req *pb_chat.RemoveMemberRequest) (*pb_chat.ActionResponse, error) {
	if s.config.GroupProvider == nil {
		return nil, status.Error(codes.Unimplemented, "group provider not configured")
	}
	return s.config.GroupProvider.RemoveMember(ctx, req)
}

func (s *server) UpdateGroup(ctx context.Context, req *pb_chat.UpdateGroupRequest) (*pb_chat.UpdateGroupResponse, error) {
	if s.config.GroupProvider == nil {
		return nil, status.Error(codes.Unimplemented, "group provider not configured")
	}
	return s.config.GroupProvider.UpdateGroup(ctx, req)
}

func (s *server) LeaveGroup(ctx context.Context, req *pb_chat.LeaveGroupRequest) (*pb_chat.ActionResponse, error) {
	if s.config.GroupProvider == nil {
		return nil, status.Error(codes.Unimplemented, "group provider not configured")
	}
	return s.config.GroupProvider.LeaveGroup(ctx, req)
}

func (s *server) DeleteGroup(ctx context.Context, req *pb_chat.DeleteGroupRequest) (*pb_chat.ActionResponse, error) {
	if s.config.GroupProvider == nil {
		return nil, status.Error(codes.Unimplemented, "group provider not configured")
	}
	return s.config.GroupProvider.DeleteGroup(ctx, req)
}

func (s *server) MuteGroup(ctx context.Context, req *pb_chat.MuteGroupRequest) (*pb_chat.ActionResponse, error) {
	if s.config.GroupProvider == nil {
		return nil, status.Error(codes.Unimplemented, "group provider not configured")
	}
	return s.config.GroupProvider.MuteGroup(ctx, req)
}

func (s *server) AcceptGroupInvitation(ctx context.Context, req *pb_chat.AcceptGroupInvitationRequest) (*pb_chat.ActionResponse, error) {
	if s.config.GroupProvider == nil {
		return nil, status.Error(codes.Unimplemented, "group provider not configured")
	}
	return s.config.GroupProvider.AcceptGroupInvitation(ctx, req)
}

// --- ATTACHMENT SERVICE (ADVANCE TIER) ---

func (s *server) RegisterAttachment(ctx context.Context, req *pb_chat.RegisterAttachmentRequest) (*pb_chat.RegisterAttachmentResponse, error) {
	if s.config.AttachmentProvider == nil {
		return nil, status.Error(codes.Unimplemented, "attachment provider not configured")
	}
	return s.config.AttachmentProvider.RegisterAttachment(ctx, req)
}

func (s *server) ListAttachmentsByIds(ctx context.Context, req *pb_chat.ListAttachmentsByIdsRequest) (*pb_chat.ListAttachmentsByIdsResponse, error) {
	if s.config.AttachmentProvider == nil {
		return nil, status.Error(codes.Unimplemented, "attachment provider not configured")
	}
	return s.config.AttachmentProvider.ListAttachmentsByIds(ctx, req)
}

// Start starts the gRPC server for Revoluchat integration.
func Start(config Config) error {
	addr := fmt.Sprintf(":%d", config.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Security Interceptor
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if config.ServerKey != "" {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Error(codes.Unauthenticated, "metadata missing")
			}
			keys := md.Get("x-server-key")
			if len(keys) == 0 {
				return nil, status.Error(codes.Unauthenticated, "invalid server key")
			}
			if keys[0] != config.ServerKey {
				return nil, status.Error(codes.Unauthenticated, "invalid server key")
			}
		}
		return handler(ctx, req)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor))
	
	srv := &server{config: config}
	
	// Register User Service (Always available)
	pb_user.RegisterUserServiceServer(s, srv)

	// Register Advance Tier Services if configured
	if config.Tier == TierAdvance {
		pb_chat.RegisterMessageServiceServer(s, srv)
		pb_chat.RegisterConversationServiceServer(s, srv)
		pb_chat.RegisterCallServiceServer(s, srv)
		pb_chat.RegisterGroupServiceServer(s, srv)
		pb_chat.RegisterAttachmentServiceServer(s, srv)
		pb_admin.RegisterAdminServiceServer(s, srv)
	}

	fmt.Printf("Revoluchat Go SDK [%s tier]: gRPC server listening on %s\n", config.Tier, addr)
	return s.Serve(lis)
}
