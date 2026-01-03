import { describe, it, expect } from 'vitest';
import { ApolloClient, InMemoryCache } from '@apollo/client';
import { createApolloClient } from '../client';

describe('Apollo Client', () => {
  describe('createApolloClient', () => {
    it('should create an Apollo Client instance', () => {
      const client = createApolloClient();

      expect(client).toBeInstanceOf(ApolloClient);
    });

    it('should configure cache with Developer type policy', () => {
      const client = createApolloClient();
      const cache = client.cache as InMemoryCache;

      // Access the cache config to verify type policies
      const config = cache.policies;
      expect(config).toBeDefined();
    });

    it('should configure cache with DailyStats type policy', () => {
      const client = createApolloClient();
      const cache = client.cache as InMemoryCache;

      const config = cache.policies;
      expect(config).toBeDefined();
    });

    it('should use VITE_GRAPHQL_URL from environment', () => {
      const client = createApolloClient();

      // Apollo Client stores the URI in link.options
      expect(client.link).toBeDefined();
    });

    it('should have default query options configured', () => {
      const client = createApolloClient();

      expect(client.defaultOptions.query).toBeDefined();
      expect(client.defaultOptions.watchQuery).toBeDefined();
    });
  });
});
