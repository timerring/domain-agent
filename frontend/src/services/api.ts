import axios from 'axios'

const API_BASE = '/api'

const api = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json',
  },
})

export interface ChatResponse {
  session_id: string
  message: string
  intent: string
  action: string
  data: any
  timestamp: string
}

export interface DomainResult {
  domain: string
  available: boolean
  signatures: string[]
  score: number
  price: string
}

export const sendMessage = async (
  message: string,
  sessionId?: string
): Promise<ChatResponse> => {
  const response = await api.post('/agent/chat', {
    message,
    session_id: sessionId || '',
  })
  return response.data
}

export const checkDomains = async (domains: string[]): Promise<DomainResult[]> => {
  const response = await api.post('/domains/check', { domains })
  return response.data.results
}

export const suggestDomains = async (
  keywords: string[],
  options?: {
    tlds?: string[]
    maxLen?: number
    minLen?: number
    count?: number
  }
): Promise<any[]> => {
  const response = await api.post('/domains/suggest', {
    keywords,
    ...options,
  })
  return response.data.suggestions
}

export const getSession = async (sessionId: string): Promise<any> => {
  const response = await api.get(`/agent/session/${sessionId}`)
  return response.data
}
