import { ref, onMounted, onUnmounted, computed } from 'vue'

type Breakpoint = 'xs' | 'mobile' | 'sm' | 'md' | 'lg' | 'xl'

const breakpoints = {
  xs: 0,
  mobile: 480,
  sm: 600,
  md: 960,
  lg: 1280,
  xl: 1920,
}

export const useBreakpoints = () => {
  const width = ref(window.innerWidth)

  const onResize = () => {
    width.value = window.innerWidth
  }

  onMounted(() => {
    window.addEventListener('resize', onResize)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', onResize)
  })

  const current = computed<Breakpoint>(() => {
    if (width.value < breakpoints.mobile) return 'xs'
    if (width.value < breakpoints.sm) return 'mobile'
    if (width.value < breakpoints.md) return 'sm'
    if (width.value < breakpoints.lg) return 'md'
    if (width.value < breakpoints.xl) return 'lg'
    return 'xl'
  })

  const isMobile = computed(() => current.value === 'xs' || current.value === 'mobile')
  const isTablet = computed(() => current.value === 'sm')
  const isDesktop = computed(() => current.value === 'md' || current.value === 'lg' || current.value === 'xl')

  return {
    width,
    current,
    isMobile,
    isTablet,
    isDesktop,
  }
}
