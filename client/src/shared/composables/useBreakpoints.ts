import { ref, onMounted, onUnmounted, computed } from 'vue'

type Breakpoint = 'xs' | 'sm' | 'md' | 'lg' | 'xl'

const breakpoints = {
  xs: 0,      // small phones
  sm: 480,    // large phones
  md: 768,    // tablets
  lg: 1024,   // laptops / small desktop
  xl: 1280,   // large desktop
}

export const useBreakpoints = () => {
  const screenWidth = ref(window.innerWidth)
  const isTouchable = ref(false)

  const updateWidth = () => {
    screenWidth.value = window.innerWidth
  }

  const detectTouch = () => {
    isTouchable.value = 'ontouchstart' in window || navigator.maxTouchPoints > 0
  }

  onMounted(() => {
    window.addEventListener('resize', updateWidth)
    detectTouch()
  })
  
  onUnmounted(() => {
    window.removeEventListener('resize', updateWidth)
  })

  const breakpoint = computed<Breakpoint>(() => {
    const w = screenWidth.value
    if (w < breakpoints.sm) return 'xs'
    if (w < breakpoints.md) return 'sm'
    if (w < breakpoints.lg) return 'md'
    if (w < breakpoints.xl) return 'lg'
    return 'xl'
  })

  const isMobile = computed(() => screenWidth.value < breakpoints.md) // < tablet
  const isTablet = computed(() => screenWidth.value >= breakpoints.md && screenWidth.value < breakpoints.lg)
  const isDesktop = computed(() => screenWidth.value >= breakpoints.lg)
  const isLargeDesktop = computed(() => screenWidth.value >= breakpoints.xl)

  return {
    screenWidth,
    breakpoint,
    isMobile,
    isTablet,
    isDesktop,
    isLargeDesktop,
    isTouchable
  }
}
