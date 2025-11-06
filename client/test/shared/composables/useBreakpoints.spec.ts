import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { useBreakpoints } from '@/shared/composables/useBreakpoints'
import { flushPromises, mount } from '@vue/test-utils'
 

describe('useBreakpoints', () => {
  const originalInnerWidth = window.innerWidth
  const originalAddEventListener = window.addEventListener
  const originalRemoveEventListener = window.removeEventListener

  beforeEach(() => {
    // Mock window methods
    window.addEventListener = vi.fn()
    window.removeEventListener = vi.fn()
    vi.spyOn(window, 'addEventListener')
    vi.spyOn(window, 'removeEventListener')
    window.innerWidth = 1024 // Default to desktop
  })

  afterEach(() => {
    // Restore original implementations
    window.addEventListener = originalAddEventListener
    window.removeEventListener = originalRemoveEventListener
    window.innerWidth = originalInnerWidth
    vi.restoreAllMocks()
  })

  it('initializes with current window width', () => {
    const { screenWidth } = useBreakpoints()
    expect(screenWidth.value).toBe(1024)
  })


  describe('breakpoint detection', () => {
    const testCases = [
      { width: 400, expected: 'xs', desc: 'extra small screens (mobile)' },
      { width: 500, expected: 'sm', desc: 'small screens (large phones)' },
      { width: 800, expected: 'md', desc: 'medium screens (tablets)' },
      { width: 1100, expected: 'lg', desc: 'large screens (laptops)' },
      { width: 1400, expected: 'xl', desc: 'extra large screens (desktops)' },
    ]

    testCases.forEach(({ width, expected, desc }) => {
      it(`detects ${desc} (${width}px) as '${expected}'`, () => {
        window.innerWidth = width
        const { breakpoint } = useBreakpoints()
        expect(breakpoint.value).toBe(expected)
      })
    })
  })

  describe('device type helpers', () => {

    // mount the component to call onMounted
    const mountComposable = () => {
      return mount({
        setup() {
          return useBreakpoints()
        },
        template: '<div></div>'
      })
    }

    it('detects mobile devices', () => {
      window.innerWidth = 700 // Below md breakpoint (768)
      const wrapper = mountComposable()
      expect(wrapper.vm.screenWidth).toBe(700)
      expect(wrapper.vm.isMobile).toBe(true)
    })

    it('detects tablet devices', () => {
      window.innerWidth = 900 // Between md and lg breakpoints
      const wrapper = mountComposable()
      expect(wrapper.vm.isTablet).toBe(true)
    })

    it('detects desktop devices', () => {
      window.innerWidth = 1200 // Above lg breakpoint (1024)
      const wrapper = mountComposable()
      expect(wrapper.vm.isDesktop).toBe(true)
    })

    it('detects large desktop devices', () => {
      window.innerWidth = 1440 // Above xl breakpoint (1280)
      const wrapper = mountComposable()
      expect(wrapper.vm.isLargeDesktop).toBe(true)
    })

    it('updates screen width on window resize', () => {
      vi.spyOn(window, 'addEventListener')
      vi.spyOn(window, 'removeEventListener')

      window.innerWidth = 1024
      const wrapper = mountComposable()

      expect(window.addEventListener).toHaveBeenCalledWith('resize', expect.any(Function))
      const resizeHandler = (window.addEventListener as any).mock.calls[0][1]
 

      window.innerWidth = 600
      resizeHandler()

      expect(wrapper.vm.screenWidth).toBe(600)
    })
  })


  it('detects touch devices',  () => {
    // Mock touch support
    Object.defineProperty(window, 'ontouchstart', {
      value: {},
      writable: true,
      configurable: true
    })

   const wrapper = mount({
    setup() {
      return useBreakpoints()
    },
    template: '<div></div>'
  })

  const { isTouchable } = wrapper.vm
  expect(isTouchable).toBe(true)
  })
})
