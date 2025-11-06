import { describe, it, expect, beforeEach } from 'vitest'
import { useToast } from '@/shared/composables/useToast'
import { debounce } from 'vuetify/lib/util/helpers.mjs'

describe('useToast', () => {
  let toast: ReturnType<typeof useToast>

  beforeEach(() => {
    // Create a fresh instance for each test
    toast = useToast()
    // Reset the state
    toast.close()
  })

  it('initializes with default values', () => {
    expect(toast.state).toEqual({
      show: false,
      message: '',
      color: 'red-darken-2',
      timeout: 4000,
      location: 'top right'
    })
  })

  it('opens a toast with default values', () => {
    const message = 'Test message'
    toast.open(message)

    expect(toast.state).toMatchObject({
      show: true,
      message,
      color: 'red-darken-2',
      location: 'top right'
    })
  })

  it('opens a toast with custom color and location', () => {
    const message = 'Success!'
    const color = 'green-darken-2'
    const location = 'bottom right'

    toast.open(message, color, location as any)

    expect(toast.state).toMatchObject({
      show: true,
      message,
      color,
      location
    })
  })

  it('closes the toast', () => {
    // First open the toast
    toast.open('Test message')
    expect(toast.state.show).toBe(true)

    // Then close it
    toast.close()
    
    // Check if it's closed and reset to default location
    expect(toast.state.show).toBe(false)
    expect(toast.state.location).toBe('top right')
  })

  it('handles multiple open and close operations', () => {
    // First open
    toast.open('First message', 'blue')
    expect(toast.state.show).toBe(true)
    expect(toast.state.color).toBe('blue')

    // Close
    toast.close()
    expect(toast.state.show).toBe(false)

    // Open again with different values
    toast.open('Second message', 'green', 'bottom left' as any)
    expect(toast.state).toMatchObject({
      show: true,
      message: 'Second message',
      color: 'green',
      location: 'bottom left'
    })
  })
})
