package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Client struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
	Finish  string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func NewClient() *Client {
	return &Client{
		apiKey:  os.Getenv("VIBECODING_API_KEY"),
		baseURL: "https://vibecodingapi.ai/v1",
		model:   "gpt-4-gizmo-g-2fkFE8rbu",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) GenerateDomainIdeas(userInput string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("VIBECODING_API_KEY not set")
	}

	prompt := fmt.Sprintf(`你是一个专业的域名顾问。根据用户需求生成创意域名建议。

用户需求：%s

请生成5-10个域名建议，要求：
1. 简短易记（2-8个字符优先）
2. 有意义或有趣味性
3. 适合品牌使用
4. 包含不同的风格（正式、创意、技术感等）

请以JSON格式返回：
{
  "suggestions": [
    {
      "domain": "域名",
      "reason": "推荐理由",
      "style": "风格类型"
    }
  ],
  "summary": "总结建议"
}`, userInput)

	messages := []Message{
		{Role: "system", Content: "你是一个专业的域名顾问，擅长根据用户需求生成有创意的域名建议。"},
		{Role: "user", Content: prompt},
	}

	req := ChatRequest{
		Model:       c.model,
		Messages:    messages,
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	return c.chat(req)
}

// GenerateResponse 生成通用对话响应（不要求返回 JSON）
func (c *Client) GenerateResponse(userInput string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("VIBECODING_API_KEY not set")
	}

	messages := []Message{
		{Role: "system", Content: "你是 Domain Agent，一个专业友好的域名查询助手。你可以帮助用户查询域名可用性和生成创意域名建议。请用简短自然的语言回应用户。"},
		{Role: "user", Content: userInput},
	}

	req := ChatRequest{
		Model:       c.model,
		Messages:    messages,
		MaxTokens:   200,
		Temperature: 0.7,
	}

	return c.chat(req)
}

func (c *Client) AnalyzeUserIntent(userInput string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("VIBECODING_API_KEY not set")
	}

	prompt := fmt.Sprintf(`分析用户输入的意图，返回以下类型之一：
- "check_specific": 用户提供了具体的域名（如 google.com, abc.cn），想查询这些域名是否可用
- "generate_ideas": 用户想要域名创意建议，或者想要生成/推荐相关的域名
- "greeting": 问候语
- "general": 一般咨询

重要区分：
- 如果用户提供了完整域名格式（包含 .com/.cn/.ai 等），选择 "check_specific"
- 如果用户只提供了关键词或想法，想要生成域名建议，选择 "generate_ideas"
- 例如："查询 google.com" → check_specific
- 例如："查询有关 kitleaf 的域名" → generate_ideas
- 例如："我想要科技感的域名" → generate_ideas

用户输入：%s

只返回意图类型，不要其他内容。`, userInput)

	messages := []Message{
		{Role: "system", Content: "你是一个意图分析助手，专门分析用户的域名查询意图。请准确区分用户是想查询具体域名的可用性，还是想要生成域名建议。"},
		{Role: "user", Content: prompt},
	}

	req := ChatRequest{
		Model:       c.model,
		Messages:    messages,
		MaxTokens:   50,
		Temperature: 0.1,
	}

	return c.chat(req)
}

func (c *Client) chat(req ChatRequest) (string, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}
