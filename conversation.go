package gogpt

type ConversationsResponse struct {
	Items                   []Conversation `json:"items"`
	Total                   int            `json:"total"`
	Limit                   int            `json:"limit"`
	Offset                  int            `json:"offset"`
	HasMissingConversations bool           `json:"has_missing_conversations"`
}

type Conversation struct {
	idBasedSetItem
	ID         string `json:"id"`
	Title      string `json:"title"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

func (c Conversation) getId() string {
	return c.ID
}
