import { useState } from 'react'
import ChatInterface from './components/ChatInterface'
import DomainResults from './components/DomainResults'

function App() {
  const [results, setResults] = useState<any[]>([])

  return (
    <main className="min-h-screen pt-12 md:pt-24 pb-4 px-4 md:py-40 md:pb-24 md:px-24 max-w-screen-xl mx-auto text-sm leading-6 md:text-lg md:leading-9 overflow-y-scroll">
      {/* Header */}
      <div id="#header">
        <div>
          <h1 className="text-3xl md:text-8xl leading-10 md:leading-[9rem] text-brand pb-3 md:pb-4">
            Domain Agent
          </h1>
          <p className="pb-5 md:pb-4 text-lg md:text-xl text-gray-600 max-w-2xl">
            AI-powered domain discovery platform. Find the perfect domain for your next big idea.
          </p>
          <ul className="pb-5 md:pb-4 text-sm md:text-base text-gray-500 max-w-xl space-y-2">
            <li className="flex items-center gap-2">
              <span className="text-brand">•</span>
              <span>Intelligent domain generation</span>
            </li>
            <li className="flex items-center gap-2">
              <span className="text-brand">•</span>
              <span>Real-time availability checking</span>
            </li>
            <li className="flex items-center gap-2">
              <span className="text-brand">•</span>
              <span>Creative suggestions powered by AI</span>
            </li>
          </ul>
        </div>
      </div>

      {/* Main Content */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 md:gap-16 mt-12 md:mt-20">
        <ChatInterface onResults={setResults} />
        <DomainResults results={results} />
      </div>
    </main>
  )
}

export default App
