package types

import "time"

// ChatRequest 对话请求
type ChatRequest struct {
	SessionID string                 `json:"session_id"`
	Message   string                 `json:"message" binding:"required"`
	Context   map[string]interface{} `json:"context"`
}

// ChatResponse 对话响应
type ChatResponse struct {
	SessionID string                 `json:"session_id"`
	Message   string                 `json:"message"`
	Intent    string                 `json:"intent"`
	Action    string                 `json:"action"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// Session 会话信息
type Session struct {
	ID        string                 `json:"id"`
	Messages  []Message              `json:"messages"`
	Context   map[string]interface{} `json:"context"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// Message 消息
type Message struct {
	Role      string    `json:"role"` // user, assistant, system
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// CheckDomainsRequest 检查域名请求
type CheckDomainsRequest struct {
	Domains []string `json:"domains" binding:"required"`
}

// DomainResult 域名检查结果
type DomainResult struct {
	Domain     string   `json:"domain"`
	Available  bool     `json:"available"`
	Signatures []string `json:"signatures"`
	Score      float64  `json:"score"`
	Price      string   `json:"price"`
}

// SuggestDomainsRequest 域名建议请求
type SuggestDomainsRequest struct {
	Keywords []string `json:"keywords" binding:"required"`
	TLDs     []string `json:"tlds"`
	MaxLen   int      `json:"max_len"`
	MinLen   int      `json:"min_len"`
	Count    int      `json:"count"`
}

// DomainSuggestion 域名建议
type DomainSuggestion struct {
	Domain       string  `json:"domain"`
	Score        float64 `json:"score"`
	Reason       string  `json:"reason"`
	Length       int     `json:"length"`
	Memorability float64 `json:"memorability"`
}
