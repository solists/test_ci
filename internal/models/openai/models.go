package openai

type GetQueryRequest struct {
	UserID   int64
	Messages []PromptMessage
}
type PromptMessage struct {
	Message string
}
type GetQueryResponse struct {
	Result string
}
