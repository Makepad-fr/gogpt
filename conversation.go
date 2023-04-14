package gogpt

import (
	"github.com/Makepad-fr/gogpt/internal"
	"github.com/google/uuid"
)

type ConversationHistoryItem struct {
	idBasedItem
	ID         string `json:"id"`
	Title      string `json:"title"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

func (c ConversationHistoryItem) getId() string {
	return c.ID
}

type Conversation struct {
	Title             string                          `json:"title"`
	CreateTime        float64                         `json:"create_time"`
	UpdateTime        float64                         `json:"update_time"`
	Mapping           map[string]internal.MappingNode `json:"mapping"`
	ModerationResults []interface{}                   `json:"moderation_results"`
	CurrentNode       string                          `json:"current_node"`
}

type ConversationHistoryResponse struct {
	Items                   []ConversationHistoryItem `json:"items"`
	Total                   int                       `json:"total"`
	Limit                   int                       `json:"limit"`
	Offset                  int                       `json:"offset"`
	HasMissingConversations bool                      `json:"has_missing_conversations"`
}

type ConversationResponse struct {
	Message        internal.Message `json:"message"`
	ConversationID string           `json:"conversation_id"`
	Error          *string          `json:"error"`
}

type GenerateConversationTitleResponse struct {
	Title string `json:"title"`
}

type TextModerationResponse struct {
	Blocked      bool   `json:"blocked"`
	Flagged      bool   `json:"flagged"`
	ModerationId string `json:"moderation_id"`
}

type conversationResponseConsumer func(event ConversationResponse)

func createMessageRequestInExistingConversation(message, model, conversationUUID string, timeZoneOffset int) (*internal.NewMessageRequest, error) {
	messageRequest, err := createMessageRequestForNewConversation(message, model, timeZoneOffset)
	if err != nil {
		return nil, err
	}
	messageRequest.ConversationId = conversationUUID
	return messageRequest, nil
}

func createMessageRequestForNewConversation(message, model string, timeZoneOffset int) (*internal.NewMessageRequest, error) {
	messageUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	parentUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &internal.NewMessageRequest{
		Action: "next",
		Messages: []internal.Message{
			{
				ID: messageUUID.String(),
				Author: internal.Author{
					Role: "user",
				},
				Content: internal.Content{
					ContentType: "text",
					Parts:       []string{message},
				},
			},
		},
		Model:             model,
		TimezoneOffsetMin: timeZoneOffset,
		ParentMessageID:   parentUUID.String(),
	}, nil
}
