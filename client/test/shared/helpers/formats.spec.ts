import { describe, it, expect } from 'vitest'
import { humanizeNumberFormat, formatJustDate } from '@/shared/helpers/formats'

describe('formats helper', () => {
  describe('humanizeNumberFormat', () => {
    it('formats numbers less than 1000 correctly', () => {
      expect(humanizeNumberFormat(42)).toBe('42')
      expect(humanizeNumberFormat(999)).toBe('999')
      expect(humanizeNumberFormat(0)).toBe('0')
    })

    it('formats thousands with K suffix', () => {
      expect(humanizeNumberFormat(1500)).toBe('1.5K')
      expect(humanizeNumberFormat(25000, 0)).toBe('25K')
      expect(humanizeNumberFormat(123456, 2)).toBe('123.46K')
    })

    it('formats millions with M suffix', () => {
      expect(humanizeNumberFormat(1500000)).toBe('1.5M')
      expect(humanizeNumberFormat(25000000, 0)).toBe('25M')
      expect(humanizeNumberFormat(123456789, 3)).toBe('123.457M')
    })

    it('formats billions with B suffix', () => {
      expect(humanizeNumberFormat(1500000000)).toBe('1.5B')
      expect(humanizeNumberFormat(25000000000, 0)).toBe('25B'
      )
    })

    it('formats trillions with T suffix when allowed', () => {
      expect(humanizeNumberFormat(1500000000000, 1, true)).toBe('1.5T')
      expect(humanizeNumberFormat(25000000000000, 2, true)).toBe('25.00T')
    })
  })

  describe('formatJustDate', () => {
    it('formats dates in YYYY-MM-DD format', () => {
      expect(formatJustDate(new Date('2023-05-15T12:00:00'))).toBe('2023-05-15')
      expect(formatJustDate(new Date('2023-12-01T00:00:00'))).toBe('2023-12-01')
      expect(formatJustDate(new Date('2000-01-01T23:59:59'))).toBe('2000-01-01')
    })

    it('pads single digit months and days with leading zeros', () => {
      expect(formatJustDate(new Date('2023-01-02T00:00:00'))).toBe('2023-01-02')
      expect(formatJustDate(new Date('2023-11-05T12:34:56'))).toBe('2023-11-05')
    })
  })
})
