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
    REVOLUCHAT SDK INITIALIZER
=========================================
`

const setupTemplate = `package main

import (
	"context"
	"fmt"
	"log"

	"github.com/oririfai/revoluchat-go-sdk/revoluchat"
	{{ if eq .Tier "advance" }}
	pb_chat "github.com/oririfai/revoluchat-go-sdk/proto/chat_v1"
	pb_user "github.com/oririfai/revoluchat-go-sdk/proto/user_v1"
	pb_admin "github.com/oririfai/revoluchat-go-sdk/proto/admin_v1"
	{{ end }}
)

// InitRevoluchat sets up the Revoluchat SDK based on your configuration.
// You can move this entire file and code to another folder or file for initialize this package
// You can use your env variable to change config value
func InitRevoluchat() {
	config := revoluchat.Config{
		Tier:      revoluchat.Tier{{ .TierTitle }},
		GRPCPort:  {{ .Port }},
		ServerKey: "{{ .ServerKey }}",
		UserProvider: func(ctx context.Context, id string) (*revoluchat.User, error) {
			// TODO: Fetch user from your database
			return &revoluchat.User{
				ID:        id,
				UUID:      "", // UUID if available
				Name:      "Example User",
				Phone:     "",
				Status:    "", // active or suspended (other status will be suspended)
				IsKYC:     false,
				AvatarURL: "",
			}, nil
		},
		{{ if eq .Tier "advance" }}
		// Advance Tier Providers
		MessageProvider:      &ChatHandler{},
		ConversationProvider: &ChatHandler{},
		CallProvider:         &CallHandler{},
		GroupProvider:        &GroupHandler{},
		AttachmentProvider:   &AttachmentHandler{},
		ContactProvider:      &ContactHandler{},
		AdminProvider:        &AdminHandler{},
		{{ end }}
	}

	go func() {
		if err := revoluchat.Start(config); err != nil {
			log.Fatalf("Failed to start Revoluchat SDK: %v", err)
		}
	}()
}

{{ if eq .Tier "advance" }}
// --- REVOLUCHAT ADVANCE TIER HANDLER CONTRACTS ---

type ChatHandler struct{}

// MessageProvider - Handles chat message persistence and delivery statuses.

// InsertMessage saves a new chat message to your database.
// This is triggered when a user sends a message. You should extract the sender, conversation, and content,
// persist it to your local database, and return the saved message details with a new unique message ID.
func (h *ChatHandler) InsertMessage(ctx context.Context, req *pb_chat.InsertMessageRequest) (*pb_chat.InsertMessageResponse, error) {
	// TODO: Implement database insertion logic.
	// 1. Extract request parameters (e.g., req.SenderId, req.ConversationId, req.Body, req.Type).
	// 2. Persist the message record to your database.
	// 3. Return the populated *pb_chat.InsertMessageResponse.
	return nil, fmt.Errorf("unimplemented: InsertMessage is used to save a newly sent message in your database (Sender ID: %s, Conversation ID: %s)", req.SenderId, req.ConversationId)
}

// ListMessages retrieves historical messages in a conversation/group.
// This is triggered when a user opens a chat room or scrolls up to load historical chat logs.
// You should fetch messages from your database with pagination support.
func (h *ChatHandler) ListMessages(ctx context.Context, req *pb_chat.ListMessagesRequest) (*pb_chat.ListMessagesResponse, error) {
	// TODO: Implement message history retrieval logic.
	// 1. Query historical messages filtering by req.ConversationId.
	// 2. Apply pagination limits/before_id using req.Limit and req.BeforeId.
	// 3. Map fetched records to *pb_chat.ListMessagesResponse and return.
	return nil, fmt.Errorf("unimplemented: ListMessages is used to fetch historical messages for a specific conversation (Conversation ID: %s)", req.ConversationId)
}

