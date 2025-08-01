import React from 'react';

function App() {
  return (
    <div className="App">
      <header>
        <h1>CICD Status Notifier Bot Dashboard</h1>
        <p>Welcome to the CI/CD monitoring dashboard!</p>
      </header>
      <main>
        <section>
          <h2>System Status</h2>
          <div style={{ padding: '20px', border: '1px solid #ccc', borderRadius: '8px', margin: '20px 0' }}>
            <p>✅ Backend API: Running</p>
            <p>✅ Database: Connected</p>
            <p>✅ Notifications: Active</p>
          </div>
        </section>
        <section>
          <h2>Recent Activity</h2>
          <div style={{ padding: '20px', border: '1px solid #ccc', borderRadius: '8px', margin: '20px 0' }}>
            <p>No recent activities to display.</p>
          </div>
        </section>
      </main>
    </div>
  );
}

export default App;