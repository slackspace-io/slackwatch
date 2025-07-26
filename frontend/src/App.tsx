import React from 'react';
import { Routes, Route, Link, Navigate } from 'react-router-dom';
import { Home } from './components/Home';
import { SettingsPage } from './components/SettingsPage';
import { RefreshAll } from './components/RefreshAll';

const NotFound: React.FC = () => (
  <div>
    <h1>404 - Not Found</h1>
    <p>The page you are looking for does not exist.</p>
    <Link to="/">Go back to Home</Link>
  </div>
);

const App: React.FC = () => {
  return (
    <div className="app-container">
      <header className="app-header">
        <h1>Slackwatch</h1>
        <nav>
          <ul>
            <li><Link to="/">Home</Link></li>
            <li><Link to="/settings">Settings</Link></li>
          </ul>
        </nav>
      </header>

      <main className="app-content">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/settings" element={<SettingsPage />} />
          <Route path="/refresh-all" element={<RefreshAll />} />
          <Route path="/404" element={<NotFound />} />
          <Route path="*" element={<Navigate to="/404" replace />} />
        </Routes>
      </main>

      <footer className="app-footer">
        <p>Slackwatch - Version 0.1.0</p>
      </footer>
    </div>
  );
};

export default App;