// MarkRead marks a message as read by the recipient.
// This is triggered when a recipient views/reads a message, updating the read status indicators (double blue checks) on the sender's screen.
func (h *ChatHandler) MarkRead(ctx context.Context, req *pb_chat.MarkReadRequest) (*pb_chat.MarkReadResponse, error) {
	// TODO: Implement message read status update logic.
	// 1. Locate the message record matching req.MessageId.
	// 2. Update its status to 'read' and record the read timestamp.
	// 3. Return a successful *pb_chat.MarkReadResponse.
	return nil, fmt.Errorf("unimplemented: MarkRead is used to mark a message as read (Message ID: %s, User ID: %s)", req.MessageId, req.UserId)
}

// MarkDelivered marks a message as successfully delivered to the recipient's device.
// This is triggered when a recipient's device receives a message push, updating the delivery status indicators (double gray checks).
func (h *ChatHandler) MarkDelivered(ctx context.Context, req *pb_chat.MarkDeliveredRequest) (*pb_chat.MarkDeliveredResponse, error) {
	// TODO: Implement message delivery status update logic.
	// 1. Locate the message record matching req.MessageId.
	// 2. Update its status to 'delivered' and record the delivery timestamp.
	// 3. Return a successful *pb_chat.MarkDeliveredResponse.
	return nil, fmt.Errorf("unimplemented: MarkDelivered is used to mark a message as successfully delivered to the recipient's device (Message ID: %s, User ID: %s)", req.MessageId, req.UserId)
}

// DeleteMessage deletes a specific message.
// This is triggered when a user deletes a single message from their chat interface (e.g., delete for me or delete for everyone).
func (h *ChatHandler) DeleteMessage(ctx context.Context, req *pb_chat.DeleteMessageRequest) (*pb_chat.DeleteMessageResponse, error) {
	// TODO: Implement message deletion logic.
	// 1. Locate the message matching req.MessageId and verify the user has permission to delete it.
	// 2. Perform a hard-delete, soft-delete, or replace the content with a 'deleted message' placeholder.
	// 3. Return a successful *pb_chat.DeleteMessageResponse.
	return nil, fmt.Errorf("unimplemented: DeleteMessage is used to delete a specific message (Message ID: %s, User ID: %s)", req.MessageId, req.UserId)
}

// BulkDeleteMessages deletes multiple messages at once.
// This is triggered when a user clears chat history or deletes multiple selected messages simultaneously.
func (h *ChatHandler) BulkDeleteMessages(ctx context.Context, req *pb_chat.BulkDeleteMessagesRequest) (*pb_chat.BulkDeleteMessagesResponse, error) {
	// TODO: Implement bulk message deletion logic.
	// 1. Iterate over req.MessageIds and perform a bulk delete query.
	// 2. Verify deletion permissions for the requesting user (req.UserId).
	// 3. Return a successful *pb_chat.BulkDeleteMessagesResponse.
	return nil, fmt.Errorf("unimplemented: BulkDeleteMessages is used to delete a batch of messages at once (User ID: %s, Message Count: %d)", req.UserId, len(req.MessageIds))
}

// ConversationProvider - Handles chat session management (direct and group).

// CreateConversation initializes a new direct or group chat session.
// This is triggered when a user initiates a new direct message conversation with another user.
func (h *ChatHandler) CreateConversation(ctx context.Context, req *pb_chat.CreateConversationRequest) (*pb_chat.CreateConversationResponse, error) {
	// TODO: Implement conversation initialization logic.
	// 1. Create a new conversation record between req.UserAId and req.UserBId in your database.
	// 2. Check if a conversation between these users already exists.
	// 3. Return the conversation details in *pb_chat.CreateConversationResponse.
	return nil, fmt.Errorf("unimplemented: CreateConversation is used to initialize a new chat session between %s and %s (App ID: %s)", req.UserAId, req.UserBId, req.AppId)
}

// ListConversations lists all active conversations for a user.
// This is triggered to populate the main chat dashboard/list screen showing all active chats, last message preview, and unread badge count.
func (h *ChatHandler) ListConversations(ctx context.Context, req *pb_chat.ListConversationsRequest) (*pb_chat.ListConversationsResponse, error) {
	// TODO: Implement active conversation listing logic.
	// 1. Query all active conversations that the user req.UserId is a member of.
	// 2. Join/fetch the last message, unread message count, and participant details.
	// 3. Map to *pb_chat.ListConversationsResponse and return.
	return nil, fmt.Errorf("unimplemented: ListConversations is used to retrieve the active chat list for a user (User ID: %s)", req.UserId)
}

