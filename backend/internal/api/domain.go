package api

import (
	"net/http"
	"domain-agent/backend/internal/scanner"
	"domain-agent/backend/internal/types"
	"github.com/gin-gonic/gin"
)

func RegisterDomainRoutes(r *gin.RouterGroup) {
	domainGroup := r.Group("/domains")
	{
		domainGroup.POST("/check", handleCheckDomains)
		domainGroup.POST("/suggest", handleSuggestDomains)
	}
}

// handleCheckDomains 批量检查域名可用性
func handleCheckDomains(c *gin.Context) {
	var req types.CheckDomainsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := scanner.CheckDomains(req.Domains)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"total":   len(results),
	})
}

// handleSuggestDomains 根据关键词生成域名建议
func handleSuggestDomains(c *gin.Context) {
	var req types.SuggestDomainsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	suggestions := scanner.GenerateSuggestions(req)

	c.JSON(http.StatusOK, gin.H{
		"suggestions": suggestions,
		"count":       len(suggestions),
	})
}
