import React from 'react'

interface DomainResult {
  domain: string
  available: boolean
  score?: number
  signatures?: string[]
  reason?: string
}

interface Props {
  results: DomainResult[]
}

export default function DomainResults({ results }: Props) {
  const [filter, setFilter] = React.useState<'all' | 'available' | 'taken'>('all')

  const filteredResults = results.filter(result => {
    if (filter === 'all') return true
    if (filter === 'available') return result.available
    if (filter === 'taken') return !result.available
    return true
  })

  const availableCount = results.filter(r => r.available).length
  const takenCount = results.filter(r => !r.available).length

  return (
    <div className="bg-white border border-gray-200 rounded-lg overflow-hidden h-[600px] flex flex-col">
      <div className="border-b border-gray-200 p-6">
        <h2 className="text-lg font-medium text-brand">Domain Results</h2>
        <p className="text-sm text-brand-light mt-1">
          {results.length > 0 ? `${results.length} domains found` : 'No results yet'}
        </p>
        
        {/* Filter Tabs */}
        {results.length > 0 && (
          <div className="mt-4 flex gap-2">
            <button
              onClick={() => setFilter('all')}
              className={`px-4 py-2 text-sm font-medium rounded-lg transition-colors ${
                filter === 'all'
                  ? 'bg-brand text-white'
                  : 'bg-gray-100 text-brand-light hover:bg-gray-200'
              }`}
            >
              All ({results.length})
            </button>
            <button
              onClick={() => setFilter('available')}
              className={`px-4 py-2 text-sm font-medium rounded-lg transition-colors ${
                filter === 'available'
                  ? 'bg-green-600 text-white'
                  : 'bg-gray-100 text-brand-light hover:bg-gray-200'
              }`}
            >
              Available ({availableCount})
            </button>
            <button
              onClick={() => setFilter('taken')}
              className={`px-4 py-2 text-sm font-medium rounded-lg transition-colors ${
                filter === 'taken'
                  ? 'bg-gray-600 text-white'
                  : 'bg-gray-100 text-brand-light hover:bg-gray-200'
              }`}
            >
              Taken ({takenCount})
            </button>
          </div>
        )}
      </div>

      <div className="flex-1 overflow-y-auto p-6">
        {results.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full text-brand-light">
            <div className="w-16 h-16 border-2 border-brand-light rounded-full flex items-center justify-center mb-4">
              <svg
                className="w-8 h-8"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={1.5}
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
            </div>
            <p className="text-sm font-medium">No domains yet</p>
            <p className="text-xs mt-2 text-center">Start a conversation to discover domains</p>
          </div>
        ) : (
          <div className="space-y-4">
            {filteredResults.map((result, idx) => (
              <div
                key={idx}
                className={`p-4 border rounded-lg transition-all hover:border-brand ${
                  result.available
                    ? 'border-green-200 bg-green-50/30'
                    : 'border-gray-200 bg-gray-50/30'
                }`}
              >
                <div className="flex items-center justify-between">
                  <div className="flex-1">
                    <div className="group relative">
                      <h3 className={`font-mono text-sm font-medium transition-colors ${
                        result.available ? 'text-brand' : 'text-gray-400'
                      }`}>
                        {result.domain}
                      </h3>
                      {result.reason && (
                        <div className="absolute bottom-full left-0 mb-2 hidden group-hover:block z-10">
                          <div className="bg-gray-900 text-white text-xs rounded-lg p-3 max-w-xs shadow-lg">
                            <div className="font-medium mb-1">Reason:</div>
                            <div>{result.reason}</div>
                            <div className="absolute top-full left-4 -mt-1">
                              <div className="border-4 border-transparent border-t-gray-900"></div>
                            </div>
                          </div>
                        </div>
                      )}
                    </div>
                  </div>
                  <div className="ml-4">
                    {result.available ? (
                      <span className="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium text-green-700 bg-green-100">
                        Available
                      </span>
                    ) : (
                      <div className="text-right">
                        <div className="text-xs font-medium text-brand-light mb-1">Taken - Verified by:</div>
                        <div className="flex flex-col gap-1">
                          {result.signatures && result.signatures.length > 0 ? (
                            result.signatures.map((sig, i) => {
                              let label = sig
                              if (sig === 'DNS_NS' || sig === 'DNS_A' || sig === 'DNS_MX') {
                                label = 'DNS records (NS, A, MX)'
                              } else if (sig === 'WHOIS') {
                                label = 'WHOIS information'
                              } else if (sig === 'SSL') {
                                label = 'SSL certificate'
                              }
                              return (
                                <span
                                  key={i}
                                  className="text-xs px-2 py-0.5 bg-gray-100 text-brand-light rounded"
                                >
                                  {label}
                                </span>
                              )
                            })
                          ) : (
                            <span className="text-xs px-2 py-0.5 bg-gray-100 text-brand-light rounded">
                              Domain registered
                            </span>
                          )}
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
