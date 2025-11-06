import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import CrossHairDetails from '@/features/tickers/components/CrossHairDetails.vue'
import type { StockHLOC } from '@/shared/models/recomendations'

describe('CrossHairDetails.vue', () => {
  const vuetify = createVuetify({ components, directives })
  
  const mockStockData: StockHLOC = {
    date: new Date('2023-11-05'),
    open: 150.25,
    high: 152.50,
    low: 149.75,
    close: 151.80,
    volume: 1200000,
    change: 2.5,
    changePercent: 1.67,
    time: '2023-11-05',
    vwap: 151.80
  }

  
  const createWrapper = (props = {}) => {
    return mount(CrossHairDetails, {
      props: {
        data: { ...mockStockData, ...props },
      },
      global: {
        plugins: [vuetify],
      },
    })
  }

  
  it('renders the component with correct icon for positive change', () => {
    const wrapper = createWrapper()
    const icon = wrapper.find('i')
    expect(icon.classes()).toContain('mdi-trending-up')
  })

  it('renders the component with correct icon for negative change', () => {
    const wrapper = createWrapper({ change: -1.5, changePercent: -1.0 })
    const icon = wrapper.find('i')
    expect(icon.classes()).toContain('mdi-trending-down')
  })

  it('displays all price fields with correct formatting', () => {
    const wrapper = createWrapper()
    const fields = wrapper.findAll('.tw-flex-row')
    
    expect(fields.length).toBe(7) // volumen, open, high, low, close, change, changePercent
  })

  it('applies correct text color class based on price change', () => {
    const wrapper = createWrapper()
    const positiveChangeElements = wrapper.findAll('.tw-text-green-400')
    
    // Should have at least one element with positive change class
    expect(positiveChangeElements.length).toBeGreaterThan(0)
    
    // Test negative change
    const negativeWrapper = createWrapper({ change: -1.5, changePercent: -1.0 })
    const negativeChangeElements = negativeWrapper.findAll('.tw-text-red-400')
    
    // Should have at least one element with negative change class
    expect(negativeChangeElements.length).toBeGreaterThan(0)
  })

  it('formats large numbers correctly', () => {
    const wrapper = createWrapper({
      volume: 1500000,
      high: 1000000,
      low: 500000,
      open: 750000,
      close: 800000,
      change: 50000,
      changePercent: 6.67
    })
    
    const fields = wrapper.findAll('.tw-flex-row')
    const fieldTexts = fields.map(field => field.text())
    
    expect(fieldTexts[0]).toContain('V:1.50M')
    expect(fieldTexts[1]).toContain('O:750.00K$')
    expect(fieldTexts[2]).toContain('C:800.00K$')
    expect(fieldTexts[3]).toContain('H:1.00M$')
    expect(fieldTexts[4]).toContain('L:500.00K$')
    expect(fieldTexts[5]).toContain('CH:50.00K')
    expect(fieldTexts[6]).toContain('%CH:6.67%')
  })
})
