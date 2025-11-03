import vuetify from '@/plugins/vuetify'
import ErrorFetchData from '@/shared/components/ErrorFetchData.vue'
import { mount,} from '@vue/test-utils'
import { describe, it, expect } from 'vitest'

describe("ErrorComponent", () => {
  const errorMessage = "test error"
  const wrapper = mount(ErrorFetchData, {
    global: {
      plugins: [vuetify], 
    },
    props: {
      error: errorMessage,
    },
  })

  it("Must exist an error message", () => {  
    expect(wrapper.text()).toContain(errorMessage)
  })

  it("Must be have a button try again", () => {
    expect(wrapper.find('[data-testid="try-again-button"]').text()).toBe("Try Again")
  })

  it("Must emit click event when button is clicked", () => {
    const button = wrapper.find('[data-testid="try-again-button"]')
    button.trigger('click')
    expect(wrapper.emitted().click).toBeTruthy()
  })
})