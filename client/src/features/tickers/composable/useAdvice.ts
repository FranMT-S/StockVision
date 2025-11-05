import { getAdviceActionColor, getAdviceActionIcon } from "@/shared/helpers/adviceActions"
import { computed, ref, watch } from "vue"

/** 
 * Use this composable to get the action, icon and color of the advice
 */
export const useAdvice = (initialAction: string) => {
  const action = ref(initialAction.toUpperCase())
  const actionUpper = computed(() => action.value.toUpperCase())

  watch(action, (newValue) => {
    const upper = newValue.toUpperCase()
    if (newValue !== upper) {
      action.value = upper
    }
  })

  const icon = computed(() => {
    return getAdviceActionIcon(actionUpper.value)
  })

  const color = computed(() => {
    return getAdviceActionColor(actionUpper.value)
  })

  return {
    action,
    actionUpper,
    icon,
    color
  }
}
