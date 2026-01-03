import { loadConfig } from './index';

describe('Config', () => {
  describe('loadConfig', () => {
    const originalEnv = process.env;

    beforeEach(() => {
      // Reset environment before each test
      jest.resetModules();
      process.env = { ...originalEnv };
    });

    afterAll(() => {
      process.env = originalEnv;
    });

    it('should load default configuration values', () => {
      const config = loadConfig();

      expect(config.port).toBe(4000);
      expect(config.nodeEnv).toBe('test'); // Jest sets NODE_ENV to 'test'
      expect(config.simulatorUrl).toBe('http://localhost:8080');
      expect(config.pollIntervalMs).toBe(60000);
      expect(config.logLevel).toBe('info');
      expect(config.queryComplexityLimit).toBe(1000);
      expect(config.dataloaderBatchSize).toBe(100);
    });

    it('should override defaults with environment variables', () => {
      process.env.PORT = '5000';
      process.env.NODE_ENV = 'production';
      process.env.SIMULATOR_URL = 'http://example.com:9000';
      process.env.POLL_INTERVAL_MS = '30000';

      const config = loadConfig();

      expect(config.port).toBe(5000);
      expect(config.nodeEnv).toBe('production');
      expect(config.simulatorUrl).toBe('http://example.com:9000');
      expect(config.pollIntervalMs).toBe(30000);
    });

    it('should throw error for invalid numeric environment variables', () => {
      process.env.PORT = 'invalid';

      expect(() => loadConfig()).toThrow('Environment variable PORT must be a valid number');
    });
  });
});
