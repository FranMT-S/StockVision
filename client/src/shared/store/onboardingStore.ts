import { defineStore } from "pinia";
import { ref } from "vue";
import { customFetch } from "../services/customFetch";
import { API_CONFIG } from "../constants/api";
import { Onboarding } from "../models/onboarding";

export const useOnboardingStore = defineStore('onboarding', () => {
  const overviewData = ref<Onboarding | null>(null)
  const abortControllerUpdate = ref<AbortController | null>(null)
  const abortControllerGet = ref<AbortController | null>(null)

  /**
   * update the onboarding data of the user
   * 
   * @param step 
   * @param done 
   */
  const updateOnboarding = async (step: number, done: boolean) => {
    if(abortControllerUpdate.value && !abortControllerUpdate.value.signal.aborted){
      abortControllerUpdate.value.abort()
    }
    
    abortControllerUpdate.value = new AbortController();
    const res = await customFetch<Onboarding>(API_CONFIG.ENDPOINTS.Onboarding(), {
      method: 'PATCH',
      body: JSON.stringify({ overviewStep: step, overviewDone: done }),
      signal: abortControllerUpdate.value.signal
    })
    
    if (!res.ok) {
      console.error(res.error)
      return
    }

    overviewData.value = res.data
  }

  /**
   * fetch the onboarding data from the server
   */
  const fetchOnboarding = async () => {
   
    if(abortControllerGet.value && !abortControllerGet.value.signal.aborted){
      abortControllerGet.value.abort()
    }

    abortControllerGet.value = new AbortController();
    const response = await customFetch<Onboarding>(API_CONFIG.ENDPOINTS.Onboarding(), {
      method: 'GET',
      signal: abortControllerGet.value.signal
    })

    if (!response.ok) {
      // set default to avoid show overview tour infinite
      console.error(response.error)
      overviewData.value = { id: 1, overviewStep: 1, overviewDone: true }
      return
    }

    overviewData.value = response.data
  }

  return {
    overviewData,
    updateOnboarding,
    fetchOnboarding
  }
})
