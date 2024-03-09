// src/config/index.ts
import type { Config } from './types'

let config: Config | undefined;
const env = import.meta.env.VITE_APP_ENV;

if (env === 'development') {
  const { default: devConfig } = await import('./development');
  config = devConfig;
} else if (env === 'production') {
  const { default: prodConfig } = await import('./production');
  config = prodConfig;
} else {
  console.error('Unknown environment');
}

export default config;