// GetConversation fetches details of a specific conversation.
// This is triggered when opening a specific conversation room to load its metadata, list of participants, and settings.
func (h *ChatHandler) GetConversation(ctx context.Context, req *pb_chat.GetConversationRequest) (*pb_chat.GetConversationResponse, error) {
	// TODO: Implement conversation detail retrieval logic.
	// 1. Query conversation details matching req.ConversationId.
	// 2. Fetch the participant list and details for this conversation.
	// 3. Return the details in *pb_chat.GetConversationResponse.
	return nil, fmt.Errorf("unimplemented: GetConversation is used to fetch details of a specific conversation (Conversation ID: %s)", req.ConversationId)
}

// DeleteConversation removes a conversation from a user's active list.
// This is triggered when a user chooses to delete or hide a chat room from their local chat list view.
func (h *ChatHandler) DeleteConversation(ctx context.Context, req *pb_chat.DeleteConversationRequest) (*pb_chat.ActionResponse, error) {
	// TODO: Implement conversation deletion/hiding logic.
	// 1. Find the conversation membership records for the user (req.UserId) matching the IDs in req.Ids.
	// 2. Soft-delete the memberships or flag them as hidden/deleted.
	// 3. Return a successful *pb_chat.ActionResponse.
	return nil, fmt.Errorf("unimplemented: DeleteConversation is used to remove or hide conversations (IDs: %v) for user %s", req.Ids, req.UserId)
}

// ArchiveConversation moves a conversation to the user's archives.
// This is triggered when a user archives a chat room to declutter their active chat list.
func (h *ChatHandler) ArchiveConversation(ctx context.Context, req *pb_chat.ArchiveConversationRequest) (*pb_chat.ActionResponse, error) {
	// TODO: Implement conversation archiving logic.
	// 1. Find the conversation memberships matching the IDs in req.Ids for user req.UserId.
	// 2. Update the status of these conversation memberships to 'archived'.
	// 3. Return a successful *pb_chat.ActionResponse.
	return nil, fmt.Errorf("unimplemented: ArchiveConversation is used to archive conversations (IDs: %v) for user %s", req.Ids, req.UserId)
}

// UnarchiveConversation restores an archived conversation.
// This is triggered when a user unarchives a chat room to bring it back to their active chat list.
func (h *ChatHandler) UnarchiveConversation(ctx context.Context, req *pb_chat.UnarchiveConversationRequest) (*pb_chat.ActionResponse, error) {
	// TODO: Implement conversation unarchiving logic.
	// 1. Find the archived conversation memberships matching the IDs in req.Ids for user req.UserId.
	// 2. Update their status back to active/unarchived.
	// 3. Return a successful *pb_chat.ActionResponse.
	return nil, fmt.Errorf("unimplemented: UnarchiveConversation is used to restore archived conversations (IDs: %v) to active status for user %s", req.Ids, req.UserId)
}

// ListArchivedConversations lists all archived conversations for a user.
// This is triggered when a user opens their archived chats section to view their archived chat history.
func (h *ChatHandler) ListArchivedConversations(ctx context.Context, req *pb_chat.ListConversationsRequest) (*pb_chat.ListConversationsResponse, error) {
	// TODO: Implement archived conversation listing logic.
	// 1. Query conversation memberships for the user that have the 'archived' status.
	// 2. Fetch the last message and metadata for each.
	// 3. Map to *pb_chat.ListConversationsResponse and return.
	return nil, fmt.Errorf("unimplemented: ListArchivedConversations is used to fetch all archived chat rooms for a user (User ID: %s)", req.UserId)
}

type CallHandler struct{}

// CallProvider - Handles VoIP call sessions and logs.

