import { db, checkDbConnection, disconnectDb } from '../client';

describe('Database Client', () => {
  afterAll(async () => {
    await disconnectDb();
  });

  describe('db singleton', () => {
    it('should export a PrismaClient instance', () => {
      expect(db).toBeDefined();
      expect(typeof db.$connect).toBe('function');
      expect(typeof db.$disconnect).toBe('function');
    });
  });

  describe('checkDbConnection', () => {
    it('should return a boolean', async () => {
      const result = await checkDbConnection();
      expect(typeof result).toBe('boolean');
    });

    // This test will pass when database is running
    it.skip('should return true when database is connected', async () => {
      const result = await checkDbConnection();
      expect(result).toBe(true);
    });
  });

  describe('disconnectDb', () => {
    it('should disconnect without errors', async () => {
      await expect(disconnectDb()).resolves.not.toThrow();
    });
  });
});
