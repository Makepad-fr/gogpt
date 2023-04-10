package gogpt

type Author struct {
	Role     string                 `json:"role"`
	Metadata map[string]interface{} `json:"metadata"`
}

type Content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type Message struct {
	ID         string                 `json:"id"`
	Author     Author                 `json:"author"`
	CreateTime float64                `json:"create_time"`
	Content    Content                `json:"content"`
	EndTurn    bool                   `json:"end_turn,omitempty"`
	Weight     float64                `json:"weight"`
	Metadata   map[string]interface{} `json:"metadata"`
	Recipient  string                 `json:"recipient"`
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
