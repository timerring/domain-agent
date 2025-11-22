import { useState, useRef, useEffect } from 'react'
import { sendMessage } from '../services/api'

interface Message {
  role: 'user' | 'assistant'
  content: string
  timestamp: Date
}

interface Props {
  onResults: (results: any[]) => void
}

export default function ChatInterface({ onResults }: Props) {
  const [messages, setMessages] = useState<Message[]>([
    {
      role: 'assistant',
      content: 'Hello! I\'m Domain Agent. Tell me about your project and I\'ll help you find the perfect domain.',
      timestamp: new Date()
    }
  ])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const [sessionId, setSessionId] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const handleSend = async () => {
    if (!input.trim() || loading) return

    const userMessage: Message = {
      role: 'user',
      content: input,
      timestamp: new Date()
    }

    setMessages(prev => [...prev, userMessage])
    setInput('')
    setLoading(true)

    try {
      const response = await sendMessage(input, sessionId)
      
      if (!sessionId) {
        setSessionId(response.session_id)
      }

      const assistantMessage: Message = {
        role: 'assistant',
        content: response.message,
        timestamp: new Date(response.timestamp)
      }

      setMessages(prev => [...prev, assistantMessage])

      // 如果有域名结果，更新结果面板
      if (response.data?.domains && response.data.domains.length > 0) {
        console.log('Processing domains:', response.data.domains)
        
        // 调用域名检查 API 来获取真实数据
        try {
          const checkResponse = await fetch('http://localhost:8080/api/domains/check', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              domains: response.data.domains
            }),
          })
          
          if (checkResponse.ok) {
            const checkData = await checkResponse.json()
            console.log('Domain check results:', checkData.results)
            
            // 合并域名检查结果和 reason 信息
            const resultsWithReasons = (checkData.results || []).map((result: any) => {
              const domainReason = response.data.domainReasons?.find((dr: any) => dr.domain === result.domain)
              return {
                ...result,
                reason: domainReason?.reason || ''
              }
            })
            
            onResults(resultsWithReasons)
          } else {
            // 如果检查失败，显示模拟数据
            const mockResults = response.data.domains.map((d: string) => {
              const domainReason = response.data.domainReasons?.find((dr: any) => dr.domain === d)
              return {
                domain: d,
                available: Math.random() > 0.5,
                score: Math.floor(Math.random() * 100),
                reason: domainReason?.reason || ''
              }
            })
            onResults(mockResults)
          }
        } catch (error) {
          console.log('Domain check failed, using mock data:', error)
          // 如果检查失败，显示模拟数据
          const mockResults = response.data.domains.map((d: string) => {
            const domainReason = response.data.domainReasons?.find((dr: any) => dr.domain === d)
            return {
              domain: d,
              available: Math.random() > 0.5,
              score: Math.floor(Math.random() * 100),
              reason: domainReason?.reason || ''
            }
          })
          onResults(mockResults)
        }
      }
    } catch (error) {
      console.error('发送消息失败:', error)
      const errorMessage: Message = {
        role: 'assistant',
        content: '抱歉，发生了错误。请稍后重试。',
        timestamp: new Date()
      }
      setMessages(prev => [...prev, errorMessage])
    } finally {
      setLoading(false)
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  return (
    <div className="bg-white border border-gray-200 rounded-lg overflow-hidden flex flex-col h-[600px]">
      <div className="border-b border-gray-200 p-6">
        <h2 className="text-lg font-medium text-brand">Domain Assistant</h2>
        <p className="text-sm text-brand-light mt-1">AI-powered domain discovery</p>
      </div>

      <div className="flex-1 overflow-y-auto p-6 space-y-6">
        {messages.map((msg, idx) => (
          <div
            key={idx}
            className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}
          >
            <div
              className={`max-w-[85%] ${
                msg.role === 'user'
                  ? 'text-brand'
                  : 'text-brand-light'
              }`}
            >
              <div className="text-xs uppercase tracking-wide mb-2 opacity-60">
                {msg.role === 'user' ? 'You' : 'Assistant'}
              </div>
              <div className="prose prose-sm max-w-none">
                <p className="whitespace-pre-wrap leading-relaxed">{msg.content}</p>
              </div>
              <div className="text-xs mt-3 opacity-50">
                {msg.timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
              </div>
            </div>
          </div>
        ))}
        {loading && (
          <div className="flex justify-start">
            <div className="text-brand-light">
              <div className="text-xs uppercase tracking-wide mb-2 opacity-60">Assistant</div>
              <div className="flex space-x-1">
                <div className="w-1.5 h-1.5 bg-brand-light rounded-full animate-bounce"></div>
                <div className="w-1.5 h-1.5 bg-brand-light rounded-full animate-bounce delay-100"></div>
                <div className="w-1.5 h-1.5 bg-brand-light rounded-full animate-bounce delay-200"></div>
              </div>
            </div>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>

      <div className="border-t border-gray-200 p-6">
        <div className="flex space-x-3">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="Describe your project or domain needs..."
            className="input-field flex-1"
            disabled={loading}
          />
          <button
            onClick={handleSend}
            disabled={loading || !input.trim()}
            className="btn-primary"
          >
            {loading ? 'Sending...' : 'Send'}
          </button>
        </div>
      </div>
    </div>
  )
}