// InitiateCall logs the initialization of a new voice/video call.
// This is triggered when a user starts an outgoing voice or video call.
func (h *CallHandler) InitiateCall(ctx context.Context, req *pb_chat.InitiateCallRequest) (*pb_chat.InitiateCallResponse, error) {
	// TODO: Implement call logging logic.
	// 1. Create a new call log entry in your database (e.g., caller, callee, type, status initiated).
	// 2. Generate a unique call session ID and set the start time.
	// 3. Return the call details in *pb_chat.InitiateCallResponse.
	return nil, fmt.Errorf("unimplemented: InitiateCall is used to log a new voice/video call session (Caller ID: %s, Receiver ID: %s)", req.CallerId, req.ReceiverId)
}

// UpdateCallStatus updates the status of a call session (e.g. ringing, connected, completed).
// This is triggered when a call's state changes (ringing, accepted, ended, missed, rejected, busy).
func (h *CallHandler) UpdateCallStatus(ctx context.Context, req *pb_chat.UpdateCallStatusRequest) (*pb_chat.UpdateCallStatusResponse, error) {
	// TODO: Implement call status update logic.
	// 1. Find the active call log record using req.CallId.
	// 2. Update the status and optionally update the end time or duration if the call is completed.
	// 3. Return *pb_chat.UpdateCallStatusResponse.
	return nil, fmt.Errorf("unimplemented: UpdateCallStatus is used to update the state of a call (Call ID: %s, Status: %s)", req.CallId, req.Status)
}

// GetCall fetches metadata of a specific call session.
// This is triggered to retrieve historical metadata or logs of a specific call.
func (h *CallHandler) GetCall(ctx context.Context, req *pb_chat.GetCallRequest) (*pb_chat.GetCallResponse, error) {
	// TODO: Implement call detail retrieval logic.
	// 1. Query the call record matching req.CallId from the database.
	// 2. Return the populated *pb_chat.GetCallResponse.
	return nil, fmt.Errorf("unimplemented: GetCall is used to retrieve detailed metadata of a call session (Call ID: %s)", req.CallId)
}

// ListCallHistory retrieves the call history log for a specific user.
// This is triggered when a user opens their 'Calls' history tab/screen to view past incoming, outgoing, and missed calls.
func (h *CallHandler) ListCallHistory(ctx context.Context, req *pb_chat.ListCallHistoryRequest) (*pb_chat.ListCallHistoryResponse, error) {
	// TODO: Implement call history retrieval logic.
	// 1. Query all call logs where the user req.UserId is either the caller or receiver.
	// 2. Apply pagination parameter req.Limit.
	// 3. Map to *pb_chat.ListCallHistoryResponse and return.
	return nil, fmt.Errorf("unimplemented: ListCallHistory is used to retrieve a user's call history list (User ID: %s)", req.UserId)
}

// DeleteCallHistory removes selected records from a user's call log.
// This is triggered when a user deletes one or more call logs from their call history tab.
func (h *CallHandler) DeleteCallHistory(ctx context.Context, req *pb_chat.DeleteCallHistoryRequest) (*pb_chat.DeleteCallHistoryResponse, error) {
	// TODO: Implement call history deletion logic.
	// 1. Loop through req.Ids and remove or hide them from the user's call logs in the database.
	// 2. Return a successful *pb_chat.DeleteCallHistoryResponse.
	return nil, fmt.Errorf("unimplemented: DeleteCallHistory is used to delete records from a user's call logs (User ID: %s, Call Count: %d)", req.UserId, len(req.Ids))
}

type GroupHandler struct{}

// GroupProvider - Handles group chat room management.

// CreateGroup creates a new group chat room.
// This is triggered when a user initiates group creation, defining name, initial members list, and avatar.
func (h *GroupHandler) CreateGroup(ctx context.Context, req *pb_chat.CreateGroupRequest) (*pb_chat.CreateGroupResponse, error) {
	// TODO: Implement group creation logic.
	// 1. Create a new group room record with req.Name, req.Description, and req.AvatarUrl.
	// 2. Insert group membership records for the creator (as admin) and the initial members in req.MemberIds / req.AdminIds.
	// 3. Return the new group details in *pb_chat.CreateGroupResponse.
	return nil, fmt.Errorf("unimplemented: CreateGroup is used to create a new group room (Group Name: %s, App ID: %s)", req.Name, req.AppId)
}

