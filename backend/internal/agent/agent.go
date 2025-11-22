package agent

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"domain-agent/backend/internal/llm"
	"domain-agent/backend/internal/types"
)

var (
	sessions  = make(map[string]*types.Session)
	mu        sync.RWMutex
	llmClient *llm.Client
)

func init() {
	llmClient = llm.NewClient()
}

// ProcessMessage 处理用户消息
func ProcessMessage(req types.ChatRequest) (*types.ChatResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	// 获取或创建会话
	session, exists := sessions[req.SessionID]
	if !exists {
		session = &types.Session{
			ID:        req.SessionID,
			Messages:  []types.Message{},
			Context:   make(map[string]interface{}),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		sessions[req.SessionID] = session
	}

	// 添加用户消息
	session.Messages = append(session.Messages, types.Message{
		Role:      "user",
		Content:   req.Message,
		Timestamp: time.Now(),
	})

	// 分析意图 - 优先使用 LLM，如果失败则回退到规则匹配
	intent := analyzeIntentWithLLM(req.Message)
	if intent == "" {
		intent = analyzeIntent(req.Message)
	}

	// 生成响应
	response := generateResponse(intent, req.Message, session)

	// 添加助手消息
	session.Messages = append(session.Messages, types.Message{
		Role:      "assistant",
		Content:   response.Message,
		Timestamp: time.Now(),
	})

	session.UpdatedAt = time.Now()

	return response, nil
}

// GetSession 获取会话
func GetSession(sessionID string) (*types.Session, error) {
	mu.RLock()
	defer mu.RUnlock()

	session, exists := sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}

	return session, nil
}

// analyzeIntentWithLLM 使用 LLM 分析意图
func analyzeIntentWithLLM(message string) string {
	intent, err := llmClient.AnalyzeUserIntent(message)
	if err != nil {
		fmt.Printf("LLM intent analysis failed: %v\n", err)
		return ""
	}

	// 清理返回的意图
	intent = strings.TrimSpace(intent)
	intent = strings.ToLower(strings.Trim(intent, `"'`))

	// 验证意图是否有效
	validIntents := map[string]bool{
		"check_specific": true,
		"generate_ideas": true,
		"greeting":       true,
		"general":        true,
	}

	if validIntents[intent] {
		return intent
	}

	return ""
}

// generateResponseWithLLM 使用 LLM 生成智能响应
func generateResponseWithLLM(message string) string {
	response, err := llmClient.GenerateDomainIdeas(message)
	if err != nil {
		fmt.Printf("LLM response generation failed: %v\n", err)
		return "我来为你生成一些创意域名..."
	}

	return response
}

// analyzeIntent 使用 AI 分析用户意图
func analyzeIntent(message string) string {
	// 使用 LLM 进行智能意图分析
	intent, err := llmClient.AnalyzeUserIntent(message)
	if err != nil {
		fmt.Printf("AI intent analysis failed: %v, falling back to keyword matching\n", err)
		// 如果 AI 分析失败，回退到简单的关键词匹配
		return analyzeIntentWithKeywords(message)
	}

	// 清理 AI 返回的结果（去除多余空格和引号）
	intent = strings.TrimSpace(intent)
	intent = strings.Trim(intent, "\"'")
	intent = strings.ToLower(intent)

	// 验证返回的意图是否有效
	validIntents := map[string]bool{
		"check_specific": true,
		"generate_ideas": true,
		"greeting":       true,
		"general":        true,
	}

	if validIntents[intent] {
		fmt.Printf("AI analyzed intent: %s for message: %s\n", intent, message)
		return intent
	}

	// 如果 AI 返回了无效的意图，使用关键词匹配作为后备
	fmt.Printf("AI returned invalid intent: %s, using keyword matching\n", intent)
	return analyzeIntentWithKeywords(message)
}

// analyzeIntentWithKeywords 使用关键词匹配分析意图（作为后备方案）
func analyzeIntentWithKeywords(message string) string {
	msgLower := strings.ToLower(message)

	// 检查是否包含具体域名
	if strings.Contains(msgLower, ".com") ||
		strings.Contains(msgLower, ".cn") ||
		strings.Contains(msgLower, ".ai") ||
		strings.Contains(msgLower, ".io") {
		return "check_specific"
	}

	// 检查是否是创意请求
	keywords := []string{"想要", "帮我", "推荐", "建议", "找", "查找", "需要", "suggest", "need", "recommend"}
	for _, keyword := range keywords {
		if strings.Contains(msgLower, keyword) {
			return "generate_ideas"
		}
	}

	// 检查是否是问候
	greetings := []string{"你好", "hi", "hello", "嗨"}
	for _, greeting := range greetings {
		if strings.Contains(msgLower, greeting) {
			return "greeting"
		}
	}

	return "general"
}

