import { reactive } from 'vue'
import { Anchor } from 'vuetify/lib/types.mjs';


const defaultState = {
  show: false,
  message: '',
  color: 'red-darken-2',
  timeout: 4000,
  location: 'top right' as Anchor
}

const state = reactive(defaultState)

/** 
 * Use this composable to show a toast message, toast is a global a unique instance for all aplication
 * 
 * @param message - The message to show
 * @param color - The color of the toast
 * @param location - The location of the toast
 */

function open(message: string, color: string = 'red-darken-2', location: Anchor = 'top right') {
  state.message = message
  state.color = color
  state.location = location
  state.show = true
}

function close() {
  state.show = false
  state.location = defaultState.location
  state.color = defaultState.color
  state.message = defaultState.message
  state.timeout = defaultState.timeout
  state.show = defaultState.show
}

export function useToast() {
  return { state, open, close }
}
