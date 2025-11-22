package api

import (
	"net/http"
	"domain-agent/backend/internal/agent"
	"domain-agent/backend/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterAgentRoutes(r *gin.RouterGroup) {
	agentGroup := r.Group("/agent")
	{
		agentGroup.POST("/chat", handleChat)
		agentGroup.GET("/session/:id", handleGetSession)
		agentGroup.GET("/stream", handleStream)
	}
}

// handleChat 处理对话请求
func handleChat(c *gin.Context) {
	var req types.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果没有 session_id，创建新会话
	if req.SessionID == "" {
		req.SessionID = uuid.New().String()
	}

	// 处理消息
	response, err := agent.ProcessMessage(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// handleGetSession 获取会话信息
func handleGetSession(c *gin.Context) {
	sessionID := c.Param("id")
	
	session, err := agent.GetSession(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// handleStream WebSocket 流式响应
func handleStream(c *gin.Context) {
	agent.HandleWebSocket(c.Writer, c.Request)
}
