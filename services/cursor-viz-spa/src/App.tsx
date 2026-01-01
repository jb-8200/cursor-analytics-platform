import { useState, useEffect } from 'react'
import './App.css'

interface HealthStatus {
  cursorSim: boolean;
  analyticsCore: boolean;
}

function App() {
  const [health, setHealth] = useState<HealthStatus>({
    cursorSim: false,
    analyticsCore: false,
  })

  useEffect(() => {
    // Check cursor-sim health
    fetch('http://localhost:8080/v1/health')
      .then(() => setHealth(prev => ({ ...prev, cursorSim: true })))
      .catch(() => setHealth(prev => ({ ...prev, cursorSim: false })))

    // Check analytics-core health (GraphQL)
    fetch('http://localhost:4000/graphql', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ query: '{ health { status } }' }),
    })
      .then(() => setHealth(prev => ({ ...prev, analyticsCore: true })))
      .catch(() => setHealth(prev => ({ ...prev, analyticsCore: false })))
  }, [])

  return (
    <div className="app">
      <h1>Cursor Analytics Dashboard</h1>
      <p className="version">Version: 0.0.1-p0 (Scaffolding)</p>

      <div className="status-grid">
        <div className="status-card">
          <h2>cursor-sim</h2>
          <div className={`status ${health.cursorSim ? 'healthy' : 'unavailable'}`}>
            {health.cursorSim ? '✓ Running' : '✗ Unavailable'}
          </div>
          <p className="port">Port: 8080</p>
        </div>

        <div className="status-card">
          <h2>cursor-analytics-core</h2>
          <div className={`status ${health.analyticsCore ? 'healthy' : 'unavailable'}`}>
            {health.analyticsCore ? '✓ Running' : '✗ Unavailable'}
          </div>
          <p className="port">Port: 4000</p>
        </div>

        <div className="status-card">
          <h2>cursor-viz-spa</h2>
          <div className="status healthy">✓ Running</div>
          <p className="port">Port: 3000</p>
        </div>
      </div>

      <div className="notice">
        <h3>P0 Scaffolding Complete</h3>
        <p>All services are scaffolded but not yet implemented.</p>
        <p>Next step: Implement features following SPEC.md files.</p>
      </div>
    </div>
  )
}

export default App