// GetGroup retrieves metadata and members of a group.
// This is triggered when opening group info or settings to view members, admins, description, and group configurations.
func (h *GroupHandler) GetGroup(ctx context.Context, req *pb_chat.GetGroupRequest) (*pb_chat.GetGroupResponse, error) {
	// TODO: Implement group metadata retrieval logic.
	// 1. Fetch the group record matching req.GroupId.
	// 2. Query active member list, their roles (e.g., admin, member), and profile details.
	// 3. Map to *pb_chat.GetGroupResponse and return.
	return nil, fmt.Errorf("unimplemented: GetGroup is used to retrieve metadata and member list of a group (Group ID: %s)", req.GroupId)
}

// AddMembers adds new participants to a group.
// This is triggered when a group admin or member adds new people to an existing group chat room.
func (h *GroupHandler) AddMembers(ctx context.Context, req *pb_chat.AddMembersRequest) (*pb_chat.ActionResponse, error) {
	// TODO: Implement group member addition logic.
	// 1. Validate that the requester has permission to add members.
	// 2. Create new membership records in the database for the user IDs in req.UserIds.
	// 3. Return a successful *pb_chat.ActionResponse.
	return nil, fmt.Errorf("unimplemented: AddMembers is used to add participants to a group (Group ID: %s, Users Added Count: %d)", req.GroupId, len(req.UserIds))
}

// RemoveMember removes/kicks a participant from a group.
// This is triggered when an admin kicks a member from the group.
func (h *GroupHandler) RemoveMember(ctx context.Context, req *pb_chat.RemoveMemberRequest) (*pb_chat.ActionResponse, error) {
	// TODO: Implement group member removal logic.
	// 1. Validate that the requester is a group admin or owner.
	// 2. Delete or set the status of the group membership record for req.UserId to 'kicked/removed'.
	// 3. Return a successful *pb_chat.ActionResponse.
	return nil, fmt.Errorf("unimplemented: RemoveMember is used by admins to kick a participant from a group (Group ID: %s, User ID: %s)", req.GroupId, req.UserId)
}

// UpdateGroup updates group metadata (e.g. group name, avatar).
// This is triggered when group details (name, description, settings, or avatar) are changed by an authorized member.
func (h *GroupHandler) UpdateGroup(ctx context.Context, req *pb_chat.UpdateGroupRequest) (*pb_chat.UpdateGroupResponse, error) {
	// TODO: Implement group update logic.
	// 1. Check permissions and locate the group record via req.GroupId.
	// 2. Update fields like req.Name, req.Description, req.AvatarUrl, or req.IsLocked.
	// 3. Return the updated group details in *pb_chat.UpdateGroupResponse.
	return nil, fmt.Errorf("unimplemented: UpdateGroup is used to update group metadata (Group ID: %s, New Name: %s)", req.GroupId, req.Name)
}

// LeaveGroup handles a member voluntarily leaving a group.
// This is triggered when a user decides to leave a group chat room on their own.
func (h *GroupHandler) LeaveGroup(ctx context.Context, req *pb_chat.LeaveGroupRequest) (*pb_chat.ActionResponse, error) {
	// TODO: Implement leave group logic.
	// 1. Find the membership record matching req.GroupId and the user calling it.
	// 2. Delete or set status of this membership record to 'left'.
	// 3. Return a successful *pb_chat.ActionResponse.
	return nil, fmt.Errorf("unimplemented: LeaveGroup is used when a user voluntarily exits a group chat room (Group ID: %s, App ID: %s)", req.GroupId, req.AppId)
}

// DeleteGroup disbands a group chat room.
// This is triggered when a group owner/creator deletes the group completely, removing all members and deleting the group room.
func (h *GroupHandler) DeleteGroup(ctx context.Context, req *pb_chat.DeleteGroupRequest) (*pb_chat.ActionResponse, error) {
	// TODO: Implement group deletion logic.
	// 1. Validate that the requester has authority (owner) to delete the group.
	// 2. Mark the group as deleted and remove or archive all membership records.
	// 3. Return a successful *pb_chat.ActionResponse.
	return nil, fmt.Errorf("unimplemented: DeleteGroup is used to disband a group completely (Group ID: %s, App ID: %s)", req.GroupId, req.AppId)
}

