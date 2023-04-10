package gogpt

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
