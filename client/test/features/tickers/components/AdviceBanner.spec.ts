import AdviceBanner from '@/features/tickers/components/AdviceBanner.vue'
import vuetify from '@/plugins/vuetify'
import { mount } from '@vue/test-utils'
import { useAdvice } from '@/features/tickers/composable/useAdvice'
import { computed, ref } from 'vue'

describe('renders advice correctly', () => {

  vi.mock('@/features/tickers/composable/useAdvice', () => ({
    useAdvice: vi.fn() // mock "vacío" inicialmente
  }))

  const cases = [
    {
      advice: 'BUY.Buy Stock',
      tittle: 'Advice',
      icon: 'mdi-cart-plus',
      color: 'green',
      expectedAction: 'BUY',
      expectedDescription: 'Buy Stock',
      mockUseAdvice: {
        action: ref('BUY'),
        icon: computed<"mdi-cart-plus" | "mdi-trending-down" | "mdi-hand-back-right" | "">(() => 'mdi-cart-plus'),
        color: computed(() => { return { class:'tw-text-green-500', vuetify:'green' }}),
      }
    },
    {
      advice: 'SELL.Sell Stock',
      tittle: 'Advice',
      icon: 'mdi-trending-down',
      color: 'red',
      expectedAction: 'SELL',
      expectedDescription: 'Sell Stock',
      mockUseAdvice: {
        action: ref('SELL'),
        icon: computed<"mdi-cart-plus" | "mdi-trending-down" | "mdi-hand-back-right" | "">(() => 'mdi-trending-down'),
        color: computed(() => { return { class:'tw-text-red-500', vuetify:'red' }})
      }
    },
    {
      advice: 'HOLD.Hold Stock',
      tittle: 'Advice',
      icon: 'mdi-hand-back-right',
      color: 'gray',
      expectedAction: 'HOLD',
      expectedDescription: 'Hold Stock',
      mockUseAdvice: {
        action: ref('HOLD'),
        icon: computed<"mdi-cart-plus" | "mdi-trending-down" | "mdi-hand-back-right" | "">(() => 'mdi-hand-back-right'),
        color: computed(() => { return { class:'tw-text-gray-500', vuetify:'gray' }})
      }
    }
  ]

  for(const ca of cases){
    vi.mocked(useAdvice).mockImplementationOnce((_: string) => {
      return {
        ...ca.mockUseAdvice,
        actionUpper: computed(() => ca.mockUseAdvice.action.value.toUpperCase())
      }
    })

    const wrapper = mount(AdviceBanner, {
      props: {
        advice: ca.advice,
        tittle: ca.tittle
      },
      global:{
        plugins:[vuetify]
      }
    })
    
    it(`test ${ca.advice}`, () => {
      expect(wrapper.find('[data-test="action"]').text()).toContain(ca.expectedAction)
      expect(wrapper.find('[data-test="description"]').text()).toContain(ca.expectedDescription)
      expect(wrapper.find('[data-test="icon"]').classes()).toContain(ca.icon)
    })

    vi.clearAllMocks()
  }
})
