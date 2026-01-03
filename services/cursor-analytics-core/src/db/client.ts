import { PrismaClient } from '../generated/prisma';

/**
 * Database client singleton for cursor-analytics-core
 *
 * This ensures a single PrismaClient instance is used throughout the application
 * to avoid connection pool exhaustion.
 *
 * Usage:
 *   import { db } from './db/client';
 *   const developers = await db.developer.findMany();
 */

// Global type augmentation for development hot-reloading
// Prevents multiple PrismaClient instances in development
declare global {
  var __prismaClient: PrismaClient | undefined;
}

const createPrismaClient = (): PrismaClient => {
  return new PrismaClient({
    log:
      process.env.NODE_ENV === 'development'
        ? ['query', 'error', 'warn']
        : ['error'],
  });
};

// In production, always create a new instance
// In development, reuse global instance to avoid connection exhaustion during hot-reload
export const db = global.__prismaClient ?? createPrismaClient();

if (process.env.NODE_ENV !== 'production') {
  global.__prismaClient = db;
}

/**
 * Gracefully disconnect from database
 * Call this during application shutdown
 */
export const disconnectDb = async (): Promise<void> => {
  await db.$disconnect();
};

/**
 * Health check function to verify database connectivity
 * Returns true if connected, false otherwise
 */
export const checkDbConnection = async (): Promise<boolean> => {
  try {
    await db.$queryRaw`SELECT 1`;
    return true;
  } catch (error) {
    console.error('Database connection check failed:', error);
    return false;
  }
};

/**
 * Refresh the daily_stats materialized view
 * Should be called periodically (e.g., every 5 minutes)
 */
export const refreshDailyStats = async (): Promise<void> => {
  await db.$executeRaw`SELECT refresh_daily_stats()`;
};
