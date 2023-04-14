package internal

type TextModerationRequestBody struct {
	ConversationId string `json:"conversation_id"`
	Input          string `json:"input"`
	MessageID      string `json:"message_id"`
	Model          string `json:"model"`
}

type Content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type GenerateConversationTitleRequestBody struct {
	MessageId string `json:"message_id"`
}

type NewMessageRequest struct {
	Action            string    `json:"action"`
	Messages          []Message `json:"messages"`
	ParentMessageID   string    `json:"parent_message_id"`
	Model             string    `json:"model"`
	TimezoneOffsetMin int       `json:"timezone_offset_min"`
	ConversationId    string    `json:"conversation_id,omitempty"`
}

type Author struct {
	Role     string                 `json:"role"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type MappingNode struct {
	ID       string   `json:"id"`
	Message  *Message `json:"message,omitempty"`
	Parent   string   `json:"parent,omitempty"`
	Children []string `json:"children"`
}

type Message struct {
	ID         string                 `json:"id"`
	Author     Author                 `json:"author"`
	CreateTime float64                `json:"create_time,omitempty"`
	UpdateTime *float64               `json:"update_time,omitempty"`
	Content    Content                `json:"content"`
	EndTurn    *bool                  `json:"end_turn,omitempty"`
	Weight     float64                `json:"weight,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Recipient  string                 `json:"recipient,omitempty"`
}
