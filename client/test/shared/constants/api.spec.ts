

import { API_CONFIG,API_URL } from '@/shared/constants/api';
import { describe, it, expect } from 'vitest';

describe("Api Routes", () => {
  
  const apiConfig = API_CONFIG
  
  describe("List Ticker Route", () => {
    it("should have default values", () => {
      expect(apiConfig.ENDPOINTS.List()).toBe(API_URL +'/api/v1/tickers?page=1&size=10&sort=asc')
    })
    it("should have the correct endpoints", () => {
      expect(apiConfig.ENDPOINTS.List({q:'test', page:10, size:20, sort:'desc'})).toBe(API_URL + '/api/v1/tickers?page=10&size=20&sort=desc&q=test')
    })
    it("shoul size must be 10 and page 1", () => {
      expect(apiConfig.ENDPOINTS.List({q:'', page:-1, size:-1, sort:'desc'})).toBe(API_URL + '/api/v1/tickers?page=1&size=10&sort=desc')
    })

    it("should be default values if q is undefined", () => {
      expect(apiConfig.ENDPOINTS.List({q:undefined, page:undefined, size:undefined, sort:'desc'})).toBe(API_URL + '/api/v1/tickers?page=1&size=10&sort=desc')
    })

    it("should be default values if sort is undefined and q=test", () => {
      expect(apiConfig.ENDPOINTS.List({q:'test', page:undefined, size:undefined, sort:undefined})).toBe(API_URL + '/api/v1/tickers?page=1&size=10&sort=asc&q=test')
    })
  })

  describe("Overview Route", () => {
    it("should be default values if sort is undefined and q=test", () => {
      expect(apiConfig.ENDPOINTS.Overview('test',new Date('2025-10-27:00:00:00'))).toBe(API_URL + '/api/v1/tickers/test/overview?from=2025-10-27')
    })
  })

  describe("Logo Route", () => {
    it("should be default values if sort is undefined and q=test", () => {
      expect(apiConfig.ENDPOINTS.Logo('AAPL')).toBe(API_URL + '/api/v1/tickers/AAPL/logo')
    })
  })

  describe("Predictions Route", () => {
    it("should be default values if sort is undefined and q=test", () => {
      expect(apiConfig.ENDPOINTS.Predictions('AAPL')).toBe(API_URL + '/api/v1/tickers/AAPL/predictions')
    })
  })

  describe("HistoricalPrices Route", () => {
    it("should be default values if sort is undefined and q=test", () => {
      expect(apiConfig.ENDPOINTS.HistoricalPrices('AAPL')).toBe(API_URL + '/api/v1/tickers/AAPL/historical')
    })

    it("should be default values if sort is undefined and q=test", () => {
      expect(apiConfig.ENDPOINTS.HistoricalPrices('AAPL',new Date('2025-10-27:00:00:00'))).toBe(API_URL + '/api/v1/tickers/AAPL/historical?from=2025-10-27')
    })
  })
});
