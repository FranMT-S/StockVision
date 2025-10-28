import { beforeEach } from 'vitest'
import { setActivePinia } from 'pinia'
import { createPinia } from 'pinia'

beforeEach(() => {
  setActivePinia(createPinia())
})