// MuteGroup mutes notifications of a group for a user.
// This is triggered when a user chooses to mute chat alert notifications for a specific group.
func (h *GroupHandler) MuteGroup(ctx context.Context, req *pb_chat.MuteGroupRequest) (*pb_chat.ActionResponse, error) {
	// TODO: Implement group muting logic.
	// 1. Update the user's group membership settings to set notifications to muted (req.Mute).
	// 2. Return a successful *pb_chat.ActionResponse.
	return nil, fmt.Errorf("unimplemented: MuteGroup is used to mute chat alerts for a user (Group ID: %s, Mute: %t)", req.GroupId, req.Mute)
}

// AcceptGroupInvitation handles a user accepting a group invitation.
// This is triggered when a user accepts an invitation or request to join a group chat room.
func (h *GroupHandler) AcceptGroupInvitation(ctx context.Context, req *pb_chat.AcceptGroupInvitationRequest) (*pb_chat.ActionResponse, error) {
	// TODO: Implement invitation acceptance logic.
	// 1. Validate the group invitation status.
	// 2. Update the user's group membership record to active.
	// 3. Return a successful *pb_chat.ActionResponse.
	return nil, fmt.Errorf("unimplemented: AcceptGroupInvitation is used when a user accepts a pending group invite (Group ID: %s, App ID: %s)", req.GroupId, req.AppId)
}

type AttachmentHandler struct{}

// AttachmentProvider - Handles file attachments registration and retrieval.

// RegisterAttachment registers a newly uploaded attachment in your database.
// This is triggered when a file is uploaded, to record its storage URL, size, and metadata in the database.
func (h *AttachmentHandler) RegisterAttachment(ctx context.Context, req *pb_chat.RegisterAttachmentRequest) (*pb_chat.RegisterAttachmentResponse, error) {
	// TODO: Implement attachment registration logic.
	// 1. Create a new record in your attachment table.
	// 2. Save req.StorageKey, req.MimeType, req.Size, and req.Checksum, and associate it with a message/sender if needed.
	// 3. Return the registered attachment info in *pb_chat.RegisterAttachmentResponse.
	return nil, fmt.Errorf("unimplemented: RegisterAttachment is used to register file attachment details (Storage Key: %s, Mime Type: %s)", req.StorageKey, req.MimeType)
}

// ListAttachmentsByIds retrieves a list of attachments by their unique IDs.
// This is triggered when retrieving or rendering a batch of file attachments by their primary key IDs.
func (h *AttachmentHandler) ListAttachmentsByIds(ctx context.Context, req *pb_chat.ListAttachmentsByIdsRequest) (*pb_chat.ListAttachmentsByIdsResponse, error) {
	// TODO: Implement batch attachment retrieval logic.
	// 1. Query your database for all attachments matching the IDs in req.Ids.
	// 2. Map the results to *pb_chat.ListAttachmentsByIdsResponse.
	return nil, fmt.Errorf("unimplemented: ListAttachmentsByIds is used to retrieve metadata for a list of file attachments by their IDs (Count: %d)", len(req.Ids))
}

type ContactHandler struct{}

// ContactProvider - Handles user contact lists.

// SearchUserByPhone searches for registered users by their phone number.
// This is triggered when a user wants to check if a specific phone number is registered on the platform to start a chat.
func (h *ContactHandler) SearchUserByPhone(ctx context.Context, req *pb_user.SearchUserByPhoneRequest) (*pb_user.UserResponse, error) {
	// TODO: Implement phone number search logic.
	// 1. Query your users table for a record matching req.Phone.
	// 2. Return the user profiles inside *pb_user.UserResponse if found, otherwise return an error or empty result.
	return nil, fmt.Errorf("unimplemented: SearchUserByPhone is used to search for registered users by phone (Phone: %s)", req.Phone)
}

