import axios from 'axios';
import { Workload, Settings } from './types';

const API_URL = '/api';

export const api = {
  // Workloads
  getAllWorkloads: async (): Promise<Workload[]> => {
    const response = await axios.get(`${API_URL}/workloads`);
    return response.data;
  },

  updateWorkload: async (workload: Workload): Promise<void> => {
    await axios.post(`${API_URL}/workloads/update`, workload);
  },

  upgradeWorkload: async (workload: Workload): Promise<void> => {
    await axios.post(`${API_URL}/workloads/upgrade`, workload);
  },

  refreshAll: async (): Promise<void> => {
    await axios.post(`${API_URL}/workloads/refresh-all`);
  },

  // Settings
  getSettings: async (): Promise<Settings> => {
    const response = await axios.get(`${API_URL}/settings`);
    return response.data;
  },

  getNextScheduleTime: async (): Promise<string> => {
    const response = await axios.get(`${API_URL}/settings/next-schedule-time`);
    // Ensure we're returning a string, not an object
    if (typeof response.data === 'object') {
      console.error('Received object instead of string for next schedule time:', response.data);
      return JSON.stringify(response.data);
    }
    return response.data;
  }
};