// generateResponse 生成响应
func generateResponse(intent, message string, session *types.Session) *types.ChatResponse {
	response := &types.ChatResponse{
		SessionID: session.ID,
		Intent:    intent,
		Timestamp: time.Now(),
		Data:      make(map[string]interface{}),
	}

	switch intent {
	case "greeting":
		// 使用 AI 生成友好的问候响应
		llmResponse, err := llmClient.GenerateResponse(message)
		if err != nil {
			response.Message = "你好！我是 Domain Agent，专业的域名查询助手。我可以帮你查询域名和生成创意域名建议。"
		} else {
			response.Message = llmResponse
		}
		response.Action = "none"

	case "check_specific":
		// 提取域名
		domains := extractDomains(message)
		response.Data["domains"] = domains
		response.Action = "check_domains"

		// 使用 AI 生成自然的响应
		llmResponse, err := llmClient.GenerateResponse(fmt.Sprintf("用户想查询这些域名的可用性：%v", domains))
		if err != nil {
			response.Message = "我来帮你查询这些域名的可用性..."
		} else {
			response.Message = llmResponse
		}

	case "generate_ideas":
		// 使用 LLM 生成智能响应
		llmResponse := generateResponseWithLLM(message)

		// 尝试解析 LLM 返回的 JSON
		var suggestions struct {
			Suggestions []struct {
				Domain string `json:"domain"`
				Reason string `json:"reason"`
				Style  string `json:"style"`
			} `json:"suggestions"`
			Summary string `json:"summary"`
		}

		if err := json.Unmarshal([]byte(llmResponse), &suggestions); err == nil && len(suggestions.Suggestions) > 0 {
			// 成功解析 JSON，格式化输出
			var domains []string
			var domainReasons []map[string]string
			var output strings.Builder

			output.WriteString(fmt.Sprintf("🤖 根据你的需求，我为你生成了 %d 个创意域名：\n\n", len(suggestions.Suggestions)))

			for i, s := range suggestions.Suggestions {
				output.WriteString(fmt.Sprintf("%d. **%s** - %s\n", i+1, s.Domain, s.Reason))
				domains = append(domains, s.Domain)
				domainReasons = append(domainReasons, map[string]string{
					"domain": s.Domain,
					"reason": s.Reason,
				})
			}

			if suggestions.Summary != "" {
				output.WriteString(fmt.Sprintf("\n💡 **总结**: %s", suggestions.Summary))
			}

			response.Message = output.String()
			response.Action = "generate_suggestions"
			response.Data["domains"] = domains
			response.Data["domainReasons"] = domainReasons
			response.Data["keywords"] = extractKeywords(message)
		} else {
			// 尝试提取 JSON 代码块中的内容
			jsonContent := llmResponse
			if strings.Contains(llmResponse, "```json") {
				// 提取 ```json 和 ``` 之间的内容
				start := strings.Index(llmResponse, "```json")
				if start != -1 {
					start += 7 // 跳过 "```json"
					end := strings.Index(llmResponse[start:], "```")
					if end != -1 {
						jsonContent = llmResponse[start : start+end]
					}
				}
			}

			// 再次尝试解析
			if err := json.Unmarshal([]byte(jsonContent), &suggestions); err == nil && len(suggestions.Suggestions) > 0 {
				// 成功解析 JSON，格式化输出
				var domains []string
				var domainReasons []map[string]string
				var output strings.Builder

				output.WriteString(fmt.Sprintf("🤖 根据你的需求，我为你生成了 %d 个创意域名：\n\n", len(suggestions.Suggestions)))

				for i, s := range suggestions.Suggestions {
					output.WriteString(fmt.Sprintf("%d. **%s** - %s\n", i+1, s.Domain, s.Reason))
					domains = append(domains, s.Domain)
					domainReasons = append(domainReasons, map[string]string{
						"domain": s.Domain,
						"reason": s.Reason,
					})
				}

				if suggestions.Summary != "" {
					output.WriteString(fmt.Sprintf("\n💡 **总结**: %s", suggestions.Summary))
				}

				response.Message = output.String()
				response.Action = "generate_suggestions"
				response.Data["domains"] = domains
				response.Data["domainReasons"] = domainReasons
				response.Data["keywords"] = extractKeywords(message)
			} else {
				// JSON 解析仍然失败，直接返回 LLM 响应
				response.Message = llmResponse
				response.Action = "generate_suggestions"
				response.Data["keywords"] = extractKeywords(message)
			}
		}

	default:
		// 使用 AI 生成通用响应
		llmResponse, err := llmClient.GenerateResponse(message)
		if err != nil {
			response.Message = "我理解你想查询域名。你可以直接告诉我域名，或描述你的需求。"
		} else {
			response.Message = llmResponse
		}
		response.Action = "clarify"
	}

	return response
}

// extractDomains 从消息中提取域名
func extractDomains(message string) []string {
	domains := []string{}
	words := strings.Fields(message)

	for _, word := range words {
		word = strings.TrimSpace(word)
		if strings.Contains(word, ".") && !strings.HasPrefix(word, ".") {
			domains = append(domains, word)
		}
	}

	return domains
}

// extractKeywords 从消息中提取关键词
func extractKeywords(message string) []string {
	// 简单的关键词提取（后续可以用 NLP）
	stopWords := map[string]bool{
		"我": true, "想要": true, "一个": true, "的": true, "帮我": true,
		"找": true, "查找": true, "推荐": true, "建议": true, "域名": true,
	}

	words := strings.Fields(message)
	keywords := []string{}

	for _, word := range words {
		word = strings.TrimSpace(word)
		if len(word) > 0 && !stopWords[word] {
			keywords = append(keywords, word)
		}
	}

	return keywords
}
