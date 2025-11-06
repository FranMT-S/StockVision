import { describe, it, expect } from 'vitest'
import { getAdviceActionIcon, getAdviceActionColor } from '@/shared/helpers/adviceActions'

describe('adviceActions helper', () => {
  describe('getAdviceActionIcon', () => {
    it('returns correct icons for each action type', () => {
      expect(getAdviceActionIcon('BUY')).toBe('mdi-cart-plus')
      expect(getAdviceActionIcon('SELL')).toBe('mdi-trending-down')
      expect(getAdviceActionIcon('HOLD')).toBe('mdi-hand-back-right')
    })

    it('returns empty string for unknown action types', () => {
      expect(getAdviceActionIcon('UNKNOWN')).toBe('')
      expect(getAdviceActionIcon('')).toBe('')
      // @ts-ignore - Testing invalid input
      expect(getAdviceActionIcon(undefined)).toBe('')
    })
  })

  describe('getAdviceActionColor', () => {
    it('returns correct colors for each action type', () => {
      expect(getAdviceActionColor('BUY')).toEqual({
        class: 'tw-text-green-500',
        vuetify: 'green'
      })
      expect(getAdviceActionColor('SELL')).toEqual({
        class: 'tw-text-red-500',
        vuetify: 'red'
      })
      expect(getAdviceActionColor('HOLD')).toEqual({
        class: 'tw-text-gray-500',
        vuetify: 'gray'
      })
    })

    it('returns empty values for unknown action types', () => {
      const expected = { class: '', vuetify: '' }
      expect(getAdviceActionColor('UNKNOWN')).toEqual(expected)
      expect(getAdviceActionColor('')).toEqual(expected)
      // @ts-ignore - Testing invalid input
      expect(getAdviceActionColor(undefined)).toEqual(expected)
    })
  })
})
