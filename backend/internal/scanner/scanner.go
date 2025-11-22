package scanner

import (
	"crypto/tls"
	"domain-agent/backend/internal/types"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/likexian/whois"
)

// CheckDomains 批量检查域名可用性
func CheckDomains(domains []string) ([]types.DomainResult, error) {
	results := make([]types.DomainResult, 0, len(domains))
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 并发检查（限制并发数）
	semaphore := make(chan struct{}, 10)

	for _, domain := range domains {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := checkSingleDomain(d)

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(domain)
	}

	wg.Wait()
	return results, nil
}

// checkSingleDomain 检查单个域名（移植自 domain-scanner）
func checkSingleDomain(domain string) types.DomainResult {
	result := types.DomainResult{
		Domain:     domain,
		Available:  false, // 默认为不可用，保守策略
		Signatures: []string{},
		Score:      calculateScore(domain),
		Price:      "standard",
	}

	// 检查域名签名（DNS、WHOIS、SSL）
	signatures := checkDomainSignatures(domain)
	result.Signatures = signatures

	// 如果有任何签名，域名已被注册
	if len(signatures) > 0 {
		result.Available = false
		return result
	}

	// 最终 WHOIS 检查可用性
	available, err := checkWHOISAvailability(domain)
	if err != nil {
		fmt.Printf("WHOIS check error for %s: %v\n", domain, err)
		result.Available = false
		return result
	}

	result.Available = available
	return result
}

// checkDomainSignatures 检查域名签名（移植自 domain-scanner）
func checkDomainSignatures(domain string) []string {
	var signatures []string

	// 1. 检查 DNS NS 记录
	nsRecords, err := net.LookupNS(domain)
	if err == nil && len(nsRecords) > 0 {
		signatures = append(signatures, "DNS_NS")
	}

	// 2. 检查 DNS A 记录
	ipRecords, err := net.LookupIP(domain)
	if err == nil && len(ipRecords) > 0 {
		signatures = append(signatures, "DNS_A")
	}

	// 3. 检查 DNS MX 记录
	mxRecords, err := net.LookupMX(domain)
	if err == nil && len(mxRecords) > 0 {
		signatures = append(signatures, "DNS_MX")
	}

	// 4. 检查 WHOIS 信息
	whoisResult, err := whois.Whois(domain)
	if err == nil && whoisResult != "" {
		resultLower := strings.ToLower(whoisResult)
		// 检查已注册指标
		registeredIndicators := []string{
			"registrar:", "creation date:", "expiration date:",
			"updated date:", "name server:", "domain status:",
		}
		for _, indicator := range registeredIndicators {
			if strings.Contains(resultLower, indicator) {
				signatures = append(signatures, "WHOIS")
				break
			}
		}
	}

	// 5. 检查 SSL 证书
	conn, err := tls.DialWithDialer(&net.Dialer{
		Timeout: 5 * time.Second,
	}, "tcp", domain+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err == nil {
		defer conn.Close()
		state := conn.ConnectionState()
		if len(state.PeerCertificates) > 0 {
			signatures = append(signatures, "SSL")
		}
	}

	return signatures
}

// checkWHOISAvailability 通过 WHOIS 检查域名可用性（移植自 domain-scanner）
func checkWHOISAvailability(domain string) (bool, error) {
	result, err := whois.Whois(domain)
	if err != nil {
		return false, err
	}

	if result == "" {
		return false, nil
	}

	resultLower := strings.ToLower(result)

	// 检查可用指标
	availableIndicators := []string{
		"status: free",
		"not found",
		"no match",
		"status: available",
		"no data found",
		"is available",
		"not registered",
		"no entries found",
	}

	for _, indicator := range availableIndicators {
		if strings.Contains(resultLower, indicator) {
			return true, nil
		}
	}

	// 检查不可用指标
	unavailableIndicators := []string{
		"registrar:",
		"creation date:",
		"expiration date:",
		"domain status:",
		"name server:",
	}

	for _, indicator := range unavailableIndicators {
		if strings.Contains(resultLower, indicator) {
			return false, nil
		}
	}

	// 保守策略：无法确定时默认为不可用
	return false, nil
}

// calculateScore 计算域名评分
func calculateScore(domain string) float64 {
	name := strings.Split(domain, ".")[0]
	score := 100.0

	// 长度评分（越短越好）
	if len(name) <= 3 {
		score += 50
	} else if len(name) <= 5 {
		score += 30
	} else if len(name) <= 8 {
		score += 10
	} else {
		score -= float64(len(name)-8) * 2
	}

	// 易读性评分（无数字混合更好）
	hasNumber := false
	for _, c := range name {
		if c >= '0' && c <= '9' {
			hasNumber = true
			break
		}
	}
	if !hasNumber {
		score += 20
	}

	return score
}

// GenerateSuggestions 生成域名建议
func GenerateSuggestions(req types.SuggestDomainsRequest) []types.DomainSuggestion {
	suggestions := []types.DomainSuggestion{}

	// 默认 TLDs
	tlds := req.TLDs
	if len(tlds) == 0 {
		tlds = []string{".com", ".cn", ".ai", ".io", ".tech"}
	}

	// 默认长度限制
	maxLen := req.MaxLen
	if maxLen == 0 {
		maxLen = 8
	}
	minLen := req.MinLen
	if minLen == 0 {
		minLen = 3
	}

	// 生成策略
	for _, keyword := range req.Keywords {
		// 1. 直接使用关键词
		for _, tld := range tlds {
			domain := keyword + tld
			if len(keyword) >= minLen && len(keyword) <= maxLen {
				suggestions = append(suggestions, types.DomainSuggestion{
					Domain:       domain,
					Score:        calculateScore(domain),
					Reason:       "直接使用关键词",
					Length:       len(keyword),
					Memorability: 0.9,
				})
			}
		}

		// 2. 关键词缩写
		if len(keyword) > 3 {
			abbr := generateAbbreviation(keyword)
			for _, tld := range tlds {
				domain := abbr + tld
				suggestions = append(suggestions, types.DomainSuggestion{
					Domain:       domain,
					Score:        calculateScore(domain),
					Reason:       fmt.Sprintf("'%s' 的缩写", keyword),
					Length:       len(abbr),
					Memorability: 0.7,
				})
			}
		}
	}

	// 3. 组合关键词
	if len(req.Keywords) >= 2 {
		combined := strings.Join(req.Keywords[:2], "")
		if len(combined) <= maxLen {
			for _, tld := range tlds {
				domain := combined + tld
				suggestions = append(suggestions, types.DomainSuggestion{
					Domain:       domain,
					Score:        calculateScore(domain),
					Reason:       "关键词组合",
					Length:       len(combined),
					Memorability: 0.8,
				})
			}
		}
	}

	// 限制返回数量
	count := req.Count
	if count == 0 {
		count = 20
	}
	if len(suggestions) > count {
		suggestions = suggestions[:count]
	}

	return suggestions
}

// generateAbbreviation 生成缩写
func generateAbbreviation(word string) string {
	if len(word) <= 3 {
		return word
	}

	// 简单策略：取首字母 + 部分字母
	abbr := string(word[0])
	if len(word) > 3 {
		abbr += string(word[len(word)/2])
	}
	if len(word) > 5 {
		abbr += string(word[len(word)-1])
	}

	return abbr
}
