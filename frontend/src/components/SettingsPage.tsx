import React, { useState, useEffect } from 'react';
import { Settings } from '../types';
import { api } from '../api';

export const SettingsPage: React.FC = () => {
  const [settings, setSettings] = useState<Settings | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchSettings = async () => {
      try {
        const data = await api.getSettings();
        setSettings(data);
      } catch (err) {
        setError('Failed to fetch settings');
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchSettings();
  }, []);

  if (isLoading) return <div>Loading settings...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!settings) return <div>No settings found</div>;

  return (
    <div className="settings-page">
      <div className="settings-section">
        <div className="settings-section-header">System Settings</div>
        <div className="settings-item">
          <span className="settings-item-key">Schedule: </span>
          <span className="settings-item-value">{settings.system.schedule}</span>
        </div>
        <div className="settings-item">
          <span className="settings-item-key">Data Directory: </span>
          <span className="settings-item-value">{settings.system.data_dir}</span>
        </div>
        <div className="settings-item">
          <span className="settings-item-key">Run at Startup: </span>
          <span className="settings-item-value">{settings.system.run_at_startup.toString()}</span>
        </div>
      </div>

      {settings.gitops && settings.gitops.length > 0 && (
        <div className="settings-section">
          <div className="settings-section-header">Gitops Settings</div>
          {settings.gitops.map((gitops, index) => (
            <div key={index}>
              <div className="settings-item">
                <span className="settings-item-key">Name: </span>
                <span className="settings-item-value">{gitops.name}</span>
              </div>
              <div className="settings-item">
                <span className="settings-item-key">Repository URL: </span>
                <span className="settings-item-value">{gitops.repository_url}</span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
