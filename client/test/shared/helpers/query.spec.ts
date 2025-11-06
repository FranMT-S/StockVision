import { describe, it, expect, vi, beforeEach } from 'vitest'
import { normalizePageNumber, normalizePageSize, isValidQuery, buildQueryString } from '@/shared/helpers/query'

describe('query helper', () => {
  describe('normalizePageNumber', () => {
    it('returns default page number (1) for invalid inputs', () => {
      expect(normalizePageNumber(undefined)).toBe(1)
      expect(normalizePageNumber(0)).toBe(1)
      expect(normalizePageNumber(-1)).toBe(1)
      expect(normalizePageNumber(NaN)).toBe(1)
    })

    it('returns the provided page number when valid', () => {
      expect(normalizePageNumber(1)).toBe(1)
      expect(normalizePageNumber(5)).toBe(5)
      expect(normalizePageNumber(100)).toBe(100)
    })
  })

  describe('normalizePageSize', () => {
    it('returns default page size for invalid inputs', () => {
      // Default size is typically 10, but we'll use the actual constant
      const defaultSize = 10 // This should match your DEFAULT_PAGINATION.SIZE
      expect(normalizePageSize(undefined)).toBe(defaultSize)
      expect(normalizePageSize(0)).toBe(defaultSize)
      expect(normalizePageSize(-1)).toBe(defaultSize)
      expect(normalizePageSize(NaN)).toBe(defaultSize)
    })

    it('returns the provided page size when valid', () => {
      expect(normalizePageSize(5)).toBe(5)
      expect(normalizePageSize(20)).toBe(20)
      expect(normalizePageSize(100)).toBe(100)
    })
  })

  describe('isValidQuery', () => {
    it('returns false for invalid queries', () => {
      expect(isValidQuery(undefined)).toBe(false)
      expect(isValidQuery('')).toBe(false)
      expect(isValidQuery('   ')).toBe(false)
      expect(isValidQuery('undefined')).toBe(false)
    })

    it('returns true for valid queries', () => {
      expect(isValidQuery('test')).toBe(true)
      expect(isValidQuery('  test  ')).toBe(true)
      expect(isValidQuery('test query')).toBe(true)
    })
  })

  describe('buildQueryString', () => {
    it('returns empty string for empty params', () => {
      expect(buildQueryString({})).toBe('')
    })

    it('builds query string correctly', () => {
      const params = {
        page: 1,
        size: 10,
        query: 'test',
        sort: 'name',
        order: 'asc'
      }
      const result = buildQueryString(params)
      
      // The order of parameters might vary, so we'll check if all key-value pairs are present
      expect(result).toContain('page=1')
      expect(result).toContain('size=10')
      expect(result).toContain('query=test')
      expect(result).toContain('sort=name')
      expect(result).toContain('order=asc')
      expect(result.startsWith('?')).toBe(true)
      expect(result.split('&').length).toBe(5)
    })

    it('handles special characters in values', () => {
      const params = {
        search: 'test@example.com',
        filter: 'status=active&type=user'
      }
      const result = buildQueryString(params)
      
      expect(result).toContain('search=test%40example.com')
      expect(result).toContain('filter=status%3Dactive%26type%3Duser')
    })
  })
})
