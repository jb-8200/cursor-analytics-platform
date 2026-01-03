/**
 * Configuration module for cursor-analytics-core
 * Manages environment variables and application settings
 */

export interface Config {
  // Server configuration
  port: number;
  nodeEnv: string;

  // Database configuration
  databaseUrl: string;

  // cursor-sim integration
  simulatorUrl: string;
  pollIntervalMs: number;

  // Logging
  logLevel: string;
  logFormat: string;

  // Performance
  queryComplexityLimit: number;
  dataloaderBatchSize: number;
}

function getEnvVar(key: string, defaultValue?: string): string {
  const value = process.env[key];
  if (value === undefined) {
    if (defaultValue !== undefined) {
      return defaultValue;
    }
    throw new Error(`Missing required environment variable: ${key}`);
  }
  return value;
}

function getEnvVarAsNumber(key: string, defaultValue?: number): number {
  const value = process.env[key];
  if (value === undefined) {
    if (defaultValue !== undefined) {
      return defaultValue;
    }
    throw new Error(`Missing required environment variable: ${key}`);
  }
  const parsed = parseInt(value, 10);
  if (isNaN(parsed)) {
    throw new Error(`Environment variable ${key} must be a valid number`);
  }
  return parsed;
}

export function loadConfig(): Config {
  return {
    // Server configuration
    port: getEnvVarAsNumber('PORT', 4000),
    nodeEnv: getEnvVar('NODE_ENV', 'development'),

    // Database configuration
    databaseUrl: getEnvVar(
      'DATABASE_URL',
      'postgresql://cursor:cursor_dev@localhost:5432/cursor_analytics'
    ),

    // cursor-sim integration
    simulatorUrl: getEnvVar('SIMULATOR_URL', 'http://localhost:8080'),
    pollIntervalMs: getEnvVarAsNumber('POLL_INTERVAL_MS', 60000),

    // Logging
    logLevel: getEnvVar('LOG_LEVEL', 'info'),
    logFormat: getEnvVar('LOG_FORMAT', 'json'),

    // Performance
    queryComplexityLimit: getEnvVarAsNumber('QUERY_COMPLEXITY_LIMIT', 1000),
    dataloaderBatchSize: getEnvVarAsNumber('DATALOADER_BATCH_SIZE', 100),
  };
}

// Export singleton config instance
export const config = loadConfig();
