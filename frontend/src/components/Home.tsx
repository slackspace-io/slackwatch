import React, { useState, useEffect } from 'react';
import { Workload } from '../types';
import { api } from '../api';
import { WorkloadCard } from './WorkloadCard';
import { SystemInfoCard } from './SystemInfoCard';
import { Link } from 'react-router-dom';

export const Home: React.FC = () => {
  const [workloads, setWorkloads] = useState<Workload[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchWorkloads = async () => {
    setIsLoading(true);
    try {
      const data = await api.getAllWorkloads();
      setWorkloads(data);
      setError(null);
    } catch (err) {
      setError('Failed to fetch workloads');
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchWorkloads();
  }, []);

  if (isLoading) return <div>Loading workloads...</div>;

  if (error) return <div>Error: {error}</div>;

  if (workloads.length === 0) {
    return (
      <div>
        <div>No workloads found</div>
        <br />
        <br />
        <Link to="/refresh-all">Click to Refresh All</Link>
      </div>
    );
  }

  return (
    <div className="workloads-page">
      <SystemInfoCard workloads={workloads} />

      {workloads.map((workload, index) => (
        <WorkloadCard
          key={index}
          workload={workload}
          onUpdate={fetchWorkloads}
        />
      ))}
    </div>
  );
};
