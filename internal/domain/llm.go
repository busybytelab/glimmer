package domain

// Usage represents the token usage and cost information for an LLM request
type Usage struct {
	LlmModelName     string  `json:"llmModelName"`
	CacheHit         bool    `json:"cacheHit"`
	Cost             float64 `json:"cost"`
	PromptTokens     int     `json:"promptTokens"`
	CompletionTokens int     `json:"completionTokens"`
	TotalTokens      int     `json:"totalTokens"`
}
