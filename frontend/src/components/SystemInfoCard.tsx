import React, { useState, useEffect } from 'react';
import { Workload } from '../types';
import { api } from '../api';
import { Link } from 'react-router-dom';

interface SystemInfoCardProps {
  workloads: Workload[];
}

export const SystemInfoCard: React.FC<SystemInfoCardProps> = ({ workloads }) => {
  const [nextSchedule, setNextSchedule] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchNextSchedule = async () => {
      try {
        const time = await api.getNextScheduleTime();
        // Make sure we're getting a string, not an object
        if (typeof time === 'object') {
          console.error('Received object instead of string for next schedule time:', time);
          setNextSchedule(JSON.stringify(time));
        } else {
          setNextSchedule(time);
        }
      } catch (err) {
        setError('Failed to fetch next schedule time');
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchNextSchedule();
  }, []);

  return (
    <div className="system-info-card">
      <div className="system-info">System Info</div>
      <div className="system-info-entry">Watched Workloads: {workloads.length}</div>

      {isLoading ? (
        <div className="system-info-entry">Loading next run time...</div>
      ) : error ? (
        <div className="system-info-entry">Error: {error}</div>
      ) : (
        <>
          <div className="system-info-entry">Next Run: {nextSchedule}</div>
          <div className="system-info-entry">
            <Link to="/refresh-all">Click to Run Now</Link>
          </div>
        </>
      )}
    </div>
  );
};
