import { reactive } from 'vue'
import { Anchor } from 'vuetify/lib/types.mjs';

const state = reactive({
  show: false,
  message: '',
  color: 'red-darken-2',
  timeout: 4000,
  location: 'top right' as Anchor
})

function open(message: string, color: string = 'red-darken-2', location: Anchor = 'top right') {
  state.message = message
  state.color = color
  state.location = location
  state.show = true
}

function close() {
  state.show = false
  state.location = 'top right' as Anchor
}

export function useToast() {
  return { state, open, close }
}
