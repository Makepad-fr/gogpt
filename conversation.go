package gogpt

import "github.com/google/uuid"

type Author struct {
	Role     string                 `json:"role"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type Content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type Message struct {
	ID         string                 `json:"id"`
	Author     Author                 `json:"author"`
	CreateTime float64                `json:"create_time,omitempty"`
	Content    Content                `json:"content"`
	EndTurn    bool                   `json:"end_turn,omitempty"`
	Weight     float64                `json:"weight,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Recipient  string                 `json:"recipient,omitempty"`
}

type MappingNode struct {
	ID       string   `json:"id"`
	Message  *Message `json:"message,omitempty"`
	Parent   string   `json:"parent,omitempty"`
	Children []string `json:"children"`
}

type Conversation struct {
	Title             string                 `json:"title"`
	CreateTime        float64                `json:"create_time"`
	UpdateTime        float64                `json:"update_time"`
	Mapping           map[string]MappingNode `json:"mapping"`
	ModerationResults []interface{}          `json:"moderation_results"`
	CurrentNode       string                 `json:"current_node"`
}

type ConversationHistoryResponse struct {
	Items                   []ConversationHistoryItem `json:"items"`
	Total                   int                       `json:"total"`
	Limit                   int                       `json:"limit"`
	Offset                  int                       `json:"offset"`
	HasMissingConversations bool                      `json:"has_missing_conversations"`
}

type ConversationHistoryItem struct {
	idBasedSetItem
	ID         string `json:"id"`
	Title      string `json:"title"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

func (c ConversationHistoryItem) getId() string {
	return c.ID
}

type NewMessageRequest struct {
	Action            string    `json:"action"`
	Messages          []Message `json:"messages"`
	ParentMessageID   string    `json:"parent_message_id"`
	Model             string    `json:"model"`
	TimezoneOffsetMin int       `json:"timezone_offset_min"`
	ConversationId    string    `json:"conversation_id,omitempty"`
}

func createMessageRequestInExistingConversation(message, model, conversationUUID string) (*NewMessageRequest, error) {
	messageRequest, err := createMessageRequestForNewConversation(message, model)
	if err != nil {
		return nil, err
	}
	messageRequest.ConversationId = conversationUUID
	return messageRequest, nil
}

func createMessageRequestForNewConversation(message, model string) (*NewMessageRequest, error) {
	messageUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	parentUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &NewMessageRequest{
		Action: "next",
		Messages: []Message{
			{
				ID: messageUUID.String(),
				Author: Author{
					Role: "user",
				},
				Content: Content{
					ContentType: "text",
					Parts:       []string{message},
				},
			},
		},
		Model:             model,
		TimezoneOffsetMin: -120,
		ParentMessageID:   parentUUID.String(),
	}, nil
}
