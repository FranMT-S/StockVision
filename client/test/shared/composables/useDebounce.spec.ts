import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { useDebounce } from '@/shared/composables/useDebounce'

describe('useDebounce', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.restoreAllMocks()
    vi.clearAllTimers()
  })

  it('calls the function after the specified delay', () => {
    const { debounced } = useDebounce()
    const mockFn = vi.fn()
    const delay = 1000

    debounced(mockFn, delay)
    
    // Should not be called immediately
    expect(mockFn).not.toHaveBeenCalled()
    
    // Fast-forward time by delay
    vi.advanceTimersByTime(delay)
    
    // Should be called after delay
    expect(mockFn).toHaveBeenCalledTimes(1)
  })

  it('cancels previous calls when called multiple times', () => {
    const { debounced } = useDebounce()
    const mockFn = vi.fn()
    const delay = 1000

    // First call
    debounced(mockFn, delay)
    
    // Second call before first delay completes
    debounced(mockFn, delay)
    
    // Fast-forward time by delay
    vi.advanceTimersByTime(delay)
    
    // Should only be called once
    expect(mockFn).toHaveBeenCalledTimes(1)
  })

  it('can cancel pending debounced calls', () => {
    const { debounced, cancel } = useDebounce()
    const mockFn = vi.fn()
    const delay = 1000

    debounced(mockFn, delay)
    cancel()
    
    // Fast-forward time by delay
    vi.advanceTimersByTime(delay)
    
    // Should not be called after cancel
    expect(mockFn).not.toHaveBeenCalled()
  })

  it('handles multiple independent instances', () => {
    const debouncer1 = useDebounce()
    const debouncer2 = useDebounce()
    
    const mockFn1 = vi.fn()
    const mockFn2 = vi.fn()
    const delay = 1000

    debouncer1.debounced(mockFn1, delay)
    debouncer2.debounced(mockFn2, delay)
    
    // Cancel first instance
    debouncer1.cancel()
    
    // Fast-forward time by delay
    vi.advanceTimersByTime(delay)
    
    // Only second function should be called
    expect(mockFn1).not.toHaveBeenCalled()
    expect(mockFn2).toHaveBeenCalledTimes(1)
  })

  it('passes the correct arguments to the debounced function', () => {
    const { debounced } = useDebounce()
    const mockFn = vi.fn()
    const delay = 1000
    const testArg = 'test argument'

    const debouncedFn = () => mockFn(testArg)
    debounced(debouncedFn, delay)
    
    vi.advanceTimersByTime(delay)
    
    expect(mockFn).toHaveBeenCalledWith(testArg)
  })
})