// AddContact adds a user to the requester's contacts.
// This is triggered when a user adds another user as a contact (creating a personal contact relationship).
func (h *ContactHandler) AddContact(ctx context.Context, req *pb_user.AddContactRequest) (*pb_user.ActionResponse, error) {
	// TODO: Implement contact creation logic.
	// 1. Save a new contact relationship record (req.OwnerId and req.ContactId) in the database.
	// 2. Return a successful *pb_user.ActionResponse.
	return nil, fmt.Errorf("unimplemented: AddContact is used to add another user to the requester's personal contact list (Owner ID: %s, Contact ID: %s)", req.OwnerId, req.ContactId)
}

// ListContacts retrieves the contacts list of a user.
// This is triggered when a user opens their contacts list tab/screen.
func (h *ContactHandler) ListContacts(ctx context.Context, req *pb_user.ListContactsRequest) (*pb_user.ListContactsResponse, error) {
	// TODO: Implement contact list retrieval logic.
	// 1. Query all contact relationships where the owner is req.UserId.
	// 2. Join or load the contact users' profile details (name, avatar, online status).
	// 3. Return the mapped contact list in *pb_user.ListContactsResponse.
	return nil, fmt.Errorf("unimplemented: ListContacts is used to retrieve the personal contact list of a user (User ID: %s)", req.UserId)
}

// RemoveContact removes a user from the requester's contacts.
// This is triggered when a user deletes/unfriends a contact.
func (h *ContactHandler) RemoveContact(ctx context.Context, req *pb_user.RemoveContactRequest) (*pb_user.ActionResponse, error) {
	// TODO: Implement contact removal logic.
	// 1. Delete the contact relationship record matching req.OwnerId and req.ContactId.
	// 2. Return a successful *pb_user.ActionResponse.
	return nil, fmt.Errorf("unimplemented: RemoveContact is used to remove a user from the requester's contact list (Owner ID: %s, Contact ID: %s)", req.OwnerId, req.ContactId)
}

type AdminHandler struct{}

// AdminProvider - Handles platform administrative operations.

// ListUsers lists all registered users for administration dashboard.
// This is triggered when an administrator views the users directory in the admin panel. Supports pagination/filters.
func (h *AdminHandler) ListUsers(ctx context.Context, req *pb_admin.ListUsersRequest) (*pb_admin.ListUsersResponse, error) {
	// TODO: Implement admin user list logic.
	// 1. Query the users table applying req.Query and status filters.
	// 2. Apply pagination parameters req.Limit and req.Page.
	// 3. Return the user list in *pb_admin.ListUsersResponse.
	return nil, fmt.Errorf("unimplemented: ListUsers is used by administrators to list and filter registered users (Limit: %d, Page: %d)", req.Limit, req.Page)
}

// SuspendUser suspends a user's account with optional duration/reason.
// This is triggered when an administrator blocks/suspends a user's account for violation or platform moderation.
func (h *AdminHandler) SuspendUser(ctx context.Context, req *pb_admin.SuspendUserRequest) (*pb_admin.ActionResponse, error) {
	// TODO: Implement user suspension logic.
	// 1. Locate the user matching req.UserId and set their status to suspended.
	// 2. Save suspension details such as reason or duration.
	// 3. Disconnect any active sockets/sessions for this user if applicable.
	// 4. Return a successful *pb_admin.ActionResponse.
	return nil, fmt.Errorf("unimplemented: SuspendUser is used by administrators to suspend/block a user account (User ID: %s, Reason: %s)", req.UserId, req.Reason)
}

// UnsuspendUser reactivates a suspended user's account.
// This is triggered when an administrator restores a suspended user's account access.
func (h *AdminHandler) UnsuspendUser(ctx context.Context, req *pb_admin.UnsuspendUserRequest) (*pb_admin.ActionResponse, error) {
	// TODO: Implement user activation/unsuspension logic.
	// 1. Locate the user matching req.UserId and restore their status to active.
	// 2. Return a successful *pb_admin.ActionResponse.
	return nil, fmt.Errorf("unimplemented: UnsuspendUser is used by administrators to reactivate a suspended user account (User ID: %s)", req.UserId)
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

	fmt.Printf("\n✅ Successfully generated file %s!\n", fileName)
	fmt.Println("Follow these steps to complete setup:")
	fmt.Println("1. Implement the providers in the generated file.")
	fmt.Printf("2. Call InitRevoluchat() in your main.go.\n\n")
}
